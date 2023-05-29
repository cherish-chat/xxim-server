package ip2region

import "testing"

func TestIp2Region(t *testing.T) {
	Init("https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb")
	t.Logf("%s", Ip2Region("114.114.114.114"))
}
