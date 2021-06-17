package menus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mana/src/config"
	"mana/src/models"
	"net/http"
	"strconv"
)

var log = config.Log()

func GetMenus(c *gin.Context) {
	uid,_ := c.Get("uid")
	userId := fmt.Sprintf("%s", uid)
	res := models.SelectMenuList(userId)
	msg := models.NewResMessage("200", "successfully")
	returns := models.NewReturns(res, msg)
	c.JSON(http.StatusOK, &returns)
}

func GetMenusAll(c *gin.Context) {
	rid, _ := strconv.ParseInt(c.DefaultQuery("rid","0" ), 10, 64)
	// 查询菜单列表
	menuList := models.FindByMenuListAll()
	// 查询已授权的id
	authorizedId := models.FindByAuthorizedId(rid)

	var items map[string]interface{}
	items = make(map[string]interface{}, 0)
	items["menu_list"] = menuList
	items["authorized"] = authorizedId


	msg := models.NewResMessage("200", "successfully")
	returns := models.NewReturns(items, msg)
	c.JSON(http.StatusOK, &returns)
}
