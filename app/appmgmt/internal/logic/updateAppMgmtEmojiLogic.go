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

type UpdateAppMgmtEmojiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtEmojiLogic {
	return &UpdateAppMgmtEmojiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtEmojiLogic) UpdateAppMgmtEmoji(in *pb.UpdateAppMgmtEmojiReq) (*pb.UpdateAppMgmtEmojiResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Emoji{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtEmoji.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtEmojiResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["name"] = in.AppMgmtEmoji.Name
		updateMap["cover"] = in.AppMgmtEmoji.Cover
		updateMap["type"] = in.AppMgmtEmoji.Type
		updateMap["staticUrl"] = in.AppMgmtEmoji.StaticUrl
		updateMap["animatedUrl"] = in.AppMgmtEmoji.AnimatedUrl
		updateMap["sort"] = in.AppMgmtEmoji.Sort
		updateMap["isEnable"] = in.AppMgmtEmoji.IsEnable
	}
	if len(updateMap) > 0 {
		err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			if in.AppMgmtEmoji.Cover {
				// 更新当前的封面
				err = tx.Model(&appmgmtmodel.EmojiGroup{}).Where("name = ?", in.AppMgmtEmoji.Group).Update("coverId", in.AppMgmtEmoji.Id).Error
				if err != nil {
					l.Errorf("更新失败: %v", err)
					return err
				}
			} else {
				if model.Cover {
					// update coverId = ""
					err = tx.Model(&appmgmtmodel.EmojiGroup{}).Where("name = ?", in.AppMgmtEmoji.Group).Update("coverId", "").Error
					if err != nil {
						l.Errorf("更新失败: %v", err)
						return err
					}
				}
			}
			// 把其他的封面都设置为false
			err := tx.Model(model).Where("`group` = ?", in.AppMgmtEmoji.Group).Update("cover", false).Error
			if err != nil {
				l.Errorf("更新失败: %v", err)
				return err
			}
			err = tx.Model(model).Where("id = ?", in.AppMgmtEmoji.Id).Updates(updateMap).Error
			if err != nil {
				l.Errorf("更新失败: %v", err)
				return err
			}
			return nil
		})
		if err != nil {
			return &pb.UpdateAppMgmtEmojiResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtEmojiResp{}, nil
}
