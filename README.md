<h1 align="center">
<img src="./docs/81841388.png" width="400" height='' alt="Noa Himesaka">
<br>
Noa Himesaka
</h1>

<h3 align="center">
Download books from <a href="https://book.sfacg.com/">sfacg</a> and
<a href="https://app.hbooker.com/">hbooker</a> to read them.
</h3>

[简中](./docs/README_zh-CN.md) | [繁中](./docs/README_zh-TW.md)

- - -
<br>

## **About downloading SFACG VIP books**

- WeChat API cannot download VIP chapters, and due to recent updates by the SFACG programmers, the new chapter API only
  returns images instead of text, making it impossible to download VIP chapters.
- To implement VIP chapter downloads, you need to enable the SFACG Android API, which can be done by modifying the App
  variable in the `main.go` file and setting it from `false` to `true`.
  <br><br>

- - -

## **Functions**

- The script implements download functions for SFACG [`Android`/`WeChat`] API and hbooker Android API.
- You can log in to your account and save your cookies in a `config.json` file.
- Input the book ID or URL to download the book to a local directory.
- Input the URL to download the book text from the URL.
- Supports downloading EPUB files from SFACG and hbooker.
- Search for books by keyword and download the search results.
- [ Warning ] The new version of book cache is incompatible with older versions of book cache.
  <br><br>

- - -

## Sign in to your Ciweimao account

- To use this script, you need to log in with your account and obtain your `token`.
- The new version of hbooker adds GEETEST verification, which will be triggered if you enter incorrect information or
  log in multiple times.
- The IP address may need to log in again after a few hours to avoid triggering the verification process. You can try
  changing the IP to avoid it.
  <br><br>

- - -

## Accessing the API using tokens

- Use tokens to access the API and bypass login
- Third-party captcha GEETEST has been added to the Ciweimao official server.
- The Ciweimao login is protected by GEETEST, which seems impossible to circumvent.
- You can capture packets of the Ciweimao Android App to obtain the account and login_token for logging in.

<br><br>

- - -

## **Example**

``` bash
NAME:
   pineapple-backups - https://github.com/AlexiaVeronica/pineapple-backups

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   V.1.8.4

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

## Thank you for JetBrains License
  [<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains" width="128"/>](https://www.jetbrains.com)

- [Goland](https://www.jetbrains.com/go/)

-  I would like to express my sincere gratitude and appreciation to JetBrains for providing me with the license to use and leverage software related to open source tools.
  
-  Your software has helped me to accomplish my tasks faster and more efficiently in my daily work, while also providing me with appropriate quality assurance. This is important for my personal growth and business development.

-   Thank you for your outstanding contribution to promoting the development of open-source software community. I will continue to use and support your company's software in my work and leisure time, and actively participate in contributing to the open-source community.

-   Once again, thank you JetBrains for your support!

-   Best regards, Alexia 

## **Disclaimers**

- This tool is for educational purposes only. Please delete it from your computer within 24 hours after downloading.
- Please respect the copyright and do not distribute the crawled books yourself.
- The authors or copyright holders shall not be liable for any claim, damages, or other liability, whether in an action
  of contract, tort, or otherwise, arising from, out of, or in connection with the software or the use or other dealings
  in the software, including but not limited to the use of the software for illegal purposes. The author is not
  responsible for any legal consequences.
- If you have any questions, please contact me via GitHub issues or email.