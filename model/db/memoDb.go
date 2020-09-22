package db

import (
	// フォーマットI/O
	"fmt"

	// Go言語のORM
	"github.com/jinzhu/gorm"

	// エンティティ(データベースのテーブルの行に対応)
	entity "git.heroku.com/golangginvuesample2.git/model/entity"
)

// DB接続する
func open() *gorm.DB {
	DBMS := "postgres"

	// ローカル
	// USER := "postgres"
	// PASS := "admin"
	// PROTOCOL := "localhost:5432"
	// DBNAME := "postgres"
	// CONNECT := "postgres://" + USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?sslmode=disable"

	// heroku
	CONNECT := "postgres://heroku postges path" + "?sslmode=disable"

	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}

	// DBエンジンを「InnoDB」に設定
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// 詳細なログを表示
	db.LogMode(true)

	// 登録するテーブル名を単数形にする（デフォルトは複数形）
	db.SingularTable(true)

	// マイグレーション（テーブルが無い時は自動生成）
	db.AutoMigrate(&entity.Memo{})
	db.AutoMigrate(&entity.User{})

	fmt.Println("db connected: ", &db)
	return db
}

// FindAllMemos は メモテーブルのレコードを全件取得する
func FindAllMemos(userID string) []entity.Memo {
	memos := []entity.Memo{}

	db := open()
	// select
	db.Order("ID asc").Where("User_ID = ?", userID).Find(&memos)

	// defer 関数がreturnする時に実行される
	defer db.Close()

	return memos
}

// FindMemo は メモテーブルのレコードを１件取得する
func FindMemo(productID int) []entity.Memo {
	memo := []entity.Memo{}

	db := open()
	// select
	db.First(&memo, productID)
	defer db.Close()

	return memo
}

// InsertMemo は メモテーブルにレコードを追加する
func InsertMemo(registerMemo *entity.Memo) {
	db := open()
	// insert
	db.Create(&registerMemo)
	defer db.Close()
}

// UpdateStateMemo は メモテーブルの指定したレコードの状態を変更する
func UpdateStateMemo(memoID int, memoState int) {
	memo := []entity.Memo{}

	db := open()
	// update
	db.Model(&memo).Where("ID = ?", memoID).Update("State", memo)
	defer db.Close()
}

// DeleteMemo は メモテーブルの指定したレコードを削除する
func DeleteMemo(memoID int) {
	memo := []entity.Memo{}

	db := open()
	// delete
	db.Delete(&memo, memoID)
	defer db.Close()
}

// UpsetUser は ユーザー情報の追加
func UpsetUser(registerUser *entity.User) {
	db := open()
	// insert
	db.Create(&registerUser)
	defer db.Close()
}

// FindUserName は ユーザーIDに対応するユーザー名を取得する
func FindUserName(userID string) string {

	users := []entity.User{}

	db := open()
	// select
	db.Where("id = ?", userID).Find(&users)

	// defer 関数がreturnする時に実行される
	defer db.Close()

	user := users[0].Name

	return user
}
