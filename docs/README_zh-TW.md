<h1 align="center">
  <img src="./81841388.png" width="400" height='' alt="Noa Himesaka">
  <br>Noa Himesaka<br>  
</h1>
<h3 align="center">
    下載 <a href=“https://book.sfacg.com/”>菠萝包</a> 和 
    <a href=“https://app.hbooker.com/”>刺蝟貓</a> 的小說到本地閱讀.
</h3> 

## **功能**
- 通過 sfacg wechat Api 和 刺蝟貓 Android Api實現下載功能
- 登錄菠萝包帳戶並將cookies保存到本地檔 ```config.json```
- 輸入圖書id或url，並將圖書下載到本地目錄
- 輸入url，並從url提取書籍id下載書籍文本
- 支援從菠萝包和刺蝟貓下載epub電子書
- 按關鍵字搜索書籍，並下載搜尋結果
- [ **警告** ] 新版本圖書快取與舊版本圖書緩存不相容。
- 
## **文件树**
``` 
├─cache
├─cfg
├─cover
├─docs
├─epub
│  └─internal
│      └─storage
│          ├─memory
│          └─osfs
├─save
├─src
│  ├─boluobao
│  ├─hbooker
│  │  └─Encrypt
│  └─https
└─struct
    ├─book_info
    ├─hbooker_structs
    └─sfacg_structs
```
## **示例**

- --app=```<type[sfacg / cat]>```
- --account=```<account>```
- --password=```<password>```
- --download=```<type[bid / url]>```
- --search=```<关键词>```
- --show  < 查看 config.json 文件 >
- 
## **免責聲明**
- 此工具僅用於學習。 請在下載後24小時內將其從計算機中刪除。
- 請尊重版權，請勿自行傳播爬蟲圖書，在任何情況下，作者或版權持有人均不對任何索賠負責
- 損害賠償或其他責任，無論是在合同訴訟中，因軟體或軟體的使用或其他交易而產生的侵權行為或其他行為，作者均不承擔責任。

 
