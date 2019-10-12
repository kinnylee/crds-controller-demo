package main

import (
	"flag"
	_"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientSet "kinnylee.com/crds-controller-demo/pkg/client/clientset/versioned"
	informers "kinnylee.com/crds-controller-demo/pkg/client/informers/externalversions"
	"kinnylee.com/crds-controller-demo/pkg/signals"
	"log"
	"os/user"
	"path/filepath"
	"time"
)

var(
	kubeConfig *string
	masterUrl string
)

func init(){
	u, err := user.Current()

	if err != nil {
		panic(err.Error())
	}
	if home := u.HomeDir; home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "path to the kube config file",
		)
	} else {
		kubeConfig = flag.String("kubeconfig", "", "path to the kube config file")
	}
}

func main()  {
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Printf("faild create cluster, config: %v", err)
		panic(err.Error())
	}

	scalingClient, err := clientSet.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	/*
	factory := informers.NewSharedInformerFactory(scalingClient, 30 * time.Second)
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
	*/

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	factory := informers.NewSharedInformerFactory(scalingClient, 30 * time.Second)

	controller := NewController(kubeClient, scalingClient, factory.Control().V1().Scalings())

	stopCh := signals.SetupSignalHandler()
	go factory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		log.Fatal("error run controller: %s", err.Error())
	}
}
