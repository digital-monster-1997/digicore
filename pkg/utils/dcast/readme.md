# Cast
> 在 Golang 當中安全且簡單的從一個類型，轉型成另外一個

## 改寫
改寫自 [spf13](https://github.com/spf13/cast) 

## 說明
 Cast 提供了簡單的函數來輕鬆的將數字轉為字串，將 interface 轉換成布林
 當可以進行明顯的轉換時，Cast 會自動地執行操作，他不會用猜的。
 Ex: 當 int 的字串 例如：'8'，使用的時候只能將字串轉為 int。
 
## 為啥使用 Cast ?
1. Golang 在處理動態的資料的時候，經常要把資料從一種類型轉換成另一種
2. 從 Yaml、Toml 或者 Json 等其他缺乏完整資料型態的格式中拿到資訊

## 用法
這裡提供了兩種方法
1. To 方法，這種方法將會返回所需的類型，如果提供的輸入，不會轉會城該類型，將返回該類型的 0 或者 nil 值
2. To_E 方法，會返回相同結果，但是加上一個額外的錯誤，告訴使用者是否成功轉換，這樣就可以區分，是因為輸入零值，而返回的或是，錯誤而返回的零值

ToString
```golang
dcast.ToString("mayonegg")         // "mayonegg"
dcast.ToString(8)                  // "8"
dcast.ToString(8.31)               // "8.31"
dcast.ToString([]byte("one time")) // "one time"
dcast.ToString(nil)                // ""

var foo interface{} = "one more time"
dcast.ToString(foo)                // "one more time"
```

ToInt
```golang
dcast.ToInt(8)                  // 8
dcast.ToInt(8.31)               // 8
dcast.ToInt("8")                // 8
dcast.ToInt(true)               // 1
dcast.ToInt(false)              // 0

var eight interface{} = 8
dcast.ToInt(eight)              // 8
dcast.ToInt(nil)                // 0              // "one more time"
```