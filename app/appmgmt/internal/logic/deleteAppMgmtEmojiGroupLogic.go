package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppMgmtEmojiGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAppMgmtEmojiGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppMgmtEmojiGroupLogic {
	return &DeleteAppMgmtEmojiGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAppMgmtEmojiGroupLogic) DeleteAppMgmtEmojiGroup(in *pb.DeleteAppMgmtEmojiGroupReq) (*pb.DeleteAppMgmtEmojiGroupResp, error) {
	model := &appmgmtmodel.EmojiGroup{}
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := tx.Model(model).Where("name in (?)", in.Names).Delete(model).Error
		if err != nil {
			l.Errorf("delete err: %v", err)
			return err
		}
		// 删除所有的表情
		emojiModel := &appmgmtmodel.Emoji{}
		err = tx.Model(emojiModel).Where("`group` in (?)", in.Names).Delete(emojiModel).Error
		if err != nil {
			l.Errorf("delete err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.DeleteAppMgmtEmojiGroupResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.DeleteAppMgmtEmojiGroupResp{}, nil
}
