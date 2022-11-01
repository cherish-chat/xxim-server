package discov

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/threading"
	apiv1 "k8s.io/api/core/v1"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

// 监听k8s某个namespace下的某个service的endpoints变化

// EndpointsEventHandler is ResourceEventHandler implementation.
type EndpointsEventHandler struct {
	update    func([]string)
	endpoints map[string]lang.PlaceholderType
	lock      sync.Mutex
}

// newEndpointsEventHandler returns an EndpointsEventHandler.
func newEndpointsEventHandler(update func([]string)) *EndpointsEventHandler {
	return &EndpointsEventHandler{
		update: update,
	}
}

// OnAdd handles the endpoints add events.
func (h *EndpointsEventHandler) OnAdd(obj interface{}) {
	endpoints, ok := obj.(*apiv1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", obj)
		return
	}

	var ips []string
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			ips = append(ips, point.IP)
		}
	}

	h.update(utils.Set[string](ips))
}

// OnDelete handles the endpoints delete events.
func (h *EndpointsEventHandler) OnDelete(obj interface{}) {
	endpoints, ok := obj.(*apiv1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", obj)
		return
	}

	var ips []string
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			ips = append(ips, point.IP)
		}
	}

	h.update(utils.Set[string](ips))
}

// OnUpdate handles the endpoints update events.
func (h *EndpointsEventHandler) OnUpdate(oldObj, newObj interface{}) {
	oldEndpoints, ok := oldObj.(*apiv1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", oldObj)
		return
	}

	newEndpoints, ok := newObj.(*apiv1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", newObj)
		return
	}

	if oldEndpoints.ResourceVersion == newEndpoints.ResourceVersion {
		return
	}

	h.Update(newEndpoints)
}

func (h *EndpointsEventHandler) Update(endpoints *apiv1.Endpoints) {
	h.lock.Lock()
	defer h.lock.Unlock()

	old := h.endpoints
	h.endpoints = make(map[string]lang.PlaceholderType)
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			h.endpoints[point.IP] = lang.Placeholder
		}
	}

	if diff(old, h.endpoints) {
		h.notify()
	}
}

func (h *EndpointsEventHandler) notify() {
	var targets []string

	for k := range h.endpoints {
		targets = append(targets, k)
	}

	h.update(utils.Set[string](targets))
}

func diff(o, n map[string]lang.PlaceholderType) bool {
	if len(o) != len(n) {
		return true
	}

	for k := range o {
		if _, ok := n[k]; !ok {
			return true
		}
	}

	return false
}

const (
	resyncInterval = 5 * time.Minute
	nameSelector   = "metadata.name="
)

func MustListenEndpoints(namespace string, serviceName string, onUpdate func([]string)) *EndpointsEventHandler {
	v, err := ListenEndpoints(namespace, serviceName, onUpdate)
	if err != nil {
		logx.Errorf("failed to listen endpoints: %v", err)
		panic(err)
	}
	return v
}

func ListenEndpoints(namespace string, serviceName string, onUpdate func([]string)) (*EndpointsEventHandler, error) {
	if onUpdate == nil {
		onUpdate = func(endpoints []string) {
			//var addrs []resolver.Address
			for _, endpoint := range endpoints {
				//addrs = append(addrs, resolver.Address{
				//	Addr: fmt.Sprintf("%s:%d", endpoint, port),
				//})
				fmt.Printf("endpoint: %s\n", endpoint)
			}
		}
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	handler := newEndpointsEventHandler(onUpdate)
	inf := informers.NewSharedInformerFactoryWithOptions(cs, resyncInterval,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = nameSelector + serviceName
		}),
	)
	in := inf.Core().V1().Endpoints()
	in.Informer().AddEventHandler(handler)
	threading.GoSafe(func() {
		inf.Start(proc.Done())
	})

	endpoints, _ := cs.CoreV1().Endpoints(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})

	handler.Update(endpoints)

	return handler, nil
}

func (h *EndpointsEventHandler) Endpoints() []string {
	h.lock.Lock()
	defer h.lock.Unlock()

	var endpoints []string
	for k := range h.endpoints {
		endpoints = append(endpoints, k)
	}

	return endpoints
}
