package logic

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/cherish-chat/xxim-server/common/utils"
	"strconv"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"

	_ "github.com/aws/aws-sdk-go-v2"
)

type TranslateTextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTranslateTextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TranslateTextLogic {
	return &TranslateTextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
RespBodyModel

	{
	    "showapi_res_error":"",
	    "showapi_res_id":"6453dc270de376c93a55cdd7",
	    "showapi_res_code":0,
	    "showapi_fee_num":1,
	    "showapi_res_body":{
	        "query":"你是谁啊",
	        "translation":[
	            "Who are you?"
	        ],
	        "web":[

	        ],
	        "basic":[

	        ],
	        "fee_num":1,
	        "ret_code":0
	    }
	}
*/
type RespBodyModel struct {
	Translation []string `json:"translation"`
}

/*
RespModel
语言	代码	语言	代码	语言	代码
阿拉伯语	ara	康瓦尔语	cor	塞尔维亚语	srp
阿尔巴尼亚语	alb	克里克语	cre	世界语	epo
阿拉贡语	arg	克罗地亚语	hrv	斯洛文尼亚语	slo
艾马拉语	aym	孔卡尼语	kok	索马里语	som
奥塞梯语	oss	老挝语	lao	泰语	th
奥里亚语	ori	拉脱维亚语	lav	泰米尔语	tam
波兰语	pl	卢干达语	lug	泰卢固语	tel
巴什基尔语	bak	卢旺达语	kin	塔塔尔语	tt
白俄罗斯语	bel	罗姆语	ro	乌克兰语	ukr
保加利亚语	bul	缅甸语	bur	文达语	ven
本巴语	bem	马拉雅拉姆语	mal	西班牙语	spa
俾路支语	bal	迈蒂利语	mai	匈牙利语	hu
博杰普尔语	bho	毛利语	mao	希利盖农语	hil
楚瓦什语	chv	那不勒斯语	nea	新挪威语	nno
丹麦语	dan	南索托语	sot	修纳语	sna
掸语	sha	旁遮普语	pan	新增粤语	yus
低地德语	log	契维语	twi	英语	en
俄语	ru	瑞典语	swe	意大利语	it
法语	fra	萨摩亚语	sm	因特语	ina
梵语	san	桑海语	sol	伊博语	ibo
法罗语	fao	书面挪威语	nob	亚美尼亚语	arm
盖尔语	gla	斯瓦希里语	swa	伊朗语	ir
高棉语	hkm	塞茨瓦纳语	tn	中文(简体)	zh
古吉拉特语	guj	土耳其语	tr	中文(粤语)	yue
瓜拉尼语	grn	他加禄语	tgl	祖鲁语	zul
韩语	kor	突尼斯阿拉伯语	tua	爱尔兰语	gle
哈卡钦语	hak	瓦隆语	wln	阿尔及利亚阿拉伯语	arq
豪萨语	hau	沃洛夫语	wol	阿姆哈拉语	amh
吉尔吉斯语	kir	希伯来语	heb	阿塞拜疆语	aze
加泰罗尼亚语	cat	西弗里斯语	fry	爱沙尼亚语	est
卡拜尔语	kab	下索布语	los	奥罗莫语	orm
卡舒比语	kah	西非书面语	nqo	波斯语	per
科西嘉语	cos	宿务语	ceb	巴斯克语	baq
克林贡语	kli	新增普通话	zhs	柏柏尔语	ber
克什米尔语	kas	印地语	hi	北方萨米语	sme
拉丁语	lat	越南语	vie	比林语	bli
拉特加莱语	lag	亚齐语	ach	冰岛语	ice
林加拉语	lin	伊多语	ido	巴西语	pt_BR
卢森尼亚语	ruy	伊努克提图特语	iku	聪加语	tso
罗曼什语	roh	中文(繁体)	cht	德语	de
马来语	may	扎扎其语	zaz	德顿语	tet
马拉加斯语	mg	爪哇语	jav	菲律宾语	fil
马绍尔语	mah	奥克语	oci	弗留利语	fri
毛里求斯克里奥尔语	mau	阿肯语	aka	刚果语	kon
马耳他语	mlt	阿萨姆语	asm	格陵兰语	kal
挪威语	nor	阿斯图里亚斯语	ast	古希腊语	gra
南非荷兰语	afr	奥杰布瓦语	oji	荷兰语	nl
葡萄牙语	pt	布列塔尼语	bre	黑山语	mot
普什图语	pus	巴西葡萄牙语	pot	加利西亚语	glg
齐切瓦语	nya	邦板牙语	pam	捷克语	cs
日语	jp	北索托语	ped	卡纳达语	kan
萨丁尼亚语	srd	比斯拉马语	bis	中文(文言文)	wyw
波斯尼亚语	bos	南恩德贝莱语	nbl	乌尔都语	urd
鞑靼语	tat	尼泊尔语	nep	希腊语	el
迪维希语	div	帕皮阿门托语	pap	西里西亚语	sil
芬兰语	fin	切罗基语	chr	夏威夷语	haw
富拉尼语	ful	塞尔维亚-克罗地亚语	sec	信德语	snd
高地索布语	ups	僧伽罗语	sin	叙利亚语	syr
格鲁吉亚语	geo	斯洛伐克语	sk	巽他语	sun
古英语	eno	苏格兰语	sco	印尼语	id
胡帕语	hup	塞尔维亚语（西里尔）	src	意第绪语	yid
海地语	ht	塔吉克语	tgk	印古什语	ing
海地克里奥尔语	ht	提格利尼亚语	tir	约鲁巴语	yor
加拿大法语	frn	土库曼语	tuk	因纽特语	iu
卡努里语	kau	威尔士语	wel	中古法语	frm
科萨语	xho	立陶宛语	lit	曼克斯语	glv
克里米亚鞑靼语	cri	逻辑语	loj	孟加拉语	ben
克丘亚语	que	马拉地语	mar	马拉提语	mr
库尔德语	kur	马其顿语	mac	林堡语	lim
罗马尼亚语	rom	卢森堡语	ltz
*/
type RespModel struct {
	ShowapiResError string         `json:"showapi_res_error"`
	ShowapiResId    string         `json:"showapi_res_id"`
	ShowapiResCode  int            `json:"showapi_res_code"`
	ShowapiFeeNum   int            `json:"showapi_fee_num"`
	ShowapiResBody  *RespBodyModel `json:"showapi_res_body"`
}

func (l *TranslateTextLogic) TranslateText(in *pb.TranslateTextReq) (*pb.TranslateTextResp, error) {
	if in.Q == "" {
		return &pb.TranslateTextResp{}, nil
	}
	if in.From == "" {
		// 判断语言
		if utils.IsChinese(in.From) {
			in.From = "zh"
		} else if utils.IsEnglish(in.From) {
			in.From = "en"
		} else {
			return &pb.TranslateTextResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "不支持的语言"))}, nil
		}
	}
	if in.To == "" {
		return &pb.TranslateTextResp{}, nil
	}
	option := l.svcCtx.ConfigMgr.TranslateOption(l.ctx)
	switch option {
	case "", "0":
		//禁用
		return &pb.TranslateTextResp{Result: in.Q}, nil
	case "1":
		//万维易源
		appId := l.svcCtx.ConfigMgr.TranslateShowApiAppId(l.ctx)
		appIdI, err := strconv.Atoi(appId)
		if err != nil {
			return &pb.TranslateTextResp{}, err
		}
		sign := l.svcCtx.ConfigMgr.TranslateShowApiSign(l.ctx)
		res := ShowApiRequest("https://route.showapi.com/32-10", appIdI, sign)
		res.AddTextPara("q", in.Q)
		res.AddTextPara("from", in.From)
		res.AddTextPara("to", in.To)
		resp, err := res.Post()
		if err != nil {
			l.Errorf("post showapi err: %s", err.Error())
			return &pb.TranslateTextResp{}, err
		}
		model := &RespModel{}
		err = json.Unmarshal([]byte(resp), model)
		if err != nil {
			l.Errorf("unmarshal showapi err: %s, resp: %s", err.Error(), resp)
			return &pb.TranslateTextResp{}, err
		}
		if model.ShowapiResBody == nil {
			l.Errorf("unmarshal showapi body=nil, resp: %s", resp)
			return &pb.TranslateTextResp{Result: in.Q}, err
		}
		if len(model.ShowapiResBody.Translation) == 0 {
			l.Errorf("unmarshal showapi translation=nil, resp: %s", resp)
			return &pb.TranslateTextResp{Result: in.Q}, err
		}
		return &pb.TranslateTextResp{Result: model.ShowapiResBody.Translation[0]}, nil
	case "2":
		keyId := l.svcCtx.ConfigMgr.TranslateAmazonAccessKeyId(l.ctx)
		secretKey := l.svcCtx.ConfigMgr.TranslateAmazonSecretAccessKey(l.ctx)
		region := l.svcCtx.ConfigMgr.TranslateAmazonRegion(l.ctx)
		cfg := aws.Config{
			Region:      region,
			Credentials: credentials.NewStaticCredentialsProvider(keyId, secretKey, ""),
		}
		client := translate.NewFromConfig(cfg)
		output, err := client.TranslateText(l.ctx, &translate.TranslateTextInput{
			SourceLanguageCode: aws.String(in.From),
			TargetLanguageCode: aws.String(in.To),
			Text:               aws.String(in.Q),
			Settings:           nil,
			TerminologyNames:   nil,
		})
		if err != nil {
			l.Errorf("post showapi err: %s", err.Error())
			return &pb.TranslateTextResp{}, err
		}
		if output == nil {
			l.Errorf("unmarshal showapi body=nil, resp: %+v", output)
			return &pb.TranslateTextResp{Result: in.Q}, err
		}
		if output.TranslatedText == nil {
			l.Errorf("unmarshal showapi translation=nil, resp: %+v", output)
			return &pb.TranslateTextResp{Result: in.Q}, err
		}
		return &pb.TranslateTextResp{Result: *output.TranslatedText}, nil
	}
	return &pb.TranslateTextResp{}, nil
}
