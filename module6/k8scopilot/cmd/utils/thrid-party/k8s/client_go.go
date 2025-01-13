package k8s

import (
	"github.com/ahyang98/k8scopilot/cmd/utils"
	"path/filepath"
	"strings"

	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ClientGo struct {
	ClientSet       *kubernetes.Clientset
	DynamicClient   dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
}

var clientGo *ClientGo

func GetClientGo() (*ClientGo, error) {
	if clientGo == nil {
		kubeconfig, err := utils.GetConfig().Get("kubeconfig")
		if err != nil {
			return nil, err
		}
		clientGo, err = newClientGo(kubeconfig.(string))
		if err != nil {
			return nil, err
		}
	}
	return clientGo, nil
}

func newClientGo(kubeconfig string) (*ClientGo, error) {
	if strings.HasPrefix(kubeconfig, "~") {
		homeDir := homedir.HomeDir()
		kubeconfig = filepath.Join(homeDir, kubeconfig[1:])
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	return &ClientGo{
		ClientSet:       clientset,
		DynamicClient:   dynamicClient,
		DiscoveryClient: discoveryClient,
	}, nil
}

func GetGVR(resourceType string) (schema.GroupVersionResource, error) {
	resourceType = strings.ToLower(resourceType)
	var gvr schema.GroupVersionResource
	switch resourceType {
	case "deployment":
		gvr = schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "deployments",
		}
	case "service":
		gvr = schema.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "services",
		}
	case "pod":
		gvr = schema.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "pods",
		}
	default:
		return schema.GroupVersionResource{}, fmt.Errorf("unsupported resource type %s", resourceType)
	}
	return gvr, nil
}
