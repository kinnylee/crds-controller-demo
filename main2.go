package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientSet "kinnylee.com/crds-controller-demo/pkg/client/clientset/versioned"
	informers "kinnylee.com/crds-controller-demo/pkg/client/informers/externalversions"
	"log"
	"os/user"
	"path/filepath"
	"time"
)

func main()  {
	u, err := user.Current()

	if err != nil {
		panic(err.Error())
	}

	kubeConfig := filepath.Join(u.HomeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Printf("faild create cluster, config: %v", err)
		panic(err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	scalingClient, err := clientSet.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(scalingClient, 30 * time.Second)

	controller := NewController(kubeClient, scalingClient, factory.Control().V1().Scalings())

	stopCh := make(<-chan struct{})
	go factory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		log.Fatal("error run controller: %s", err.Error())
	}
}
