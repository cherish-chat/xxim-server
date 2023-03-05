package msgservice

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/common/discov"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"sync"
)

func NewMsgPodsMgr(config MsgPodsConfig) *MsgPodsMgr {
	c := &MsgPodsMgr{Config: config}
	c.initMsgRpc()
	return c
}

type MsgPodsMgr struct {
	msgPods               sync.Map
	endpointsEventHandler *discov.EndpointsEventHandler
	Config                MsgPodsConfig
}

type MsgPodsConfig struct {
	DiscovType string   `json:",default=endpoints"` // endpoint, k8s
	Endpoints  []string `json:",optional"`
	K8s        struct {
		Namespace string `json:",default=default"`
		Port      int    `json:",default=6700"`
	} `json:",optional"`
}

type MsgPod struct {
	PodIpPort string
	MsgService
}

func (s *MsgPodsMgr) AllMsgServices() []*MsgPod {
	var services []*MsgPod
	s.msgPods.Range(func(key, value interface{}) bool {
		services = append(services, value.(*MsgPod))
		return true
	})
	return services
}

func (s *MsgPodsMgr) initMsgRpc() {
	if s.Config.DiscovType == "k8s" {
		{
			s.endpointsEventHandler = discov.MustListenEndpoints(s.Config.K8s.Namespace, "msg-rpc-svc", func(endpoints []string) {
				endpoints = utils.UpdateSlice(endpoints, func(endpoint string) string {
					return fmt.Sprintf("%s:%d", endpoint, s.Config.K8s.Port)
				})
				for _, endpoint := range endpoints {
					if _, ok := s.msgPods.Load(endpoint); !ok {
						s.msgPods.Store(endpoint, &MsgPod{
							PodIpPort: endpoint,
							MsgService: NewMsgService(zrpc.MustNewClient(
								zrpc.RpcClientConf{
									Endpoints: []string{endpoint},
									NonBlock:  true,
								},
								utils.Zrpc.Options()...,
							)),
						})
					}
				}
				// 删除不存在的
				s.msgPods.Range(func(key, value interface{}) bool {
					exist := false
					for _, endpoint := range endpoints {
						if key == endpoint {
							exist = true
							break
						}
					}
					if !exist {
						s.msgPods.Delete(key)
					}
					return true
				})
				// 列出所有的
				count := 0
				s.msgPods.Range(func(endpoint, value interface{}) bool {
					logx.Infof("msg pod endpoint: %s", endpoint)
					count++
					return true
				})
				logx.Infof("msg pod count: %d", count)
			})
		}
	} else {
		for _, endpoint := range s.Config.Endpoints {
			s.msgPods.Store(endpoint, &MsgPod{
				PodIpPort: endpoint,
				MsgService: NewMsgService(zrpc.MustNewClient(zrpc.RpcClientConf{
					Endpoints: []string{endpoint},
					NonBlock:  true,
				})),
			})
		}
	}
}
