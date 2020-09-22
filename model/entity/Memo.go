package model

// Memo はテーブルのモデル
type Memo struct {
	ID     int    `gorm:"primary_key;not null"       json:"id"`
	UserID string `gorm:"type:varchar(200);not null" json:"userId"`
	Memo   string `gorm:"type:varchar(400)"          json:"memo"`
	State  int    `gorm:"not null"                   json:"state"`
}
