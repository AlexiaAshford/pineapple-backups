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

## **免责声明**
- 此工具仅用于学习。请在下载后24小时内将其从计算机中删除。 
- 请尊重版权，请勿自行传播爬虫图书，在任何情况下，作者或版权持有人均不对任何索赔负责 
- 损害赔偿或其他责任，无论是在合同诉讼中，因软件或软件的使用或其他交易而产生的侵权行为或其他行为，作者均不承担责任。

 
