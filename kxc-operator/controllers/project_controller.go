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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

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
	found := &corev1.Namespace{}
	err = r.Get(ctx, types.NamespacedName{Name: ProjectNamespaceName(project.Name)}, found)
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

	return ctrl.Result{}, nil
}

// namespaceForProject returns a Namespace object
func (r *ProjectReconciler) namespaceForProject(project *cloudv1alpha1.Project) (*corev1.Namespace, error) {
	ls := labelsForProject(project.Name)

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   ProjectNamespaceName(project.Name),
			Labels: ls,
		},
	}
	// Set Project instance as the owner and controller
	err := ctrl.SetControllerReference(project, namespace, r.Scheme)
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

// NamespacePrefix is the prefix for all namespaces created for projects
const projectNamespacePrefix = "kxc-proj-"

// ProjectNamespaceName returns the namespaceName.
func ProjectNamespaceName(projectName string) string {
	return projectNamespacePrefix + projectName
}

// labelsForProject returns the labels for selecting the resources
// belonging to the given project CR name.
func labelsForProject(name string) map[string]string {
	return map[string]string{"app": "kxc_project", "kxc_project_cr": name}
}

func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha1.Project{}).
		Owns(&corev1.Namespace{}).
		Complete(r)
}
