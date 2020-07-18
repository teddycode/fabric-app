package models

import (
	"bufio"
	"github.com/fabric-app/pkg/util/hash"
	"github.com/fabric-app/pkg/util/rand"
	"github.com/jinzhu/gorm"
	"os"
	"path"
	"time"
)

const HEADER_IMAGE_PATH = "./test/header/images/"

// user table structure
type User struct {
	Model
	UserName string `json:"user_name" gorm:"type:varchar(20);unique;not null"`
	Identity string `json:"identity" gorm:"type:varchar(20)"`
	Password string `json:"password" gorm:"type:varchar(20)"`
	Phone    string `json:"phone" gorm:"type:varchar(12)"`
	Email    string `json:"email" gorm:"type:varchar(20)"`
	Role     int    `json:"role" gorm:"type:int;default 1"`
	CaSecure string `json:"ca_secure" gorm:"type:varchar(20)"`
	Secret   string `json:"secret" gorm:"type:varchar(20)"`
	Address  string `json:"address" gorm:"type:varchar(50)"`
	Header   string `json:"header" gorm:"type:varchar(10)"`
}

// create user info
func NewUser(user *User) (int, error) {
	err := db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// del user
func DelUser(user *User) (int, error) {
	err := db.Delete(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// find user by id
func FindUserById(id int) (User, error) {
	var user User
	err := db.First(&user, id).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

// by name
func FindUserByName(name string) (User, error) {
	var user User
	err := db.Where("user_name = ?", name).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

// find user by email
func FindUserByEmail(e string) (User, error) {
	var user User
	err := db.Where("email = ?", e).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

// update user info
func UpdateUserInfo(newUser *User) (int, error) {
	var oldUser User
	err := db.First(&oldUser, newUser.ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	oldUser.Email = newUser.Email
	oldUser.Address = newUser.Address
	oldUser.Phone = newUser.Phone
	oldUser.ModifiedOn = int(time.Now().Unix())
	err = db.Save(&oldUser).Error
	if err != nil {
		return 0, nil
	}
	return oldUser.ID, nil
}

// update user secret
func UpdateUserSecret(user *User) (int, error) {
	var secretString string
	for {
		secretString = rand.RandStringBytesMaskImprSrcUnsafe(5)
		if user.Secret != secretString {
			break
		}
	}
	db.First(user)
	user.Secret = secretString
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// update user header
func UpdateUserheader(name, header string) (int, error) {
	var user User
	db.Where("user_name = ?", name).First(&user)
	user.Header = header
	err := db.Save(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// update user password
func UpdateUserNewPassword(user *User, newPassword string) (int, error) {
	var secretString string
	for {
		secretString = rand.RandStringBytesMaskImprSrcUnsafe(5)
		if user.Secret != secretString {
			break
		}
	}
	db.First(user)
	user.Secret = secretString
	user.Password = hash.EncodeMD5(newPassword)
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// store user header
func SaveUserHeader(username string, data []byte) (int, error) {
	path := path.Join(HEADER_IMAGE_PATH, username+".jpg")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	writer := bufio.NewWriter(file)
	count, err := writer.Write(data)
	if err != nil {
		return 0, err
	}
	writer.Flush()
	return count, err
}

// get user header file name
func GetUserHeader(username string) (string, error) {
	var user User
	err := db.Select("header").Where("user_name = ?", username).Find(&user).Error
	if err != nil {
		return "", err
	}
	return user.Header, err
}

//
//// get user header
//func GetUserHeader(username string) (*os.File, error) {
//	path1 := path.Join(HEADER_IMAGE_PATH, username+".jpg")
//	file, err := os.OpenFile(path1, os.O_RDONLY|os.O_TRUNC, 0666)
//	if err != nil || file == nil {
//		path1 = path.Join(HEADER_IMAGE_PATH, "default.jpg")
//		file, err = os.OpenFile(path1, os.O_RDONLY|os.O_TRUNC, 0666)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return file, err
//}
