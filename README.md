#detect port (golang version)

編譯好的版本 linux 限定



#[範例]：

./detect_port 5 10

./detect_port 5(secs) 10(mins)

參數有二個都是必填的

參數一：5(secs) 這個是偵測ip port時的time out設定，內網可以再短一點

參數二：10(mins) 這個是定時器，幾分鐘偵測一次 (新增的功能取消要使用crond來配合使用)



#[檔案說明]：

dp.config 可以設定寄信人，收件人依格式填入 

(寄信人要開低安全性，另外第一次在新的SERVER使用時需要把錯訊息中的網址複製出來登入寄件人的帳號去開啟確認是信任的程式和地點，詳情可以參考gmail其他應用程式使用gmail寄信)

dpconfig.json 要偵測的 IP 及PORT的列表

detect_port.go 源碼

detect_port 編好的執行檔不需要再加載套件

目前程式是多緒執行，當所有檢查都完成後才會結束一個回合。

如果有失敗的結果才會寄信，檢查的結果都正常是不會有任何反應。

現在寄信的使用者是為了測試用申請的google帳號，請替換成自已的帳號。
