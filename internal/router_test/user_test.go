package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	token string
}

func (s *UserTestSuite) SetupSuite() {
	s.token = loginDefault(s.T())
}

// ========== 用户列表 ==========

func (s *UserTestSuite) TestUserList() {
	resp := httpGet(s.T(), "/users"+buildQuery(map[string]string{"page": "1", "size": "10"}), s.token)
	assertSuccess(s.T(), resp)

	var pr PageResult
	parseJSON(s.T(), resp.Data, &pr)
	assert.GreaterOrEqual(s.T(), pr.Total, int64(1))
	assert.Len(s.T(), pr.Items, pr.Size)
}

func (s *UserTestSuite) TestUserList_WithKeyword() {
	resp := httpGet(s.T(), "/users"+buildQuery(map[string]string{"keyword": "admin", "page": "1", "size": "10"}), s.token)
	assertSuccess(s.T(), resp)

	var pr PageResult
	parseJSON(s.T(), resp.Data, &pr)
	assert.GreaterOrEqual(s.T(), pr.Total, int64(1))
}

// ========== 用户详情 ==========

func (s *UserTestSuite) TestUserGet_Success() {
	resp := httpGet(s.T(), "/users/"+AdminID, s.token)
	assertSuccess(s.T(), resp)

	var data map[string]any
	parseJSON(s.T(), resp.Data, &data)
	assert.Equal(s.T(), "admin", data["username"])
	assert.NotEmpty(s.T(), data["id"])
}

func (s *UserTestSuite) TestUserGet_NotFound() {
	resp := httpGet(s.T(), "/users/not_exist_id_12345", s.token)
	assertCode(s.T(), resp, 404)
}

// ========== 创建用户 ==========

func (s *UserTestSuite) TestUserCreate_Success() {
	resp := httpPost(s.T(), "/users", s.token, map[string]any{
		"username": "testuser_" + uniqueSuffix(),
		"password": "Test@123456",
		"nickname": "测试用户",
		"email":    "test@example.com",
		"phone":    "13800138000",
	})
	assertSuccess(s.T(), resp)
}

func (s *UserTestSuite) TestUserCreate_DuplicateUsername() {
	resp := httpPost(s.T(), "/users", s.token, map[string]string{
		"username": "admin",
		"password": "Test@123456",
	})
	assertCode(s.T(), resp, 500)
}

func (s *UserTestSuite) TestUserCreate_MissingFields() {
	// 缺少 username
	resp := httpPost(s.T(), "/users", s.token, map[string]string{
		"password": "Test@123456",
	})
	assertCode(s.T(), resp, 400)

	// 缺少 password
	resp = httpPost(s.T(), "/users", s.token, map[string]string{
		"username": "testuser_missing_pwd",
	})
	assertCode(s.T(), resp, 400)
}

// ========== 更新用户 ==========

func (s *UserTestSuite) TestUserUpdate_Success() {
	// 先创建一个用户
	createResp := httpPost(s.T(), "/users", s.token, map[string]any{
		"username": "updateuser_" + uniqueSuffix(),
		"password": "Test@123456",
		"nickname": "原始昵称",
	})
	assertSuccess(s.T(), createResp)
	userID := extractID(s.T(), createResp.Data)

	// 更新用户
	resp := httpPut(s.T(), "/users/"+userID, s.token, map[string]any{
		"nickname": "更新后昵称",
		"email":    "updated@example.com",
	})
	assertSuccess(s.T(), resp)

	// 验证更新结果
	getResp := httpGet(s.T(), "/users/"+userID, s.token)
	assertSuccess(s.T(), getResp)
	var data map[string]any
	parseJSON(s.T(), getResp.Data, &data)
	assert.Equal(s.T(), "更新后昵称", data["nickname"])
	assert.Equal(s.T(), "updated@example.com", data["email"])
}

func (s *UserTestSuite) TestUserUpdate_NotFound() {
	resp := httpPut(s.T(), "/users/not_exist_id", s.token, map[string]string{
		"nickname": "测试",
	})
	assertCode(s.T(), resp, 500) // GetByID 失败
}

// ========== 删除用户 ==========

func (s *UserTestSuite) TestUserDelete_Success() {
	// 先创建
	createResp := httpPost(s.T(), "/users", s.token, map[string]any{
		"username": "deleteuser_" + uniqueSuffix(),
		"password": "Test@123456",
	})
	assertSuccess(s.T(), createResp)
	userID := extractID(s.T(), createResp.Data)

	// 删除
	resp := httpDelete(s.T(), "/users/"+userID, s.token)
	assertSuccess(s.T(), resp)

	// 验证已被删除
	getResp := httpGet(s.T(), "/users/"+userID, s.token)
	assertCode(s.T(), getResp, 404)
}

func (s *UserTestSuite) TestUserDelete_NotFound() {
	resp := httpDelete(s.T(), "/users/not_exist_id", s.token)
	assertCode(s.T(), resp, 500)
}

// ========== 分配角色 ==========

func (s *UserTestSuite) TestUserAssignRoles_Success() {
	// 先创建用户
	createResp := httpPost(s.T(), "/users", s.token, map[string]any{
		"username": "roleuser_" + uniqueSuffix(),
		"password": "Test@123456",
	})
	assertSuccess(s.T(), createResp)
	userID := extractID(s.T(), createResp.Data)

	// 创建角色
	roleResp := httpPost(s.T(), "/roles", s.token, map[string]any{
		"name": "测试角色_" + uniqueSuffix(),
		"code": "test_role_" + uniqueSuffix(),
	})
	assertSuccess(s.T(), roleResp)
	roleID := extractID(s.T(), roleResp.Data)

	// 分配角色
	resp := httpPut(s.T(), "/users/"+userID+"/roles", s.token, map[string]any{
		"role_ids": []string{roleID},
	})
	assertSuccess(s.T(), resp)
}

// ========== 未认证访问 ==========

func (s *UserTestSuite) TestUserList_NoToken() {
	resp := httpGet(s.T(), "/users", "")
	assertCode(s.T(), resp, 401)
}

func (s *UserTestSuite) TestUserCreate_NoToken() {
	resp := httpPost(s.T(), "/users", "", map[string]string{
		"username": "noauth",
		"password": "123456",
	})
	assertCode(s.T(), resp, 401)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
