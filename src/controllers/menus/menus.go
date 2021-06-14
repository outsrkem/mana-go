package menus

import (
	"github.com/gin-gonic/gin"
	"mana/src/config"
	"mana/src/models"
	"net/http"
)

var log = config.Log()

func GetMenus(c *gin.Context) {
	/*
	   	jsonData := `{
	       "id": 101,
	       "authName": "商品管理",
	       "path": null,
	       "children": [
	           {
	               "id": 104,
	               "authName": "商品列表",
	               "path": null,
	               "children": []
	           }
	       ]
	   }`
	*/
	// 临时返回
	msg := models.NewResMessage("200", "successfully")
	res := models.SelectMenuList()
	returns := models.NewReturns(res, msg)
	c.JSON(http.StatusOK, &returns)
}
