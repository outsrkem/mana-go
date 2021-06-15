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
	atTimes := time.Now().Unix()
	atTimesStr := time.Unix(atTimes, 0).Format("2006-01-02 15:04:05")
	// 使用uuid，并去除“-”
	uuid, _ := uuid.NewV4()
	uid := strings.Replace(uuid.String(), "-", "", -1)
	nickname := name
	role, expires, inactive := 3, 2, 1
	tx, err := mysql.DB.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("事务开启失败:%v\n", err)
		return userInfo, err
	}
	// 插入用户表信息
	sqlStr1 := `INSERT INTO user (USERID, USERNAME, NICKNAME, ROLE, PASSWD, UPDATETIME, EXPIRES, INACTIVE, CREATETIME) VALUES (?,?,?,?,?,?,?,?,?);`
	_, err = tx.Exec(sqlStr1, uid, name, nickname, role, passwd, atTimesStr, expires, inactive, atTimesStr)
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
	sqlStr := `SELECT ID,USERID,USERNAME,NICKNAME,ROLE,PASSWD,EXPIRES,INACTIVE,CREATETIME,UPDATETIME FROM  user WHERE USERNAME = ?`
	//fmt.Println(sqlStr)
	var row = mysql.DB.QueryRow(sqlStr, username)
	//err := row.Scan(u.ID, u.USERID, u.USERNAME, u.NICKNAME, u.ROLE, u.PASSWD, u.UPDATETIME, u.EXPIRES, u.INACTIVE, u.CREATETIME)
	err := row.Scan(&u.ID, &u.USERID, &u.USERNAME, &u.NICKNAME, &u.ROLE, &u.PASSWD, &u.EXPIRES, &u.INACTIVE, &u.CREATETIME, &u.UPDATETIME)
	if err != nil {
		fmt.Println("asd", err.Error())

	}
	return &u, err
}

// SelectUserQueryMultiRow 查询多条
func SelectUserQueryMultiRow(id int) {
	// 1.sql
	sqlStr := `SELECT * FROM  user WHERE ID > ?`
	// 2.执行
	rows, err := mysql.DB.Query(sqlStr, id)
	if err != nil {
		fmt.Println("exec %s query failed, err:%v\n", sqlStr, err)
	}
	// 3 一定要关闭连接
	defer rows.Close()
	// 4. 循环取值
	for rows.Next() {
		var u userInfo
		err := rows.Scan(&u.ID, &u.USERID, &u.USERNAME, &u.NICKNAME, &u.ROLE, &u.PASSWD, &u.UPDATETIME, &u.EXPIRES, &u.INACTIVE, &u.CREATETIME)
		if err != nil {
			fmt.Println("error")
		}
		fmt.Printf("u:%#v\n", u)
	}

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
