package models

import (
	"mana/src/connections/database/mysql"
)

type menuList struct {
	id          string `json:"id"`
	name        string `json:"name"`
	p_code      string `json:"p_code"`
	level       string `json:"level"`
	sub_m_id    string `json:"sub_m_id"`
	sub_m_name  string `json:"sub_m_name"`
	sub_m_path  string `json:"sub_m_path"`
	sub_m_level string `json:"sub_m_level"`
}

type menuListLevel1 struct {
	id       string              `json:"id"`
	name     string              `json:"name"`
	path     string              `json:"path"`
	leafNode []map[string]string `json:"leafNode"`
}

// selectMenuLevel2
func selectMenuLevel2(id string) []map[string]interface{} {
	sqlStr := `SELECT
					sub.id AS sub_m_id,
					sub.name AS sub_m_name,
					sub.path AS sub_m_path,
					sub.level AS sub_m_level
				FROM
					menus m
				LEFT JOIN menus sub ON (m.id = sub.parent_menu_id)
				WHERE
					m.id = ?;`

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
		err := rows.Scan(&m.sub_m_id, &m.sub_m_name, &m.sub_m_path, &m.sub_m_level)
		if err != nil {
			log.Info(err.Error())
		}

		if m.sub_m_id != "" {
			item := make(map[string]interface{})
			item["id"] = m.sub_m_id
			item["name"] = m.sub_m_name
			item["path"] = m.sub_m_path
			item["leafNode"] = ""
			items = append(items, item)
		}

	}
	return items
}

// SelectMenuList 查询菜单
func SelectMenuList() *map[string]interface{} {
	// 查询一级菜单 level=1
	var items []map[string]interface{}
	items = make([]map[string]interface{}, 0)

	sqlStr_1 := `SELECT id,name,path FROM menus WHERE level =1;`
	rows_1, err := mysql.DB.Query(sqlStr_1)
	if err != nil {
		log.Info(sqlStr_1, err)
	}
	defer rows_1.Close()
	for rows_1.Next() {
		var m menuListLevel1
		err := rows_1.Scan(&m.id, &m.name, &m.path)
		if err != nil {
			log.Info(err.Error())
		}
		log.Info(&m)
		item := make(map[string]interface{})
		item["id"] = m.id
		item["name"] = m.name
		item["path"] = m.path
		leafNode := selectMenuLevel2(m.id)
		item["leafNode"] = leafNode
		items = append(items, item)
		log.Info(&items)
	}

	returns := NewResponse(items, nil)
	return &returns
}
