package memcache

import (
	"context"
	"encoding/json"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkauth "github.com/larksuite/oapi-sdk-go/v3/service/auth/v3"
)

var responseObject AppToken

type AppToken struct {
	AppAccessToken     string `json:"app_access_token"`
	Code               int    `json:"code"`
	Expire             int    `json:"expire"`
	Message            string `json:"msg"`
	TernantAccessToken string `json:"tenant_access_token"`
}

type AppAccess interface {
	GetAppAccessToken() (*AppToken, error)
}

type appAccess struct {
	appId     string `json:"app_id"`
	appSecret string `json:"app_secret"`
}

func NewAppToken() *AppToken {
	return &AppToken{
		AppAccessToken:     "",
		Code:               0,
		Expire:             1800,
		Message:            "",
		TernantAccessToken: "",
	}
}

func (a *AppToken) GetAppAccessToken() (*AppToken, error) {
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient("cli_a4b0a37dd8f8d02f", "ziCKGTkVuprRLpoV17rrzcaCkjZV5lBq", lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkauth.NewInternalAppAccessTokenReqBuilder().
		Body(larkauth.NewInternalAppAccessTokenReqBodyBuilder().
			AppId(`cli_a4b0a37dd8f8d02f`).
			AppSecret(`ziCKGTkVuprRLpoV17rrzcaCkjZV5lBq`).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Auth.AppAccessToken.Internal(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, err
	}

	json.Unmarshal(resp.RawBody, &responseObject)

	return &responseObject, nil
}
