package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

// ListNotice 列出通知
func (c *HttpClient) ListNotice(req *pb.ListNoticeReq) (resp *pb.ListNoticeResp, err error) {
	resp = &pb.ListNoticeResp{}
	err = c.Request("/v1/notice/listNotice", req, resp)
	return
}

// ListNotice 列出通知
func (c *WsClient) ListNotice(req *pb.ListNoticeReq) (resp *pb.ListNoticeResp, err error) {
	resp = &pb.ListNoticeResp{}
	err = c.Request("/v1/notice/listNotice", req, resp)
	return
}

func (c *WsClient) onNewNotice(notice *pb.Notice) {
	switch notice.ContentType {
	case pb.NoticeContentType_NewFriendRequest:
		logx.Infof("收到好友申请: ")
		listFriendApplyResp, err := c.httpClient.ListFriendApply(&pb.ListFriendApplyReq{
			Cursor: 0,
			Limit:  100,
			Filter: &pb.ListFriendApplyReq_Filter{
				Status: nil,
			},
			Option: &pb.ListFriendApplyReq_Option{
				IncludeApplyByMe: false,
			},
		})
		if err != nil {
			logx.Errorf("ListFriendApply error: %v", err)
			return
		}
		logx.Infof("listFriendApplyResp: %s", utils.AnyString(listFriendApplyResp))
	case pb.NoticeContentType_JoinNewGroup:
		logx.Infof("收到加入群组通知: ")
		//重新订阅群组
		c.KeepAlive()
	default:
		logx.Infof("onNewNotice: %s", utils.AnyString(notice))
	}
}
