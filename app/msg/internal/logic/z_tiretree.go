package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/stringx"
	"sync"
)

type ShieldWordTrieTree struct {
	svcCtx   *svc.ServiceContext
	tire     stringx.Trie
	tireLock sync.RWMutex
}

var ShieldWordTrieTreeInstance *ShieldWordTrieTree

func InitShieldWordTrieTree(svcCtx *svc.ServiceContext) {
	ShieldWordTrieTreeInstance = &ShieldWordTrieTree{svcCtx: svcCtx}
	ShieldWordTrieTreeInstance.Flush()
}

func (l *ShieldWordTrieTree) Flush() {
	var words []string
	var maxCreateAt int64
	for {
		// 1000条一次
		var shieldWords []*appmgmtmodel.ShieldWord
		err := l.svcCtx.Mysql().Model(&appmgmtmodel.ShieldWord{}).
			Where("createTime > ?", maxCreateAt).
			Order("createTime asc").
			Limit(1000).
			Find(&shieldWords).Error
		if err != nil {
			panic(err)
		}
		if len(shieldWords) == 0 {
			break
		}
		for _, shieldWord := range shieldWords {
			words = append(words, shieldWord.Word)
			maxCreateAt = shieldWord.CreateTime
		}
	}
	trie := stringx.NewTrie(words, stringx.WithMask(l.svcCtx.ConfigMgr.MessageShieldWordReplace(context.Background())))
	l.tireLock.Lock()
	l.tire = trie
	l.tireLock.Unlock()
}

func (l *ShieldWordTrieTree) Check(content string) (sentence string, found bool) {
	l.tireLock.RLock()
	defer l.tireLock.RUnlock()
	if l.tire == nil {
		return "", false
	}
	sentence, _, found = l.tire.Filter(content)
	return sentence, found
}
