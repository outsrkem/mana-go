package models

import (
    "fmt"
    "mana/src/connections/database/mysql"
    "mana/src/filters/uuid"
    "strings"
    "time"
)

// 用户表
type userInfo struct {
    ID         string `json:"Id"`         // Id
    USERID     string `json:"userid"`     // 用户id
    USERNAME   string `json:"username"`   // 用户名
    NICKNAME   string `json:"nickname"`   // 昵称
    ROLE       string `json:"role"`       // 角色
    PASSWD     string `json:"passwd"`     // 密码
    UPDATETIME string `json:"change"`     // 最近一次密码修改时间
    EXPIRES    string `json:"expires"`    // 密码过期时间
    INACTIVE   string `json:"inactive"`   // 用户状态
    CREATETIME string `json:"createtime"` // 创建时间
}

type userCenter struct {
    ID         string `json:"Id"`         // Id
    USERID     string `json:"userid"`     // 用户id,外键
    USERNAME   string `json:"username"`   // 用户名，外键
    NICKNAME   string `json:"nickname"`   // 昵称
    MOBILE     string `json:"mobile"`     //手机
    EMAIL      string `json:"email"`      // 邮箱
    DESCRIBES  string `json:"describes"`  // 描述说明
    PICTURE    string `json:"picture"`    // 头像
    CREATETIME string `json:"createtime"` // 创建时间
    UPDATETIME string `json:"updatetime"` // 最近更新时间
}

// InstUser 注册用户
func InstUser(name string, passwd string) (map[string]string, error) {
    userInfo := make(map[string]string)
    //atTimes := time.Now().UnixNano() / 1e6
    //atTimesStr := time.Unix(atTimes, 0).Format("2006-01-02 15:04:05")
    atTimesStr := time.Now().UnixNano() / 1e6
    // 使用uuid，并去除“-”
    uuid, _ := uuid.NewV4()
    uid := strings.Replace(uuid.String(), "-", "", -1)
    nickname := name
    expires, inactive := 2, 1
    tx, err := mysql.DB.Begin() // 开启事务
    if err != nil {
        if tx != nil {
            tx.Rollback() // 回滚
        }
        fmt.Printf("事务开启失败:%v\n", err)
        return userInfo, err
    }
    // 插入用户表信息
    sqlStr1 := `INSERT INTO user (USERID, PASSWD, UPDATETIME, EXPIRES, INACTIVE, CREATETIME) VALUES (?,?,?,?,?,?);`
    _, err = tx.Exec(sqlStr1, uid, passwd, atTimesStr, expires, inactive, atTimesStr)
    if err != nil {
        tx.Rollback() // 回滚
        fmt.Printf("用户表插入失败:%v\n", err)
        return userInfo, err
    }
    // 插入用户中心表信息
    sqlStr2 := `INSERT INTO user_center (USERID, USERNAME, NICKNAME, CREATETIME, UPDATETIME) VALUES (?,?,?,?,?);`
    _, err = tx.Exec(sqlStr2, uid, name, nickname, atTimesStr, atTimesStr)
    if err != nil {
        tx.Rollback() // 回滚
        fmt.Printf("用户中心表插入失败:%v\n", err)
        return userInfo, err
    }
    // 提交事务
    if err = tx.Commit(); err != nil {
        // 事务回滚
        tx.Rollback()
        fmt.Println("事务回滚...")
        return userInfo, err
    }
    userInfo["userid"] = uid
    userInfo["username"] = name
    return userInfo, err
}

// SelectUserQueryRow 查询单条
func SelectUserQueryRow(username string) (*userInfo, error) {
    var u userInfo
    sqlStr := `SELECT ue.USERID,uc.USERNAME,uc.NICKNAME,ue.PASSWD,ue.EXPIRES FROM user ue INNER JOIN user_center uc ON (ue.USERID = uc.USERID) WHERE uc.USERNAME = ?`
    //fmt.Println(sqlStr)
    var row = mysql.DB.QueryRow(sqlStr, username)
    //err := row.Scan(u.ID, u.USERID, u.USERNAME, u.NICKNAME, u.ROLE, u.PASSWD, u.UPDATETIME, u.EXPIRES, u.INACTIVE, u.CREATETIME)
    err := row.Scan(&u.USERID, &u.USERNAME, &u.NICKNAME, &u.PASSWD, &u.EXPIRES)
    if err != nil {
        fmt.Println("asd", err.Error())

    }
    return &u, err
}

// SelectUidUserQueryRow 查询单条
func SelectUidUserQueryRow(uid string) (*userInfo, error) {
    var u userInfo
    sqlStr := `SELECT ID,USERID,USERNAME,NICKNAME,ROLE,PASSWD,EXPIRES,INACTIVE,CREATETIME,UPDATETIME FROM  user WHERE USERID = ?`
    var row = mysql.DB.QueryRow(sqlStr, uid)
    err := row.Scan(&u.ID, &u.USERID, &u.USERNAME, &u.NICKNAME, &u.ROLE, &u.PASSWD, &u.EXPIRES, &u.INACTIVE, &u.CREATETIME, &u.UPDATETIME)
    if err != nil {
        fmt.Println("asd", err.Error())

    }
    return &u, err
}

// SelectByUserInfo 查询用户详细信息
func SelectByUserInfo(uid string) (*userCenter, error) {
    var u userCenter
    sqlStr := `SELECT user_center.ID, user_center.USERID,user_center.USERNAME,user_center.NICKNAME,
                user_center.MOBILE,user_center.EMAIL,user_center.DESCRIBES,user_center.PICTURE,user_center.CREATETIME,user_center.UPDATETIME
                FROM user inner join user_center on   (user.USERID=user_center.USERID) WHERE user.USERID=?`
    var row = mysql.DB.QueryRow(sqlStr, uid)
    err := row.Scan(&u.ID, &u.USERID, &u.USERNAME, &u.NICKNAME, &u.MOBILE, &u.EMAIL, &u.DESCRIBES, &u.PICTURE, &u.CREATETIME, &u.UPDATETIME)
    if err != nil {
        fmt.Println("asd", err.Error())
    }
    return &u, err
}

// GetUserLists 获取用户列表
// 后续优化，表结构
func GetUserLists() *map[string]interface{} {
    /* 查询用户信息
        SELECT
            u.ID,
            u.USERID
        FROM
            user u
        LEFT JOIN user_center uc ON (u.USERID = uc.USERID)
     */
    sqlStr := `SELECT ID,USERID,USERNAME,UPDATETIME,CREATETIME FROM user;`
    rows, err := mysql.DB.Query(sqlStr)
    if err != nil {
        log.Info(sqlStr, err)
    }
    defer rows.Close()

    var u userInfo
    var items []map[string]interface{}
    items = make([]map[string]interface{}, 0)
    for rows.Next() {
        if rows.Scan(&u.ID, &u.USERID, &u.USERNAME,&u.UPDATETIME,&u.CREATETIME) != nil {
            log.Info("Get the user list, ", err.Error())
        }

        item := make(map[string]interface{})
        item["id"] = u.ID
        item["user_id"] = u.USERID
        item["user_name"] = u.USERNAME
        item["nickname"] = u.NICKNAME
        item["create_time"] = u.CREATETIME
        item["update_time"] = u.UPDATETIME
        items = append(items, item)
    }

    returns := NewResponse(items, nil)
    return &returns
}

// 优化
type users struct {
	id         string `json:"id"`
	userId     string `json:"user_id"`
	passwd     string `json:"passwd"`
	expires    int    `json:"expires"`
	inactive   int    `json:"inactive"`
	username   string `json:"username"`
	nickname   string `json:"nickname"`
	mobile     int64  `json:"mobile"`
	email      string `json:"email"`
	describes  string `json:"describes"`
	picture    string `json:"picture"`
	createTime int64  `json:"create_time"`
	updateTime int64  `json:"update_time"`
}

// SelectUsersQueryMultiRow 查询用户表
// uId 为用户id，即字段USERID
// 如果为空，则查询所有用户，并分页
// 	SelectUsersQueryMultiRow("89b8bd3386ab46c5a906bd4e3818bbca",1,5)
func SelectUsersQueryMultiRow(uId string, page, pageSize int) []map[string]interface{} {
    n, m := (page - 1) * pageSize, pageSize

    sqlCountStr := `SELECT COUNT(*) FROM user ue INNER JOIN user_center uc ON (ue.USERID = uc.USERID) WHERE 1=1`
	// 1.sql
	sqlStr := `SELECT uc.ID AS id, ue.USERID AS user_id, ue.PASSWD AS passwd, ue.EXPIRES AS expires,
			ue.INACTIVE AS inactive, uc.USERNAME AS username, uc.NICKNAME AS nickname, uc.MOBILE AS mobile,
			uc.EMAIL AS email, uc.DESCRIBES AS describes, uc.PICTURE AS picture, uc.CREATETIME AS create_time,
			uc.UPDATETIME AS update_time FROM user ue INNER JOIN user_center uc ON (ue.USERID = uc.USERID) WHERE 1=1`
	if uId != "" {
		sqlStr = sqlStr + ` AND ue.USERID = ` + `'` + uId + `'`
        sqlCountStr += ` AND ue.USERID = ` + `'` + uId + `'`
	}
    sqlStr += ` ORDER BY id LIMIT ?, ?;`
    log.Debug(sqlStr)
    log.Debug("n=",n ," m=", m)

    // 查询总记录数
    totalRow, _ := mysql.DB.Query(sqlCountStr)
    var total, pageNum int
    for totalRow.Next() {
        err := totalRow.Scan(
            &total,
        )
        if err != nil {
            log.Error("GetKnowledgePointListTotal error", err)
            continue
        }
    }
    log.Debug("总记录数：",total," pageNum: ",pageNum)
    // 查询记录
	rows, err := mysql.DB.Query(sqlStr,n , m) // 2.执行sql
	if err != nil {
		log.Error("exec %s query failed, err:%v\n",sqlStr, err)
	}
	// 3 一定要关闭连接
	defer rows.Close()
	// 4. 循环取值
    items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var u users
        item := make(map[string]interface{})
		if nil != rows.Scan(&u.id, &u.userId, &u.passwd, &u.expires, &u.inactive, &u.username, &u.nickname,
		    &u.mobile, &u.email, &u.describes, &u.picture, &u.createTime, &u.updateTime) {
            log.Error("SelectUsersQueryMultiRow Scan error", err)
		}
        item["id"] = u.id
        item["user_id"] = u.userId
        item["passwd"] = u.passwd
        item["expires"] = u.expires
        item["inactive"] = u.inactive
        item["username"] = u.username
        item["nickname"] = u.nickname
        item["mobile"] = u.mobile
        item["email"] = u.email
        item["describes"] = u.describes
        item["picture"] = u.picture
        item["create_time"] = u.createTime
        item["update_time"] = u.updateTime
        items = append(items, item)
	}
    return items
}
