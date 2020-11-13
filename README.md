# slack-word-bot
エンジニアに響く名言(迷言)をslackに投稿する。また、ボタンで「さいこー」「いいね」「ふつう」「びみょー」を評価できるようにする。

名言はMongoDBに格納しておく

## 設定ファイル
`config.json`として作成する。

SlackへのPostのための設定、名言が格納されているMongoDBについての設定を入れる。

設定：
```
{
    "slack": {
        "webHookUrl": "your_webhook_url",
        "oauthAccessToken": "your_access_token",
        "channelID": "your_channnel_id"
    },
    "mongo": {
        "host": "mongohost.example",
        "port": "27017",
        "dbName": "word_db",
        "collectionName": "word_col"
    }
}
```

## 実行
```
$ go mod install
$ go run main.go
```