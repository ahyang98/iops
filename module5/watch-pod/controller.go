package main

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

type PodController struct {
	indexer  cache.Indexer
	queue    workqueue.TypedRateLimitingInterface[string]
	informer cache.Controller
}

func NewPodController(indexer cache.Indexer, queue workqueue.TypedRateLimitingInterface[string], informer cache.Controller) *PodController {
	return &PodController{indexer: indexer, queue: queue, informer: informer}
}

func (c *PodController) ProcessNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	err := c.syncToHandle(key)
	c.handleErr(err, key)
	return true
}

func (c *PodController) syncToHandle(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		log.Printf("fetch object with key %s from store fail for %v\n", key, err)
		return err
	}
	if !exists {
		log.Printf("pod %s does not exists  \n", key)
		return nil
	}
	pod := obj.(*v1.Pod)
	log.Printf("sync pod %s, status %s\n", pod.Name, pod.Status)
	return nil
}

func (c *PodController) handleErr(err error, key string) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		log.Printf("Retry %d for key %s\n", c.queue.NumRequeues(key), key)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	log.Printf("Dropping pod %s out of the queue: %v \n", key, err)
}
