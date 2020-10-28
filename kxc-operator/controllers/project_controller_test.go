package controllers

import (
	"context"
	"time"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Project controller", func() {
	const (
		ProjectName   = "test-project"
		NamespaceName = "kxc-proj-test-project"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a project", func() {
		var proj *cloudv1alpha1.Project

		It("Should create a matching namespace", func() {
			ctx := context.Background()
			proj = &cloudv1alpha1.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name: ProjectName,
				},
			}
			Expect(k8sClient.Create(ctx, proj)).Should(Succeed())

			namespaceLookupKey := types.NamespacedName{Name: NamespaceName}
			createdNamespace := &corev1.Namespace{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, namespaceLookupKey, createdNamespace)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdNamespace.Labels).Should(Equal(map[string]string{
				"app":        "kxc",
				"project_cr": "test-project",
			}))
		})

		AfterEach(func() {
			ctx := context.Background()
			Expect(k8sClient.Delete(ctx, proj)).Should(Succeed())
		})
	})

})
