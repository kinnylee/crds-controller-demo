package main

import (
	"fmt"
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	scalingv1 "kinnylee.com/crds-controller-demo/pkg/apis/control/v1"
	clientset "kinnylee.com/crds-controller-demo/pkg/client/clientset/versioned"
	scalingscheme "kinnylee.com/crds-controller-demo/pkg/client/clientset/versioned/scheme"
	informer "kinnylee.com/crds-controller-demo/pkg/client/informers/externalversions/control/v1"
	lister "kinnylee.com/crds-controller-demo/pkg/client/listers/control/v1"
	"time"
)

const controllerAgentName = "scaling-controller"

type Controller struct {
	kubeclientset kubernetes.Interface
	scalingclientset clientset.Interface
	scalingLister lister.ScalingLister
	scalingSyncd cache.InformerSynced
	workqueue workqueue.RateLimitingInterface
	record record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	scalingclientset clientset.Interface,
	scalingInformer informer.ScalingInformer) *Controller{
		runtime.Must(scalingscheme.AddToScheme(scheme.Scheme))
		glog.V(4).Info("Create event broadcaster")
		eventBroadcaster := record.NewBroadcaster()
		eventBroadcaster.StartLogging(glog.Infof)
		eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
			Interface: kubeclientset.CoreV1().Events(""),
		})
		recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{
			Component: controllerAgentName,
		})
		
		controller := &Controller{
			kubeclientset:    kubeclientset,
			scalingclientset: scalingclientset,
			scalingLister:    scalingInformer.Lister(),
			scalingSyncd:     scalingInformer.Informer().HasSynced,
			workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Scaling"),
			record:           recorder,
		}

		glog.Info("Setting up event handlers")

		// 添加回调函数，处理增删改查操作
		scalingInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc:    controller.enqueueScaling,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldScaling := (oldObj).(scalingv1.Scaling)
				newScaling := (newObj).(scalingv1.Scaling)
				if oldScaling.ResourceVersion == newScaling.ResourceVersion {
					return
				}
				controller.enqueueScaling(newScaling)
			},
			DeleteFunc: controller.enqueueScalingForDelete,
		})
	return controller
}

// 启动入口
func (c *Controller) Run(threadiness int, stopCh <- chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("start run controller")
	if ok := cache.WaitForCacheSync(stopCh, c.scalingSyncd); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("start worker")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker started.")
	<-stopCh
	glog.Info("worker stop.")
	return nil
}

func (c *Controller) runWorker(){
	for c.processNextWorkItem() {

	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error{
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("excepted string in workqueue, but get %#V", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("Error syncing %s:%s", key, err.Error())
		}
		c.workqueue.Forget(obj)
		glog.Info("Successfully Synced %s", key)
		return nil

	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}
	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key : %s", key))
		return nil
	}

	scaling, err := c.scalingLister.Scalings(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err){
			glog.Infof("scaling object is deleted. %s/%s", namespace, name)
			return nil
		}
		runtime.HandleError(fmt.Errorf("faild to list sacling by: %s/%s", namespace, name))
		return err
	}

	glog.Infof("这里是student对象的期望状态: %#v ...", scaling)
	glog.Infof("实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)")
	c.record.Event(scaling, v1.EventTypeNormal, "Synced", "Scaling synced successfully")
	return nil
}


// 数据先放入缓存，再放入队列
func (c * Controller) enqueueScaling(obj interface{}){
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	// 将key放入队列
	c.workqueue.AddRateLimited(key)
}

// 删除操作
func (c * Controller) enqueueScalingForDelete(obj interface{}){
	var key string
	var err error

	// 先从缓存中删除对象
	if key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	// 再将key放入队列
	c.workqueue.AddRateLimited(key)
}
