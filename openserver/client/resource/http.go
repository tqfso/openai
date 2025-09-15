package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"openserver/config"
	"openserver/middleware/auth"
	"time"

	"github.com/google/go-querystring/query"
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

/*
Example:

	type Params struct {
	    Page int `url:"page"`
	    Size int `url:"size"`
	}

p := Params{Page: 1, Size: 10}
*/

func do(method, endpoint string, param, data, resp any) error {

	zdan := config.GetZdan()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 构建请求参数

	var url string
	if param != nil {
		v, err := query.Values(param)
		if err != nil {
			return err
		}
		url = fmt.Sprintf("https://%s:%s/zresource/v1/%s?%s", zdan.ZdanHost, zdan.ZdanPort, endpoint, v.Encode())
	} else {
		url = fmt.Sprintf("https://%s:%s/zresource/v1/%s", zdan.ZdanHost, zdan.ZdanPort, endpoint)
	}

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
		return fmt.Errorf("response status code: %d", response.StatusCode)
	}

	// 解析响应体

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return fmt.Errorf("response body empty")
	}

	var resourceResp Response
	if err := json.Unmarshal(body, &resourceResp); err != nil {
		return err
	}

	if resourceResp.Data != nil && resp != nil {
		if err := json.Unmarshal(resourceResp.Data, resp); err != nil {
			return err
		}
	}

	return nil
}
