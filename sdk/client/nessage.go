package client

import "github.com/cherish-chat/xxim-server/common/pb"

// MessageBatchSend 批量发送消息
func (c *HttpClient) MessageBatchSend(req *pb.MessageBatchSendReq) (resp *pb.MessageBatchSendResp, err error) {
	resp = &pb.MessageBatchSendResp{}
	err = c.Request("/v1/message/messageBatchSend", req, resp)
	return
}

// MessageBatchSend 批量发送消息
func (c *WsClient) MessageBatchSend(req *pb.MessageBatchSendReq) (resp *pb.MessageBatchSendResp, err error) {
	resp = &pb.MessageBatchSendResp{}
	err = c.Request("/v1/message/messageBatchSend", req, resp)
	return
}
