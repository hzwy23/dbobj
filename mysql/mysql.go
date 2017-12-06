package mysql

import (
	"github.com/hzwy23/panda/logger"
	"database/sql"
	"github.com/hzwy23/panda/crypto/aes"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hzwy23/dbobj/dbhandle"
)

type mysql struct {
	db *sql.DB
}

func NewMySQL() dbhandle.DbObj {

	var err error

	o := new(mysql)

	red,err:=dbhandle.GetConfig()
	if err != nil {
		panic("cant not read ./conf/dbobj.conf.please check this file.")
	}

	tns,_ :=  red.Get("DB.tns")
	usr,_ := red.Get("DB.user")
	pad,_ := red.Get("DB.passwd")
	mc,_ := red.Get("DB.maxConn")
	maxConn := 100
	if len(mc) != 0 {
		mx, err := strconv.Atoi(mc)
		if err == nil {
			maxConn = mx
		}
	}

	if len(pad) == 24 {
		pad, err = aes.Decrypt(pad)
		if err != nil {
			logger.Error("Decrypt mysql passwd failed.")
			return nil
		}
	}

	o.db, err = sql.Open("mysql", usr+":"+pad+"@"+tns)

	if err != nil {
		logger.Error("open oracle database failed.", err.Error())
		return nil
	}
	if len(pad) != 24 {
		psd, err := aes.Encrypt(pad)
		if err != nil {
			logger.Error("decrypt passwd failed."+psd)
			return nil
		}
		psd = "\"" + psd + "\""
		red.Set("DB.passwd", psd)
	}

	// 设置连接池最大值
	o.db.SetMaxOpenConns(maxConn)
	o.db.SetConnMaxLifetime(0)
	logger.Info("create mysql dbhandle success. max connect value is:", maxConn)
	return o
}

func (this *mysql) GetErrorCode(err error) string {
	ret := err.Error()
	if n := strings.Index(ret, ":"); n > 0 {
		return strings.TrimSpace(ret[:n])
	} else {
		logger.Error("this error information is not mysql return info")
		return ""
	}
}

func (this *mysql) GetErrorMsg(err error) string {
	ret := err.Error()
	if n := strings.Index(ret, ":"); n > 0 {
		return strings.TrimSpace(ret[n+1:])
	} else {
		logger.Error("this error information is not mysql return info")
		return ""
	}
}

func (this *mysql) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewMySQL().(*mysql); ok {
				this.db = val.db
			}
			return this.db.Query(sql, args...)
		}
	}
	return rows, err
}

func (this *mysql) Exec(sql string, args ...interface{}) (sql.Result, error) {
	result, err := this.db.Exec(sql, args...)
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewMySQL().(*mysql); ok {
				this.db = val.db
			}
			return this.db.Exec(sql, args...)
		}
	}
	return result, err
}

func (this *mysql) Begin() (*sql.Tx, error) {
	tx, err := this.db.Begin()
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewMySQL().(*mysql); ok {
				this.db = val.db
			}
			return this.db.Begin()
		}
	}
	return tx, err
}

func (this *mysql) Prepare(sql string) (*sql.Stmt, error) {
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewMySQL().(*mysql); ok {
				this.db = val.db
			}
			return this.db.Prepare(sql)
		}
	}
	return stmt, err
}

func (this *mysql) QueryRow(sql string, args ...interface{}) *sql.Row {
	if this.db.Ping() != nil {
		logger.Warn("Connection is broken")
		if val, ok := NewMySQL().(*mysql); ok {
			this.db = val.db
		}
	}
	return this.db.QueryRow(sql, args...)
}

func init() {
	dbhandle.Register("mysql", NewMySQL)
}
