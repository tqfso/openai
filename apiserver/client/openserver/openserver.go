package openserver

import (
	"apiserver/config"
	"bytes"
	"common"
	"common/logger"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/form"
)

var (
	transport *http.Transport
)

func init() {
	transport = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

}

func Get(ctx context.Context, endpoint string, param, resp any) error {
	return do(ctx, "GET", endpoint, param, nil, resp)
}

func Post(ctx context.Context, endpoint string, data, resp any) error {
	return do(ctx, "POST", endpoint, nil, data, resp)
}

func do(ctx context.Context, method, endpoint string, param, data, resp any) error {

	zdan := config.GetZdan()
	path, err := url.JoinPath(zdan.OpenBaseURL, endpoint)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	// 构建请求参数

	if param != nil {
		encoder := form.NewEncoder()
		v, err := encoder.Encode(param)
		if err != nil {
			return err
		}
		path = fmt.Sprintf("%s?%s", path, v.Encode())
	}

	logger.Debug("OpenServer Access", logger.String("method", method), logger.String("path", path))

	// 构建请求体

	var reader io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, path, reader)
	if err != nil {
		return err
	}

	// 构建请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetZdan().ApiServerKey))
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
		return &common.Error{Code: common.InnerAccessError, Msg: standResp.Msg}
	}

	if standResp.Data != nil && resp != nil {
		if err := json.Unmarshal(standResp.Data, resp); err != nil {
			return err
		}
	}

	return nil
}
