package xsms

import (
	"encoding/json"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
	"github.com/zeromicro/go-zero/core/logx"
)

type TencentSmsSender struct {
	Config *pb.GetServerAllConfigResp_TencentSmsConfig
	client *sms.Client
}

var singletonTencentSmsSender *TencentSmsSender

func NewTencentSmsSender(config *pb.GetServerAllConfigResp_TencentSmsConfig) (*TencentSmsSender, error) {
	if singletonTencentSmsSender != nil {
		return singletonTencentSmsSender, nil
	}
	t := &TencentSmsSender{Config: config}
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
	 * 这里采用的是从环境变量读取的方式，需要在环境变量中先设置这两个值。
	 * 你也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
	 * 以免泄露密钥对危及你的财产安全。
	 * SecretId、SecretKey 查询: https://console.cloud.tencent.com/cam/capi */
	credential := common.NewCredential(config.SecretId, config.SecretKey)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()
	/* SDK默认使用POST方法。
	 * 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"
	/* SDK有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5
	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"
	if config.Region == "" {
		config.Region = "ap-guangzhou"
	}
	client, err := sms.NewClient(credential, config.Region, cpf)
	if err != nil {
		logx.Errorf("NewClient error, %v", err)
		return nil, err
	}
	t.client = client
	singletonTencentSmsSender = t
	return t, nil
}

func (l *TencentSmsSender) SendMsg(mobiles []string, args ...interface{}) error {
	if len(mobiles) == 0 {
		return nil
	}
	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	 * 你可以直接查询SDK源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用IDE进行开发，可以方便的跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK采用的是指针风格指定参数，即使对于基本类型你也需要用指针来对参数赋值。
	 * SDK提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台: https://console.cloud.tencent.com/smsv2
	 * 腾讯云短信小助手: https://cloud.tencent.com/document/product/382/3773#.E6.8A.80.E6.9C.AF.E4.BA.A4.E6.B5.81 */

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr(l.Config.AppId)

	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名 */
	// 签名信息可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-sign) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-sign) 的签名管理查看
	request.SignName = common.StringPtr(l.Config.Sign)

	/* 模板 ID: 必须填写已审核通过的模板 ID */
	// 模板 ID 可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-template) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-template) 的正文模板管理查看
	request.TemplateId = common.StringPtr(l.Config.TemplateId)

	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, arg.(string))
	}
	request.TemplateParamSet = common.StringPtrs(strArgs)

	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(mobiles)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := l.client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	logx.Debugf("response: %s", string(b))

	/* 当出现以下错误码时，快速解决方案参考
	 * [FailedOperation.SignatureIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.signatureincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [FailedOperation.TemplateIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.templateincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [UnauthorizedOperation.SmsSdkAppIdVerifyFail](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunauthorizedoperation.smssdkappidverifyfail-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [UnsupportedOperation.ContainDomesticAndInternationalPhoneNumber](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunsupportedoperation.containdomesticandinternationalphonenumber-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * 更多错误，可咨询[腾讯云助手](https://tccc.qcloud.com/web/im/index.html#/chat?webAppId=8fa15978f85cb41f7e2ea36920cb3ae1&title=Sms)
	 */
	return nil
}
