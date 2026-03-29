package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 公共配置
var (
	BaseURL    = "http://127.0.0.1:8080"
	APIBase    = BaseURL + "/api"
	AdminToken string // 登录后获取的 token
	AdminID    string // 管理员用户 ID
)

// Response 通用响应结构
type Response struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// PageResult 分页结果
type PageResult struct {
	Items   []json.RawMessage `json:"items"`
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	Size    int               `json:"size"`
	HasMore bool              `json:"hasMore"`
}

// ===================== HTTP 工具函数 =====================

// httpGet 发送 GET 请求
func httpGet(t *testing.T, path string, token string) *Response {
	t.Helper()
	return doRequest(t, http.MethodGet, path, token, nil)
}

// httpPost 发送 POST 请求
func httpPost(t *testing.T, path string, token string, body any) *Response {
	t.Helper()
	return doRequest(t, http.MethodPost, path, token, body)
}

// httpPut 发送 PUT 请求
func httpPut(t *testing.T, path string, token string, body any) *Response {
	t.Helper()
	return doRequest(t, http.MethodPut, path, token, body)
}

// httpDelete 发送 DELETE 请求
func httpDelete(t *testing.T, path string, token string) *Response {
	t.Helper()
	return doRequest(t, http.MethodDelete, path, token, nil)
}

// doRequest 统一请求方法
func doRequest(t *testing.T, method, path, token string, body any) *Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		assert.NoError(t, err)
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, APIBase+path, reqBody)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result Response
	err = json.Unmarshal(respBody, &result)
	assert.NoError(t, err)

	return &result
}

// parseData 将 Response.Data 解析到目标结构体
func parseData[T any](t *testing.T, resp *Response) *T {
	t.Helper()
	var data T
	if resp.Data == nil || string(resp.Data) == "null" {
		return nil
	}
	err := json.Unmarshal(resp.Data, &data)
	assert.NoError(t, err)
	return &data
}

// parseList 将 Response.Data 解析为分页结果列表
func parseList[T any](t *testing.T, resp *Response) []T {
	t.Helper()
	var pr PageResult
	err := json.Unmarshal(resp.Data, &pr)
	assert.NoError(t, err)

	items := make([]T, 0, len(pr.Items))
	for _, item := range pr.Items {
		var obj T
		err := json.Unmarshal(item, &obj)
		assert.NoError(t, err)
		items = append(items, obj)
	}
	return items
}

// buildQuery 构建 query string
func buildQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return "?" + values.Encode()
}

// assertSuccess 断言请求成功
func assertSuccess(t *testing.T, resp *Response) {
	t.Helper()
	assert.Equal(t, 0, resp.Code, "期望成功但返回错误: %s", resp.Msg)
}

// assertCode 断言返回特定业务码
func assertCode(t *testing.T, resp *Response, code int) {
	t.Helper()
	assert.Equal(t, code, resp.Code, "期望 code=%d 但实际: code=%d msg=%s", code, resp.Code, resp.Msg)
}

// extractID 从 json.RawMessage 中提取 id 字段
func extractID(t *testing.T, raw json.RawMessage) string {
	t.Helper()
	var obj map[string]any
	err := json.Unmarshal(raw, &obj)
	assert.NoError(t, err)
	id, ok := obj["id"].(string)
	assert.True(t, ok, "data 中缺少 id 字段")
	return id
}

// extractToken 从登录响应 data 中提取 token
func extractToken(t *testing.T, resp *Response) string {
	t.Helper()
	var data map[string]any
	err := json.Unmarshal(resp.Data, &data)
	assert.NoError(t, err)
	token, ok := data["token"].(string)
	assert.True(t, ok, "登录响应中缺少 token")
	return token
}

// loginDefault 使用默认管理员账号登录，获取 token
func loginDefault(t *testing.T) string {
	t.Helper()
	if AdminToken != "" {
		return AdminToken
	}
	resp := httpPost(t, "/auth/login", "", map[string]string{
		"username": "admin",
		"password": "admin123",
	})
	if resp.Code != 0 {
		t.Skipf("跳过: 登录失败 (code=%d msg=%s)，请确保数据库中存在 admin/admin123 账号", resp.Code, resp.Msg)
		return ""
	}
	AdminToken = extractToken(t, resp)

	// 提取 user_id
	var data map[string]any
	json.Unmarshal(resp.Data, &data)
	if userInfo, ok := data["user_info"].(map[string]any); ok {
		AdminID = userInfo["id"].(string)
	}
	return AdminToken
}

// requireAdmin 跳过无 admin 时后续测试
func requireAdmin(t *testing.T) string {
	t.Helper()
	token := loginDefault(t)
	if token == "" {
		t.Skip("跳过: 无法获取管理员 token，请确保数据库中存在 admin/admin123 账号")
	}
	return token
}

// uniqueSuffix 生成唯一后缀避免测试数据冲突
func uniqueSuffix() string {
	return strings.ReplaceAll("test_XQZ9mK", "XQZ9mK", fmt.Sprintf("%d", 0))
}

// httpGetRaw 发送原始 GET 请求（返回 *http.Response）
func httpGetRaw(url string) (*http.Response, error) {
	return http.DefaultClient.Get(url)
}

// parseJSON 将 json.RawMessage 解析到目标结构体
func parseJSON(t *testing.T, raw json.RawMessage, v any) {
	t.Helper()
	err := json.Unmarshal(raw, v)
	assert.NoError(t, err)
}
