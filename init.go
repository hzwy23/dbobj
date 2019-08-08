package dbobj

import (
	"github.com/hzwy23/dbobj/dbhandle"

	_ "github.com/hzwy23/dbobj/mysql"
)

func init() {
	conf, err := dbhandle.GetConfig()
	if err != nil {
		panic("init database failed." + err.Error())
	}
	Default, err = conf.Get("DB.type")
	if err != nil {
		panic("get default database type failed." + err.Error())
	}
	InitDB(Default)
}
