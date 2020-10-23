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
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloud.kubexcloud.com,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloud.kubexcloud.com,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;

func (r *AppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("app", req.NamespacedName)

	// Fetch the App instance
	app := &cloudv1alpha1.App{}
	err := r.Get(ctx, req.NamespacedName, app)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("App resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get app")
		return ctrl.Result{}, err
	}

	// Check if the namespace exists
	namespace := &corev1.Namespace{}
	err = r.Get(ctx, types.NamespacedName{Name: ProjectNamespaceName(AppProjectName(app))}, namespace)
	if err != nil {
		log.Error(err, "Failed to get Namespace")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	dep := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: app.Name, Namespace: namespace.Name}, dep)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForApp(app)
		if err != nil {
			log.Error(err, "Failed to build new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	targetDep, err := r.deploymentForApp(app)
	if err != nil {
		log.Error(err, "Failed to build target Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		return ctrl.Result{}, err
	}

	// Ensure the deployment replicas are the same
	if *dep.Spec.Replicas != *targetDep.Spec.Replicas {
		log.Info("Updating deployment size",
			"Deployment.Namespace", dep.Namespace,
			"Deployment.Name", dep.Name,
			"old", *dep.Spec.Replicas,
			"new", *targetDep.Spec.Replicas)

		dep.Spec.Replicas = targetDep.Spec.Replicas
		err = r.Update(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// ensure the deployment containers are the same
	if !r.containersEqual(dep.Spec.Template.Spec.Containers, targetDep.Spec.Template.Spec.Containers) {
		dep.Spec.Template.Spec.Containers = targetDep.Spec.Template.Spec.Containers

		log.Info("Updating deployment containers",
			"Deployment.Namespace", dep.Namespace,
			"Deployment.Name", dep.Name)

		err = r.Update(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to update Deployment containers", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Check if the service already exists, if not create a new one
	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: app.Name, Namespace: namespace.Name}, svc)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc, err := r.serviceForApp(app)
		if err != nil {
			log.Error(err, "Failed to build new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}

		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	targetSvc, err := r.serviceForApp(app)
	if err != nil {
		log.Error(err, "Failed to build target Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		return ctrl.Result{}, err
	}

	// check service ports are equal
	if !r.servicePortsEqual(svc.Spec.Ports, targetSvc.Spec.Ports) {
		svc.Spec.Ports = targetSvc.Spec.Ports

		log.Info("Updating service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)

		err = r.Update(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to update Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// update status if necessary
	if app.Status.AvailableReplicas != dep.Status.AvailableReplicas || app.Status.UnavailableReplicas != dep.Status.UnavailableReplicas {
		app.Status.AvailableReplicas = dep.Status.AvailableReplicas
		app.Status.UnavailableReplicas = dep.Status.UnavailableReplicas

		err := r.Status().Update(ctx, app)
		if err != nil {
			log.Error(err, "Failed to update app status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *AppReconciler) containersEqual(containers, targetContainers []corev1.Container) bool {
	if len(containers) != len(targetContainers) {
		return false
	}

	for i, c := range containers {
		targetC := targetContainers[i]
		if c.Image != targetC.Image {
			return false
		}

		if c.Name != targetC.Name {
			return false
		}

		if strings.Join(c.Command, " ") != strings.Join(targetC.Command, " ") {
			return false
		}

		if len(c.Ports) != len(targetC.Ports) {
			return false
		}

		for i, p := range c.Ports {
			targetP := targetC.Ports[i]

			if p.ContainerPort != targetP.ContainerPort || p.Protocol != targetP.Protocol {
				return false
			}
		}
	}

	return true
}

func AppProjectName(app *cloudv1alpha1.App) string {
	return app.ObjectMeta.Labels[projectCRKey]
}

// deploymentForApp returns an app Deployment object
func (r *AppReconciler) deploymentForApp(app *cloudv1alpha1.App) (*appsv1.Deployment, error) {
	projectName := AppProjectName(app)
	namespaceName := ProjectNamespaceName(projectName)
	labels := LabelsForApp(projectName, app.Name)
	replicas := app.Spec.Replicas

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: namespaceName,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{},
				},
			},
		},
	}

	containers := []corev1.Container{}
	for _, c := range app.Spec.Containers {
		container := corev1.Container{
			Image:   c.Image,
			Name:    c.Name,
			Command: c.Command,
			Ports:   []corev1.ContainerPort{},
		}

		for _, p := range c.Ports {
			container.Ports = append(container.Ports, corev1.ContainerPort{
				Protocol:      p.Protocol,
				ContainerPort: p.Number,
			})
		}

		containers = append(containers, container)
	}
	dep.Spec.Template.Spec.Containers = containers

	// Set app instance as the owner and controller
	err := ctrl.SetControllerReference(app, dep, r.Scheme)
	if err != nil {
		return nil, err
	}
	return dep, nil
}

// LabelsForApp returns the labels for an app
func LabelsForApp(projectName, appName string) map[string]string {
	return map[string]string{
		"app":        appName,
		projectCRKey: projectName,
	}
}

// serviceForApp returns an app Service object
func (r *AppReconciler) serviceForApp(app *cloudv1alpha1.App) (*corev1.Service, error) {
	projectName := AppProjectName(app)
	namespaceName := ProjectNamespaceName(projectName)
	labels := LabelsForApp(projectName, app.Name)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: namespaceName,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports:    []corev1.ServicePort{},
		},
	}

	// expose all ports defined in app containers
	for _, c := range app.Spec.Containers {
		for _, p := range c.Ports {
			port := corev1.ServicePort{
				Protocol: p.Protocol,
				Port:     p.Number,
				Name:     fmt.Sprintf("%s-%d", c.Name, p.Number),
			}

			svc.Spec.Ports = append(svc.Spec.Ports, port)
		}
	}

	// Set app instance as the owner and controller
	err := ctrl.SetControllerReference(app, svc, r.Scheme)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (r *AppReconciler) servicePortsEqual(ports, targetPorts []corev1.ServicePort) bool {
	if len(ports) != len(targetPorts) {
		return false
	}

	for i, p := range ports {
		targetP := targetPorts[i]

		if p.Port != targetP.Port || p.Protocol != targetP.Protocol || p.Name != targetP.Name {
			return false
		}
	}

	return true
}

func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha1.App{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
