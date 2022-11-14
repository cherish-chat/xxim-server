package connservice

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/common/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"sync"
)

func NewConnPodsMgr(config ConnPodsConfig) *ConnPodsMgr {
	c := &ConnPodsMgr{Config: config}
	c.initConnRpc()
	return c
}

type ConnPodsMgr struct {
	connPods              sync.Map
	endpointsEventHandler *discov.EndpointsEventHandler
	Config                ConnPodsConfig
}

type ConnPodsConfig struct {
	DiscovType string   `json:",default=endpoints"` // endpoint, k8s
	Endpoints  []string `json:",optional"`
	K8s        struct {
		Namespace string `json:",default=default"`
		Port      int    `json:",default=6700"`
	} `json:",optional"`
}

type ConnPod struct {
	PodIp string
	ConnService
}

func (s *ConnPodsMgr) AllConnServices() []*ConnPod {
	var services []*ConnPod
	s.connPods.Range(func(key, value interface{}) bool {
		services = append(services, value.(*ConnPod))
		return true
	})
	return services
}

func (s *ConnPodsMgr) initConnRpc() {
	if s.Config.DiscovType == "k8s" {
		{
			s.endpointsEventHandler = discov.MustListenEndpoints(s.Config.K8s.Namespace, "conn", func(endpoints []string) {
				for _, endpoint := range endpoints {
					if _, ok := s.connPods.Load(endpoint); !ok {
						s.connPods.Store(endpoint, &ConnPod{
							PodIp: fmt.Sprintf("%s:%d", endpoint, s.Config.K8s.Port),
							ConnService: NewConnService(zrpc.MustNewClient(zrpc.RpcClientConf{
								Endpoints: []string{endpoint},
								NonBlock:  true,
							})),
						})
					}
				}
				// 删除不存在的
				s.connPods.Range(func(key, value interface{}) bool {
					exist := false
					for _, endpoint := range endpoints {
						if key == endpoint {
							exist = true
							break
						}
					}
					if !exist {
						s.connPods.Delete(key)
					}
					return true
				})
				// 列出所有的
				s.connPods.Range(func(endpoint, value interface{}) bool {
					fmt.Println("conn pod endpoint:", endpoint)
					return true
				})
			})
		}
	} else {
		for _, endpoint := range s.Config.Endpoints {
			s.connPods.Store(endpoint, &ConnPod{
				PodIp: endpoint,
				ConnService: NewConnService(zrpc.MustNewClient(zrpc.RpcClientConf{
					Endpoints: []string{endpoint},
					NonBlock:  true,
				})),
			})
		}
	}
}

func (s *ConnPodsMgr) ConnByIp(ip string) ConnService {
	if v, ok := s.connPods.Load(ip); ok {
		return v.(*ConnPod).ConnService
	}
	return nil
}
