package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	resp, err := httpGetRaw(BaseURL + "/health")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthLogin_Success(t *testing.T) {
	resp := httpPost(t, "/auth/login", "", map[string]string{
		"username": "admin",
		"password": "admin123",
	})
	if resp.Code != 0 {
		t.Skipf("跳过: 登录失败 (code=%d msg=%s)，请确保数据库中存在 admin/admin123 账号", resp.Code, resp.Msg)
	}
	assertSuccess(t, resp)
	assert.NotNil(t, resp.Data)

	// 验证返回了 token 和 user_info
	var data map[string]any
	parseJSON(t, resp.Data, &data)
	assert.NotEmpty(t, data["token"])

	userInfo, ok := data["user_info"].(map[string]any)
	if !ok {
		t.Fatal("user_info 字段缺失或类型错误")
	}
	assert.Equal(t, "admin", userInfo["username"])
	assert.NotEmpty(t, userInfo["id"])
	assert.NotEmpty(t, userInfo["roles"])
	assert.NotEmpty(t, userInfo["permissions"])
}

func TestAuthLogin_WrongPassword(t *testing.T) {
	resp := httpPost(t, "/auth/login", "", map[string]string{
		"username": "admin",
		"password": "wrong_password",
	})
	assertCode(t, resp, 400)
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	resp := httpPost(t, "/auth/login", "", map[string]string{
		"username": "notexist_user",
		"password": "123456",
	})
	assertCode(t, resp, 400)
}

func TestAuthLogin_EmptyBody(t *testing.T) {
	resp := httpPost(t, "/auth/login", "", nil)
	assertCode(t, resp, 400)
}

func TestAuthLogin_MissingFields(t *testing.T) {
	// 缺少 password
	resp := httpPost(t, "/auth/login", "", map[string]string{
		"username": "admin",
	})
	assertCode(t, resp, 400)

	// 缺少 username
	resp = httpPost(t, "/auth/login", "", map[string]string{
		"password": "admin123",
	})
	assertCode(t, resp, 400)
}

func TestGetUserInfo_Success(t *testing.T) {
	token := loginDefault(t)

	resp := httpGet(t, "/user/info", token)
	assertSuccess(t, resp)

	var data map[string]any
	parseJSON(t, resp.Data, &data)
	assert.Equal(t, "admin", data["username"])
	assert.NotEmpty(t, data["id"])
	assert.NotEmpty(t, data["roles"])
	assert.NotEmpty(t, data["permissions"])
}

func TestGetUserInfo_NoToken(t *testing.T) {
	resp := httpGet(t, "/user/info", "")
	assertCode(t, resp, 401)
}

func TestGetUserInfo_InvalidToken(t *testing.T) {
	resp := httpGet(t, "/user/info", "invalid_token_value")
	assertCode(t, resp, 401)
}
