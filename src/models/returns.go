package models

import "time"

// 最外层结构
func NewReturns(result interface{}, msg interface{}) map[string]interface{} {
	returns := make(map[string]interface{})

	returns["response"] = result
	returns["metaInfo"] = msg

	return returns
}

// response 结构，
// 包含 items 和 pageInfo 两部分信息，
// 如果没有 pageInfo 可传递 nil。
func NewResponse(items, pageInfo interface{}) map[string]interface{} {
	returns := make(map[string]interface{})

	returns["items"] = items
	if pageInfo != nil {
		returns["pageInfo"] = pageInfo
	}
	return returns
}

// response 中的页码信息,可选
func NewPageInfo(page, pageSize, pageNum, total int) map[string]int {
	pageInfo := make(map[string]int)

	pageInfo["page"] = page
	pageInfo["pageSize"] = pageSize
	pageInfo["pageNum"] = pageNum
	pageInfo["total"] = total
	return pageInfo
}

// 返回体的统一消息结构
type resMessage struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
}

func NewResMessage(code string, msg string) *resMessage {
	return &resMessage{
		Code:        code,
		Msg:         msg,
		RequestTime: time.Now().UnixNano() / 1e6,
	}
}
