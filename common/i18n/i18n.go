package i18n

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

func NewI18N(mongodb *xmgo.Client) *I18N {
	m := &I18N{
		c:           mongodb.Collection(&Language{}),
		LanguageMap: map[string]map[string]string{},
	}
	m.init()
	return m
}

type I18N struct {
	c           *qmgo.Collection
	LanguageMap map[string]map[string]string
}

type Language struct {
	Language string `bson:"language" json:"language"` // 语言类型 zh_CN en_US
	Key      string `bson:"key" json:"key"`           // 语言key
	Value    string `bson:"value" json:"value"`       // 语言值
}

func (m *Language) CollectionName() string {
	return "language"
}

func (m *Language) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key:          []string{"language", "key"},
		IndexOptions: opts.Index().SetUnique(true),
	}, {
		Key: []string{"language"},
	}})
	return nil
}

func (l *I18N) init() {
	var languageList []Language
	err := l.c.Find(context.Background(), bson.M{}).All(&languageList)
	if err != nil {
		logx.Errorf("init language error: %v", err)
		panic(err)
	}
	for _, language := range languageList {
		if _, ok := l.LanguageMap[language.Language]; !ok {
			l.LanguageMap[language.Language] = map[string]string{}
		}
		l.LanguageMap[language.Language][language.Key] = language.Value
	}
}

func (l *I18N) T(lang string, key string) (value string) {
	if _, ok := l.LanguageMap[lang]; !ok {
		return key
	}
	if v, ok := l.LanguageMap[lang][key]; ok {
		return v
	}
	return key
}
