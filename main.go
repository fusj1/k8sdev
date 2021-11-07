package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func HomeDir() string {
	if homedir := os.Getenv("HOME"); homedir != "" {
		return homedir
	} else {
		return os.Getenv("USERPROFILE")
	}
}

func main() {
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
	// 初始化一个工厂函数
	informerfactory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	// 初始化一个deploments的 informer
	deployinformer := informerfactory.Apps().V1().Deployments()

	// 创建 Informer（相当于注册到工厂中去，这样下面启动的时候就会去 List & Watch 对应的资源）
	informer := deployinformer.Informer()
	// 启动一个deploylistener
	deploylistener := deployinformer.Lister()
	// 注册事件处理函数
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			fmt.Println("Add deployments:", deploy.Name)
		},
		UpdateFunc: func(OldObj, NewObj interface{}) {
			olddeploy := OldObj.(*appsv1.Deployment)
			newdeploy := NewObj.(*appsv1.Deployment)
			fmt.Printf("%s->%s\n", olddeploy.Name, newdeploy.Name)
		},
		DeleteFunc: func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			fmt.Println("Delete deployments:", deploy.Name)
		},
	})
	stopper := make(chan struct{})
	// 启动工厂函数监听
	informerfactory.Start(stopper)
	defer close(stopper)
	informerfactory.WaitForCacheSync(stopper)
	deploys, err := deploylistener.Deployments("default").List(labels.Everything())
	if err != nil {
		panic(err.Error())
	}
	for idx, deploy := range deploys {
		fmt.Printf("%d -->%s", idx+1, deploy.Name)
	}
	<-stopper
}
