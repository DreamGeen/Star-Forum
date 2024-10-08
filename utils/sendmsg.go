package utils

import (
	"context"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
	"math/rand/v2"
	"star/app/storage/cached"
	"star/constant/settings"
	"star/constant/str"
	"strconv"
	"time"
)

func HandleSendSms(ctx context.Context, phone string, templateCode string) error {
	if err := sendMsg(ctx, phone, templateCode); err != nil {
		Logger.Error("send message error",
			zap.Error(err))
		return str.ErrSendSmsError
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
func sendMsg(ctx context.Context, phone, templateCode string) error {
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
	key := "captcha:" + phone
	cached.WriteWithOvertime(ctx, key, code, 5*time.Minute)
	return nil
}

// 生成六位数验证码
func generateCode() string {
	return strconv.Itoa(rand.IntN(899999) + 100000)
}
