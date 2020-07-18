package transh

import (
	"github.com/fabric-app/models"
	"github.com/jinzhu/gorm"
	"os"
	"path"
)

const PICTURE_PATH = "/root/Desktop/DataSources/lzawt/video_pics"

// user table structure
type Picture struct {
	ID   int    `json:"id"  gorm:"PRIMARY_KEY,AUTO_INCREMENT"`
	Name string `json:"name"  gorm:"type:varchar(20)"`
	Hash string `json:"hash" gorm:"type:varchar(64)"`
	Size int    `json:"size"`
	Date string `json:"date" gorm:"type:date"`
	Type string `json:"type" gorm:"type:varchar(10)"`
}

// create pic info
func NewPic(pic *Picture) (int, error) {
	err := models.db.Create(pic).Error
	if err != nil {
		return 0, err
	}
	return pic.ID, err
}

// update user info
func UpdatePicHash(name, hash string) (int, error) {
	var oldPic Picture
	err := models.db.Where("name = ?", name).First(&oldPic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	oldPic.Hash = hash
	err = models.db.Save(oldPic).Error
	if err != nil {
		return 0, nil
	}
	return oldPic.ID, nil
}

// get pictures file
func GetPicFile(point, date, name string) (*os.File, error) {
	path := path.Join(PICTURE_PATH, point, date, name+".jpg")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
