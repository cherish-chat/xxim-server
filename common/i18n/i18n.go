package i18n

import (
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

func (l *I18N) init() {
	l.mysql.AutoMigrate(&Language{})
	var languageList []Language
	// 查询所有
	err := l.mysql.Find(&languageList).Error
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
