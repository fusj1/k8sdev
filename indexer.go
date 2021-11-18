package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	NameSpaceIndexName = "namespace"
	NodeNameIndexName  = "node"
)

func NameSpaceIndexFunc(obj interface{}) ([]string, error) {
	// 直接断言成了一个clientset
	m, err := meta.Accessor(obj)
	if err != nil {
		panic(err.Error())
	}
	return []string{m.GetNamespace()}, nil
}
func NodeNameIndexFunc(obj interface{}) ([]string, error) {
	// 断言成了一个pod
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return []string{""}, nil
	}
	return []string{pod.Spec.NodeName}, nil
}

func indexer() {
	indexers := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		NameSpaceIndexName: NameSpaceIndexFunc,
		NodeNameIndexName:  NodeNameIndexFunc,
	})
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "node1",
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod2",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "node1",
		},
	}
	pod3 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod3",
			Namespace: "kube-system",
		},
		Spec: v1.PodSpec{
			NodeName: "node2",
		},
	}
	_ = indexers.Add(pod1)
	_ = indexers.Add(pod2)
	_ = indexers.Add(pod3)

	pods, err := indexers.ByIndex(NameSpaceIndexName, "default")
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods {
		fmt.Println(pod.(*v1.Pod).Name)
	}
	pods, err = indexers.ByIndex(NodeNameIndexName, "node1")
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods {
		fmt.Println(pod.(*v1.Pod).Name)

	}
}
