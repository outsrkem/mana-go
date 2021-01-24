package user

import (
	"github.com/gin-gonic/gin"
	"mana/src/config"
	"mana/src/filters/util"
	"mana/src/models"
	"net/http"
	"time"
)

// 日志
var _log = config.Log()

// 注册用户的body结构
type userRegisterInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 用户注册
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

// 用户登录
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
	err := util.PasswordAuthentication(loginPassword, result.PASSWD)
	if err != nil {
		_log.Error("登录错误", err)
		msg := models.NewResMessage("401", "Logon failed")
		c.JSON(http.StatusUnauthorized, msg)
	} else {
		var user models.UserLoginStruct
		meta := &user.MetaInfo
		resp := &user.Response

		// 生成token
		token := util.EncodeAuthToken(result.USERID, result.USERNAME, result.ROLE)

		// 构造返回数据
		meta.RequestTime = time.Now().UnixNano()
		meta.Msg = "login successfully"
		meta.Code = "200"
		resp.Userid = result.USERID
		resp.Username = result.USERNAME
		resp.Nickname = result.NICKNAME
		resp.Role = result.ROLE
		resp.Expires = result.EXPIRES
		resp.Token = token
		c.JSON(http.StatusOK, &user)
	}
}

// 获取用户信息

//func FindByUserinfo(c *gin.Context) {
//	//fmt.Println(c.MustGet("uid").(string))
//	//fmt.Println(c.MustGet("role").(string))
//	// 若有高权限的token，则可以查询其他用户，则此处需要传递用户id，后面优化
//	//uid := c.MustGet("uid").(string)
//	uid := c.Param("uid") // 获取路径参数
//	result, err := models.SelectUidUserQueryRow(uid)
//	if err != nil {
//		_log.Error("用户信息查询异常", err)
//		var user models.MetaInfo
//		user.RequestTime = time.Now().UnixNano()
//		user.Msg = "Query exception"
//		user.Code = "1"
//		c.JSON(http.StatusUnauthorized, &user)
//	} else {
//		var user models.UserLoginStruct
//		meta := &user.MetaInfo
//		resp := &user.Response
//
//		// 构造返回数据
//		meta.RequestTime = time.Now().UnixNano()
//		meta.Msg = "successfully"
//		meta.Code = "200"
//		resp.Userid = result.USERID
//		resp.Username = result.USERNAME
//		resp.Nickname = result.NICKNAME
//		resp.Role = result.ROLE
//		resp.Expires = result.EXPIRES
//		c.JSON(http.StatusOK, &user)
//	}
//}
// 查询用户信息
func FindByUserinfo(c *gin.Context) {
	result, err := models.SelectByUserInfo(c.Param("uid"))
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
