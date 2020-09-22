package model

// User はユーザー情報のモデル
type User struct {
	ID   string `gorm:"primary_key;"type:varchar(200); not null"   json:"id"`
	Name string `gorm:"type:varchar(200);not null" 				json:"name"`
}
