<h1 align="center">
  <img src="./81841388.png" width="400" height='' alt="Noa Himesaka">
  <br>Noa Himesaka<br>  
</h1>
<h3 align="center">
    下载 <a href="https://book.sfacg.com/">菠萝包</a> 和 
    <a href="https://app.hbooker.com/">刺猬猫</a> 的小说到本地阅读. 

</h3>

## 关于下载sfacg vip书籍 

---

- 微信API无法下载vip章节，因为sfacg程序员更新了章节API返回值，新的API无法获取文本，只能获取图片，因此无法下载vip章节。

- 您需要启用sfacg Android API来实现vip章节下载，您可以修改`main.go`文件中的`App`变量，并将`false`设置为`true`以实现API切换。



<br><br>


## **功能**



- 菠萝包[`Android`/`WeChat`]刺猬猫安卓接口实现了下载功能

- 登录您的帐户并将cookie保存到`config.json`

- 输入图书id或url并将图书下载到本地目录

- 输入url并从url下载书籍文本

- 支持从sfacg和hbooker下载epub

- 按关键字搜索书籍，并下载搜索结果

- [**警告**]新版本图书缓存与旧版本图书缓存不兼容。


<br><br>



## 登录您的ciweimao帐户

- - -

- 登录您的帐户以获取使用此脚本的“令牌”

- hboker新版本增加了GEETEST验证，如果您输入错误信息或多次登录，将触发GEETEST校验。

- IP地址可能需要在几小时后再次登录以避免触发验证，您可以尝试更改IP以避免触发确认。


<br><br>



## API访问通过令牌实现。

- - -

- **采用token访问api，绕过登录**

- 第三方captcha geetest已添加到ciweimao官方服务器。

- ciweimao登录受到geetest的保护，这似乎是不可能规避的。

- 您可以提供`刺猬猫 Android`应用程序的数据包捕获以获取`account`和`login token`登录。


<br><br>

## **Example**
- - -
``` bash
NAME:
   pineapple-backups - https://github.com/VeronicaAlexia/pineapple-backups

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   V.1.7.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -a value, --app value       cheng app type (default: "cat")
   -d value, --download value  book id
   -t, --token                 input hbooker token
   -m value, --max value       change max thread number (default: 16)
   -u value, --user value      input account name
   -p value, --password value  input password
   --update                    update book
   -s value, --search value    search book by keyword
   -l, --login                 login local account
   -e, --epub                  start epub
   --help, -h                  show help
   --version, -v               print the version

```

## **Disclaimers**

- This tool is for learning only. Please delete it from your computer within 24 hours after downloading.
- Please respect the copyright and do not spread the crawled books by yourself.
- In no event shall the authors or copyright holders be liable for any claim damages or other liability, whether in an
  action of contract tort or otherwise, arising from, out of or in connection with the software or the use or other
  dealings in the software , including but not limited to the use of the software for illegal purposes,author is not
  responsible for any legal consequences.
- If you have any questions, please contact me by github issues or email.