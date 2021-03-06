package auth

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/tidwall/gjson"
    "mana/src/models"
    "net/http"
)

// UpdateUserRole 更新角色权限
// {"user_id":"5b4b0238d6e04c319c966aac6cd813a1","role_id":[1002,1003]}
func UpdateUserRole(c *gin.Context) {
    raw := make(map[string]interface{})
    c.ShouldBind(&raw)
    data, _ := json.Marshal(raw)
    userId := gjson.Get(string(data), "user_id").String() // 转化为string

    // 获取角色转化为切片
    roleId := gjson.Get(string(data), "role_id").String()
    roleIdList := make([]int64, 0)
    json.Unmarshal([]byte(roleId), &roleIdList)

    if nil != models.UpdateUserRoles(userId, &roleIdList) {
        msg := models.NewResMessage("500", "The user failed to bind the role")
        c.JSON(http.StatusInternalServerError, msg)
        log.Error("UpdateUserRole error, The user failed to bind the role")
        return
    }
    msg := models.NewResMessage("200", "Successful.")
    c.JSON(http.StatusOK, &msg)
}

func GetUserList(c *gin.Context) {
    res := models.GetUserLists()
    msg := models.NewResMessage("200", "successfully")
    returns := models.NewReturns(res, msg)
    c.JSON(http.StatusOK, &returns)
}


func GetUserRoleList(c *gin.Context) {
    uid := c.DefaultQuery("uid","0" )
    // 查询用户列表
    roleList := models.FindByRoleList()
    // 查询已授权的角色id
    authorizedId := models.FindByAuthorizedRoleId(uid)

    var items map[string]interface{}
    items = make(map[string]interface{}, 0)
    //fmt.Println(*roleList["items"])
    items["role_list"] = roleList
    items["authorized"] = authorizedId


    msg := models.NewResMessage("200", "successfully")
    returns := models.NewReturns(items, msg)
    c.JSON(http.StatusOK, &returns)
}
