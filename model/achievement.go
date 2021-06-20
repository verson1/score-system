package model

type Achievement struct {
	Id      int64 `gorm:"colum:id,auto"`
	Name    string
	Class   int32
	Score   int32
	Subject string
}

func (*Achievement) TableName() string {
	return "achievement"
}
