package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/didil/kubexcloud/kxc-api/requests"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
)

type AppSvc interface {
	Create(ctx context.Context, projectName string, reqData *requests.CreateApp) error
}

type AppService struct {
	k8sSvc K8sSvc
}

func NewAppService(k8sSvc K8sSvc) *AppService {
	return &AppService{
		k8sSvc: k8sSvc,
	}
}

func (svc *AppService) validateApp(reqData *requests.CreateApp) error {
	if reqData.Name == "" {
		return fmt.Errorf("name is required")
	}

	if errs := validationutils.IsDNS1123Label(reqData.Name); len(errs) > 0 {
		return fmt.Errorf("name: %s", strings.Join(errs, "."))
	}

	return nil
}

func (svc *AppService) Create(ctx context.Context, projectName string, reqData *requests.CreateApp) error {
	client := svc.k8sSvc.Client()

	err := svc.validateApp(reqData)
	if err != nil {
		return fmt.Errorf("app invalid: %v", err)
	}

	app := &cloudv1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: reqData.Name,
		},
		Spec: cloudv1alpha1.AppSpec{
			Replicas:   reqData.Replicas,
			Containers: []cloudv1alpha1.Container{},
		},
	}

	for _, c := range reqData.Containers {
		container := cloudv1alpha1.Container{
			Image:   c.Image,
			Name:    c.Name,
			Command: c.Command,
			Ports:   []cloudv1alpha1.Port{},
		}

		for _, p := range c.Ports {
			container.Ports = append(container.Ports, cloudv1alpha1.Port{
				Number:   p.Number,
				Protocol: corev1.Protocol(p.Protocol),
			})
		}

		app.Spec.Containers = append(app.Spec.Containers, container)
	}

	err = client.Create(ctx, app)
	if err != nil {
		return fmt.Errorf("create app: %v", err)
	}
	return nil
}
