package models

import (
	"fmt"
	"mana/src/connections/database/mysql"
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
		RoleType:   1,
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

// DeleteRole 删除角色
func DeleteRole()  {
	//sqlStr := `DELETE FROM role WHERE id in(1033,1033,1035,1036)`
	rolrIdList := make([]string, 0)
	rolrIdList = append(rolrIdList, "1042")
	rolrIdList = append(rolrIdList, "1043")
	rolrIdList = append(rolrIdList, "1044")
	rolrIdList = append(rolrIdList, "1045")

	id := ""
	for _, v := range rolrIdList {
		if len(id) == 0 {
			id = id + v
			continue
		}
		id = id + "," + v
	}
	sqlStr := `DELETE FROM role WHERE id in(` + id + `)`
	fmt.Println(sqlStr)
	ret, _ := mysql.DB.Exec(sqlStr)
	n, _ := ret.RowsAffected()
	fmt.Println(n)
}
