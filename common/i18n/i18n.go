package i18n

import (
	_ "embed"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

func NewI18N(mysql *gorm.DB) *I18N {
	m := &I18N{
		mysql:       mysql,
		LanguageMap: map[string]map[string]string{},
	}
	m.init()
	return m
}

type I18N struct {
	mysql       *gorm.DB
	LanguageMap map[string]map[string]string
}

type Language struct {
	Language string `bson:"language" json:"language"` // 语言类型 zh_CN en_US
	Key      string `bson:"key" json:"key"`           // 语言key
	Value    string `bson:"value" json:"value"`       // 语言值
}

func (m *Language) TableName() string {
	return "language"
}

//go:embed i18n.json
var i18nJson []byte

func (l *I18N) init() {
	languageMap := make(map[string]map[string]string)
	err := json.Unmarshal(i18nJson, &languageMap)
	if err != nil {
		logx.Errorf("init language error: %v", err)
		panic(err)
	}
	l.LanguageMap = languageMap
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
