package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/clientcmd"
	clientSet "kinnylee.com/crds-controller-demo/pkg/client/clientset/versioned"
	informers "kinnylee.com/crds-controller-demo/pkg/client/informers/externalversions"
	"log"
	"os/user"
	"path/filepath"
	"time"
)

func main()  {
	client, err := newKubeClient()
	if err != nil {
		log.Fatalf("new kube client error: %v", err)
	}

	factory := informers.NewSharedInformerFactory(client, 30 * time.Second)
	informer := factory.Control().V1().Scalings()
	lister := informer.Lister()

	stopCh := make(chan struct{})
	factory.Start(stopCh)

	for{
		ret, err := lister.List(labels.Everything())
		if err != nil {
			log.Printf("list error: %v", err)
		} else {
			for _, scaling := range ret{
				log.Println(scaling)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func newKubeClient()(clientSet.Interface, error) {
	u, err := user.Current()

	if err != nil {
		panic(err.Error())
	}

	kubeConfig := filepath.Join(u.HomeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("faild create cluster, config: %v", err)
	}
	cli, err := clientSet.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("faild create custom kube client: %v", err)
	}
	return cli, nil
}
