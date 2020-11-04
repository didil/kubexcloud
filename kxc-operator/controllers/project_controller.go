/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumslimmetav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	ciliumpolicyapi "github.com/cilium/cilium/pkg/policy/api"
	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloud.kubexcloud.com,resources=projects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloud.kubexcloud.com,resources=projects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete

func (r *ProjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("project", req.NamespacedName)

	// Fetch the Project instance
	project := &cloudv1alpha1.Project{}
	err := r.Get(ctx, req.NamespacedName, project)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Project resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Project")
		return ctrl.Result{}, err
	}

	// Check if the namespace already exists, if not create a new one
	namespace := &corev1.Namespace{}
	err = r.Get(ctx, types.NamespacedName{Name: ProjectNamespaceName(project.Name)}, namespace)
	if err != nil && errors.IsNotFound(err) {
		// Define new
		namespace, err := r.namespaceForProject(project)
		if err != nil {
			log.Error(err, "Failed to build new Namespace", "Namespace.Name", namespace.Name)
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Namespace", "Namespace.Name", namespace.Name)
		err = r.Create(ctx, namespace)
		if err != nil {
			log.Error(err, "Failed to create new Namespace", "Namespace.Name", namespace.Name)
			return ctrl.Result{}, err
		}

		// created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Namespace")
		return ctrl.Result{}, err
	}

	// create cilium network policy
	policy := &ciliumv2.CiliumNetworkPolicy{}
	err = r.Get(ctx, types.NamespacedName{Name: ciliumNetworkPolicyName, Namespace: namespace.Name}, policy)
	if err != nil && errors.IsNotFound(err) {
		policy, err := r.ciliumNetworkPolicyForProject(project)
		if err != nil {
			log.Error(err, "Failed to build new CiliumNetworkPolicy", "CiliumNetworkPolicy.Namespace", namespace.Name, "CiliumNetworkPolicy.Name", ciliumNetworkPolicyName)
			return ctrl.Result{}, err
		}

		log.Info("Creating a new CiliumNetworkPolicy", "CiliumNetworkPolicy.Namespace", namespace.Name, "CiliumNetworkPolicy.Name", ciliumNetworkPolicyName)
		err = r.Create(ctx, policy)
		if err != nil {
			log.Error(err, "Failed to create new CiliumNetworkPolicy", "CiliumNetworkPolicy.Namespace", namespace.Name, "CiliumNetworkPolicy.Name", ciliumNetworkPolicyName)
			return ctrl.Result{}, err
		}

		// created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get CiliumNetworkPolicy")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// namespaceForProject returns a Namespace object
func (r *ProjectReconciler) namespaceForProject(project *cloudv1alpha1.Project) (*corev1.Namespace, error) {
	labels := LabelsForNamespace(project.Name)

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   ProjectNamespaceName(project.Name),
			Labels: labels,
		},
	}
	// Set Project instance as the owner and controller
	err := ctrl.SetControllerReference(project, namespace, r.Scheme)
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

func ingressControllerNamespaceName() string {
	n := os.Getenv("INGRESS_NAMESPACE")
	if n == "" {
		n = "kube-system"
	}

	return n
}

const ciliumNetworkPolicyName = "allow-from-same-namespace-and-ingress"

func (r *ProjectReconciler) ciliumNetworkPolicyForProject(project *cloudv1alpha1.Project) (*ciliumv2.CiliumNetworkPolicy, error) {
	namespaceName := ProjectNamespaceName(project.Name)
	ingressControllerNamespaceName := ingressControllerNamespaceName()

	policy := &ciliumv2.CiliumNetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ciliumNetworkPolicyName,
			Namespace: namespaceName,
		},
		Specs: ciliumpolicyapi.Rules{
			&ciliumpolicyapi.Rule{
				EndpointSelector: ciliumpolicyapi.EndpointSelector{
					LabelSelector: &ciliumslimmetav1.LabelSelector{
						MatchLabels: map[string]string{},
					},
				},
				// allow ingress from same namesapce
				Ingress: []ciliumpolicyapi.IngressRule{
					ciliumpolicyapi.IngressRule{
						FromEndpoints: []ciliumpolicyapi.EndpointSelector{
							ciliumpolicyapi.NewESFromK8sLabelSelector("k8s.", &ciliumslimmetav1.LabelSelector{
								MatchLabels: map[string]string{
									"io.kubernetes.pod.namespace": namespaceName,
								},
							}),
						},
					},
				},
			},
			&ciliumpolicyapi.Rule{
				EndpointSelector: ciliumpolicyapi.EndpointSelector{
					LabelSelector: &ciliumslimmetav1.LabelSelector{
						MatchLabels: map[string]string{},
					},
				},
				// allow ingress from ingress namespace
				Ingress: []ciliumpolicyapi.IngressRule{
					ciliumpolicyapi.IngressRule{
						FromEndpoints: []ciliumpolicyapi.EndpointSelector{
							ciliumpolicyapi.NewESFromK8sLabelSelector("k8s.", &ciliumslimmetav1.LabelSelector{
								MatchLabels: map[string]string{
									"io.kubernetes.pod.namespace": ingressControllerNamespaceName,
								},
							}),
						},
					},
				},
			},
		},
	}

	// Set Project instance as the owner and controller
	err := ctrl.SetControllerReference(project, policy, r.Scheme)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// projectNamespacePrefix is the prefix for all namespaces created for projects
const projectNamespacePrefix = "kxc-proj-"

// ProjectNamespaceName returns the namespaceName.
func ProjectNamespaceName(projectName string) string {
	return projectNamespacePrefix + projectName
}

// LabelsForNamespace returns the labels for a namespace
func LabelsForNamespace(projectName string) map[string]string {
	return map[string]string{"app": "kxc", projectCRKey: projectName}
}

const projectCRKey = "project_cr"

const userAccountCRKey = "user_account_cr"

// LabelsForProject returns the labels for a project
func LabelsForProject(userName string) map[string]string {
	return map[string]string{"app": "kxc", userAccountCRKey: userName}
}

// ProjectUserName returns the username for a project
func ProjectUserName(project *cloudv1alpha1.Project) string {
	if project == nil {
		return ""
	}

	return project.Labels[userAccountCRKey]
}

func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha1.Project{}).
		Owns(&corev1.Namespace{}).
		Complete(r)
}
