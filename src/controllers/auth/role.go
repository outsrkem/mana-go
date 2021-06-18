package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"mana/src/config"
	"mana/src/models"
	"net/http"
	"regexp"
)

var log = config.Log()

func GetRoleList(c *gin.Context) {

	res := models.FindByRoleList()
	msg := models.NewResMessage("200", "successfully")
	returns := models.NewReturns(res, msg)
	c.JSON(http.StatusOK, &returns)
}

// AddRole 添加角色
func AddRole(c *gin.Context) {
	var r = models.NewRoleList()

	role := make(map[string]interface{})
	c.ShouldBind(&role)

	data, err := json.Marshal(role)
	if err != nil {
		msg := models.NewResMessage("406", "JSON serialization error.")
		c.JSON(http.StatusNotAcceptable, &msg)
		return
	}

	r.RoleName = gjson.Get(string(data), "role_name").String()
	matched, _ := regexp.MatchString("^([A-Za-z_]{3,20})$", r.RoleName)
	if !matched {
		msg := models.NewResMessage("400", "The role name does not meet the requirements: , ^([A-Za-z]{3,20})$ ")
		c.JSON(http.StatusBadRequest, msg)
		log.Error("AddRole. role_name: ", r.RoleName)
		return
	}
	r.RoleDesc = gjson.Get(string(data), "role_desc").String()

	item := make(map[string]int64)
	roleId, _ := models.InsertTheRole(r)
	item["id"] = roleId

	msg := models.NewResMessage("200", "successfully")
	returns := models.NewReturns(item, msg)
	c.JSON(http.StatusOK, &returns)
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context)  {
	roleId := make(map[string]interface{})
	c.ShouldBind(&roleId)

	data, _ := json.Marshal(roleId)
	idList := gjson.Get(string(data), "role_id").String()
	// [1001,1003,1002]

	RoleIdListSlice := make([]int64, 0)
	json.Unmarshal([]byte(idList), &RoleIdListSlice)

	result, err := models.DeleteRoles(&RoleIdListSlice)
	if err != nil {
		msg := models.NewResMessage("400", "Delete fail.")
		c.JSON(http.StatusBadRequest, &msg)
		return
	}

	item := make(map[string]int64)
	item["result"] = result
	msg := models.NewResMessage("200", "Successful.")
	returns := models.NewReturns(item, msg)
	c.JSON(http.StatusOK, &returns)
}

// UpdateRolePermission 更新角色权限
func UpdateRolePermission(c *gin.Context) {
	raw := make(map[string]interface{})
	c.ShouldBind(&raw)
	data, _ := json.Marshal(raw)
	roleId := gjson.Get(string(data), "role_id").String() // 转化为string

	// 获取权限转化为切片
	permissionId := gjson.Get(string(data), "permission_id").String()
	permissionIdList := make([]int64, 0)
	json.Unmarshal([]byte(permissionId), &permissionIdList)

	if nil != models.UpdateRolePermissions(roleId, &permissionIdList) {
		msg := models.NewResMessage("500", "Failed to update permissions")
		c.JSON(http.StatusInternalServerError, msg)
		log.Error("UpdateRolePermission error")
		return
	}
	msg := models.NewResMessage("200", "Successful.")
	c.JSON(http.StatusOK, &msg)
}