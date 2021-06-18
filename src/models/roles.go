package models

import (
	"fmt"
	"mana/src/connections/database/mysql"
	"strconv"
	"time"
)

type roleList struct {
	Id         int    `json:"role_id"`
	RoleName   string `json:"role_name"`
	RoleState  int    `json:"role_state"`
	RoleType   int    `json:"role_type"`
	RoleDesc   string `json:"role_desc"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

func NewRoleList() *roleList {
	return &roleList{
		RoleState:  1,
		RoleType:   2,
		CreateTime: time.Now().UnixNano() / 1e6,
		UpdateTime: time.Now().UnixNano() / 1e6,
	}
}

// FindByRoleList 查询所有角色
func FindByRoleList() *map[string]interface{} {

	sqlStr := `SELECT id, r_name, r_state, r_type, r_desc, create_time, update_time FROM role;`
	rows, err := mysql.DB.Query(sqlStr)
	if err != nil {
		log.Info(sqlStr, err)
	}
	defer rows.Close()

	var r roleList
	var items []map[string]interface{}
	items = make([]map[string]interface{}, 0)
	for rows.Next() {
		if rows.Scan(&r.Id, &r.RoleName, &r.RoleState, &r.RoleType, &r.RoleDesc, &r.CreateTime, &r.UpdateTime) != nil {
			log.Info("Get the role list, ", err.Error())
		}

		item := make(map[string]interface{})
		item["role_id"] = r.Id
		item["role_name"] = r.RoleName
		item["role_state"] = r.RoleState
		item["role_type"] = r.RoleType
		item["role_desc"] = r.RoleDesc
		item["create_time"] = r.CreateTime
		item["update_time"] = r.UpdateTime
		items = append(items, item)
	}
	returns := NewResponse(items, nil)
	return &returns
}

// InsertTheRole 插入角色
func InsertTheRole(r *roleList) (int64, error) {
	sqlStr := `INSERT INTO role (r_name, r_state, r_type, r_desc, create_time, update_time ) VALUE (?,?,?,?,?,?)`
	ret, err := mysql.DB.Exec(sqlStr, r.RoleName, r.RoleState, r.RoleType, r.RoleDesc, r.CreateTime, r.UpdateTime)
	if err != nil {
		log.Error("insert role failed, ", err)
		return -1, err
	}
	id, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		log.Error("get last Insert role ID failed, err:", err)
		return -1, err
	}
	return id, err
}

// DeleteRoles 删除角色
func DeleteRoles(idList *[]int64) (int64, error) {
	//sqlStr := `DELETE FROM role WHERE id in(1033,1033,1035,1036)`
	if len(*idList) < 1 {
		return -1, fmt.Errorf("The slice cannot be empty. ")
	}
	id := ""
	for _, v := range *idList {
		if len(id) == 0 {
			id = id + strconv.FormatInt(v, 10)
			continue
		}
		id = id + "," + strconv.FormatInt(v, 10)
	}

	sqlStr := `DELETE FROM role WHERE id in(` + id + `)`

	log.Debug("DeleteRoles, ", sqlStr)
	ret, err := mysql.DB.Exec(sqlStr)
	if err != nil {
		log.Error("Delete role failed, ", err)
		return -1, err
	}
	n, _ := ret.RowsAffected()
	if err != nil {
		log.Error("RowsAffected returns the number of rows affected by an update failed, ", err)
	}

	return n, nil
}

// 更新角色权限，清空权限，新增权限
func UpdateRolePermissions(roleId string, perList *[]int64) error {
	sqlStr := `SELECT id from role_menu WHERE id=0;`
	if len(*perList) != 0 {
		val := ""
		createTime := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		// INSERT INTO role_menu (mid, rid) VALUES (2, 0000),(2,0000)
		for _, v := range *perList {
			if len(val) == 0 {
				val = val + `(` + strconv.FormatInt(v, 10) + `,` + roleId + `,` + createTime + `)`
				continue
			}
			val = val + `,` + `(` + strconv.FormatInt(v, 10) + `,` + roleId + `,` + createTime + `)`

		}

		sqlStr = `INSERT INTO role_menu (mid, rid, create_time) VALUES ` + val + `;`
	}

	tx, err := mysql.DB.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		log.Error("UpdateRolePermission， 开启事务失败", err)
		return err
	}
	// 删除原有的角色权限
	_, err = tx.Exec(`DELETE FROM role_menu WHERE rid=?;`, roleId)
	if err != nil {
		tx.Rollback()
		log.Error("UpdateRolePermission，删除原有的角色权限失败 ", err)
		return err
	}
	// 添加新的角色权限
	_, err = tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		log.Error("UpdateRolePermission，添加新的角色权限失败 ", err, sqlStr)
		return err
	}
	// 提交事务
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Error("UpdateRolePermission，提交事务失败，事务回滚...", err)
		return err
	}
	return nil
}

// UpdateUserRoles 更新用户角色
func UpdateUserRoles(uid string, roleList *[]int64) error {
	sqlStr := `SELECT id from role_menu WHERE id=0;`
	if len(*roleList) != 0 {
		val := ""
		createTime := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		for _, v := range *roleList {
			if len(val) == 0 {
				val = val + `(` + strconv.FormatInt(v, 10) + `,` + `'` + uid + `'` + `,` + createTime + `)`
				continue
			}
			val = val + `,` + `(` + strconv.FormatInt(v, 10) + `,` + `'` + uid + `'` + `,` + createTime + `)`

		}

		sqlStr = `INSERT INTO role_user (rid, userid, create_time) VALUES ` + val + `;`
	}
	log.Debug(sqlStr)
	log.Debug(uid)
	log.Debug(roleList)
	tx, err := mysql.DB.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		log.Error("UpdateUserRole， 开启事务失败", err)
		return err
	}
	// 删除原有的角色权限
	_, err = tx.Exec(`DELETE FROM role_user WHERE userid=?;`, uid)
	if err != nil {
		tx.Rollback()
		log.Error("UpdateUserRole，删除用户绑定的角色失败 ", err)
		return err
	}
	// 添加新的角色权限
	_, err = tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		log.Error("UpdateUserRole，添加新的角色失败 ", err, sqlStr)
		return err
	}
	// 提交事务
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Error("UpdateUserRole，提交事务失败，事务回滚...", err)
		return err
	}
	return nil
}


