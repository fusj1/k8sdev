package main

import (
	"flag"
	"os"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func HHomeDir() string {
	if homedir := os.Getenv("HOME"); homedir != "" {
		return homedir
	} else {
		return os.Getenv("USERPROFILE")
	}
}

func NewClientSet() (kubernetes.Clientset, error) {
	var config *rest.Config
	var kubeconfig *string
	var err error
	if h := HomeDir(); h != "" {
		kubeconfig = flag.String("kubeconfig", path.Join(h, ".kube", "config"), "(可选的)kubernetes config")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubernetes config")
	}
	flag.Parse()
	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}

	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return kubernetes.Clientset{}, err
	}
	return *clientset, nil
}
