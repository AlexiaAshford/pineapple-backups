<h1 align="center">
  <img src="./81841388.png" width="400" height='' alt="Noa Himesaka">
  <br>Noa Himesaka<br>  
</h1>
<h3 align="center">
    下载 <a href="https://book.sfacg.com/">菠萝包</a> 和 
    <a href="https://app.hbooker.com/">刺猬猫</a> 的小说到本地阅读. 

</h3> 

## **功能**

- 通过 sfacg wechat Api 和 刺猬猫 Android Api实现下载功能
- 登录菠萝包帐户并将cookies保存到本地文件 ```config.json```
- 输入图书id或url，并将图书下载到本地目录
- 输入url，并从url提取书籍id下载书籍文本
- 支持从菠萝包和刺猬猫下载epub电子书
- 按关键字搜索书籍，并下载搜索结果
- [ **警告** ] 新版本图书缓存与旧版本图书缓存不兼容。

## **示例**

- --app=```<type[sfacg / cat]>```
- --account=```<account>```
- --password=```<password>```
- --download=```<type[bid / url]>```
- --search=```<关键词>```
- --show  < 查看 config.json 文件 >

## **免责声明**

- 此工具仅用于学习。请在下载后24小时内将其从计算机中删除。
- 请尊重版权，请勿自行传播爬虫图书，在任何情况下，作者或版权持有人均不对任何索赔负责
- 损害赔偿或其他责任，无论是在合同诉讼中，因软件或软件的使用或其他交易而产生的侵权行为或其他行为，作者均不承担责任。

## **文件树**

``` 
C:.
│  .gitignore
│  config.json
│  go.mod
│  go.sum
│  LICENSE
│  main.go
│  README.md
│  
├─.idea
│      workspace.xml
│
├─cache
├─config
│      command.go
│      config.go
│      file.go
│      msg.go
│      thread.go
│      tool.go
│ 
├─docs
│      81841388.png
│      84782349.png
│      README_zh-CN.md
│      README_zh-TW.md
│
├─epub
│  │  dirinfo.go
│  │  epub.go
│  │  fetchmedia.go
│  │  fs.go
│  │  pkg.go
│  │  toc.go
│  │  write.go
│  │  xhtml.go
│  │
│  └─internal
│      └─storage
│          │  storage.go
│          │
│          ├─memory
│          │      file.go
│          │      fs.go
│          │
│          └─osfs
│                  fs.go
│
├─save
├─src
│  │  book.go
│  │  bookshelf.go
│  │  catalogue.go
│  │  login.go
│  │  progressbar.go
│  │  search.go
│  │
│  ├─boluobao
│  │      api.go
│  │
│  ├─hbooker
│  │  │  api.go
│  │  │  Geetest.go
│  │  │  UrlConstants.go
│  │  │
│  │  └─Encrypt
│  │          decode.go
│  │          Encrypt.go
│  │
│  └─https
│          Header.go
│          param.go
│          request.go
│          urlconstant.go
│
└─struct
    │  command.go
    │  config.go
    │
    ├─book_info
    │      book_info.go
    │
    ├─hbooker_structs
    │  │  chapter.go
    │  │  config.go
    │  │  content.go
    │  │  detail.go
    │  │  geetest.go
    │  │  key.go
    │  │  login.go
    │  │  recommend.go
    │  │  search.go
    │  │
    │  ├─bookshelf
    │  │      bookshelf.go
    │  │
    │  └─division
    │          division.go
    │
    └─sfacg_structs
        │  account.go
        │  book.go
        │  catalogue.go
        │  content.go
        │  login.go
        │  search.go
        │
        └─bookshelf
                bookshelf.go

```
