package services

import (
	"fmt"
	"os"
	"path/filepath"

	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sSvc interface {
	Client() client.Client
}

type K8sService struct {
	client client.Client
}

func NewK8sService() (*K8sService, error) {
	svc := &K8sService{}

	// init runtime scheme
	scheme := runtime.NewScheme()
	// add k8s scheme
	err := clientgoscheme.AddToScheme(scheme)
	if err != nil {
		return nil, fmt.Errorf("clientgoscheme: %v", err)
	}
	// add crd schemes
	err = cloudv1alpha1.AddToScheme(scheme)
	if err != nil {
		return nil, fmt.Errorf("cloudv1alpha1: %v", err)
	}

	client, err := svc.initK8sClient(scheme)
	if err != nil {
		return nil, fmt.Errorf("initK8sClient: %v", err)
	}

	svc.client = client

	return svc, nil
}

func (svc *K8sService) initK8sClient(scheme *runtime.Scheme) (client.Client, error) {
	config, err := svc.getInClusterConfig()
	if err != nil {
		return nil, err
	}

	if config == nil {
		config, err = svc.getOutOfClusterConfig()
		if err != nil {
			return nil, err
		}

	}

	client, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (svc *K8sService) getInClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (svc *K8sService) getOutOfClusterConfig() (*rest.Config, error) {
	var kubeconfig string
	if kubeconfigenv := os.Getenv("KUBECONFIG"); kubeconfigenv != "" {
		kubeconfig = kubeconfigenv
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	if kubeconfig == "" {
		return nil, fmt.Errorf("out of cluster: cloud not find a kubeconfig")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (svc *K8sService) Client() client.Client {
	return svc.client
}
