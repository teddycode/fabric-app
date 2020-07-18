package models

import (
	"github.com/jinzhu/gorm"
)

// farm type
type FarmType struct {
	ID    int    `json:"id" gorm:"primary_key"`
	Value string `json:"value" gorm:"type:varchar(20)"`
}

// new  record
func NewFarmType(t *FarmType) (int, error) {
	err := db.Create(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, err
}

// query types
func QueryFarmTypes() ([]FarmType, error) {
	var types []FarmType
	err := db.Find(&types).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return types, err
	}
	return types, err
}

// del record
func DelFarmType(t *FarmType) (int, error) {
	err := db.Delete(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, err
}

// query user farm type
func CountUserFarmNum(user string) (int64, error) {
	var cnt int64
	err := db.Model(&Transaction{}).Where("type = ? and point = ?","f", user).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return cnt, err
}
