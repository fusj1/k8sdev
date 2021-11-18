package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func informerv1() {
	clientset, err := NewClientSet()
	if err != nil {
		panic(err.Error())
	}
	// 初始化 informer factory（为了测试方便这里设置每30s重新 List 一次）
	informerFactory := informers.NewSharedInformerFactory(&clientset, time.Second*30)
	// 对 Deployment 监听
	deploymentinformer := informerFactory.Apps().V1().Deployments()
	// 创建informer并注册到相应的工厂当中去（这样下面就会去listener和watcher相应的资源）
	informer := deploymentinformer.Informer()
	// 创建listener，注册事件处理函数
	infromerListener := deploymentinformer.Lister()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    AddFunc,
		UpdateFunc: UpdateFunc,
		DeleteFunc: DeleteFunc,
	})
	stopCh := make(chan struct{})
	// 启动工厂程序开始listener and watcher 对象
	informerFactory.Start(stopCh)

	// 从本地的缓存中获取数据
	deployments, err := infromerListener.Deployments("default").List(labels.Everything())
	for idx, deploy := range deployments {
		fmt.Printf("%s ->%s", idx+1, deploy)
	}
	<-stopCh
}
func AddFunc(obj interface{}) {
	deploy := obj.(*v1.Deployment)
	fmt.Println("Add deploy:", deploy.Name)
}

func UpdateFunc(oldObj, newObj interface{}) {
	olddeploy := oldObj.(*v1.Deployment)
	newdeploy := newObj.(*v1.Deployment)
	fmt.Printf("Update deploy for %s ->%s", olddeploy.Name, newdeploy.Name)
}

func DeleteFunc(obj interface{}) {
	deployobj := obj.(*v1.Deployment)
	fmt.Println("Delete deploy", deployobj)
}
