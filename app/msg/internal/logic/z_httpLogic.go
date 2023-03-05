package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type HttpLogic struct {
	svcCtx *svc.ServiceContext
}

func (l *HttpLogic) Start() {
	// 监听80端口
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/syncMsgPart", func(w http.ResponseWriter, r *http.Request) {
		err := l.SyncMsgPart()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":80", nil)
}

func (l *HttpLogic) SyncMsgPart() error {
	// 读取 msg 表中的数据，同步到 msg_part 表中
	var latestMsgId string
	for {
		var msgs []*msgmodel.Msg
		err := l.svcCtx.Mysql().Model(&msgmodel.Msg{}).Where("id > ?", latestMsgId).Limit(1000).Find(&msgs).Error
		if err != nil {
			logx.Errorf("SyncMsgPart err: %v", err)
			return err
		}
		if len(msgs) == 0 {
			break
		}
		latestMsgId = msgs[len(msgs)-1].ServerMsgId
		err = msgmodel.InsertManyMsg(context.Background(), l.svcCtx.Mysql(), msgs)
		if err != nil {
			logx.Errorf("SyncMsgPart err: %v", err)
			return err
		}
	}
	return nil
}

func NewHttpLogic(svcCtx *svc.ServiceContext) *HttpLogic {
	return &HttpLogic{svcCtx: svcCtx}
}
