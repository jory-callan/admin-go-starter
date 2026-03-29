package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RoleTestSuite struct {
	suite.Suite
	token string
}

func (s *RoleTestSuite) SetupSuite() {
	s.token = loginDefault(s.T())
}

// ========== 角色列表 ==========

func (s *RoleTestSuite) TestRoleList() {
	resp := httpGet(s.T(), "/roles"+buildQuery(map[string]string{"page": "1", "size": "10"}), s.token)
	assertSuccess(s.T(), resp)

	var pr PageResult
	parseJSON(s.T(), resp.Data, &pr)
	assert.GreaterOrEqual(s.T(), pr.Total, int64(1))
}

func (s *RoleTestSuite) TestRoleList_WithKeyword() {
	resp := httpGet(s.T(), "/roles"+buildQuery(map[string]string{"keyword": "admin", "page": "1", "size": "10"}), s.token)
	assertSuccess(s.T(), resp)

	var pr PageResult
	parseJSON(s.T(), resp.Data, &pr)
	assert.GreaterOrEqual(s.T(), pr.Total, int64(1))
}

// ========== 角色详情 ==========

func (s *RoleTestSuite) TestRoleGet_Success() {
	// 先从列表获取一个角色 ID
	listResp := httpGet(s.T(), "/roles"+buildQuery(map[string]string{"page": "1", "size": "1"}), s.token)
	assertSuccess(s.T(), listResp)

	var pr PageResult
	parseJSON(s.T(), listResp.Data, &pr)
	assert.NotEmpty(s.T(), pr.Items)

	roleID := extractID(s.T(), pr.Items[0])

	resp := httpGet(s.T(), "/roles/"+roleID, s.token)
	assertSuccess(s.T(), resp)

	var data map[string]any
	parseJSON(s.T(), resp.Data, &data)
	assert.NotEmpty(s.T(), data["id"])
	assert.NotEmpty(s.T(), data["name"])
	assert.NotEmpty(s.T(), data["code"])
}

func (s *RoleTestSuite) TestRoleGet_NotFound() {
	resp := httpGet(s.T(), "/roles/not_exist_id", s.token)
	assertCode(s.T(), resp, 404)
}

// ========== 创建角色 ==========

func (s *RoleTestSuite) TestRoleCreate_Success() {
	code := "test_create_role_" + uniqueSuffix()
	resp := httpPost(s.T(), "/roles", s.token, map[string]any{
		"name":        "测试创建角色",
		"code":        code,
		"description": "用于测试的角色",
		"sort":        100,
		"status":      1,
	})
	assertSuccess(s.T(), resp)

	// 验证创建结果
	var data map[string]any
	parseJSON(s.T(), resp.Data, &data)
	createdRoleID := data["id"].(string)

	getResp := httpGet(s.T(), "/roles/"+createdRoleID, s.token)
	assertSuccess(s.T(), getResp)
	var role map[string]any
	parseJSON(s.T(), getResp.Data, &role)
	assert.Equal(s.T(), "测试创建角色", role["name"])
	assert.Equal(s.T(), code, role["code"])
}

func (s *RoleTestSuite) TestRoleCreate_DuplicateCode() {
	resp := httpPost(s.T(), "/roles", s.token, map[string]string{
		"name": "重复角色",
		"code": "admin",
	})
	assertCode(s.T(), resp, 500)
}

func (s *RoleTestSuite) TestRoleCreate_MissingFields() {
	resp := httpPost(s.T(), "/roles", s.token, map[string]string{
		"name": "缺少编码",
	})
	assertCode(s.T(), resp, 400)
}

// ========== 更新角色 ==========

func (s *RoleTestSuite) TestRoleUpdate_Success() {
	// 先创建
	createResp := httpPost(s.T(), "/roles", s.token, map[string]any{
		"name": "待更新角色_" + uniqueSuffix(),
		"code": "to_update_" + uniqueSuffix(),
	})
	assertSuccess(s.T(), createResp)
	roleID := extractID(s.T(), createResp.Data)

	// 更新
	resp := httpPut(s.T(), "/roles/"+roleID, s.token, map[string]any{
		"name":        "已更新角色",
		"description": "更新后的描述",
	})
	assertSuccess(s.T(), resp)

	// 验证
	getResp := httpGet(s.T(), "/roles/"+roleID, s.token)
	assertSuccess(s.T(), getResp)
	var data map[string]any
	parseJSON(s.T(), getResp.Data, &data)
	assert.Equal(s.T(), "已更新角色", data["name"])
	assert.Equal(s.T(), "更新后的描述", data["description"])
}

// ========== 删除角色 ==========

func (s *RoleTestSuite) TestRoleDelete_Success() {
	// 先创建
	createResp := httpPost(s.T(), "/roles", s.token, map[string]any{
		"name": "待删除角色_" + uniqueSuffix(),
		"code": "to_delete_" + uniqueSuffix(),
	})
	assertSuccess(s.T(), createResp)
	roleID := extractID(s.T(), createResp.Data)

	// 删除
	resp := httpDelete(s.T(), "/roles/"+roleID, s.token)
	assertSuccess(s.T(), resp)

	// 验证已被删除
	getResp := httpGet(s.T(), "/roles/"+roleID, s.token)
	assertCode(s.T(), getResp, 404)
}

// ========== 分配权限 ==========

func (s *RoleTestSuite) TestRoleAssignPermissions_Success() {
	// 创建角色
	roleResp := httpPost(s.T(), "/roles", s.token, map[string]any{
		"name": "权限角色_" + uniqueSuffix(),
		"code": "perm_role_" + uniqueSuffix(),
	})
	assertSuccess(s.T(), roleResp)
	roleID := extractID(s.T(), roleResp.Data)

	// 创建权限
	permResp := httpPost(s.T(), "/permissions", s.token, map[string]any{
		"name":   "测试按钮权限_" + uniqueSuffix(),
		"code":   "test:btn_" + uniqueSuffix(),
		"type":   2,
		"status": 1,
	})
	assertSuccess(s.T(), permResp)
	permID := extractID(s.T(), permResp.Data)

	// 分配权限
	resp := httpPut(s.T(), "/roles/"+roleID+"/permissions", s.token, map[string]any{
		"permission_ids": []string{permID},
	})
	assertSuccess(s.T(), resp)

	// 验证角色已包含权限
	getResp := httpGet(s.T(), "/roles/"+roleID, s.token)
	assertSuccess(s.T(), getResp)
	var role map[string]any
	parseJSON(s.T(), getResp.Data, &role)
	permissions := role["permissions"].([]any)
	assert.GreaterOrEqual(s.T(), len(permissions), 1)
}

// ========== 未认证访问 ==========

func (s *RoleTestSuite) TestRoleList_NoToken() {
	resp := httpGet(s.T(), "/roles", "")
	assertCode(s.T(), resp, 401)
}

func (s *RoleTestSuite) TestRoleCreate_NoToken() {
	resp := httpPost(s.T(), "/roles", "", map[string]string{
		"name": "test",
		"code": "test",
	})
	assertCode(s.T(), resp, 401)
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleTestSuite))
}
