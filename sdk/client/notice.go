package client

import "github.com/cherish-chat/xxim-server/common/pb"

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
