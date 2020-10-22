package services

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sSvc interface {
	Clientset() *kubernetes.Clientset
}

type K8sService struct {
	clientset *kubernetes.Clientset
}

func NewK8sService() (*K8sService, error) {
	svc := &K8sService{}
	clientset, err := svc.initK8sClientSet()
	if err != nil {
		return nil, fmt.Errorf("initK8sClientSet: %v", err)
	}

	svc.clientset = clientset

	return svc, nil
}

func (svc *K8sService) initK8sClientSet() (*kubernetes.Clientset, error) {
	clientset, err := svc.initInClusterK8sClientSet()
	if err != rest.ErrNotInCluster && err != nil {
		return nil, err
	}
	if err == nil {
		// ok
		return clientset, nil
	}

	clientset, err = svc.initOutOfClusterK8sClientSet()
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (svc *K8sService) initInClusterK8sClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (svc *K8sService) initOutOfClusterK8sClientSet() (*kubernetes.Clientset, error) {
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

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (svc *K8sService) Clientset() *kubernetes.Clientset {
	return svc.clientset
}
