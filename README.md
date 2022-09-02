<h1 align="center">
  <img src="./docs/81841388.png" width="400" height='' alt="Noa Himesaka">
  <br>Noa Himesaka<br>  
</h1>
<h3 align="center">
    Download books from the <a href="https://book.sfacg.com/">sfacg</a> and 
    <a href="https://app.hbooker.com/">hbooker</a> to read them. 
    
</h3> 

[简中](./docs/README_zh-CN.md) | [繁中](./docs/README_zh-TW.md)

## **Functions**
 - Download function is implemented for sfacg WeChat Api and hbooker Android API
 - Login your account and save cookies to a ```config.json```
 - Input the book id or url and download the book to the local directory
 - Input url and download book text from the url
 - Support download epub from sfacg and hbooker
 - Search books by keyword,and download the search result
 - [ **warning** ] New version book cache is incompatible with older book cache.

## **File tree**
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
## **Example**

- --app=```<type[sfacg / cat]>```
- --account=```<account>```
- --password=```<password>```
- --download=```<type[bid / url]>``` 
- --search=```<keyword>```
- --show  < show the config.json file >
 
## **Disclaimers**
- This tool is for learning only. Please delete it from your computer within 24 hours after downloading.
- Please respect the copyright and do not spread the crawled books by yourself.
- In no event shall the authors or copyright holders be liable for any claim
- damages or other liability, whether in an action of contract
- tort or otherwise, arising from, out of or in connection with the software or the use or other dealings in the
  software.

 
