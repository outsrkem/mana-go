package menus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mana/src/config"
	"mana/src/models"
	"net/http"
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
