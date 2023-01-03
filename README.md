<h1 align="center">
  <img src="./docs/81841388.png" width="400" height='' alt="Noa Himesaka">
  <br>Noa Himesaka<br>  
</h1>
<h3 align="center">
    Download books from the <a href="https://book.sfacg.com/">sfacg</a> and 
    <a href="https://app.hbooker.com/">hbooker</a> to read them. 

</h3> 
 

[简中](./docs/README_zh-CN.md) | [繁中](./docs/README_zh-TW.md)

- - -
<br>
 

## **About download sfacg vip books**    
- - -
- Wechat API can't download vip chapters, because the sfacg programmer updates the chapter api return value, the new api can't get the text, only get the picture, so you can't download vip chapters.
- you need to enable sfacg Android API to implement vip chapter download, you can modify the `App` variable in the `main.go` file and set `false` to `true` to implement the api switch. 


<br><br>

## **Functions**
- - -
- Download function is implemented for sfacg [`Android`/`WeChat`] Api and hbooker Android API 
- Login your account and save cookies to a ```config.json```
- Input the book id or url and download the book to the local directory
- Input url and download book text from the url
- Support download epub from sfacg and hbooker
- Search books by keyword,and download the search result
- [ **warning** ] New version book cache is incompatible with older book cache.

<br><br>


## Sign in to your ciweimao Account 
- - -
  - Login your account to get your `token` to use this script
  - hbooker new version add GEETEST verification, if you enter the wrong information or log in multiple times, GEETEST verification will be triggered.
  - IP address may need to log in again after a few hours to avoid triggering verification, you can try to change the IP to avoid triggering verification.

<br><br>


## API access is achieved through token.
- - -
  - **Adopt token to access api, bypass login**
  - third party captcha geetest has been adding to the ciweimao official server.
  - ciweimao login is protected by geetest, which seems impossible to circumvent.
  - you can **`Packet Capture`** of the `ciweimao Android App` to get the `account` and `login_token` login.

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