package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"mana/src/config"
	"mana/src/models"
	"net/http"
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
	r.RoleDesc = gjson.Get(string(data), "role_desc").String()

	item := make(map[string]int64)
	roleId, _ := models.InsertTheRole(r)
	item["id"] = roleId

	msg := models.NewResMessage("200", "successfully")
	returns := models.NewReturns(item, msg)
	c.JSON(http.StatusOK, &returns)
}