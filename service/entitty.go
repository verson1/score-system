package service

import "score-system/model"

type OrderFiled int

const (
	Subject OrderFiled = 1
	Score   OrderFiled = 2
	Class   OrderFiled = 3
)

type QueryCondition struct {
	Start   int    `form:"start"`
	Limit   int    `form:"limit"`
	Name    string `form:"name"`
	Class   int    `form:"class"`
	Subject string `form:"subject"`
	Score   int    `form:"score"`
	// 排序类别：1-asc 升序，2-desc 降序
	OrderType int `form:"order_type"`
	// 排序字段：score-分数 subject-科目 class-班级
	OrderFiled OrderFiled `form:"order_filed"`
}

type QueryRes struct {
	// 查询到的记录数量
	Count        int64                `json:"count"`
	Achievements []*model.Achievement `json:"achievements"`
}
