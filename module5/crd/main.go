package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func main() {
	kind := handleArgs()
	config := initCfg()
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	discovery := clientSet.Discovery()
	apiGroupResources, err := restmapper.GetAPIGroupResources(discovery)
	if err != nil {
		panic(err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(apiGroupResources)
	groupVersionKind := schema.GroupVersionKind{
		Group:   "aiops.geektime.com",
		Version: "v1alpha1",
		Kind:    kind,
	}
	mapping, err := mapper.RESTMapping(groupVersionKind.GroupKind(), groupVersionKind.Version)
	if err != nil {
		panic(err)
	}
	resourceInterface := dynamicClient.Resource(mapping.Resource).Namespace(corev1.NamespaceDefault)
	resources, err := resourceInterface.List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, resource := range resources.Items {
		fmt.Printf("Name: %s, Namespace %s \n", resource.GetName(), resource.GetNamespace())
		domain, found, err := unstructured.NestedString(resource.Object, "spec", "domain")
		if err != nil {
			continue
		}
		if found {
			fmt.Println("domain:", domain)
		}
	}

}

func handleArgs() string {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s get <resource>\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[1]
	kind := os.Args[2]

	if command != "get" {
		fmt.Println("unsupported command:", command)
		os.Exit(1)
	}
	return kind
}

func initCfg() *rest.Config {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig user path")
	} else {
		kubeconfig = flag.String("kubeconfig", "/etc/kubernetes/admin.conf", "etc kube config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	return config
}
