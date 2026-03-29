package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PermissionTestSuite struct {
	suite.Suite
	token string
}

func (s *PermissionTestSuite) SetupSuite() {
	s.token = loginDefault(s.T())
}

// ========== 权限列表 ==========

func (s *PermissionTestSuite) TestPermissionList() {
	resp := httpGet(s.T(), "/permissions"+buildQuery(map[string]string{"page": "1", "size": "10"}), s.token)
	assertSuccess(s.T(), resp)

	var pr PageResult
	parseJSON(s.T(), resp.Data, &pr)
	assert.GreaterOrEqual(s.T(), pr.Total, int64(1))
}

func (s *PermissionTestSuite) TestPermissionList_WithKeyword() {
	resp := httpGet(s.T(), "/permissions"+buildQuery(map[string]string{"keyword": "system", "page": "1", "size": "10"}), s.token)
	assertSuccess(s.T(), resp)

	var pr PageResult
	parseJSON(s.T(), resp.Data, &pr)
	assert.GreaterOrEqual(s.T(), pr.Total, int64(1))
}

// ========== 权限树 ==========

func (s *PermissionTestSuite) TestPermissionTree() {
	resp := httpGet(s.T(), "/permissions/tree", s.token)
	assertSuccess(s.T(), resp)

	var tree []map[string]any
	parseJSON(s.T(), resp.Data, &tree)
	assert.NotEmpty(s.T(), tree)
}

// ========== 权限详情 ==========

func (s *PermissionTestSuite) TestPermissionGet_Success() {
	// 先从列表获取
	listResp := httpGet(s.T(), "/permissions"+buildQuery(map[string]string{"page": "1", "size": "1"}), s.token)
	assertSuccess(s.T(), listResp)

	var pr PageResult
	parseJSON(s.T(), listResp.Data, &pr)
	assert.NotEmpty(s.T(), pr.Items)

	permID := extractID(s.T(), pr.Items[0])

	resp := httpGet(s.T(), "/permissions/"+permID, s.token)
	assertSuccess(s.T(), resp)

	var data map[string]any
	parseJSON(s.T(), resp.Data, &data)
	assert.NotEmpty(s.T(), data["id"])
	assert.NotEmpty(s.T(), data["name"])
	assert.NotEmpty(s.T(), data["code"])
}

func (s *PermissionTestSuite) TestPermissionGet_NotFound() {
	resp := httpGet(s.T(), "/permissions/not_exist_id", s.token)
	assertCode(s.T(), resp, 404)
}

// ========== 创建权限 ==========

func (s *PermissionTestSuite) TestPermissionCreate_MenuType() {
	code := "test:menu:create_" + uniqueSuffix()
	resp := httpPost(s.T(), "/permissions", s.token, map[string]any{
		"name":        "测试菜单权限",
		"code":        code,
		"description": "测试创建的菜单权限",
		"type":        1, // 菜单
		"sort":        200,
		"status":      1,
	})
	assertSuccess(s.T(), resp)

	var data map[string]any
	parseJSON(s.T(), resp.Data, &data)
	createdID := data["id"].(string)

	getResp := httpGet(s.T(), "/permissions/"+createdID, s.token)
	assertSuccess(s.T(), getResp)
	var perm map[string]any
	parseJSON(s.T(), getResp.Data, &perm)
	assert.Equal(s.T(), "测试菜单权限", perm["name"])
	assert.Equal(s.T(), code, perm["code"])
	assert.Equal(s.T(), float64(1), perm["type"])
}

func (s *PermissionTestSuite) TestPermissionCreate_ButtonType() {
	resp := httpPost(s.T(), "/permissions", s.token, map[string]any{
		"name": "测试按钮权限_" + uniqueSuffix(),
		"code": "test:btn:create_" + uniqueSuffix(),
		"type": 2, // 按钮
	})
	assertSuccess(s.T(), resp)
}

func (s *PermissionTestSuite) TestPermissionCreate_ApiType() {
	resp := httpPost(s.T(), "/permissions", s.token, map[string]any{
		"name": "测试接口权限_" + uniqueSuffix(),
		"code": "test:api:create_" + uniqueSuffix(),
		"type": 3, // 接口
	})
	assertSuccess(s.T(), resp)
}

func (s *PermissionTestSuite) TestPermissionCreate_DuplicateCode() {
	resp := httpPost(s.T(), "/permissions", s.token, map[string]string{
		"name": "重复权限",
		"code": "system:user:list",
	})
	assertCode(s.T(), resp, 500)
}

func (s *PermissionTestSuite) TestPermissionCreate_MissingFields() {
	resp := httpPost(s.T(), "/permissions", s.token, map[string]string{
		"name": "缺少编码",
	})
	assertCode(s.T(), resp, 400)
}

// ========== 更新权限 ==========

func (s *PermissionTestSuite) TestPermissionUpdate_Success() {
	// 先创建
	createResp := httpPost(s.T(), "/permissions", s.token, map[string]any{
		"name": "待更新权限_" + uniqueSuffix(),
		"code": "test:update_" + uniqueSuffix(),
		"type": 2,
	})
	assertSuccess(s.T(), createResp)
	permID := extractID(s.T(), createResp.Data)

	// 更新
	resp := httpPut(s.T(), "/permissions/"+permID, s.token, map[string]any{
		"name":        "已更新权限",
		"description": "更新后的描述",
		"status":      1,
	})
	assertSuccess(s.T(), resp)

	// 验证
	getResp := httpGet(s.T(), "/permissions/"+permID, s.token)
	assertSuccess(s.T(), getResp)
	var data map[string]any
	parseJSON(s.T(), getResp.Data, &data)
	assert.Equal(s.T(), "已更新权限", data["name"])
	assert.Equal(s.T(), "更新后的描述", data["description"])
}

// ========== 删除权限 ==========

func (s *PermissionTestSuite) TestPermissionDelete_Success() {
	// 先创建
	createResp := httpPost(s.T(), "/permissions", s.token, map[string]any{
		"name": "待删除权限_" + uniqueSuffix(),
		"code": "test:delete_" + uniqueSuffix(),
		"type": 2,
	})
	assertSuccess(s.T(), createResp)
	permID := extractID(s.T(), createResp.Data)

	// 删除
	resp := httpDelete(s.T(), "/permissions/"+permID, s.token)
	assertSuccess(s.T(), resp)

	// 验证已被删除
	getResp := httpGet(s.T(), "/permissions/"+permID, s.token)
	assertCode(s.T(), getResp, 404)
}

// ========== 未认证访问 ==========

func (s *PermissionTestSuite) TestPermissionList_NoToken() {
	resp := httpGet(s.T(), "/permissions", "")
	assertCode(s.T(), resp, 401)
}

func (s *PermissionTestSuite) TestPermissionTree_NoToken() {
	resp := httpGet(s.T(), "/permissions/tree", "")
	assertCode(s.T(), resp, 401)
}

func (s *PermissionTestSuite) TestPermissionCreate_NoToken() {
	resp := httpPost(s.T(), "/permissions", "", map[string]string{
		"name": "test",
		"code": "test:code",
	})
	assertCode(s.T(), resp, 401)
}

func TestPermissionSuite(t *testing.T) {
	suite.Run(t, new(PermissionTestSuite))
}
