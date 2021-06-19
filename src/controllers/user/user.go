package user

import (
	"github.com/gin-gonic/gin"
	"mana/src/config"
	"mana/src/filters/util"
	"mana/src/models"
	"net/http"
)

// 日志
var _log = config.Log()

// 注册用户的body结构
type userRegisterInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// InstUser 用户注册
func InstUser(c *gin.Context) {

	var userRegisterInfo userRegisterInfo

	if err := c.BindJSON(&userRegisterInfo); err != nil {
		_log.Error("用户注册获取数据error", err)
	}
	password := userRegisterInfo.Password
	username := userRegisterInfo.Username
	// 用户名或密码不能为空
	if password == "" || username == "" {
		_log.Error("注册错误,用户名或密码为空")
		msg := models.NewResMessage("403", "The user name or password cannot be empty")
		c.JSON(http.StatusForbidden, msg)
		return
	}
	// 加密密码
	encodePassword, _ := util.PasswordBcrypt(password)
	// 把用户信息插入到数据库
	userInfo, err := models.InstUser(username, encodePassword)
	if err != nil {
		// 插入失败
		_log.Error("把用户信息插入到数据库失败")
		msg := models.NewResMessage("500", "internal error")
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	msg := models.NewResMessage("201", "registered successfully")
	returns := models.NewReturns(userInfo, msg)
	c.JSON(http.StatusCreated, returns)
}

// Login 用户登录
func Login(c *gin.Context) {
	var userLoginInfo models.UserLoginInfo
	if err := c.BindJSON(&userLoginInfo); err != nil {
		_log.Error("用户登录获取数据error", err)
	}
	username := userLoginInfo.Username
	loginPassword := userLoginInfo.Password
	// 用户名或密码不能为空
	if loginPassword == "" || username == "" {
		_log.Error("登录错误,用户名或密码为空")
		msg := models.NewResMessage("403", "The user name or password cannot be empty")
		c.JSON(http.StatusForbidden, msg)
		return
	}

	result, _ := models.SelectUserQueryRow(username)
	_log.Info("登录用户===> ", username)

	// 校验密码
	if !util.PasswordAuthentication(loginPassword, result.PASSWD) {
		_log.Error("密码校验失败")
		msg := models.NewResMessage("401", "Logon failed")
		c.JSON(http.StatusUnauthorized, msg)
		return
	}
	items := make(map[string]interface{}, 0)
	items["token"] = util.EncodeAuthToken(result.USERID, result.USERNAME, result.ROLE)  // 生成token
	items["userid"] = result.USERID
	items["username"] = result.USERNAME
	items["nickname"] = result.NICKNAME
	items["expires"] = result.EXPIRES
	msg := models.NewResMessage("200", "login Successfully")
	returns := models.NewReturns(items, msg)
	c.JSON(http.StatusOK, &returns)
}

// FindByUserinfo 查询用户信息
func FindByUserinfo(c *gin.Context) {
	uid := c.Param("uid")
	// 参数校验
	if !util.RegexpMatchString(uid,"^([0-9]|[a-f]){32}$") {
		msg := models.NewResMessage("400", "The parameter ID must be an integer, ^([0-9]|[a-f]){32}$ ")
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	result, err := models.SelectByUserInfo(uid)
	if err != nil {
		_log.Error("用户信息查询异常", err)
		msg := models.NewResMessage("404", "Query exception")
		c.JSON(http.StatusOK, msg)
		return
	}
	msg := models.NewResMessage("200", "successfully")
	returns := models.NewReturns(result, msg)
	c.JSON(http.StatusOK, returns)
}
