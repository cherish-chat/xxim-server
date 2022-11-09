package msgmodel

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

// ConvIdSeparator 会话id之间的分隔符
const ConvIdSeparator = ":"

type (
	Msg struct {
		ServerMsgId    string          `bson:"_id"`                   // 服务端生成的消息id convId+seq
		ConvId         string          `bson:"conv_id"`               // 会话id // 单聊：sender_id + receiver_id // 群聊：group_id
		ClientMsgId    string          `bson:"clientMsgId"`           // 客户端生成的消息id
		ClientTime     int64           `bson:"clientTime"`            // 客户端发送消息的时间 13位时间戳
		ServerTime     int64           `bson:"serverTime"`            // 服务端接收到消息的时间 13位时间戳
		Sender         string          `bson:"sender"`                // 发送者id
		SenderInfo     string          `bson:"senderInfo"`            // 发送者信息
		SenderConvInfo string          `bson:"senderConvInfo"`        // 发送者在会话中的信息
		Receiver       MsgReceiver     `bson:"receiver"`              // 接收者id (单聊时为对方id, 群聊时为群id)
		AtUsers        []string        `bson:"atUsers"`               // 强提醒用户id列表 用户不在线时，会收到离线推送，除非用户屏蔽了该会话 如果需要提醒所有人，可以传入"all"
		ContentType    ContentType     `bson:"contentType"`           // 消息内容类型
		Content        []byte          `bson:"content"`               // 消息内容
		Seq            int64           `bson:"seq"`                   // 消息序号 会话内唯一且递增
		Options        MsgOptions      `bson:"options"`               // 消息选项
		OfflinePush    *MsgOfflinePush `bson:"offlinePush,omitempty"` // 离线推送
		Ext            []byte          `bson:"ext"`                   // 扩展字段
	}
	BatchMsg struct {
		Id          string   `bson:"_id"` // 批量id
		Msg         *Msg     `bson:"msg"` // 原始消息
		UserIdList  []string `bson:"userIdList"`
		GroupIdList []string `bson:"groupIdList"`
	}
	ContentType = pb.ContentType
	MsgReceiver struct {
		UserId  string `bson:"userId"`  // 单聊时为对方的userId
		GroupId string `bson:"groupId"` // 群聊时为群组id
	}
	MsgOptions struct {
		OfflinePush      bool `bson:"offlinePush"`      // 是否需要离线推送
		StorageForServer bool `bson:"storageForServer"` // 服务端是否需要保存消息
		StorageForClient bool `bson:"storageForClient"` // 客户端是否需要保存消息
		UnreadCount      bool `bson:"unreadCount"`      // 消息是否需要计入未读数
		NeedDecrypt      bool `bson:"needDecrypt"`      // 是否需要解密 （端对端加密技术，服务端无法解密）
		UpdateConv       bool `bson:"updateConv"`       // 是否需要重新渲染会话
	}
	MsgOfflinePush struct {
		Title   string `bson:"title"`   // 离线推送标题
		Content string `bson:"content"` // 离线推送内容
		Payload string `bson:"payload"` // 离线推送自定义字段
	}
)

func (m *BatchMsg) CollectionName() string {
	return "batch_msg"
}

func (m *BatchMsg) Indexes(c *qmgo.Collection) error {
	return nil
}

func (m *Msg) CollectionName() string {
	return "msg"
}

func (m *Msg) Indexes(c *qmgo.Collection) error {
	c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key: []string{"clientMsgId"},
	}})
	return nil
}

func NewMsgFromPb(in *pb.MsgData) *Msg {
	userId, groupId := "", ""
	if in.Receiver != nil {
		if in.Receiver.UserId != nil {
			userId = *in.Receiver.UserId
		}
		if in.Receiver.GroupId != nil {
			groupId = *in.Receiver.GroupId
		}
	}
	return &Msg{
		ServerMsgId:    in.ServerMsgId,
		ConvId:         in.ConvId,
		ClientMsgId:    in.ClientMsgId,
		ClientTime:     utils.AnyToInt64(in.ClientTime),
		ServerTime:     utils.AnyToInt64(in.ServerTime),
		Sender:         in.Sender,
		SenderInfo:     in.SenderInfo,
		SenderConvInfo: in.SenderConvInfo,
		Receiver: MsgReceiver{
			UserId:  userId,
			GroupId: groupId,
		},
		AtUsers:     utils.AnyMakeSlice(in.AtUsers),
		ContentType: in.ContentType,
		Content:     in.Content,
		Seq:         utils.AnyToInt64(in.Seq),
		Options: MsgOptions{
			OfflinePush:      in.Options.OfflinePush,
			StorageForServer: in.Options.StorageForServer,
			StorageForClient: in.Options.StorageForClient,
			UnreadCount:      in.Options.UnreadCount,
			NeedDecrypt:      in.Options.NeedDecrypt,
			UpdateConv:       in.Options.UpdateConv,
		},
		OfflinePush: &MsgOfflinePush{
			Title:   in.OfflinePush.Title,
			Content: in.OfflinePush.Content,
			Payload: in.OfflinePush.Payload,
		},
		Ext: in.Ext,
	}
}

func (m *Msg) NotFound(serverId string) {
	m.ServerMsgId = serverId
	convId, seq := ParseConvServerMsgId(serverId)
	m.ConvId = convId
	m.Seq = seq
}

func (m *Msg) AutoConvId() *Msg {
	if m.Receiver.GroupId == "" {
		// 单聊
		m.ConvId = SingleConvId(m.Sender, m.Receiver.UserId)
	} else {
		// 群聊
		m.ConvId = m.Receiver.GroupId
	}
	return m
}

func (m *Msg) SetSeq(seq int64) *Msg {
	m.Seq = seq
	m.ServerMsgId = ServerMsgId(m.ConvId, seq)
	return m
}

func (m *Msg) Check() *Msg {
	if m.ServerTime == 0 {
		m.ServerTime = time.Now().UnixMilli()
	}
	if m.ClientTime == 0 {
		m.ClientTime = m.ServerTime
	}
	if m.ClientMsgId == "" {
		m.ClientMsgId = m.ServerMsgId
	}
	return m
}

func (m *Msg) ToMsgData() *pb.MsgData {
	offlinePush := m.OfflinePush
	if offlinePush == nil {
		offlinePush = &MsgOfflinePush{}
	}
	return &pb.MsgData{
		ServerMsgId:    m.ServerMsgId,
		ConvId:         m.ConvId,
		ClientMsgId:    m.ClientMsgId,
		ClientTime:     utils.AnyToString(m.ClientTime),
		ServerTime:     utils.AnyToString(m.ServerTime),
		Sender:         m.Sender,
		SenderInfo:     m.SenderInfo,
		SenderConvInfo: m.SenderConvInfo,
		Receiver: &pb.MsgData_Receiver{
			UserId:  &m.Receiver.UserId,
			GroupId: &m.Receiver.GroupId,
		},
		AtUsers:     m.AtUsers,
		ContentType: m.ContentType,
		Content:     m.Content,
		Seq:         utils.AnyToString(m.Seq),
		Options: &pb.MsgData_Options{
			OfflinePush:      m.Options.OfflinePush,
			StorageForServer: m.Options.StorageForServer,
			StorageForClient: m.Options.StorageForClient,
			UnreadCount:      m.Options.UnreadCount,
			NeedDecrypt:      m.Options.NeedDecrypt,
			UpdateConv:       m.Options.UpdateConv,
		},
		OfflinePush: &pb.MsgData_OfflinePush{
			Title:   offlinePush.Title,
			Content: offlinePush.Content,
			Payload: offlinePush.Payload,
		},
		Ext: m.Ext,
	}
}

func (m *Msg) ExpireSeconds() int {
	return xredis.ExpireMinutes(5)
}

func (m *Msg) SinglePb(seq int64, uid string) *pb.MsgData {
	convId := SingleConvId(uid, m.Sender)
	return &pb.MsgData{
		ClientMsgId:    BatchMsgClientMsgId(m.ClientMsgId, convId),
		ServerMsgId:    ServerMsgId(convId, seq),
		ClientTime:     utils.AnyToString(m.ClientTime),
		ServerTime:     utils.AnyToString(m.ServerTime),
		Sender:         m.Sender,
		SenderInfo:     m.SenderInfo,
		SenderConvInfo: m.SenderConvInfo,
		Receiver:       &pb.MsgData_Receiver{UserId: &uid},
		ConvId:         convId,
		AtUsers:        m.AtUsers,
		ContentType:    m.ContentType,
		Content:        m.Content,
		Seq:            utils.AnyToString(seq),
		Options: &pb.MsgData_Options{
			OfflinePush:      m.Options.OfflinePush,
			StorageForServer: m.Options.StorageForServer,
			StorageForClient: m.Options.StorageForClient,
			UnreadCount:      m.Options.UnreadCount,
			NeedDecrypt:      m.Options.NeedDecrypt,
			UpdateConv:       m.Options.UpdateConv,
		},
		OfflinePush: &pb.MsgData_OfflinePush{
			Title:   m.OfflinePush.Title,
			Content: m.OfflinePush.Content,
			Payload: m.OfflinePush.Payload,
		},
		Ext: m.Ext,
	}
}

func (m *Msg) GroupPb(seq int64, groupId string) *pb.MsgData {
	convId := groupId
	return &pb.MsgData{
		ClientMsgId:    BatchMsgClientMsgId(m.ClientMsgId, convId),
		ServerMsgId:    ServerMsgId(convId, seq),
		ClientTime:     utils.AnyToString(m.ClientTime),
		ServerTime:     utils.AnyToString(m.ServerTime),
		Sender:         m.Sender,
		SenderInfo:     m.SenderInfo,
		SenderConvInfo: m.SenderConvInfo,
		Receiver:       &pb.MsgData_Receiver{GroupId: &groupId},
		ConvId:         convId,
		AtUsers:        m.AtUsers,
		ContentType:    m.ContentType,
		Content:        m.Content,
		Seq:            utils.AnyToString(seq),
		Options: &pb.MsgData_Options{
			OfflinePush:      m.Options.OfflinePush,
			StorageForServer: m.Options.StorageForServer,
			StorageForClient: m.Options.StorageForClient,
			UnreadCount:      m.Options.UnreadCount,
			NeedDecrypt:      m.Options.NeedDecrypt,
			UpdateConv:       m.Options.UpdateConv,
		},
		OfflinePush: &pb.MsgData_OfflinePush{
			Title:   m.OfflinePush.Title,
			Content: m.OfflinePush.Content,
			Payload: m.OfflinePush.Payload,
		},
		Ext: m.Ext,
	}
}

func ServerMsgId(convId string, seq int64) string {
	return convId + ConvIdSeparator + strconv.FormatInt(seq, 10)
}

func SingleConvId(id1 string, id2 string) string {
	if id1 < id2 {
		return id1 + ConvIdSeparator + id2
	}
	return id2 + ConvIdSeparator + id1
}

func ParseSingleServerMsgId(serverMsgId string) (convId string, seq int64) {
	arr := strings.Split(serverMsgId, ConvIdSeparator)
	if len(arr) == 3 {
		convId = arr[0] + ConvIdSeparator + arr[1]
		seq, _ = strconv.ParseInt(arr[2], 10, 64)
	}
	return
}

func BatchMsgClientMsgId(clientMsgId string, convId string) string {
	return clientMsgId + ConvIdSeparator + convId
}

func ParseGroupServerMsgId(serverMsgId string) (groupId string, seq int64) {
	arr := strings.Split(serverMsgId, ConvIdSeparator)
	if len(arr) == 2 {
		groupId = arr[0]
		seq, _ = strconv.ParseInt(arr[1], 10, 64)
	}
	return
}

func ParseConvServerMsgId(serverMsgId string) (convId string, seq int64) {
	arr := strings.Split(serverMsgId, ConvIdSeparator)
	if len(arr) == 2 {
		convId = arr[0]
		seq, _ = strconv.ParseInt(arr[1], 10, 64)
	} else if len(arr) == 3 {
		convId = arr[0] + ConvIdSeparator + arr[1]
		seq, _ = strconv.ParseInt(arr[2], 10, 64)
	}
	return
}

func ConvIsGroup(convId string) bool {
	return !strings.Contains(convId, ConvIdSeparator)
}

func ParseSingleConvId(convId string) (string, string) {
	arr := strings.Split(convId, ConvIdSeparator)
	if len(arr) == 2 {
		return arr[0], arr[1]
	}
	return "", ""
}

func MsgFromMongo(
	ctx context.Context,
	rc *zedis.Redis,
	collection *qmgo.Collection,
	ids []string,
) (msgList []*Msg, err error) {
	if len(ids) == 0 {
		return make([]*Msg, 0), nil
	}
	xtrace.StartFuncSpan(ctx, "FindMsgByIds", func(ctx context.Context) {
		err = collection.Find(ctx, bson.M{
			"_id": bson.M{"$in": ids},
		}).All(&msgList)
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("GetSingleMsgListBySeq failed, err: %v", err)
		return nil, err
	}
	msgMap := make(map[string]*Msg)
	for _, msg := range msgList {
		msgMap[msg.ServerMsgId] = msg
		// 存入redis
		redisMsg, _ := json.Marshal(msg)
		err = rc.SetexCtx(ctx, rediskey.MsgKey(msg.ServerMsgId), string(redisMsg), msg.ExpireSeconds())
		if err != nil {
			logx.WithContext(ctx).Errorf("redis Setex error: %v", err)
			continue
		}
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, ok := msgMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		// 占位符写入redis
		for _, id := range notFoundIds {
			err = rc.SetexCtx(ctx, rediskey.MsgKey(id), xredis.NotFound, xredis.ExpireMinutes(5))
			if err != nil {
				logx.WithContext(ctx).Errorf("redis Setex error: %v", err)
				continue
			}
		}
	}
	return msgList, nil
}

func FlushMsgCache(ctx context.Context, rc *zedis.Redis, ids []string) error {
	var err error
	if len(ids) > 0 {
		xtrace.StartFuncSpan(ctx, "DeleteCache", func(ctx context.Context) {
			redisKeys := utils.UpdateSlice(ids, func(v string) string {
				return rediskey.MsgKey(v)
			})
			_, err = rc.DelCtx(ctx, redisKeys...)
		})
	}
	return err
}
