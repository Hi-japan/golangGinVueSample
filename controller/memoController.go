package controller

import (
	// 文字列と基本データ型の変換パッケージ

	"net/http"
	strconv "strconv"
	"strings"

	// Gin
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	oauthapi "google.golang.org/api/oauth2/v2"

	// エンティティ(データベースのテーブルの行に対応)
	entity "git.heroku.com/golangginvuesample2.git/model/entity"

	// DBアクセス用モジュール
	db "git.heroku.com/golangginvuesample2.git/model/db"
)

var conf = &oauth2.Config{
	ClientID:     "your client ID",
	ClientSecret: "your client secret key",
	Scopes:       []string{oauthapi.UserinfoProfileScope},
	Endpoint:     google.Endpoint,

	// ローカル
	//RedirectURL:  "http://localhost:8080/auth/callback/google",
	// heroku
	RedirectURL: "https://golangginvuesample2.herokuapp.com/auth/callback/google",
}

// 商品の購入状態を定義
const (
	// NotCompleted は 未完了
	NotCompleted = 0

	// Completed は 完了
	Completed = 1
)

// FetchAllMemos は ユーザーのメモをすべて取得する
func FetchAllMemos(c *gin.Context) {

	userID := c.Query("userId")

	resultProducts := db.FindAllMemos(userID)

	// URLへのアクセスに対してJSONを返す
	c.JSON(200, resultProducts)
}

// FindMemo は 指定したIDのメモを取得する
func FindMemo(c *gin.Context) {
	memoIDStr := c.Query("productID")

	memoID, _ := strconv.Atoi(memoIDStr)

	result := db.FindMemo(memoID)

	// URLへのアクセスに対してJSONを返す
	c.JSON(200, result)
}

// AddMemo は メモをDBへ登録する
func AddMemo(c *gin.Context) {
	memoStr := c.PostForm("productMemo")
	userID := c.PostForm("userId")

	var memo = entity.Memo{
		Memo:   memoStr,
		State:  NotCompleted,
		UserID: userID,
	}

	db.InsertMemo(&memo)
}

// ChangeStateMemo は 商品情報の状態を変更する
func ChangeStateMemo(c *gin.Context) {
	reqMemoID := c.PostForm("productID")
	reqMemoState := c.PostForm("productState")

	memoID, _ := strconv.Atoi(reqMemoID)
	memoState, _ := strconv.Atoi(reqMemoState)
	changeState := NotCompleted

	// 商品状態が未購入の場合
	if memoState == NotCompleted {
		changeState = Completed
	} else {
		changeState = NotCompleted
	}

	db.UpdateStateMemo(memoID, changeState)
}

// DeleteMemo は メモをDBから削除する
func DeleteMemo(c *gin.Context) {
	memoIDStr := c.PostForm("productID")

	memoID, _ := strconv.Atoi(memoIDStr)

	db.DeleteMemo(memoID)
}

// Login はログイン
func Login(c *gin.Context) {
	// urlの分解。区切り文字は "/"
	segs := strings.Split(c.FullPath(), "/")
	action := segs[2] //処理内容。login あるいは callback.

	switch action {

	/* loginページから遷移した場合 */
	case "login":
		authUrl := conf.AuthCodeURL("yourStateUUID", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusSeeOther, authUrl)

	/* プロバイダの元での認証を終え、callbackされた場合 */
	case "callback":

		code := c.Request.URL.Query()["code"]
		if code == nil || len(code) == 0 {
			//fmt.Fprintln("Invalid Parameter")
		}

		ctx := context.Background()
		tok, err := conf.Exchange(ctx, code[0])
		if err != nil {
			//fmt.Fprintf(w, "OAuth Error:%v", err)
		}
		//APIクライアント作成
		client := conf.Client(ctx, tok)
		//Userinfo APIをGet
		svr, err := oauthapi.New(client)
		ui, err := svr.Userinfo.Get().Do()

		var user = entity.User{
			ID:   ui.Id,
			Name: ui.Name,
		}

		db.UpsetUser(&user)

		// 抜き出した情報をクッキーに保存
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "userId",
			Value: ui.Id,
			Path:  "/memoapp",
		})

		// ログイン後の画面に遷移する
		c.Redirect(http.StatusSeeOther, "/memoapp")
	}
}

// GetUserName はユーザーIDに対応するユーザー名を取得する
func GetUserName(c *gin.Context) {

	userID := c.Query("userId")
	userName := db.FindUserName(userID)

	// URLへのアクセスに対してJSONを返す
	c.JSON(200, userName)
}
