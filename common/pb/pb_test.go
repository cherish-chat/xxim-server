package pb

import (
	_ "embed"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/protobuf/proto"
	"testing"
)

//go:embed data.uu
var testData []byte

func TestParseBytes(t *testing.T) {
	m := &RequestBody{}
	err := proto.Unmarshal(testData, m)
	if err != nil {
		t.Fatal(err)
	}
	req := &EditGroupInfoReq{}
	err = proto.Unmarshal(m.Data, req)
	if err != nil {
		t.Fatal(err)
	}

	{
		req := &EditGroupInfoReq{
			CommonReq:          nil,
			GroupId:            "100068",
			Name:               nil,
			Avatar:             nil,
			Introduction:       nil,
			AllMute:            nil,
			MemberCanAddFriend: nil,
			CanAddMember:       utils.AnyPtr(true),
		}
		data, _ := proto.Marshal(req)
		m := &RequestBody{
			ReqId:  "client:p0:req:22",
			Method: "/v1/group/editGroupInfo",
			Data:   data,
		}
		bytes, _ := proto.Marshal(m)
		t.Log(string(bytes))
	}
}
