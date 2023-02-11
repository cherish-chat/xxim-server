package pb

func (x *SendMsgListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *BatchGetMsgListByConvIdReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
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

func (x *GetUserHomeReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetUserSettingsReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SetUserSettingsReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *RequestAddFriendReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *AcceptAddFriendReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
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

func (x *BatchGetConvSeqReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetMyFriendEventListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *AckNoticeDataReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *UpdateUserInfoReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetMyGroupListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SetGroupMemberInfoReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetGroupMemberInfoReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *ApplyToBeGroupMemberReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *HandleGroupApplyReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *KickGroupMemberReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *ReadMsgReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *EditMsgReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *KeepAliveReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SetCxnParamsReq) SetCommonReq(r *CommonReq) {
}

func (x *SetCxnParamsReq) GetCommonReq() *CommonReq {
	return &CommonReq{}
}

func (x *SetCxnParamsReq) Validate() error {
	return nil
}

func (x *SetCxnParamsResp) GetCommonResp() *CommonResp {
	return &CommonResp{}
}

func (x *SetUserParamsReq) SetCommonReq(r *CommonReq) {
}

func (x *SetUserParamsReq) GetCommonReq() *CommonReq {
	return &CommonReq{}
}

func (x *SetUserParamsReq) Validate() error {
	return nil
}

func (x *SetUserParamsResp) GetCommonResp() *CommonResp {
	return &CommonResp{
		Code: 0,
		Msg:  nil,
		Data: nil,
	}
}

func (x *RegisterReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}
