package memcache

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)
type appToken struct {
	AppAccessToken     string `json:"app_access_token"`
    Code               int    `json:"code"`
    Expire             int    `json:"expire"`
    Message            string `json:"msg"`
    TernantAccessToken string `json:"tenant_access_token"`
}

type AppAccess interface {
	GetAppAccessToken() error
}

type appAccess struct {
	appId     string `json:"app_id"`
	appSecret string `json:"app_secret"`
}

func NewAppToken() *appToken {
	return &appToken{
		AppAccessToken:     "",
		Code: 0,
		Expire: 1800,
		Message: "",
		TernantAccessToken: "",
	}
}

func (a * appToken) GetAppAccessToken()  error {
	appInfo := appAccess{
        appId:     "cli_a4b0a37dd8f8d02f",
        appSecret: "ziCKGTkVuprRLpoV17rrzcaCkjZV5lBq",
    }
    postBody, _ := json.Marshal(appInfo)
    responseBody := bytes.NewBuffer(postBody)
    res, err := http.Post("https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", "application/json", responseBody)

    if err != nil {
        log.Println(err)
        return err
    }

    defer res.Body.Close()

    var responseObject appToken
    if err := json.NewDecoder(res.Body).Decode(&responseObject); err != nil {
        log.Println(err)
        return err
    }

    log.Println("app access: ", responseObject.Expire)
	return nil
}