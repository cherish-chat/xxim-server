package logic

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/mr"
	"strconv"
	"sync"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchTranslateTextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchTranslateTextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchTranslateTextLogic {
	return &BatchTranslateTextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchTranslateTextLogic) BatchTranslateText(in *pb.BatchTranslateTextReq) (*pb.BatchTranslateTextResp, error) {
	if in.Q == "" {
		return &pb.BatchTranslateTextResp{}, nil
	}
	if in.From == "" {
		// 判断语言
		if utils.IsChinese(in.From) {
			in.From = "zh"
		} else if utils.IsEnglish(in.From) {
			in.From = "en"
		} else {
			return &pb.BatchTranslateTextResp{CommonResp: pb.NewToastErrorResp("不支持的语言")}, nil
		}
	}
	if len(in.ToList) == 0 {
		return &pb.BatchTranslateTextResp{}, nil
	}
	option := l.svcCtx.ConfigMgr.TranslateOption(l.ctx)
	switch option {
	case "", "0":
		//禁用
		return &pb.BatchTranslateTextResp{Results: map[string]string{}}, nil
	case "1":
		//万维易源
		appId := l.svcCtx.ConfigMgr.TranslateShowApiAppId(l.ctx)
		appIdI, err := strconv.Atoi(appId)
		if err != nil {
			return &pb.BatchTranslateTextResp{}, err
		}
		sign := l.svcCtx.ConfigMgr.TranslateShowApiSign(l.ctx)
		var fs []func() error
		resultMap := sync.Map{}
		for _, to := range in.ToList {
			to := to
			fs = append(fs, func() error {
				res := ShowApiRequest("https://route.showapi.com/32-10", appIdI, sign)
				res.AddTextPara("q", in.Q)
				res.AddTextPara("from", in.From)
				res.AddTextPara("to", to)
				resp, err := res.Post()
				if err != nil {
					l.Errorf("post showapi err: %s", err.Error())
					return err
				}
				model := &RespModel{}
				err = json.Unmarshal([]byte(resp), model)
				if err != nil {
					l.Errorf("unmarshal showapi err: %s, resp: %s", err.Error(), resp)
					return err
				}
				if model.ShowapiResBody == nil {
					l.Errorf("unmarshal showapi body=nil, resp: %s", resp)
					return err
				}
				if len(model.ShowapiResBody.Translation) == 0 {
					l.Errorf("unmarshal showapi translation=nil, resp: %s", resp)
					return err
				}
				resultMap.Store(to, model.ShowapiResBody.Translation[0])
				return nil
			})
		}
		err = mr.Finish(fs...)
		if err != nil {
			return &pb.BatchTranslateTextResp{}, err
		}
		results := make(map[string]string)
		resultMap.Range(func(key, value interface{}) bool {
			results[key.(string)] = value.(string)
			return true
		})
		return &pb.BatchTranslateTextResp{Results: results}, nil
	case "2":
		keyId := l.svcCtx.ConfigMgr.TranslateAmazonAccessKeyId(l.ctx)
		secretKey := l.svcCtx.ConfigMgr.TranslateAmazonSecretAccessKey(l.ctx)
		region := l.svcCtx.ConfigMgr.TranslateAmazonRegion(l.ctx)
		cfg := aws.Config{
			Region:      region,
			Credentials: credentials.NewStaticCredentialsProvider(keyId, secretKey, ""),
		}
		client := translate.NewFromConfig(cfg)
		var fs []func() error
		resultMap := sync.Map{}
		for _, to := range in.ToList {
			to := to
			fs = append(fs, func() error {
				output, err := client.TranslateText(l.ctx, &translate.TranslateTextInput{
					SourceLanguageCode: aws.String(in.From),
					TargetLanguageCode: aws.String(to),
					Text:               aws.String(in.Q),
					Settings:           nil,
					TerminologyNames:   nil,
				})
				if err != nil {
					l.Errorf("post showapi err: %s", err.Error())
					return err
				}
				if output == nil {
					l.Errorf("unmarshal showapi body=nil, resp: %+v", output)
					return err
				}
				if output.TranslatedText == nil {
					l.Errorf("unmarshal showapi translation=nil, resp: %+v", output)
					return err
				}
				resultMap.Store(to, *output.TranslatedText)
				return nil
			})
		}
		err := mr.Finish(fs...)
		results := make(map[string]string)
		resultMap.Range(func(key, value interface{}) bool {
			results[key.(string)] = value.(string)
			return true
		})
		return &pb.BatchTranslateTextResp{Results: results}, err
	}
	return &pb.BatchTranslateTextResp{}, nil
}
