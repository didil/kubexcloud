package services

import (
	"context"
	"fmt"

	"github.com/didil/kubexcloud/kxc-api/requests"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
	client := svc.k8sSvc.Client()

	err := svc.validateProject(reqData)
	if err != nil {
		return fmt.Errorf("project invalid: %v", err)
	}

	proj := &cloudv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: reqData.Name,
		},
	}

	err = client.Get(ctx, types.NamespacedName{Name: proj.Name}, &cloudv1alpha1.Project{})
	if err == nil {
		return fmt.Errorf("project already exists: %v", proj.Name)
	} else if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("get project: %v", err)
	}

	err = client.Create(ctx, proj)
	if err != nil {
		return fmt.Errorf("create project: %v", err)
	}
	return nil
}
