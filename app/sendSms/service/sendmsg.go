package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/rand/v2"
	logger "star/app/sendSms/logger"
	"strconv"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"star/app/sendSms/dao"
	"star/proto/sendSms/sendSmsPb"
	"star/settings"
	"star/utils"
)

type SendSmsSrv struct {
}

var (
	sendSmsSrvIns *SendSmsSrv
	once          sync.Once
)

func GetSendSmsSrv() *SendSmsSrv {
	once.Do(func() {
		sendSmsSrvIns = &SendSmsSrv{}
	})
	return sendSmsSrvIns
}

func (s *SendSmsSrv) HandleSendSms(ctx context.Context, req *sendSmsPb.SendRequest, resp *sendSmsPb.EmptySendResponse) error {
	if err := sendMsg(req.Phone, req.TemplateCode); err != nil {
		logger.SendSmsLogger.Error("发送短信失败", zap.Error(err))
		return utils.ErrSendSmsFailed
	}
	return nil
}

// 使用AK和SK初始化账号Client
func createClient() (client *dysmsapi.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(settings.Conf.AliyunConfig.AccessKeyId),
		AccessKeySecret: tea.String(settings.Conf.AliyunConfig.AccessKeySecret),
	}
	client, err = dysmsapi.NewClient(config)
	return
}

// 发送短信
func sendMsg(phone, templateCode string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	//生成验证码
	code := generateCode()
	templateParam := fmt.Sprintf(`{"code":"%s"}`, code)

	//tea.string()取地址
	sendMsg := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),                  //手机号
		SignName:      tea.String(settings.Conf.SignName), //签名
		TemplateCode:  tea.String(templateCode),           //模版code
		TemplateParam: tea.String(templateParam),          //短信模板变量对应的实际值
	}
	resp, err := client.SendSms(sendMsg)
	if err != nil {
		return err
	}
	if *(resp.Body.Code) != "OK" {
		err = errors.New(*(resp.Body.Message))
		return err
	}
	//将验证码储存在redis中
	redis.SaveCaptcha(code, phone)
	return nil
}

// 生成六位数验证码
func generateCode() string {
	return strconv.Itoa(rand.IntN(899999) + 100000)
}
