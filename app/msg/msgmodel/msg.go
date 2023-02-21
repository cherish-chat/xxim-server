package msgmodel

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"time"
)

type (
	Msg struct {
		// 服务端生成的消息id convId+seq
		ServerMsgId string `bson:"_id" gorm:"column:id;primary_key;type:char(128);"`
		// 会话id // 单聊：sender_id + receiver_id // 群聊：group_id
		ConvId string `bson:"convId" gorm:"column:convId;type:char(96);index;index:servertime_convid;"`
		// 客户端生成的消息id
		ClientMsgId string `bson:"clientMsgId" gorm:"column:clientMsgId;type:char(128);index;"`
		// 客户端发送消息的时间 13位时间戳
		ClientTime int64 `bson:"clientTime" gorm:"column:clientTime;type:bigint;index,sort:desc;;"`
		// 服务端接收到消息的时间 13位时间戳 index DESC
		ServerTime int64 `bson:"serverTime" gorm:"column:serverTime;type:bigint;index,sort:desc;index:servertime_convid;"`
		// 发送者id
		SenderId string `bson:"senderId" gorm:"column:senderId;type:char(32);index;"`
		// 发送者信息
		SenderInfo []byte `bson:"senderInfo" gorm:"column:senderInfo;type:blob;"`
		// 强提醒用户id列表 用户不在线时，会收到离线推送，除非用户屏蔽了该会话 如果需要提醒所有人，可以传入"all"
		AtUsers xorm.SliceString `bson:"atUsers" gorm:"column:atUsers;type:JSON;"`
		// 消息内容类型
		ContentType ContentType `bson:"contentType" gorm:"column:contentType;type:tinyint;"`
		// 消息内容
		Content []byte `bson:"content" gorm:"column:content;type:blob;"`
		// 消息序号 会话内唯一且递增
		Seq int64 `bson:"seq" gorm:"column:seq;type:bigint;index;"`
		// 消息选项
		Options MsgOptions `bson:"options" gorm:"column:options;type:JSON;"`
		// 离线推送
		OfflinePush *MsgOfflinePush `bson:"offlinePush,omitempty" gorm:"column:offlinePush;type:JSON;"`
		// 扩展字段
		Ext []byte `bson:"ext" gorm:"column:ext;type:blob;"`
		// internal
		internal MsgInternal `bson:"-" gorm:"-"`
	}
	MsgInternal struct {
		NotFound bool
	}
	ContentType = pb.ContentType
	MsgOptions  struct {
		OfflinePush       bool `bson:"offlinePush"`       // 是否需要离线推送
		StorageForServer  bool `bson:"storageForServer"`  // 服务端是否需要保存消息
		StorageForClient  bool `bson:"storageForClient"`  // 客户端是否需要保存消息
		UpdateUnreadCount bool `bson:"updateUnreadCount"` // 消息是否需要计入未读数
		NeedDecrypt       bool `bson:"needDecrypt"`       // 是否需要解密 （端对端加密技术，服务端无法解密）
		UpdateConvMsg     bool `bson:"updateConvMsg"`     // 是否需要重新渲染会话
	}
	MsgOfflinePush struct {
		Title   string `bson:"title"`   // 离线推送标题
		Content string `bson:"content"` // 离线推送内容
		Payload string `bson:"payload"` // 离线推送自定义字段
	}
)

func (m *Msg) TableName() string {
	return "msg"
}

func NewMsgFromPb(in *pb.MsgData) *Msg {
	if in.Options == nil {
		in.Options = &pb.MsgData_Options{}
	}
	if in.OfflinePush == nil {
		in.OfflinePush = &pb.MsgData_OfflinePush{}
	}
	return &Msg{
		ServerMsgId: in.ServerMsgId,
		ConvId:      in.ConvId,
		ClientMsgId: in.ClientMsgId,
		ClientTime:  utils.AnyToInt64(in.ClientTime),
		ServerTime:  utils.AnyToInt64(in.ServerTime),
		SenderId:    in.SenderId,
		SenderInfo:  in.SenderInfo,
		AtUsers:     utils.AnyMakeSlice(in.AtUsers),
		ContentType: ContentType(in.ContentType),
		Content:     in.Content,
		Seq:         utils.AnyToInt64(in.Seq),
		Options: MsgOptions{
			OfflinePush:       in.Options.OfflinePush,
			StorageForServer:  in.Options.StorageForServer,
			StorageForClient:  in.Options.StorageForClient,
			UpdateUnreadCount: in.Options.UpdateUnreadCount,
			NeedDecrypt:       in.Options.NeedDecrypt,
			UpdateConvMsg:     in.Options.UpdateConvMsg,
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
	convId, seq := pb.ParseConvServerMsgId(serverId)
	m.ConvId = convId
	m.Seq = seq
	m.internal.NotFound = true
}

func (m *Msg) IsNotFound() bool {
	return m.internal.NotFound
}

func (m *Msg) SetSeq(seq int64) *Msg {
	m.Seq = seq
	m.ServerMsgId = pb.ServerMsgId(m.ConvId, seq)
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
		ServerMsgId: m.ServerMsgId,
		ConvId:      m.ConvId,
		ClientMsgId: m.ClientMsgId,
		ClientTime:  utils.AnyToString(m.ClientTime),
		ServerTime:  utils.AnyToString(m.ServerTime),
		SenderId:    m.SenderId,
		SenderInfo:  m.SenderInfo,
		AtUsers:     m.AtUsers,
		ContentType: int32(m.ContentType),
		Content:     m.Content,
		Seq:         utils.AnyToString(m.Seq),
		Options: &pb.MsgData_Options{
			OfflinePush:       m.Options.OfflinePush,
			StorageForServer:  m.Options.StorageForServer,
			StorageForClient:  m.Options.StorageForClient,
			UpdateUnreadCount: m.Options.UpdateUnreadCount,
			NeedDecrypt:       m.Options.NeedDecrypt,
			UpdateConvMsg:     m.Options.UpdateConvMsg,
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

func GetMsgTableNameById(id string) string {
	convId, seq := pb.ParseConvServerMsgId(id)
	md5 := utils.Md5(convId)
	// 取后2位
	tableName := fmt.Sprintf("msg_%s", md5[len(md5)-2:])
	// 每100000 seq分一个表
	tableName = fmt.Sprintf("%s_%d", tableName, seq/100000)
	return tableName
}

var tableNameExists = make(map[string]bool)

func CreateMsgTable(tx *gorm.DB, tableName string) error {
	if len(tableNameExists) == 0 {
		// 获取所有表名
		rows, err := tx.Raw("show tables").Rows()
		if err != nil {
			logx.Errorf("CreateMsgTable error: %v", err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				logx.Errorf("CreateMsgTable error: %v", err)
				return err
			}
			tableNameExists[name] = true
		}
	}
	if _, ok := tableNameExists[tableName]; ok {
		return nil
	}
	err := tx.Table(tableName).Migrator().CreateTable(&Msg{})
	if err != nil {
		logx.Errorf("CreateMsgTable error: %v", err)
		return err
	}
	tableNameExists[tableName] = true
	return nil
}

func InsertManyMsg(ctx context.Context, tx *gorm.DB, models []*Msg) error {
	if len(models) == 0 {
		return nil
	}
	// 分表
	var tableModels = make(map[string][]*Msg)
	for _, model := range models {
		tableName := GetMsgTableNameById(model.ServerMsgId)
		if _, ok := tableModels[tableName]; !ok {
			tableModels[tableName] = make([]*Msg, 0)
		}
		tableModels[tableName] = append(tableModels[tableName], model)
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		for tableName, models := range tableModels {
			err := CreateMsgTable(tx, tableName)
			if err != nil {
				return err
			}
			err = tx.Table(tableName).Create(models).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func MsgFromMysql(
	ctx context.Context,
	rc *zedis.Redis,
	tx *gorm.DB,
	ids []string,
) (msgList []*Msg, err error) {
	if len(ids) == 0 {
		return make([]*Msg, 0), nil
	}
	xtrace.StartFuncSpan(ctx, "FindMsgByIds", func(ctx context.Context) {
		// tableNameIds
		var tableNameIds = make(map[string][]string)
		for _, id := range ids {
			tableName := GetMsgTableNameById(id)
			if _, ok := tableNameIds[tableName]; !ok {
				tableNameIds[tableName] = make([]string, 0)
			}
			tableNameIds[tableName] = append(tableNameIds[tableName], id)
		}
		// 查询
		for tableName, ids := range tableNameIds {
			var tmpMsgList []*Msg
			err = tx.Table(tableName).Where("id in (?)", ids).Limit(len(ids)).Find(&tmpMsgList).Error
			if err != nil {
				return
			}
			msgList = append(msgList, tmpMsgList...)
		}
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

func (m Msg) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Msg) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}

func (m MsgOptions) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MsgOptions) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}

func (m MsgOfflinePush) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MsgOfflinePush) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}
