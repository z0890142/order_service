# 啟動方式
於根目錄下執行 make build，進行 docker compose build，在進行 make start
如沒有安裝 make 環境可直接參考 makefile 內的指令。

## makefile 啟動方式
```
make build-image
make start
```

## docker 啟動方式
```
docker-compose build  --no-cache
docker-compose up -d  
```

## 產生文件

於根目錄執行以下指令，執行完後位於根目錄下的 /docs 資料夾裡
```
make generate_doc
```
---

# 資料庫選擇方式
考量到patient與doctor的資料都屬於結構性較強，不易變動的資料類型，所以才用 postgreSQL 進行儲存。

醫囑的多樣性較高且資料格式可能會根據不同得情況有不一樣的格式，考量到未來的擴充性決定採用 MongoDB 。

# 使用方式

打開 WEB 頁面進行登入，分別有兩組帳號其登入資訊如下

```
username: doctor1
password: passowrd1
```

```
username: doctor2
password: passowrd2
```
# Demo



https://github.com/z0890142/order_service/assets/22657048/c4303acf-917b-491e-9f2b-3e1c9eaae3a4

