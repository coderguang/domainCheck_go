domainCheck_go
====

[![Build Status](https://travis-ci.org/coderguang/domainCheck_go.svg?branch=master)](https://travis-ci.org/coderguang/domainCheck_go)
![](https://img.shields.io/badge/language-golang-orange.svg)
[![codebeat badge](https://codebeat.co/badges/4c9ab03b-4424-48e3-8d1f-66a5350374e9)](https://codebeat.co/projects/github-com-coderguang-domaincheck_go-master)
[![](https://img.shields.io/badge/wp-@royalchen-blue.svg)](https://www.royalchen.com)


## what it can do
* auto check domain **regist status**
* auto update domain **expiry time**
* **notice by email** when a high value can be registed

## require
require | version
------ | ------
go | 1.11
GameEngine_go | [v1.0.0](https://github.com/coderguang/GameEngine_go)
go-sql-driver/mysql | [v1.4.1 ](https://github.com/go-sql-driver/mysql)


## how to start
### 1. clone repository 
```shell
git clone git@github.com:coderguang/domainCheck_go.git domainCheck_go
```

### 2. import a null sql table to you mysql,sql file in **_sql/domain_info.sql_** 

### 3. config your *database* and *email* infomation in **config/config.json**
```json
{
    "Dburl":"localhost",  //your mysql database url
    "Dbport":"3306",  //database port
    "Dbuser":"test", //database user
    "Dbpwd":"test", //database password
    "Dbname":"domain", //database name
    "DoroutinesNum":5, //how many coroutine to check domain(only first run use)
    "EmailFrom":"",  //notice email from,empty mean wouldn't notice by email
    "EmailPwd":"", //notice email password
    "EmailTo":"royalchen@royalchen.com", //receiver 
    "Smtp":"smtps://smtp.exmail.qq.com", //smtp addr
    "SmtpPort":"465"  //smtp port
}
```

### 4. procedure for execution
```shell
go run main.go config/config.json
```
### 5. if no problem,the output message should like this:
 ![run img](https://github.com/coderguang/img/blob/master/domainCheck_go/run.png)
 
 enter any key ,it would show command,if **first** run,enter **["StartGenDataAndCheck"]** to create domain file and start run,if not first run,just enter **["StartCheckByDb"]** for scanner by database data.
  after enter **["StartGenDataAndCheck"]**,output message should like this
  ![run_ok img](https://github.com/coderguang/img/blob/master/domainCheck_go/run_ok.png)
  
### 6. all log will write in _log/xx-xx-xx.log_
   you can also get domain info in you database,like below:
   ![domain](https://github.com/coderguang/img/blob/master/domainCheck_go/domain.png)
   

## About me

**Author** | _[royalchen](https://www.royalchen.com)_
---------- | -----------------
email  | royalchen@royalchen.com
qq  | royalchen@royalchen.com
website | [www.royalchen.com](https://www.royalchen.com)
  
 
 
 

