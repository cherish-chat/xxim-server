package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"regexp"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUsersByKeywordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUsersByKeywordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUsersByKeywordLogic {
	return &SearchUsersByKeywordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUsersByKeywordLogic) SearchUsersByKeyword(in *pb.SearchUsersByKeywordReq) (*pb.SearchUsersByKeywordResp, error) {
	found := &usermodel.User{}
	// 判断 keyword 是不是只包含 英文、数字、下划线
	if regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(in.Keyword) && len(in.Keyword) >= 6 && len(in.Keyword) <= 20 {
		// 使用id查询
		_ = xorm.DetailByWhere(l.svcCtx.Mysql(), found, xorm.Where("id = ?", in.Keyword))
	} else {
		// 使用昵称查询
		_ = xorm.DetailByWhere(l.svcCtx.Mysql(), found, xorm.Where("nickname LIKE ?", in.Keyword+"%"))
	}
	if found.Id == "" {
		// 没有找到
		return &pb.SearchUsersByKeywordResp{}, nil
	}
	resp := found.BaseInfo()
	latestConn, err := l.svcCtx.ImService().GetUserLatestConn(l.ctx, &pb.GetUserLatestConnReq{UserId: found.Id})
	if err != nil {
		l.Errorf("get user latest conn failed, err: %v", err)
	} else {
		resp.IpRegion = latestConn.IpRegion
	}
	return &pb.SearchUsersByKeywordResp{Users: []*pb.UserBaseInfo{resp}}, nil
}
