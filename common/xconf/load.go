package xconf

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"strings"
)

// MustLoad loads config into v from path, exits on error.
func MustLoad(path string, v any, opts ...conf.Option) {
	// path是否以http://或者https://开头
	// 如果是，则调用LoadRemote方法
	// 否则调用Load方法
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		if err := LoadRemote(path, v); err != nil {
			log.Fatalf("error: config file %s, %s", path, err.Error())
		}
		return
	}
	if err := conf.Load(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

// LoadRemote loads config into v from remote path.
func LoadRemote(path string, v any) error {
	bytes, err := utils.File.Download(path)
	if err != nil {
		return err
	}
	filename := utils.File.FilenameFromUrl(path)
	// 存储到本地
	if err := utils.File.Save(filename, bytes, 0644); err != nil {
		return err
	}
	return conf.Load(filename, v)
}
