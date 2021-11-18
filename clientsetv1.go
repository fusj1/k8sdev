package main

import (
	"context"
	"flag"
	"fmt"
	"path"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func clientsetv1() kubernetes.Clientset {
	var config *rest.Config
	var ctx context.Context
	var err error

	// *string是为了符合flag.String的返回类型
	var kubeconfig *string
	// 获取配置文件的方法
	if home := HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", path.Join(home, ".kube", "config"), "（可选的）kubeconfig配置文件")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件的路径")
	}
	// 获取config的方法，或者利用已经得到的配置文件的地址去初始化配置文件
	if config, err = rest.InClusterConfig(); err != nil {
		// 从配置文件初始化Config文件的时候需要用clientcmd库
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//利用clientset获取deployment
	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, v1.ListOptions{})
	for idx, deploy := range deployments.Items {
		fmt.Printf("%v-->%v", idx, deploy)
	}
	return *clientset
}

// // 获取本地的家目录
// func HomeDir() string {
// 	// 得到当前的家目录的位置,如果报错则证明当前的环境是windows环境
// 	if home := os.Getenv("home"); home != "" {
// 		return home
// 	} else {
// 		return os.Getenv("USERPROFILE")
// 	}
// }
