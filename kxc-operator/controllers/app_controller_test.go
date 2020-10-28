package controllers

import (
	"context"
	"time"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("App controller", func() {
	const (
		ProjectName    = "test-proj-for-app"
		NamespaceName  = "kxc-proj-test-proj-for-app"
		AppName        = "test-app"
		DeploymentName = "test-app"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating an app", func() {
		var proj *cloudv1alpha1.Project
		var app *cloudv1alpha1.App

		It("Should create a deployment, service and ingress", func() {
			ctx := context.Background()
			proj = &cloudv1alpha1.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name: ProjectName,
				},
			}
			Expect(k8sClient.Create(ctx, proj)).Should(Succeed())

			// wait for namespace creation
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: NamespaceName}, &corev1.Namespace{})
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			app = &cloudv1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: NamespaceName,
					Labels:    LabelsForApp(ProjectName, AppName),
				},
				Spec: cloudv1alpha1.AppSpec{
					Replicas: 2,
					Containers: []cloudv1alpha1.Container{
						cloudv1alpha1.Container{
							Image:   "busybox",
							Name:    "test-container",
							Command: []string{"sleep", "60"},
							Ports: []cloudv1alpha1.Port{
								cloudv1alpha1.Port{
									Number:           9123,
									Protocol:         "TCP",
									ExposeExternally: true,
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, app)).Should(Succeed())

			// check deployment

			createdDeployment := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: DeploymentName, Namespace: NamespaceName}, createdDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(*createdDeployment.Spec.Replicas).Should(Equal(int32(2)))

			Expect(createdDeployment.Spec.Selector).Should(Equal(&metav1.LabelSelector{
				MatchLabels: LabelsForApp(ProjectName, AppName),
			}))
			Expect(createdDeployment.Spec.Template.ObjectMeta.Labels).Should(Equal(LabelsForApp(ProjectName, AppName)))
			Expect(createdDeployment.Spec.Template.Spec.Containers).Should(HaveLen(1))
			cont := createdDeployment.Spec.Template.Spec.Containers[0]
			Expect(cont.Image).Should(Equal("busybox"))
			Expect(cont.Name).Should(Equal("test-container"))
			Expect(cont.Command).Should(Equal([]string{"sleep", "60"}))
			Expect(cont.Ports).Should(Equal([]corev1.ContainerPort{
				corev1.ContainerPort{
					Protocol:      "TCP",
					ContainerPort: 9123,
				},
			}))

			Expect(createdDeployment.Labels).Should(Equal(map[string]string{
				"app":        AppName,
				"project_cr": ProjectName,
			}))

			Expect(createdDeployment.OwnerReferences).Should(HaveLen(1))
			Expect(createdDeployment.OwnerReferences[0].UID).Should(Equal(app.UID))

			// check service

			createdService := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: DeploymentName, Namespace: NamespaceName}, createdService)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdService.Labels).Should(Equal(map[string]string{
				"app":        AppName,
				"project_cr": ProjectName,
			}))

			Expect(createdService.Spec.Selector).Should(Equal(LabelsForApp(ProjectName, AppName)))
			Expect(createdService.Spec.Type).Should(Equal(corev1.ServiceTypeClusterIP))
			Expect(createdService.Spec.Ports).Should(Equal([]corev1.ServicePort{
				corev1.ServicePort{
					Protocol:   "TCP",
					Port:       9123,
					TargetPort: intstr.FromInt(9123),
					Name:       "test-container-9123",
				},
			}))

			Expect(createdService.OwnerReferences).Should(HaveLen(1))
			Expect(createdService.OwnerReferences[0].UID).Should(Equal(app.UID))

			// add ingress test

			createdIngress := &netv1beta1.Ingress{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: DeploymentName, Namespace: NamespaceName}, createdIngress)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdIngress.Labels).Should(Equal(map[string]string{
				"app":        AppName,
				"project_cr": ProjectName,
			}))

			Expect(createdIngress.Spec.Rules).Should(Equal([]netv1beta1.IngressRule{
				netv1beta1.IngressRule{
					Host: "test-app.127.0.0.1.xip.io",
					IngressRuleValue: netv1beta1.IngressRuleValue{
						HTTP: &netv1beta1.HTTPIngressRuleValue{
							Paths: []netv1beta1.HTTPIngressPath{
								netv1beta1.HTTPIngressPath{
									Backend: netv1beta1.IngressBackend{
										ServiceName: app.Name,
										ServicePort: intstr.FromInt(9123),
									},
								},
							},
						},
					},
				},
			}))

			Expect(createdIngress.OwnerReferences).Should(HaveLen(1))
			Expect(createdIngress.OwnerReferences[0].UID).Should(Equal(app.UID))

			// check app status

			Eventually(func() string {
				app := &cloudv1alpha1.App{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: AppName, Namespace: NamespaceName}, app)
				if err != nil {
					return ""
				}
				return app.Status.ExternalURL
			}, timeout, interval).Should(Equal("http://test-app.127.0.0.1.xip.io/"))
		})

		AfterEach(func() {
			ctx := context.Background()
			Expect(k8sClient.Delete(ctx, app)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, proj)).Should(Succeed())
		})
	})

})
