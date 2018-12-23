## KARI

ConoHa API V1 オブジェクトストレージ ファイル削除ツールです

goroutineを使って超高速にオブジェクトファイルを削除します

旧ConoHa(2015年5月17日以前のConoHa)でオブジェクトストレージを使ってた人向けです

## 設定ファイル

confフォルダにある.jsonファイルをリネームして使ってください

```
mv ./conf/conoha_api_v1_key_sample.json ./conf/conoha_api_v1_key.json
```


### 設定ファイルサンプル

自分のアカウントでConoHaにログインして　https://cp.conoha.jp/Account/API/ の情報を転記してください

```
{
  "auth_url" : "https://ident-r1nd1001.cnode.jp/v2.0/tokens",
  "tenantName" : "11111",
  "username" : "11111",
  "password" : "APIのパスワード",
  "endPoint" : "https://objectstore-r1nd1001.cnode.jp/v1/XXXXXX"
}
```

## ビルド

```
go build kasumi.go
```


## 実行

```
 ./kasumi /container_name_path
```

例えばオブジェクトストレージのルートに「image」というコンテナがある場合

```
 ./kasumi /image
```

と実行すると中身のオブジェクトファイルを削除します
