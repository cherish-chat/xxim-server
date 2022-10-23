package xtdmq

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	v20200217 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/zeromicro/go-zero/core/logx"
)

func GetApiClient(topicName string, config TDMQConfig) *v20200217.Client {
	tencentClient, err := v20200217.NewClient(common.NewCredential(config.SecretId, config.SecretKey), config.Region, profile.NewClientProfile())
	if err != nil {
		logx.Errorf("Could not instantiate tencent client: %v", err)
		panic(err)
	}
	// 判断topic是否存在
	{
		request := v20200217.NewDescribeTopicsRequest()
		request.EnvironmentId = common.StringPtr(config.Namespace)
		request.TopicName = common.StringPtr(topicName)
		request.ClusterId = common.StringPtr(config.ClusterId)
		request.Filters = nil
		request.Offset = nil
		request.Limit = common.Uint64Ptr(1)
		response, err := tencentClient.DescribeTopics(request)
		if err != nil {
			logx.Errorf("Could not DescribeTopics: %v", err)
			panic(err)
		}
		if *response.Response.TotalCount == 0 {
			logx.Infof("topic not exist: %s. auto create", topicName)
			// 创建topic
			request := v20200217.NewCreateTopicRequest()
			request.EnvironmentId = common.StringPtr(config.Namespace)
			request.TopicName = common.StringPtr(topicName)
			request.Partitions = common.Uint64Ptr(2) // 默认2个分区
			request.Remark = common.StringPtr("auto create by xxim-server")
			request.ClusterId = common.StringPtr(config.ClusterId)
			request.PulsarTopicType = common.Int64Ptr(3) // 持久分区
			_, err := tencentClient.CreateTopic(request)
			if err != nil {
				logx.Errorf("Could not CreateTopic: %v", err)
				panic(err)
			}
		}
	}
	return tencentClient
}
