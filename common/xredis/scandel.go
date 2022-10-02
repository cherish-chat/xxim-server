package xredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	scriptScanDel = `
local function scan(key)
    local cursor = 0
    local keynum = 0

    repeat
        local res = redis.call("scan", cursor, "match", key,'COUNT',ARGV[1])

        if (res ~= nil and #res >= 0) then
            redis.replicate_commands()
            cursor = tonumber(res[1])
            local ks = res[2]
            keynum = #ks
            for i=1,keynum,1 do
                local k = tostring(ks[i])
                redis.call("del", k)
            end
        end
    until (cursor <= 0)

    return keynum
end

local a = #KEYS
local b = 1
local total = 0
while (b <= a)
do
    total = total + scan(KEYS[b])
    b = b + 1
end

return total
`
)

var (
	scriptScanDelSha = ``
)

func ScanDel(rc redis.UniversalClient, ctx context.Context, patternKeys ...string) error {
	if scriptScanDelSha == "" {
		var err error
		scriptScanDelSha, err = redis.NewScript(scriptScanDel).Load(ctx, rc).Result()
		if err != nil {
			panic(err)
		}
	}
	result, err := rc.EvalSha(ctx, scriptScanDelSha, patternKeys, 1000).Result()
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("ScanDel result: %v", result)
	return nil
}
