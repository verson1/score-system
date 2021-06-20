package dao

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"score-system/model"
)

type Dao struct {
	orm *gorm.DB
}

type OrderFiled string
type OrderType string

const (
	ByClass   OrderFiled = "class"
	ByScore   OrderFiled = "score"
	BySubject OrderFiled = "subject"
)
const (
	ByDESC OrderType = "DESC"
	ByASC  OrderType = "ASC"
)

type Condition struct {
	Start      int
	Limit      int
	Name       string
	Class      int
	Score      int
	Subject    string
	OrderFiled OrderFiled
	OrderType  OrderType
}

func NewDao(orm *gorm.DB) *Dao {
	return &Dao{
		orm: orm,
	}
}

func (d *Dao) BatchInsertAchievement(achievements []*model.Achievement) error {
	err := d.orm.CreateInBatches(achievements, len(achievements)).Error
	if err != nil {
		return errors.Wrap(err, "batch insert fail")
	}
	return nil
}

func (d *Dao) GetAchievementOrderByFiled(condition *Condition) (int64, []*model.Achievement, error) {
	if condition == nil {
		return 0, nil, errors.New("query condition not allow is nil")
	}

	db := d.orm.Model(&model.Achievement{})
	if condition.Name != "" {
		db = db.Where("name = ?", condition.Name)
	}
	if condition.Class != 0 {
		db = db.Where("class = ?", condition.Class)
	}
	if condition.Subject != "" {
		db = db.Where("subject = ?", condition.Subject)
	}

	var res []*model.Achievement
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	if count == 0 || count < int64(condition.Start) {
		return 0, res, err
	}

	if condition.OrderType != "" && condition.OrderFiled != "" {
		db = db.Order(fmt.Sprintf("%s %s", condition.OrderFiled, condition.OrderType))
	}

	// 查询条数如果没设置，默认查10条
	if condition.Limit == 0 {
		condition.Limit = 10
	}
	db.Offset(condition.Start).Limit(condition.Limit)

	err = db.Scan(&res).Error
	if err != nil {
		return 0, res, err
	}
	return count, res, nil
}
