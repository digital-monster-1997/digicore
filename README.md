# Digicore

## 簡介
『 數碼核 』 是一款以治理為導向的微服務框架核心 ，擅長讓微服務中的狀態可視化

## 支援特性
1. 註冊服務
2. 發現服務
3. 連線方式
    * grpc (服務端以及客戶端，參數可自由設定)
    * http1.1/http2
    * websocket
4. 資料庫支援的有
    * orm
    * redis
    * mongodb
5. 熱加載 config
6. corn job 任務
7. queue(rocketMQ)
8. 日誌
9. 異步任務池
10. 事件總線（這個不清楚，再研究
11. Prometheus/pprof 監控
12. 優雅重啟
13. 工具類 command tool
14. 全域變量註冊
15. 在線應用程式附載均衡
16. RPC 健康檢查
17. 接入授權
18. ghz 壓力測試工具
19. 在線服務限流
20. 在線服務熔斷，異常接入 sentry
21. watch 服務在線狀態
22. cache 多級暫存


## 文件
請看 [文件](https://digicore.30cm.net) 來檢視詳細文件

## 快速開始
```golang
func main(){
	fmt.Pringln("Quick Start")
}
```

## 更多範例請看
