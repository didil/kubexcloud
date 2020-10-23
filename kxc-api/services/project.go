package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/didil/kubexcloud/kxc-api/requests"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
)

type ProjectSvc interface {
	Create(ctx context.Context, reqData *requests.CreateProject) error
	Get(ctx context.Context, projectName string) (*cloudv1alpha1.Project, error)
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
		return fmt.Errorf("name is required")
	}

	if errs := validationutils.IsDNS1123Label(reqData.Name); len(errs) > 0 {
		return fmt.Errorf("name: %s", strings.Join(errs, "."))
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

func (svc *ProjectService) Get(ctx context.Context, projectName string) (*cloudv1alpha1.Project, error) {
	client := svc.k8sSvc.Client()

	proj := &cloudv1alpha1.Project{}
	err := client.Get(ctx, types.NamespacedName{Name: projectName}, proj)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project: %v", err)
	}

	return proj, nil
}
