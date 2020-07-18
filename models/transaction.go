package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

const (
	HOUR_TIMESTAMP = 60 * 60
	DAY_TIMESTAMP  = 24 * HOUR_TIMESTAMP
	WEEK_TIMESTAMP = 7 * DAY_TIMESTAMP
	MOTH_TIMESTAMP = 30 * DAY_TIMESTAMP
)

// transaction
type Transaction struct {
	ID        int       `json:"id" gorm:"PRIMARY_KEY"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type" gorm:"type:varchar(2)"`
	Hash      string    `json:"hash" gorm:"type:varchar(64)"`
	Point     string    `json:"point" gorm:"type:varchar(20)"`
}

type PointCnter struct {
	Point string `json:"point"`
	Cnt   int    `json:"cnt"`
}

type TxCnter struct {
	Unit  string `json:"unit"`
	Value int64  `json:"value"`
}

func (v PointCnter) TableName() string {
	return "transactions"
}

func (v TxCnter) TableName() string {
	return "transactions"
}

// new tx record
func NewTx(tx *Transaction) (int, error) {
	err := db.Create(tx).Error
	if err != nil {
		return 0, err
	}
	return tx.ID, err
}

// count tx number by time period
//func countTxNumByTimePeriod(s, e int64) (int64, error) {
//	var count int64
//	err := db.Model(&Transaction{}).Where("timestamp >= ? and timestamp < ?", s, e).Count(&count).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		return 0, err
//	}
//	return count, err
//}

// count tx number by day  with last 12h
func CountTxNumByDay() ([]TxCnter, error) {
	var sql string
	txCnts := []TxCnter{}
	var t = time.Now() // 从当前时间开始
	d, _ := time.ParseDuration("1h")
	nt := t.Truncate(d).Add(d) // 截取hour
	d, _ = time.ParseDuration("-1h")
	for i := 0; i < 24; i++ {
		lt := nt.Add(d)
		sql += fmt.Sprintf("SELECT '%d' AS unit, COUNT(id) AS value FROM transactions WHERE TIMESTAMP <= '%s' and TIMESTAMP > '%s' ",
			lt.Hour(), nt.Format("2006-01-02 15:04:05"), lt.Format("2006-01-02 15:04:05"))
		if i != 23 {
			sql += " UNION "
		}
		nt = lt
	}
	//fmt.Println(sql)
	err := db.Raw(sql).Find(&txCnts).Error
	if err != nil {
		return txCnts, err
	}
	return txCnts, err
}

// count tx number by week
func CountTxNumByWeek() ([]TxCnter, error) {
	var sql string
	txCnts := []TxCnter{}
	var t = time.Now() // 从当前时间开始
	d, _ := time.ParseDuration("24h")
	nt := t.Truncate(d).Add(d) // 截取day
	d, _ = time.ParseDuration("-24h")
	for i := 0; i < 7; i++ {
		lt := nt.Add(d)
		sql += fmt.Sprintf("SELECT '%s' AS unit, COUNT(id) AS value FROM transactions WHERE TIMESTAMP <= '%s' and TIMESTAMP > '%s' ",
			lt.Weekday().String(), nt.Format("2006-01-02 15:04:05"), lt.Format("2006-01-02 15:04:05"))
		if i != 6 {
			sql += " UNION "
		}
		nt = lt
	}
	//fmt.Println(sql)
	err := db.Raw(sql).Find(&txCnts).Error
	if err != nil {
		return txCnts, err
	}
	return txCnts, err
}

// count tx number by moth
func CountTxNumByMoth() ([]TxCnter, error) {
	var sql string
	txCnts := []TxCnter{}
	var t = time.Now() // 从当前时间开始
	d, _ := time.ParseDuration("24h")
	nt := t.Truncate(d).Add(d) // 截取day
	d, _ = time.ParseDuration("-24h")
	for i := 0; i < 30; i++ {
		lt := nt.Add(d)
		sql += fmt.Sprintf("SELECT '%d' AS unit, COUNT(id) AS value FROM transactions WHERE TIMESTAMP <= '%s' and TIMESTAMP > '%s' ",
			nt.Day(), nt.Format("2006-01-02 15:04:05"), lt.Format("2006-01-02 15:04:05"))
		if i != 29 {
			sql += " UNION "
		}
		nt = lt
	}
//	fmt.Println(sql)
	err := db.Raw(sql).Find(&txCnts).Error
	if err != nil {
		return txCnts, err
	}
	return txCnts, err
}

// count tx number by moth
func CountTxNumByYear() ([]TxCnter, error) {
	var sql string
	txCnts := []TxCnter{}
	var mStr = time.Now().Format("01") // 从当前月份开始
	m, _ := strconv.ParseInt(mStr, 10, 64)
	var yStr = time.Now().Format("2006") // 从当前年份开始
	y, _ := strconv.ParseInt(yStr, 10, 64)
	for i := 0; i < 12; i++ {
		sql += fmt.Sprintf("select '%d' as unit, COUNT(id) as value from transactions where year(timestamp) = '%d' and month(timestamp) = '%d'", m, y, m)
		if i != 11 {
			sql += " UNION "
		}
		m--
		if m < 0 {
			m = 12
			y--
		}
	}
//	fmt.Println(sql)
	err := db.Raw(sql).Find(&txCnts).Error
	if err != nil {
		return txCnts, err
	}
	return txCnts, err
}

// count tx number by point
func CountTxNums() (int64, error) {
	var count int64
	err := db.Model(&Transaction{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

// get all points
func GetAllPoints() ([]PointCnter, error) {
	var points []PointCnter
	err := db.Model(&Transaction{}).Select("point, count(id) as cnt").Where("type = ? or type = ?", "s", "p").Group("point").Order("cnt DESC").Find(&points).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return points, err
}
