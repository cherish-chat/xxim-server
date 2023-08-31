package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"strconv"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGroupMemberLogic {
	return &SearchGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchGroupMember 搜索群成员
func (l *SearchGroupMemberLogic) SearchGroupMember(in *pb.SearchGroupMemberReq) (*pb.SearchGroupMemberResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{
			Page: 1,
			Size: 20,
		}
	}
	groupMemberModel := &groupmodel.GroupMember{}
	groupMemberTableName := groupMemberModel.TableName()
	userModel := &usermodel.User{}
	userTableName := userModel.TableName()
	/*
		SELECT gm.createTime, gm.groupId, gm.role, gm.remark, gm.unbanTime,
		       u.id, u.nickname, u.avatar
		FROM ${groupMemberTableName} as gm
		INNER JOIN ${userTableName} as u ON gm.userId = u.id
		WHERE gm.groupId = ? AND (u.nickname LIKE ? OR u.id IN (${orInUserIds})) AND gm.role IN (${mustRoles})
		ORDER BY u.id ASC
		LIMIT ?, ?
	*/
	args := []interface{}{in.GroupId}
	sqlBuilder := `SELECT gm.createTime, gm.groupId, gm.role, gm.remark, gm.unbanTime,
	       u.id, u.nickname, u.avatar
	FROM %s as gm
	INNER JOIN %s as u ON gm.userId = u.id
	WHERE gm.groupId = ? `
	if len(in.Keyword) > 0 && len(in.OrInUserIds) > 0 {
		sqlBuilder += `AND (u.nickname LIKE ? OR u.id LIKE ? OR u.id IN (%s)) `
		args = append(args, "%"+in.Keyword+"%")
		args = append(args, "%"+in.Keyword+"%")
		args = append(args, in.OrInUserIds)
	} else if len(in.Keyword) > 0 && len(in.OrInUserIds) == 0 {
		sqlBuilder += `AND (u.nickname LIKE ? OR u.id LIKE ?) `
		args = append(args, "%"+in.Keyword+"%")
		args = append(args, "%"+in.Keyword+"%")
	} else if len(in.Keyword) == 0 && len(in.OrInUserIds) > 0 {
		sqlBuilder += `AND u.id IN (%s) `
		args = append(args, in.OrInUserIds)
	}
	if len(in.MustRoles) > 0 {
		sqlBuilder += `AND gm.role IN (?) `
		args = append(args, in.MustRoles)
	}
	sqlBuilder += `ORDER BY u.id ASC LIMIT `
	sqlBuilder += strconv.Itoa(int((in.Page.Page - 1) * in.Page.Size))
	sqlBuilder += `, `
	sqlBuilder += strconv.Itoa(int(in.Page.Size))
	sql := fmt.Sprintf(sqlBuilder, groupMemberTableName, userTableName)
	type Result struct {
		CreateTime int64               `gorm:"column:createTime"`
		GroupId    string              `gorm:"column:groupId"`
		Role       groupmodel.RoleType `gorm:"column:role"`
		Remark     string              `gorm:"column:remark"`
		UnbanTime  int64               `gorm:"column:unbanTime"`
		Id         string              `gorm:"column:id"`
		Nickname   string              `gorm:"column:nickname"`
		Avatar     string              `gorm:"column:avatar"`
	}
	var results []Result
	l.Infof("sql: %s, args: %v", sql, args)
	err := l.svcCtx.Mysql().Raw(sql, args...).Scan(&results).Error
	if err != nil {
		l.Errorf("SearchGroupMember error: %v", err)
		return nil, err
	}
	var resp pb.SearchGroupMemberResp
	for _, result := range results {
		resp.GroupMemberList = append(resp.GroupMemberList, &pb.GroupMemberInfo{
			GroupId:      result.GroupId,
			MemberId:     result.Id,
			Remark:       result.Remark,
			Role:         pb.GroupRole(result.Role),
			UnbanTime:    result.UnbanTime,
			UserBaseInfo: &pb.UserBaseInfo{Id: result.Id, Nickname: result.Nickname, Avatar: result.Avatar},
		})
	}
	return &resp, nil
}
