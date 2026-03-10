package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type accessTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
}

func httpGet(url string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// 简单 POST JSON 请求
func httpPostJSON(url string, body interface{}) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	return client.Do(req)
}
func GetAccessToken() (string, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", Config.AppKey, Config.AppSecret)
	resp, err := httpGet(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data accessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.ErrCode != 0 {
		return "", fmt.Errorf("gettoken error: %d %s", data.ErrCode, data.ErrMsg)
	}
	return data.AccessToken, nil
}
