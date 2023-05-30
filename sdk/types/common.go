package types

import "github.com/cherish-chat/xxim-server/common/pb"

type RequestHeader = pb.RequestHeader
type ResponseHeader = pb.ResponseHeader
type Platform = pb.Platform
type I18NLanguage = pb.I18NLanguage
type EncodingProto = pb.EncodingProto

const (
	I18NLanguage_UNSET_LANGUAGE      I18NLanguage = 0  // 未设置
	I18NLanguage_Afrikaans           I18NLanguage = 1  // 南非荷兰语	af
	I18NLanguage_Albanian            I18NLanguage = 2  // 阿尔巴尼亚语	sq
	I18NLanguage_Amharic             I18NLanguage = 3  // 阿姆哈拉语	am
	I18NLanguage_Arabic              I18NLanguage = 4  // 阿拉伯语	ar
	I18NLanguage_Armenian            I18NLanguage = 5  // 亚美尼亚语	hy
	I18NLanguage_Azerbaijani         I18NLanguage = 6  // 阿塞拜疆语	az
	I18NLanguage_Bengali             I18NLanguage = 7  // 孟加拉语	bn
	I18NLanguage_Bosnian             I18NLanguage = 8  // 波斯尼亚语	bs
	I18NLanguage_Bulgarian           I18NLanguage = 9  // 保加利亚语	bg
	I18NLanguage_Catalan             I18NLanguage = 10 // 加泰罗尼亚语	ca
	I18NLanguage_Chinese_Simplified  I18NLanguage = 11 // 简体中文 	zh
	I18NLanguage_Chinese_Traditional I18NLanguage = 12 // 繁体中文 	zh-TW
	I18NLanguage_Croatian            I18NLanguage = 13 // 克罗地亚语	hr
	I18NLanguage_Czech               I18NLanguage = 14 // 捷克语	cs
	I18NLanguage_Danish              I18NLanguage = 15 // 丹麦语	da
	I18NLanguage_Dari                I18NLanguage = 16 // 波斯语	fa-AF
	I18NLanguage_Dutch               I18NLanguage = 17 // 荷兰语	nl
	I18NLanguage_English             I18NLanguage = 18 // 英语	en
	I18NLanguage_Estonian            I18NLanguage = 19 // 爱沙尼亚语	et
	I18NLanguage_Farsi_Persian       I18NLanguage = 20 // 波斯语 	fa
	I18NLanguage_Filipino_Tagalog    I18NLanguage = 21 // 菲律宾语	tl
	I18NLanguage_Finnish             I18NLanguage = 22 // 芬兰语	fi
	I18NLanguage_French              I18NLanguage = 23 // 法语	fr
	I18NLanguage_French_Canada       I18NLanguage = 24 // 法语（加拿大）	fr-CA
	I18NLanguage_Georgian            I18NLanguage = 25 // 格鲁吉亚语	ka
	I18NLanguage_German              I18NLanguage = 26 // 德语	de
	I18NLanguage_Greek               I18NLanguage = 27 // 希腊语	el
	I18NLanguage_Gujarati            I18NLanguage = 28 // 古吉拉特语	gu
	I18NLanguage_Haitian_Creole      I18NLanguage = 29 // 海地克里奥尔语 	ht
	I18NLanguage_Hausa               I18NLanguage = 30 // 豪萨语	ha
	I18NLanguage_Hebrew              I18NLanguage = 31 // 希伯来语	he
	I18NLanguage_Hindi               I18NLanguage = 32 // 印地语	hi
	I18NLanguage_Hungarian           I18NLanguage = 33 // 匈牙利语	hu
	I18NLanguage_Icelandic           I18NLanguage = 34 // 冰岛语	is
	I18NLanguage_Indonesian          I18NLanguage = 35 // 印度尼西亚语	id
	I18NLanguage_Irish               I18NLanguage = 36 // 爱尔兰语	ga
	I18NLanguage_Italian             I18NLanguage = 37 // 意大利语	it
	I18NLanguage_Japanese            I18NLanguage = 38 // 日语	ja
	I18NLanguage_Kannada             I18NLanguage = 39 // 卡纳达语	kn
	I18NLanguage_Kazakh              I18NLanguage = 40 // 哈萨克语	kk
	I18NLanguage_Korean              I18NLanguage = 41 // 韩语	ko
	I18NLanguage_Latvian             I18NLanguage = 42 // 拉脱维亚语	lv
	I18NLanguage_Lithuanian          I18NLanguage = 43 // 立陶宛语	lt
	I18NLanguage_Macedonian          I18NLanguage = 44 // 马其顿语	mk
	I18NLanguage_Malay               I18NLanguage = 45 // 马来语	ms
	I18NLanguage_Malayalam           I18NLanguage = 46 // 马拉雅拉姆语	ml
	I18NLanguage_Maltese             I18NLanguage = 47 // 马耳他语	mt
	I18NLanguage_Marathi             I18NLanguage = 48 // 马拉地语	mr
	I18NLanguage_Mongolian           I18NLanguage = 49 // 蒙古语	mn
	I18NLanguage_Norwegian_Bokmal    I18NLanguage = 50 // 挪威语 	no
	I18NLanguage_Pashto              I18NLanguage = 51 // 普什图语	ps
	I18NLanguage_Polish              I18NLanguage = 52 // 波兰语	pl
	I18NLanguage_Portuguese_Brazil   I18NLanguage = 53 // 葡萄牙语（巴西）	pt
	I18NLanguage_Portuguese_Portugal I18NLanguage = 54 // 葡萄牙语（葡萄牙）	pt-PT
	I18NLanguage_Punjabi             I18NLanguage = 55 // 旁遮普语	pa
	I18NLanguage_Romanian            I18NLanguage = 56 // 罗马尼亚语	ro
	I18NLanguage_Russian             I18NLanguage = 57 // 俄语	ru
	I18NLanguage_Serbian             I18NLanguage = 58 // 塞尔维亚语	sr
	I18NLanguage_Sinhala             I18NLanguage = 59 // 僧伽罗语	si
	I18NLanguage_Slovak              I18NLanguage = 60 // 斯洛伐克语	sk
	I18NLanguage_Slovenian           I18NLanguage = 61 // 斯洛文尼亚语	sl
	I18NLanguage_Somali              I18NLanguage = 62 // 索马里语	so
	I18NLanguage_Spanish             I18NLanguage = 63 // 西班牙语	es
	I18NLanguage_Spanish_Mexico      I18NLanguage = 64 // 西班牙语（墨西哥）	es-MX
	I18NLanguage_Swahili             I18NLanguage = 65 // 斯瓦希里语	sw
	I18NLanguage_Swedish             I18NLanguage = 66 // 瑞典语	sv
	I18NLanguage_Tamil               I18NLanguage = 67 // 泰米尔语	ta
	I18NLanguage_Telugu              I18NLanguage = 68 // 泰卢固语	te
	I18NLanguage_Thai                I18NLanguage = 69 // 泰语	th
	I18NLanguage_Turkish             I18NLanguage = 70 // 土耳其语	tr
	I18NLanguage_Ukrainian           I18NLanguage = 71 // 乌克兰语	uk
	I18NLanguage_Urdu                I18NLanguage = 72 // 乌尔都语	ur
	I18NLanguage_Uzbek               I18NLanguage = 73 // 乌兹别克语	uz
	I18NLanguage_Vietnamese          I18NLanguage = 74 // 越南语	vi
	I18NLanguage_Welsh               I18NLanguage = 75 // 威尔士语	cy
)
const (
	EncodingProto_PROTOBUF EncodingProto = 0 // protobuf
	EncodingProto_JSON     EncodingProto = 1 // json
)
