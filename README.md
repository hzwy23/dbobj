## 介绍
dbobj 提供一个操作数据库的工具包，通过直接写sql的方式操作数据库。目前支持oracle,mysql,mariadb数据库

### 获取方法：

> go get github.com/hzwy23/dbobj

## 使用方法:

### oracle数据库 

1. 选择的是oracle数据库,请按照go-oci8包的要求配置pkgconfig和oracle instantclient.

2. oci8.pc在vendor/github.com/mattn/go-oci8中.请按照要求,修改oci8.pc文件

3. 修改dbobj根目录中init.go文件，**将第5行的注释去掉**。

3. 请设置环境变量.HBIGDATA_HOME.这个变量中创建目录conf.然后将dbobj中的asofdate.conf复制到conf中.

### mysql，mariadb数据库
1. 请设置环境变量.HBIGDATA_HOME.这个变量中创建目录conf.然后将dbobj中的asofdate.conf复制到conf中.

### 创建目录

```shell
    export HBIGDATA_HOME=/opt/go/hcloud
    mkdir $HBIGDATA_HOME/conf
    cp asofdate.conf $HBIGDATA_HOME/conf
```

### 工程目录样式:
```
$HBIGDATA_HOME

            --------github.com

            ------------hzwy23

            ----------------dbobj

            ----conf

            --------asofdate.conf
            
            main.go
```

在指定的配置文件目录中创建配置文件,配置文件名称指定为:asofdate.conf,在文件中输入下面信息:

#### mysql配置文件

```shell
    DB.type=mysql
    DB.tns = "tcp(localhost:3306)/bigdata"
    DB.user = root
    DB.passwd= huang
```

#### oracle配置文件

``` shell
    DB.type=oracle
    DB.tns = "192.168.1.101:1521/orcl"
    DB.user = test
    DB.passwd= huang
	
```

#### 系统启动后,会默认自动对密码进行加密.

dbobj示例工程[dbobj example](https://github.com/hzwy23/dbobj-example)

### 例子
```go
package main

import (
	"fmt"

	"github.com/hzwy23/dbobj"
)

type UserInfo struct {
	UserId   string
	UserName string
}

func GetUserDetails(userId int) ([]UserInfo, error) {
	rows, err := dbobj.Query("select user_id,user_name from dbobj_test_table where age = ?", userId)
	if err != nil {
		fmt.Println("query table failed", err)
		return nil, err
	}
	var rst []UserInfo
	err = dbobj.Scan(rows, &rst)
	if err != nil {
		fmt.Println("scan table failed", err)
		return nil, err
	}
	return rst, nil
}

func main() {
	tmp1, err := GetUserDetails(12)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success", tmp1)
}

```