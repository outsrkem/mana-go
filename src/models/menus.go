package models

import (
    "fmt"
    "mana/src/connections/database/mysql"
)

// menuList 一级菜单
type menuList struct {
    sub_m_id    string `json:"sub_m_id"`
    sub_m_name  string `json:"sub_m_name"`
    sub_m_path  string `json:"sub_m_path"`
    sub_m_level string `json:"sub_m_level"`
    sub_p_code  string `json:"sub_p_code"`
    sub_desc    string `json:"sub_desc"`
    CreateTime  int64  `json:"create_time"`
    UpdateTime  int64  `json:"update_time"`
}

// menuListLevel1 一级菜单
type menuListLevel1 struct {
    id       string              `json:"Id"`
    name     string              `json:"name"`
    path     string              `json:"path"`
    leafNode []map[string]string `json:"leafNode"`
}

type menus struct {
    id             int64 `json:"id"`
    level          string `json:"level"`
    name           string `json:"name"`
    path           string `json:"path"`
    p_code         string `json:"p_code"`
    parent_menu_id string `json:"parent_menu_id"`
    icon_name      string `json:"icon_name"`
    description    string `json:"description"`
    CreateTime     int64  `json:"create_time"`
    UpdateTime     int64  `json:"update_time"`
}

// selectMenuLevel2 根据查询二级菜单
// 根据用户权限过滤二级菜单
func selectMenuLevel2(mId, userId string) []map[string]interface{} {

    // 1.查询用户的权限
    pCodeArray := FindByMenuPermission(userId)

    // 2.查询菜单的时候匹配权限
    sqlStr := `SELECT
					sub.Id AS sub_m_id,
					sub.name AS sub_m_name,
					sub.path AS sub_m_path,
					sub.level AS sub_m_level,
					sub.p_code
				FROM
					menus m
				LEFT JOIN menus sub ON (m.Id = sub.parent_menu_id)
				WHERE
					m.Id = ?;`

    log.Debug(sqlStr)
    rows, err := mysql.DB.Query(sqlStr, mId)
    if err != nil {
        log.Info(sqlStr, err)
    }

    defer rows.Close()

    var items []map[string]interface{}
    items = make([]map[string]interface{}, 0)

    for rows.Next() {
        var m menuList
        err := rows.Scan(&m.sub_m_id, &m.sub_m_name, &m.sub_m_path, &m.sub_m_level, &m.sub_p_code)
        if err != nil {
            log.Info("Get the secondary menu, ", err.Error())
        }

        if m.sub_m_id != "" {
            if PermissionMatchCheck(pCodeArray, m.sub_p_code) != nil {
                continue
            }
            item := make(map[string]interface{})
            item["Id"] = m.sub_m_id
            item["name"] = m.sub_m_name
            item["path"] = m.sub_m_path
            item["leafNode"] = ""
            items = append(items, item)
        }

    }
    return items
}

// SelectMenuList 查询菜单
func SelectMenuList(userId string) *map[string]interface{} {
    // 查询一级菜单 level=1
    var items []map[string]interface{}
    items = make([]map[string]interface{}, 0)

    sqlStr_1 := `SELECT Id,name,path FROM menus WHERE level =1;`
    rows_1, err := mysql.DB.Query(sqlStr_1)
    if err != nil {
        log.Info(sqlStr_1, err)
    }
    defer rows_1.Close()
    for rows_1.Next() {
        var m menuListLevel1
        err := rows_1.Scan(&m.id, &m.name, &m.path)
        if err != nil {
            log.Info("Get the first level menu, ", err.Error())
        }
        item := make(map[string]interface{})

        leafNode := selectMenuLevel2(m.id, userId)
        if len(leafNode) == 0 {
            log.Debug("没有子菜单，或略一级菜单, 一级菜单id: ", m.id)
            continue
        }
        item["Id"] = m.id
        item["name"] = m.name
        item["path"] = m.path
        item["leafNode"] = leafNode
        items = append(items, item)
    }
    returns := NewResponse(items, nil)
    return &returns
}

// FindByMenuPermission 获取用户权限码列表
func FindByMenuPermission(userId string) *[]string {
    sqlStr := `SELECT
					m.p_code
				FROM
					role_user ru,
					role_menu rm,
					menus m
				WHERE
					ru.rid = rm.rid
				AND rm.mid = m.Id
				AND ru.userid = ?
				AND m.p_code IS NOT NULL;`

    rows, err := mysql.DB.Query(sqlStr, userId)
    if err != nil {
        log.Info(sqlStr, err)
    }
    defer rows.Close()

    pCodeArray := make([]string, 0)
    for rows.Next() {
        var pCode string
        err := rows.Scan(&pCode)
        if err != nil {
            log.Info("FindByMenuPermission", err.Error())
        }
        pCodeArray = append(pCodeArray, pCode)
    }
    return &pCodeArray
}

// PermissionMatchCheck 权限匹配校验，使用用户权限和菜单权限对比，用于菜单的加载
func PermissionMatchCheck(pCodeArray *[]string, pCode string) error {
    for _, value := range *pCodeArray {
        if pCode == value {
            return nil
        }
    }
    return fmt.Errorf("Lack of permissions: %s ", pCode)
}

// FindByMenuListAll 获取所有菜单列表
func FindByMenuListAll() *[]map[string]interface{} {
    // 查询一级菜单 level=1
    var items []map[string]interface{}
    items = make([]map[string]interface{}, 0)

    sqlStr := `SELECT id,name,path,p_code,description,create_time,update_time FROM menus WHERE level =1;`
    rows, err := mysql.DB.Query(sqlStr)
    if err != nil {
        log.Info(sqlStr, err)
    }
    defer rows.Close()
    for rows.Next() {
        var m menus
        err := rows.Scan(&m.id, &m.name, &m.path, &m.p_code, &m.description, &m.CreateTime, &m.UpdateTime)
        if err != nil {
            log.Info("Get the first level menu, ", err.Error())
        }
        item := make(map[string]interface{})

		children := FindByMenuLeve2ListAll(m.id)

        item["children"] = children
		item["id"] = m.id
		item["name"] = m.name
		item["path"] = m.path
		item["p_code"] = m.p_code
		item["description "] = m.description
		item["create_time"] = m.CreateTime
		item["update_time"] = m.UpdateTime

        items = append(items, item)
    }
    return &items
}

// FindByMenuLeve2ListAll 根据1级菜单获取二级菜单
func FindByMenuLeve2ListAll(id int64) []map[string]interface{} {

    sqlStr := `SELECT id,name,path,p_code,description,create_time,update_time FROM menus WHERE parent_menu_id = ?;`
    log.Debug(sqlStr)

    rows, err := mysql.DB.Query(sqlStr, id)
    if err != nil {
        log.Info(sqlStr, err)
    }
    defer rows.Close()

    var items []map[string]interface{}
    items = make([]map[string]interface{}, 0)

    for rows.Next() {
        var m menuList
        err := rows.Scan(&m.sub_m_id, &m.sub_m_name, &m.sub_m_path, &m.sub_p_code, &m.sub_desc, &m.CreateTime, &m.UpdateTime)
        if err != nil {
            log.Info("Get the secondary menu, ", err.Error())
        }

        if m.sub_m_id != "" {
            item := make(map[string]interface{})
			item["children"] = make([]map[string]interface{}, 0)
            item["id"] = m.sub_m_id
            item["name"] = m.sub_m_name
            item["path"] = m.sub_m_path
            item["p_code"] = m.sub_p_code
            item["description"] = m.sub_desc
            item["create_time"] = m.CreateTime
            item["update_time"] = m.UpdateTime
            items = append(items, item)
        }
    }
    return items
}

func FindByAuthorizedId(rid int64) *[]int64 {
	sqlStr := `SELECT mid FROM role_menu WHERE rid=?;`
	log.Debug(sqlStr, rid)

	rows, err := mysql.DB.Query(sqlStr, rid)
	if err != nil {
		log.Info(sqlStr, err)
	}
	defer rows.Close()

	items := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Info("Get the secondary menu, ", err.Error())
		}

		items = append(items, id)

	}
	return &items
}
