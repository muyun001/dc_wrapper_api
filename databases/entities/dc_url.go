package entities

import "time"

type DcUrl struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Url        string    `gorm:"type:text" json:"url"`
	Md5        string    `gorm:"type:varchar(32)" json:"md5"`
	Type       int       `gorm:"type:tinyint(1)" json:"type"`
	Status     int       `gorm:"type:tinyint(1)" json:"status"`
	CreateTime time.Time `gorm:"type:timestamp" json:"create_time"`
	UpdateTime time.Time `gorm:"type:timestamp" json:"update_time"`
}

func (DcUrl) TableName() string {
	return "urls_25"
}
