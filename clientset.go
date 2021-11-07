package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func HomeDir() string {
	if homedir := os.Getenv("HOME"); homedir != "" {
		return homedir
	} else {
		return os.Getenv("USERPROFILE")
	}
}

func clientset() {
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
		panic(err.Error())
	}
	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for idx, deploy := range deployments.Items {
		fmt.Printf("%d->%s", idx+1, deploy.Name)
	}

}
