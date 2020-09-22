new Vue({
    // 「el」プロパティーで、Vueの表示を反映する場所＝HTML要素のセレクター（id）を定義
    el: '#app',

    // data オブジェクトのプロパティの値を変更すると、ビューが反応し、新しい値に一致するように更新
    data: {
        // 利用者名
        userName: '',
        // メモ情報
        products: [],
        // メモ
        productMemo: '',
        // メモ情報の状態
        current: -1,
        // メモ情報の状態一覧
        options: [
            { value: -1, label: 'すべて' },
            { value:  0, label: '未完了' },
            { value:  1, label: '完了' }
        ],
        // true：入力済・false：未入力
        isEntered: false
    },

    // 算出プロパティ
    computed: {
        // メモの状態一覧を表示する
        labels() {
            return this.options.reduce(function (a, b) {
                return Object.assign(a, { [b.value]: b.label })
            }, {})
        },
        // 表示対象のメモを返却する
        computedMemos() {
          return this.products.filter(function (el) {
            var option = this.current < 0 ? true : this.current === el.state
            return option
          }, this)
        },
        // 入力チェック
        validate() {
            var isEnteredProductMemo = 0 < this.productMemo.length
            this.isEntered = isEnteredProductMemo
            return isEnteredProductMemo
        },

    },

    // インスタンス作成時の処理
    created: function() {
        this.doGetUserName()
        this.doFetchAllMemos()
    },

    // メソッド定義
    methods: {

        // ユーザーID取得
        getUserId() {
            const cookies = document.cookie
            const cookiesArray = cookies.split('; '); // ;で分割し配列に

            var id = null
            for (var c of cookiesArray) {
                var cArray = c.split('=');
                if (cArray[0] == 'userId') {
                    console.log(cArray);  // [key,value] 
                    id = cArray[1]
                }
            }
            return id
        },
        // ユーザー名を取得する
        doGetUserName() {

            axios.get('/getUserName', {
                params: {
                    userId: this.getUserId()
                }
            })
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('レスポンスエラー')
                    } else {
                        var resultProducts = response.data

                        // サーバから取得したメモ情報をdataに設定する
                        this.userName = resultProducts
                    }
                })
        },

        // 全てのメモを取得する
        doFetchAllMemos() {

            axios.get('/fetchAllMemos', {
                params: {
                    userId: this.getUserId()
                }
            })
            .then(response => {
                if (response.status != 200) {
                    throw new Error('レスポンスエラー')
                } else {
                    var resultProducts = response.data

                    // サーバから取得したメモ情報をdataに設定する
                    this.products = resultProducts
                }
            })
        },
        // １つのメモを取得する
        doFetchMemo(product) {
            axios.get('/fetchMemo', {
                params: {
                    productID: product.id
                }
            })
            .then(response => {
                if (response.status != 200) {
                    throw new Error('レスポンスエラー')
                } else {
                    var resultProduct = response.data

                    // 選択されたメモ情報のインデックスを取得する
                    var index = this.products.indexOf(product)

                    // spliceを使うとdataプロパティの配列の要素をリアクティブに変更できる
                    this.products.splice(index, 1, resultProduct[0])
                }
            })
        },
        // メモを登録する
        doAddMemo() {
            // サーバへ送信するパラメータ
            const params = new URLSearchParams();
            params.append('productMemo', this.productMemo)
            params.append('userId', this.getUserId())

            axios.post('/addMemo', params)
            .then(response => {
                if (response.status != 200) {
                    throw new Error('レスポンスエラー')
                } else {
                    // メモ情報を取得する
                    this.doFetchAllMemos()

                    // 入力値を初期化する
                    this.initInputValue()
                }
            })
        },
        // メモの状態を変更する
        doChangeMemoState(product) {
            // サーバへ送信するパラメータ
            const params = new URLSearchParams();
            params.append('productID', product.id)
            params.append('productState', product.state)

            axios.post('/changeStateMemo', params)
            .then(response => {
                if (response.status != 200) {
                    throw new Error('レスポンスエラー')
                } else {
                    // メモ情報を取得する
                    this.doFetchMemo(product)
                }
            })
        },
        // メモを削除する
        doDeleteMemo(product) {
            // サーバへ送信するパラメータ
            const params = new URLSearchParams();
            params.append('productID', product.id)

            axios.post('/deleteMemo', params)
            .then(response => {
                if (response.status != 200) {
                    throw new Error('レスポンスエラー')
                } else {
                    // メモ情報を取得する
                    this.doFetchAllMemos()
                }
            })
        },
        // 入力値を初期化する
        initInputValue() {
            this.current = -1
            this.productMemo = ''
        }
    }
})