package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

type getAllConnectionListReq struct {
	CommonReq *pb.CommonReq     `json:"commonReq"`
	Page      *pb.Page          `json:"page"`
	Filter    map[string]string `json:"filter"`
}

type getAllConnectionListResp struct {
	CommonResp  *pb.CommonResp                 `json:"commonResp"`
	Connections []getAllConnectionListRespItem `json:"connections"`
	Total       int64                          `json:"total"`
}

type getAllConnectionListRespItem struct {
	UserId         string `json:"userId"`
	DeviceId       string `json:"deviceId"`
	Platform       string `json:"platform"`
	ConnectTime    int64  `json:"connectTime"`
	ConnectTimeStr string `json:"connectTimeStr"`
	Ip             string `json:"ip"`
	IpRegion       string `json:"ipRegion"`
	PodIp          string `json:"podIp"`
}

// getAllConnectionList 获取所有连接列表
// @Summary 获取所有连接列表
// @Description 使用此接口获取所有连接列表
// @Tags app连接管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body getAllConnectionListReq true "请求参数"
// @Success 200 {object} getAllConnectionListResp "响应数据"
// @Router /appmgmt/get/connection/list/all [post]
func (r *AppMgrHandler) getAllConnectionList(ctx *gin.Context) {
	in := &getAllConnectionListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	var (
		userIds   []string
		platforms []string
		devices   []string
	)
	for k, v := range in.Filter {
		if v == "" {
			continue
		}
		switch k {
		case "userId":
			userIds = append(userIds, v)
		case "platform":
			platforms = append(platforms, v)
		case "deviceId":
			devices = append(devices, v)
		}
	}
	connResp, err := r.svcCtx.ImService().GetUserConn(ctx, &pb.GetUserConnReq{
		UserIds:   userIds,
		Platforms: platforms,
		Devices:   devices,
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	connParams := connResp.ConnParams
	// 使用userId正序
	sort.Slice(connParams, func(i, j int) bool {
		return connParams[i].UserId < connParams[j].UserId
	})
	var (
		connections []getAllConnectionListRespItem
		total       = int64(len(connParams))
	)
	// 分页
	if in.Page == nil {
		in.Page = &pb.Page{
			Size: 15,
			Page: 1,
		}
	}
	var (
		offset = (in.Page.Page - 1) * in.Page.Size
		pass   = 0
	)
	for _, conn := range connParams {
		if ip, ok := in.Filter["ip"]; ok && ip != "" {
			if conn.Ips != ip {
				continue
			}
		}
		if timeGte, ok := in.Filter["connectTimeGte"]; ok {
			timeGteI, err := strconv.ParseInt(timeGte, 10, 64)
			if err == nil {
				if conn.Timestamp < timeGteI {
					continue
				}
			}
		}
		if timeLte, ok := in.Filter["connectTimeLte"]; ok {
			timeLteI, err := strconv.ParseInt(timeLte, 10, 64)
			if err == nil {
				if conn.Timestamp > timeLteI {
					continue
				}
			}
		}
		// offset
		if int32(pass) < offset {
			pass++
			continue
		}
		// size
		if int32(len(connections)) >= in.Page.Size {
			break
		}
		connections = append(connections, getAllConnectionListRespItem{
			UserId:         conn.UserId,
			DeviceId:       conn.DeviceId,
			Platform:       conn.Platform,
			ConnectTime:    conn.Timestamp,
			ConnectTimeStr: utils.TimeFormat(conn.Timestamp),
			Ip:             conn.Ips,
			IpRegion:       ip2region.Ip2Region(conn.Ips).String(),
			PodIp:          conn.PodIp,
		})
	}
	handler.ReturnOk(ctx, &getAllConnectionListResp{
		Connections: connections,
		Total:       total,
	})
}

type kickoutConnectionReq struct {
	CommonReq *pb.CommonReq `json:"commonReq"`
	UserId    string        `json:"userId"`
	DeviceId  string        `json:"deviceId"`
	Platform  string        `json:"platform"`
}

type kickoutConnectionResp struct {
	CommonResp *pb.CommonResp `json:"commonResp"`
}

// kickoutConnection 踢出连接
// @Summary 踢出连接
// @Description 使用此接口踢出连接
// @Tags app连接管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body kickoutConnectionReq true "请求参数"
// @Success 200 {object} kickoutConnectionResp "响应数据"
// @Router /appmgmt/kickout/connection [post]
func (r *AppMgrHandler) kickoutConnection(ctx *gin.Context) {
	in := &kickoutConnectionReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	resp, err := r.svcCtx.ImService().KickUserConn(ctx, &pb.KickUserConnReq{GetUserConnReq: &pb.GetUserConnReq{
		UserIds:   []string{in.UserId},
		Platforms: []string{in.Platform},
		Devices:   []string{in.DeviceId},
	}})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, &kickoutConnectionResp{CommonResp: resp.GetCommonResp()})
}
