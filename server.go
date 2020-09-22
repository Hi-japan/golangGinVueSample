package main

import (
	// ロギングを行うパッケージ
	"log"
	"os"

	// HTTPを扱うパッケージ
	"net/http"

	// Gin
	"github.com/gin-gonic/gin"

	// Postgres用ドライバ
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// コントローラー
	controller "git.heroku.com/golangginvuesample2.git/controller"
)

func main() {
	// サーバーを起動する
	serve()
}

func serve() {
	router := gin.Default()

	// 静的ファイルのパスを指定
	router.Static("/view", "./view")

	// ルーティングの設定
	router.StaticFS("/memoapp", http.Dir("./view/static"))
	router.GET("/fetchAllMemos", controller.FetchAllMemos)
	router.GET("/fetchMemo", controller.FindMemo)
	router.POST("/addMemo", controller.AddMemo)
	router.POST("/changeStateMemo", controller.ChangeStateMemo)
	router.POST("/deleteMemo", controller.DeleteMemo)
	router.GET("/getUserName", controller.GetUserName)
	router.GET("/auth/login/google", controller.Login)
	router.GET("/auth/callback/google", controller.Login)

	// ローカル
	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatal("Server Run Failed.: ", err)
	// }

	// heroku
	port := os.Getenv("PORT")
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server Run Failed.: ", err)
	}
}
