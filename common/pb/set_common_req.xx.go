package pb

func (x *SendMsgListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SendMsgListReq) Path() string {
	return "msg/SendMsgListAsync"
}

func (x *BatchGetMsgListByConvIdReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *BatchGetMsgListByConvIdReq) Path() string {
	return "msg/SyncMsgList"
}

func (x *CreateGroupReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *LoginReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *ConfirmRegisterReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SearchUsersByKeywordReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SearchUsersByKeywordReq) Path() string {
	return "user/SearchUsersByKeyword"
}

func (x *GetUserHomeReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetUserSettingsReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetUserSettingsReq) Path() string {
	return "user/GetUserSettings"
}

func (x *SetUserSettingsReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SetUserSettingsReq) Path() string {
	return "user/SetUserSettings"
}

func (x *RequestAddFriendReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *RequestAddFriendReq) Path() string {
	return "relation/RequestAddFriend"
}

func (x *AcceptAddFriendReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *AcceptAddFriendReq) Path() string {
	return "relation/AcceptAddFriend"
}

func (x *RejectAddFriendReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *BlockUserReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *DeleteBlockUserReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *DeleteFriendReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SetSingleConvSettingReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetSingleConvSettingReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetFriendListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetMsgByIdReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetMsgByIdReq) Path() string {
	return "msg/GetMsgById"
}

func (x *BatchGetConvSeqReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *BatchGetConvSeqReq) Path() string {
	return "msg/BatchGetConvSeq"
}

func (x *GetMyFriendEventListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetMyFriendEventListReq) Path() string {
	return "relation/GetMyFriendEventList"
}

func (x *GetAppSystemConfigReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *AckNoticeDataReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *AckNoticeDataReq) Path() string {
	return "msg/AckNoticeData"
}
