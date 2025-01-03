package main

import (
	"flag"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"log"
	"path/filepath"
	"time"
)

func main() {

	clientSet, err := kubernetes.NewForConfig(getClusterCfg())
	if err != nil {
		panic(err)
	}
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Hour*12)
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	podInformer := informerFactory.Core().V1().Pods()
	_, err = podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onAdd(obj, queue)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			onUpdate(newObj, queue)
		},
		DeleteFunc: func(obj interface{}) {
			onDelete(obj, queue)
		},
	})
	if err != nil {
		panic(err)
	}
	podController := NewPodController(podInformer.Informer().GetIndexer(), queue, podInformer.Informer())
	stopper := make(chan struct{})
	informerFactory.Start(stopper)
	informerFactory.WaitForCacheSync(stopper)
	go func() {
		for {
			if !podController.ProcessNextItem() {
				break
			}
		}
	}()
	<-stopper
}

func onDelete(obj interface{}, queue workqueue.TypedRateLimitingInterface[string]) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Println("on delete ", key, err)
		return
	}
	queue.Add(key)
}

func onUpdate(obj interface{}, queue workqueue.TypedRateLimitingInterface[string]) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Println("on update ", key, err)
		return
	}
	queue.Add(key)
}

func onAdd(obj interface{}, queue workqueue.TypedRateLimitingInterface[string]) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Println("on add ", key, err)
		return
	}
	queue.Add(key)
}

func getClusterCfg() *rest.Config {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig user path")
	} else {
		kubeconfig = flag.String("kubeconfig", "/etc/kubernetes/admin.conf", "etc kube config")
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}
	return config
}
