package tdmq

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	v20200217 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/zeromicro/go-zero/core/logx"
)

type Mgr struct {
	Config        *Config
	tencentClient *v20200217.Client
}

func NewMgr(config *Config) *Mgr {
	m := &Mgr{Config: config}
	tencentClient, err := v20200217.NewClient(common.NewCredential(m.Config.SecretId, m.Config.SecretKey), m.Config.Region, profile.NewClientProfile())
	if err != nil {
		logx.Errorf("Could not instantiate tencent client: %v", err)
		panic(err)
	}
	m.tencentClient = tencentClient
	// 查询集群
	if m.Config.ClusterId == "" {
		request := v20200217.NewDescribeClustersRequest()
		response, err := m.tencentClient.DescribeClusters(request)
		if err != nil {
			logx.Errorf("Could not DescribeClusters: %v", err)
			panic(err)
		}
		clusterNameExist := false
		for _, cluster := range response.Response.ClusterSet {
			logx.Infof("cluster: %s %s", *cluster.ClusterId, *cluster.ClusterName)
			if *cluster.ClusterName == m.Config.ClusterName {
				clusterNameExist = true
				m.Config.ClusterId = *cluster.ClusterId
				break
			}
		}
		if !clusterNameExist {
			// 创建集群
			request := v20200217.NewCreateClusterRequest()
			request.ClusterName = common.StringPtr(m.Config.ClusterName)
			request.Remark = common.StringPtr("auto create by xxim-server")
			response, err := m.tencentClient.CreateCluster(request)
			if err != nil {
				logx.Errorf("Could not CreateCluster: %v", err)
				panic(err)
			}
			m.Config.ClusterId = *response.Response.ClusterId
			logx.Infof("create cluster: %s %s", m.Config.ClusterId, m.Config.ClusterName)
		}
	}
	// 查询命名空间
	{
		request := v20200217.NewDescribeEnvironmentsRequest()
		request.ClusterId = common.StringPtr(m.Config.ClusterId)
		response, err := m.tencentClient.DescribeEnvironments(request)
		if err != nil {
			logx.Errorf("Could not DescribeEnvironments: %v", err)
			panic(err)
		}
		namespaceExist := false
		for _, namespace := range response.Response.EnvironmentSet {
			logx.Infof("namespace: %s %s", *namespace.EnvironmentId, *namespace.NamespaceName)
			if *namespace.NamespaceName == m.Config.Namespace {
				namespaceExist = true
				m.Config.Namespace = *namespace.EnvironmentId
				break
			}
		}
		if !namespaceExist {
			// 创建命名空间
			request := v20200217.NewCreateEnvironmentRequest()
			request.EnvironmentId = common.StringPtr(m.Config.Namespace)
			request.ClusterId = common.StringPtr(m.Config.ClusterId)
			request.Remark = common.StringPtr("auto create by xxim-server")
			request.MsgTTL = common.Uint64Ptr(60 * 60 * 24 * 7)
			response, err := m.tencentClient.CreateEnvironment(request)
			if err != nil {
				logx.Errorf("Could not CreateEnvironment: %v", err)
				panic(err)
			}
			m.Config.Namespace = *response.Response.EnvironmentId
			logx.Infof("create namespace: %s %s", m.Config.Namespace, m.Config.Namespace)
		}
	}
	// token
	if m.Config.Token == "" {
		request := v20200217.NewDescribeRolesRequest()
		request.RoleName = common.StringPtr("xxim")
		request.ClusterId = common.StringPtr(m.Config.ClusterId)
		response, err := m.tencentClient.DescribeRoles(request)
		if err != nil {
			logx.Errorf("Could not DescribeRoles: %v", err)
			panic(err)
		}
		if len(response.Response.RoleSets) == 0 {
			// 创建角色
			request := v20200217.NewCreateRoleRequest()
			request.ClusterId = common.StringPtr(m.Config.ClusterId)
			request.RoleName = common.StringPtr("xxim")
			request.Remark = common.StringPtr("auto create by xxim-server")
			response, err := m.tencentClient.CreateRole(request)
			if err != nil {
				logx.Errorf("Could not CreateRole: %v", err)
				panic(err)
			}
			logx.Infof("create role: %s", *response.Response.RoleName)
		} else {
			m.Config.Token = *response.Response.RoleSets[0].Token
		}
		// 命名空间关联角色
		{
			request := v20200217.NewDescribeEnvironmentRolesRequest()
			request.ClusterId = common.StringPtr(m.Config.ClusterId)
			request.EnvironmentId = common.StringPtr(m.Config.Namespace)
			response, err := m.tencentClient.DescribeEnvironmentRoles(request)
			if err != nil {
				logx.Errorf("Could not DescribeEnvironmentRoles: %v", err)
				panic(err)
			}
			roleExist := false
			for _, role := range response.Response.EnvironmentRoleSets {
				logx.Infof("role: %s %s", *role.RoleName, *role.RoleName)
				if *role.RoleName == "xxim" {
					roleExist = true
					break
				}
			}
			if !roleExist {
				request := v20200217.NewCreateEnvironmentRoleRequest()
				request.EnvironmentId = common.StringPtr(m.Config.Namespace)
				request.RoleName = common.StringPtr("xxim")
				request.ClusterId = common.StringPtr(m.Config.ClusterId)
				request.Permissions = []*string{common.StringPtr("produce"), common.StringPtr("consume")}
				_, err := m.tencentClient.CreateEnvironmentRole(request)
				if err != nil {
					logx.Errorf("Could not CreateEnvironmentRole: %v", err)
					panic(err)
				}
			}
		}
		if m.Config.Token == "" {
			request := v20200217.NewDescribeRolesRequest()
			request.RoleName = common.StringPtr("xxim")
			request.ClusterId = common.StringPtr(m.Config.ClusterId)
			response, err := m.tencentClient.DescribeRoles(request)
			if err != nil {
				logx.Errorf("Could not DescribeRoles: %v", err)
				panic(err)
			}
			m.Config.Token = *response.Response.RoleSets[0].Token
		}
	}
	logx.Infof("config: %+v", m.Config)
	return m
}

type Config struct {
	Namespace   string // 命名空间
	ClusterId   string `json:",optional"` // 集群名称
	ClusterName string // 集群名称
	SecretId    string // 腾讯云 SecretId
	SecretKey   string // 腾讯云 SecretKey
	Region      string // 腾讯云 地域 ap-guangzhou/ap-shanghai/ap-beijing/ap-hongkong ...
	Token       string `json:",optional"` // 腾讯云 Token
}

func (m *Mgr) CreateTopic(topicName string, partition uint64) error {
	// 判断topic是否存在
	{
		request := v20200217.NewDescribeTopicsRequest()
		request.EnvironmentId = common.StringPtr(m.Config.Namespace)
		request.TopicName = common.StringPtr(topicName)
		request.ClusterId = common.StringPtr(m.Config.ClusterId)
		request.Filters = nil
		request.Offset = nil
		request.Limit = common.Uint64Ptr(20)
		response, err := m.tencentClient.DescribeTopics(request)
		if err != nil {
			logx.Errorf("Could not DescribeTopics: %v", err)
			return err
		}
		if *response.Response.TotalCount == 0 {
			logx.Infof("topic not exist: %s. auto create", topicName)
			// 创建topic
			request := v20200217.NewCreateTopicRequest()
			request.EnvironmentId = common.StringPtr(m.Config.Namespace)
			request.TopicName = common.StringPtr(topicName)
			request.Partitions = common.Uint64Ptr(partition)
			request.Remark = common.StringPtr("auto create by xxim-server")
			request.ClusterId = common.StringPtr(m.Config.ClusterId)
			request.PulsarTopicType = common.Int64Ptr(3) // 持久分区
			_, err := m.tencentClient.CreateTopic(request)
			if err != nil {
				logx.Errorf("Could not CreateTopic: %v", err)
				return err
			}
		}
	}
	return nil
}

func (m *Mgr) CreateSubscriptionBroadcast(topicName string, subscriptionName string, maxPodNum int) error {
	for i := 0; i < maxPodNum; i++ {
		subName := fmt.Sprintf("%s-%d", subscriptionName, i)
		err := m.CreateSubscription(topicName, subName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mgr) CreateSubscription(topicName string, subscriptionName string) error {
	// 判断是否创建了订阅关系
	// 列出所有消费订阅
	request := v20200217.NewDescribeSubscriptionsRequest()
	{
		request.EnvironmentId = common.StringPtr(m.Config.Namespace)
		request.TopicName = common.StringPtr(topicName)
		request.Offset = nil
		request.Limit = common.Uint64Ptr(1)
		request.SubscriptionName = common.StringPtr(subscriptionName)
		request.Filters = nil
		request.ClusterId = common.StringPtr(m.Config.ClusterId)
	}
	describeSubscriptionsResponse, err := m.tencentClient.DescribeSubscriptions(request)
	if err != nil {
		logx.Errorf("Could not describe subscriptions: %v", err)
		return err
	}
	if *describeSubscriptionsResponse.Response.TotalCount == 0 {
		// 创建订阅关系
		request := v20200217.NewCreateSubscriptionRequest()
		{
			request.EnvironmentId = common.StringPtr(m.Config.Namespace)
			request.TopicName = common.StringPtr(topicName)
			request.SubscriptionName = common.StringPtr(subscriptionName)
			request.IsIdempotent = common.BoolPtr(true)
			request.Remark = common.StringPtr("auto create by xxim-server")
			request.ClusterId = common.StringPtr(m.Config.ClusterId)
			request.AutoCreatePolicyTopic = common.BoolPtr(true)
			request.PostFixPattern = common.StringPtr("COMMUNITY")
		}
		_, err := m.tencentClient.CreateSubscription(request)
		if err != nil {
			logx.Errorf("Could not create subscription: %v", err)
			return err
		}
	}
	return nil
}
