package models

import (
	"mana/src/connections/database/mysql"
	"strconv"
	"time"
)

type resourceLinks struct {
	ID         string `json:"Id"`         // Id
	USERID     string `json:"userid"`     // 用户id
	LINKNAME   string `json:"linkname"`   // 链接名
	LINKURL    string `json:"linkurl"`    // url
	DESCRIBES  string `json:"describes"`  // 说明
	CATEGORY   string `json:"category"`   // 类别
	ACTIVATE   string `json:"activate"`   // 状态
	CREATETIME string `json:"CreateTime"` // 创建时间
	UPDATETIME string `json:"UpdateTime"` // 更新时间
}

func NewResourceLinks() *resourceLinks {
	return &resourceLinks{
		ACTIVATE:   "0",
		UPDATETIME: strconv.FormatInt(time.Now().Unix(), 10),
		CREATETIME: strconv.FormatInt(time.Now().UnixNano(), 10),
	}
}

// 查询导航链接列表
func FindByResourceLinks(pageSize, page int, activate, category string) (*map[string]interface{}, error) {
	var l resourceLinks
	n, m := (page - 1) * pageSize, pageSize
	var items []map[string]string
	items = make([]map[string]string, 0)
	// SELECT COUNT(*) FROM kube_config
	csqCountStr := `SELECT COUNT(*) FROM resource_links WHERE 1=1`
	sqlStr := `SELECT ID,USERID,LINKNAME,LINKURL,DESCRIBES,CATEGORY,ACTIVATE,CREATETIME,UPDATETIME FROM resource_links WHERE 1=1`
	// 动态sql拼接
	if activate == "" {
		sqlStr += ` AND ACTIVATE = 1`
		csqCountStr += ` AND ACTIVATE = 1`
	} else if activate != "all" {
		sqlStr += ` AND ACTIVATE = ` + activate
		csqCountStr += ` AND ACTIVATE = ` + activate
	}
	if category != "" {
		sqlStr += ` AND CATEGORY = ` + category
		csqCountStr += ` AND CATEGORY = ` + category
	}
	// 分页
	sqlStr += ` ORDER BY ID LIMIT ?, ?;`

	// 查询记录总数,用于PageInfo信息
	var total, pageNum int
	totalRow, err := mysql.DB.Query(csqCountStr)
	if err != nil {
		log.Error("GetKnowledgePointListTotal error", err)
		return nil, err
	}
	for totalRow.Next() {
		err := totalRow.Scan(
			&total,
		)
		if err != nil {
			log.Error("GetKnowledgePointListTotal error", err)
			continue
		}
	}
	// 计算页数
	if total%pageSize == 0 {
		//
		pageNum = total / pageSize
	} else {
		pageNum = total/pageSize + 1
	}
	pageInfo := NewPageInfo(page, pageSize, pageNum, total)

	// 查询记录
	rows, err := mysql.DB.Query(sqlStr, n, m)
	if err != nil {
		log.Error("exec  query failed, err", sqlStr, err)
		return nil, err
	}
	// 3 一定要关闭连接
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&l.ID, &l.USERID, &l.LINKNAME, &l.LINKURL, &l.DESCRIBES,
			&l.CATEGORY, &l.ACTIVATE, &l.CREATETIME, &l.UPDATETIME)
		if err != nil {
			log.Error("Exception reading database result", err)
		}

		createTime, _ := strconv.ParseInt(l.CREATETIME, 10, 64)
		updateTime, _ := strconv.ParseInt(l.UPDATETIME, 10, 64)
		item := make(map[string]string)
		item["id"] = l.ID
		item["name"] = l.LINKNAME
		item["content"] = l.LINKURL
		item["describes"] = l.DESCRIBES
		item["activate"] = l.ACTIVATE
		item["category"] = l.CATEGORY
		item["CreateTime"] = time.Unix(createTime, 0).Format("2006-01-02 15:04:05")
		item["UpdateTime"] = time.Unix(updateTime, 0).Format("2006-01-02 15:04:05")
		items = append(items, item)
	}

	returns := NewResponse(items, pageInfo)
	return &returns, err
}

// 根据id查询导航链接
func FindByResourceLinksTheId(id string) (*map[string]string, error) {
	var l resourceLinks
	sqlStr := `SELECT ID,USERID,LINKNAME,LINKURL,DESCRIBES,CATEGORY,ACTIVATE,CREATETIME,UPDATETIME FROM resource_links WHERE ID=?;`
	var row = mysql.DB.QueryRow(sqlStr, id)
	err := row.Scan(&l.ID, &l.USERID, &l.LINKNAME, &l.LINKURL, &l.DESCRIBES, &l.CATEGORY, &l.ACTIVATE,
		&l.CREATETIME, &l.UPDATETIME)

	if err != nil {
		log.Error("exec  query failed, ", sqlStr, "ID=", id, err)
		return nil, err
	}

	createTime, _ := strconv.ParseInt(l.CREATETIME, 10, 64)
	updateTime, _ := strconv.ParseInt(l.UPDATETIME, 10, 64)
	item := make(map[string]string)
	item["Id"] = l.ID
	item["name"] = l.LINKNAME
	item["content"] = l.LINKURL
	item["describes"] = l.DESCRIBES
	item["activate"] = l.ACTIVATE
	item["category"] = l.CATEGORY
	item["CreateTime"] = time.Unix(createTime, 0).Format("2006-01-02 15:04:05")
	item["UpdateTime"] = time.Unix(updateTime, 0).Format("2006-01-02 15:04:05")

	returns := item
	return &returns, err
}

// 插入
func InsertResourceLink(l *resourceLinks) (string, error) {
	// 插入用户表信息
	sqlStr1 := `INSERT INTO resource_links (USERID, LINKNAME, LINKURL, DESCRIBES, CATEGORY, ACTIVATE, CREATETIME, UPDATETIME)
                         VALUES (?,?,?,?,?,?,?,?);`
	ret, err := mysql.DB.Exec(sqlStr1, l.USERID, l.LINKNAME, l.LINKURL, l.DESCRIBES,
		l.CATEGORY, l.ACTIVATE, l.CREATETIME, l.UPDATETIME)

	if err != nil {
		log.Error("insert failed,", err)
		return "", err
	}

	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		log.Error("get last Insert ID failed, err:", err)
		return "", err
	}

	return strconv.FormatInt(theID, 10), err

}

// 更新
func UpdateResourceLinkToDb(l *resourceLinks) (string, error) {
	// 插入用户表信息
	sqlStr1 := `UPDATE resource_links SET LINKNAME=?,LINKURL=?,DESCRIBES=?,CATEGORY=?,ACTIVATE=?,UPDATETIME=? WHERE ID=?`
	ret, err := mysql.DB.Exec(sqlStr1, l.LINKNAME, l.LINKURL, l.DESCRIBES, l.CATEGORY, l.ACTIVATE, l.UPDATETIME, l.ID)

	if err != nil {
		log.Error("update failed,", err)
		return "", err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		log.Error("get RowsAffected failed, ", err)
		return "", err
	}

	return strconv.FormatInt(n, 10), err

}

// DeleteLink 删除
func DeleteLink(id string) (string, error) {
	sqlStr1 := `DELETE FROM resource_links WHERE ID=?`
	ret, err := mysql.DB.Exec(sqlStr1, id)
	if err != nil {
		log.Error("delete failed, err ", err)
		return "", err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		log.Error("get RowsAffected failed, err: ", err)
		return "", err
	}
	log.Info("Delete navigation links success, affected rows: ", n)
	return strconv.FormatInt(n, 10), nil
}
