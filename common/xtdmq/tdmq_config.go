package xtdmq

import (
	"os"
	"strings"
)

type TDMQConfig struct {
	Token     string // 角色Token
	VpcUrl    string // 私有网络接入地址
	Namespace string // 命名空间
	ClusterId string // 集群名称
	SecretId  string `json:",optional"` // 腾讯云 SecretId
	SecretKey string `json:",optional"` // 腾讯云 SecretKey
	Region    string `json:",optional"` // 腾讯云 地域 ap-guangzhou/ap-shanghai/ap-beijing/ap-hongkong ...
}

type TDMQProducerConfig struct {
	TopicName    string // 主题名称
	SendTimeout  int64  `json:",default=3000"` // 发送超时时间
	ProducerName string `json:",optional"`     // 生产者名称
}

type TDMQConsumerConfig struct {
	TopicName    string // 主题名称
	SubName      string // 订阅名
	ConsumerName string `json:",optional"` // 消费者名称
	// 设置consumer初始接收消息的位置，可选参数为： 0:Latest 1:Earliest
	SubInitialPosition int `json:",default=0,options=0|1"`
	// 0:独占模式（Exclusive）（可创建多个Subscription，但每个Subscription只能被一个Consumer消费，以此实现广播模式）
	//一个 Subscription 只能与一个 Consumer 关联，只有这个 Consumer 可以接收到 Topic 的全部消息，如果该 Consumer 出现故障了就会停止消费。Exclusive 订阅模式下，同一个 Subscription 里只有一个 Consumer 能消费 Topic，如果多个 Consumer 订阅则会报错，适用于全局有序消费的场景。
	//
	// 1:共享模式（Shared）
	//消息通过 round robin 轮询机制（也可以自定义）分发给不同的消费者，并且每个消息仅会被分发给一个消费者。当消费者断开连接，所有被发送给他，但没有被确认的消息将被重新安排，分发给其它存活的消费者。
	//
	// 2:灾备模式（Failover）
	//当存在多个 consumer 时，将会按字典顺序排序，第一个 consumer 被初始化为唯一接受消息的消费者。当第一个 consumer 断开时，所有的消息（未被确认和后续进入的）将会被分发给队列中的下一个 consumer。
	//
	// 3:KEY 共享模式（Key_Shared）
	//当存在多个 consumer 时，将根据消息的 key 进行分发，key 相同的消息只会被分发到同一个消费者。
	SubType           int  `json:",default=0,options=0|1|2|3"`
	EnableRetry       bool `json:",default=true"`
	ReceiverQueueSize int  `json:",default=10"`    // 消费者接收队列大小
	IsBroadcast       bool `json:",default=false"` // 是否广播模式 如果开启广播模式 并且SubType为0 会判断环境变量中POD_NAME是否存在，如果存在则会自动设置SubName为SubName-POD_INDEX，如果不存在则会报错
}

func (c TDMQConsumerConfig) GetSubName() string {
	if c.IsBroadcast && c.SubType == 0 {
		if podName := os.Getenv("POD_NAME"); podName != "" {
			tmp := strings.Split(podName, "-")
			podNum := tmp[len(tmp)-1]
			return c.SubName + "-" + podNum
		}
		panic("env:POD_NAME is not set")
	}
	return c.SubName
}

func (c TDMQConsumerConfig) GetConsumerName() string {
	if c.IsBroadcast && c.SubType == 0 {
		if podName := os.Getenv("POD_NAME"); podName != "" {
			tmp := strings.Split(podName, "-")
			podNum := tmp[len(tmp)-1]
			return c.ConsumerName + "-" + podNum
		}
		panic("env:POD_NAME is not set")
	}
	return c.ConsumerName
}

func (c TDMQConsumerConfig) GetTopicName(clusterName string, namespace string) string {
	return clusterName + "/" + namespace + "/" + c.TopicName
}

func (c TDMQProducerConfig) GetProducerName() string {
	if podName := os.Getenv("POD_NAME"); podName != "" {
		tmp := strings.Split(podName, "-")
		podNum := tmp[len(tmp)-1]
		return c.ProducerName + "-" + podNum
	}
	return c.ProducerName
}
