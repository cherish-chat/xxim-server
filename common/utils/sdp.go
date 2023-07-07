package utils

import "github.com/pion/sdp/v2"

type xSdp struct {
}

var Sdp = &xSdp{}

func (x *xSdp) GetClientIp(sd *sdp.SessionDescription) string {
	for _, m := range sd.MediaDescriptions {
		for _, a := range m.Attributes {
			if a.Key == "candidate" {
				return a.Value
			}
		}
	}
	return ""
}
