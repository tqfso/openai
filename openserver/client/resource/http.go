package resource

import (
	"bytes"
	"common"
	"common/logger"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"openserver/config"
	"openserver/middleware/auth"
	"time"

	"github.com/go-playground/form"
)

type Response struct {
	Data json.RawMessage `json:"data"` // 延迟解析
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
}

func (r *Response) IsSuccess() bool {
	return r.Code == 0
}

func Get(endpoint string, param, resp any) error {
	return do("GET", endpoint, param, nil, resp)
}

func Post(endpoint string, data, resp any) error {
	return do("POST", endpoint, nil, data, resp)
}

func do(method, endpoint string, param, data, resp any) error {

	zdan := config.GetZdan()

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 构建请求参数

	var url string
	if param != nil {
		encoder := form.NewEncoder()
		v, err := encoder.Encode(param)
		if err != nil {
			return err
		}
		url = fmt.Sprintf("https://%s/zresource/v1%s?%s", zdan.Address(), endpoint, v.Encode())
	} else {
		url = fmt.Sprintf("https://%s/zresource/v1%s", zdan.Address(), endpoint)
	}

	logger.Debug("Resource Access", logger.String("method", method), logger.String("url", url))

	// 构建请求体

	var reader io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}

	// 构建请求头

	token, err := auth.ZCloudMakeToken(zdan.CloudDmappId, zdan.CloudDmappKey)
	if err != nil {
		return err
	}

	req.Header.Set("ZCookie", token)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 获取返回状态字

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("resource status code: %d", response.StatusCode)
	}

	// 解析响应体

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return fmt.Errorf("response body empty")
	}

	var standResp Response
	if err := json.Unmarshal(body, &standResp); err != nil {
		return err
	}

	if !standResp.IsSuccess() {
		return &common.Error{Code: standResp.Code, Msg: standResp.Msg}
	}

	if standResp.Data != nil && resp != nil {
		if err := json.Unmarshal(standResp.Data, resp); err != nil {
			return err
		}
	}

	return nil
}
