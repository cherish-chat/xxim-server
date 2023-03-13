package pb

func (x *SendMsgListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *ReportUserReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *ReportGroupReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *UpdateConvSettingReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *BatchKickGroupMemberReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *EditGroupInfoReq) SetCommonReq(r *CommonReq) {
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

func (x *ResetPasswordReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetGroupApplyListReq) SetCommonReq(r *CommonReq) {
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

func (x *GetGroupHomeReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetGroupMemberListReq) SetCommonReq(r *CommonReq) {
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

func (x *GetConvSettingReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetFriendListReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *GetLatestVersionReq) SetCommonReq(r *CommonReq) {
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

func (x *UpdateUserPasswordReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SearchGroupsByKeywordReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *AppGetAllConfigReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *SendSmsReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}

func (x *VerifySmsReq) SetCommonReq(r *CommonReq) {
	x.CommonReq = r
}
