package services

import (
	"context"
	"fmt"

	"github.com/didil/kubexcloud/kxc-api/requests"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectSvc interface {
	Create(ctx context.Context, reqData *requests.CreateProject) error
}

type ProjectService struct {
	k8sSvc K8sSvc
}

func NewProjectService(k8sSvc K8sSvc) *ProjectService {
	return &ProjectService{
		k8sSvc: k8sSvc,
	}
}

func (svc *ProjectService) validateProject(reqData *requests.CreateProject) error {
	if reqData.Name == "" {
		return fmt.Errorf("project name is required")
	}

	return nil
}

func (svc *ProjectService) Create(ctx context.Context, reqData *requests.CreateProject) error {
	clientset := svc.k8sSvc.Clientset()

	err := svc.validateProject(reqData)
	if err != nil {
		return fmt.Errorf("project invalid: %v", err)
	}

	namespace := &cloudv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: reqData.Name,
		},
	}

	_, err = clientset.r.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create project: %v", err)
	}
	return nil
}
