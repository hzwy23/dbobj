// Copyright 2016 huangzhanwei. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dbobj

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hzwy23/dbobj/dbhandle"
	"github.com/hzwy23/dbobj/utils"
)

var (
	dbobj   dbhandle.DbObj
	Default = "mysql"
)

func InitDB(dbtyp string) error {
	if dbobj == nil {
		if val, ok := dbhandle.Adapter[dbtyp]; ok {
			dbobj = val()
		}
	}
	return nil
}

// get default DbObj name
// return DbObj name
func GetDefaultName() string {
	return Default
}

func Begin() (*sql.Tx, error) {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return nil, errors.New("can not connect database again.")
		}
		return dbobj.Begin()
	}
	return dbobj.Begin()
}

func Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return nil, errors.New("can not connect database again.")
		}
		return dbobj.Query(sql, args...)
	}
	return dbobj.Query(sql, args...)
}

func QueryRow(sql string, args ...interface{}) *sql.Row {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return nil
		}
		return dbobj.QueryRow(sql, args...)
	}
	return dbobj.QueryRow(sql, args...)
}

func Exec(sql string, args ...interface{}) (sql.Result, error) {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return nil, errors.New("connect database failed.")
		}
		return dbobj.Exec(sql, args...)
	}
	return dbobj.Exec(sql, args...)
}

func Prepare(sql string) (*sql.Stmt, error) {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return nil, errors.New("can not connect database again.")
		}
		return dbobj.Prepare(sql)
	}
	return dbobj.Prepare(sql)
}

func GetErrorCode(errs error) string {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return err.Error()
		}
		return dbobj.GetErrorCode(errs)
	}
	return dbobj.GetErrorCode(errs)
}

func GetErrorMsg(errs error) string {
	if dbobj == nil {
		err := InitDB(Default)
		if err != nil {
			return err.Error()
		}
		return dbobj.GetErrorMsg(errs)
	}
	return dbobj.GetErrorMsg(errs)
}

// 将查询结果赋予某个slice
func QueryForSlice(sql string, rst interface{}, args ...interface{}) error {
	rows, err := Query(sql, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return Scan(rows, rst)
}

// 将查询结果赋予某个对象
func QueryForStruct(sql string, rst interface{}, args ...interface{}) error {
	rows, err := Query(sql, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return Scan(rows, rst)
}

func PackArgs(args ...interface{}) []interface{} {
	var obj []interface{}
	return append(obj, args...)
}

// 将查询结果，赋予一系列对象
func QueryForObject(sql string, args []interface{}, result ...interface{}) error {
	rows, err := Query(sql, args...)
	if err != nil {
		return err
	}
	return ScanRow(rows, result...)
}

// 扫描单行值
// 如果有多行值，则只返回第一行
func ScanRow(rows *sql.Rows, result ...interface{}) error {
	// 获取查询结果字段数量
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return err
	}
	size := len(cols)
	values := make([]interface{}, size)

	for index, val := range result {
		obj := reflect.ValueOf(val).Elem()
		switch obj.Kind() {
		case reflect.String:
			values[index] = &sql.NullString{}
		case reflect.Float64, reflect.Float32:
			values[index] = &sql.NullFloat64{}
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int, reflect.Int8, reflect.Uint64, reflect.Uint32, reflect.Uint:
			values[index] = &sql.NullInt64{}
		case reflect.Bool:
			values[index] = &sql.NullBool{}
		default:
			values[index] = &[]byte{}
			fmt.Println("dbobj.Scan unsupported this type is :", obj.Kind())
		}
	}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return err
		}
		argLen := len(result)

		for index := 0; index < argLen; index++ {
			obj := reflect.ValueOf(result[index]).Elem()

			switch obj.Kind() {
			case reflect.String:
				obj.SetString(values[index].(*sql.NullString).String)
			case reflect.Float64, reflect.Float32:
				obj.SetFloat(values[index].(*sql.NullFloat64).Float64)
			case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int, reflect.Int8, reflect.Uint64, reflect.Uint32, reflect.Uint:
				obj.SetInt(values[index].(*sql.NullInt64).Int64)
			case reflect.Bool:
				obj.SetBool(values[index].(*sql.NullBool).Bool)
			default:
				str, _ := values[index].(*[]byte)
				obj.SetBytes(*str)
				fmt.Println("dbobj.Scan unsupported this type is :", obj.Kind())
			}
		}
		return nil
	}
	return nil
}

// Function: scan DbObj query result
// Time: 2016-06-10
// Author: huangzhanwei
// Notice: second argument of type must be valid pointer.
func Scan(rows *sql.Rows, rst interface{}) error {
	defer rows.Close()

	// 查询结果返回地址，必须是指针类型
	obj := reflect.ValueOf(rst)
	if obj.Kind() != reflect.Ptr || obj.IsNil() {
		fmt.Errorf("second must be valid pointer")
		return errors.New("second argument of type must be valid pointer")
	}

	// 查询接收方的类型，如果返回的是多行记录
	// 必须是Slice
	switch obj.Elem().Kind() {
	case reflect.Slice:
		return scanForSlice(obj, rows)
	case reflect.Struct:
		return scanForStruct(obj, rows)
	default:
		return errors.New("unsupported type. type is:" + obj.Elem().Kind().String())
	}
	return nil
}

// Count total rows by sqlText
//
func Count(sql string, args ...interface{}) int64 {
	var total int64
	err := QueryRow(sql, args...).Scan(&total)
	if err != nil {
		fmt.Errorf("%v", err)
		return 0
	}
	return total
}

func scanForStruct(obj reflect.Value, rows *sql.Rows) error {
	// 获取查询结果字段数量
	cols, err := rows.Columns()
	if err != nil {
		fmt.Errorf("%v", err)
		return err
	}
	size := len(cols)

	obj = obj.Elem()
	argList := make([]interface{}, size)

	for rows.Next() {
		for j := 0; j < size; j++ {
			switch obj.Field(j).Kind() {
			case reflect.String:
				argList[j] = &sql.NullString{}
			case reflect.Float64, reflect.Float32:
				argList[j] = &sql.NullFloat64{}
			case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int, reflect.Int8, reflect.Uint64, reflect.Uint32, reflect.Uint:
				argList[j] = &sql.NullInt64{}
			case reflect.Bool:
				argList[j] = &sql.NullBool{}
			default:
				argList[j] = &[]byte{}
				fmt.Println("dbobj.Scan unsupported this type is :", obj.Field(j).Kind())
			}
		}

		err := rows.Scan(argList...)
		if err != nil {
			fmt.Errorf("%v", err)
			return err
		}

		// 读取查询到的值，并赋值到返回变量rst中
		for index, vals := range argList {
			switch vals.(type) {
			case *sql.NullString:
				str := vals.(*sql.NullString).String
				tag := obj.Type().Field(index).Tag.Get("dateType")
				if tag != "" {
					str, _ = utils.DateFormat(str, tag)
				}
				obj.Field(index).SetString(str)
			case *sql.NullBool:
				obj.Field(index).SetBool(vals.(*sql.NullBool).Bool)
			case *sql.NullFloat64:
				obj.Field(index).SetFloat(vals.(*sql.NullFloat64).Float64)
			case *sql.NullInt64:
				obj.Field(index).SetInt(vals.(*sql.NullInt64).Int64)
			default:
				str, _ := vals.(*[]byte)
				obj.Field(index).SetBytes(*str)
				fmt.Println("this type cannot by supported now, set Type to []byte, value is: ", string(*str))
			}
		}
		return nil
	}
	return nil
}

func scanForSlice(obj reflect.Value, rows *sql.Rows) error {
	// 获取查询结果字段数量
	cols, err := rows.Columns()
	if err != nil {
		fmt.Errorf("%v", err)
		return err
	}
	size := len(cols)

	obj = obj.Elem()
	argList := make([]interface{}, size)
	var i = 0
	for rows.Next() {
		if i >= obj.Cap() {
			newcap := obj.Cap() + obj.Cap()/2
			if newcap < 4 {
				newcap = 4
			}
			newv := reflect.MakeSlice(obj.Type(), obj.Len(), newcap)
			reflect.Copy(newv, obj)
			obj.Set(newv)
		}

		if i >= obj.Len() {
			obj.SetLen(i + 1)
			if max := obj.Index(0).NumField(); max < size {
				fmt.Errorf("slice colunms less then dest ", max, size)
				return errors.New("slice colunms less then dest. numFiled is :" + strconv.Itoa(max) + ". need numFiled " + strconv.Itoa(size))
			}
		}

		if i == 0 {
			for j := 0; j < size; j++ {
				switch obj.Index(i).Field(j).Kind() {
				case reflect.String:
					argList[j] = &sql.NullString{}
				case reflect.Float64, reflect.Float32:
					argList[j] = &sql.NullFloat64{}
				case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int, reflect.Int8, reflect.Uint64, reflect.Uint32, reflect.Uint:
					argList[j] = &sql.NullInt64{}
				case reflect.Bool:
					argList[j] = &sql.NullBool{}
				default:
					argList[j] = &[]byte{}
					fmt.Println("dbobj.Scan unsupported this type is :", obj.Index(i).Field(j).Kind())
				}
			}
		}

		err := rows.Scan(argList...)
		if err != nil {
			fmt.Errorf("%v", err)
			return err
		}

		// 读取查询到的值，并赋值到返回变量rst中
		for index, vals := range argList {
			switch vals.(type) {
			case *sql.NullString:
				str := vals.(*sql.NullString).String
				tag := obj.Index(i).Type().Field(index).Tag.Get("dateType")
				if tag != "" {
					str, _ = utils.DateFormat(str, tag)
				}
				obj.Index(i).Field(index).SetString(str)
			case *sql.NullBool:
				obj.Index(i).Field(index).SetBool(vals.(*sql.NullBool).Bool)
			case *sql.NullFloat64:
				obj.Index(i).Field(index).SetFloat(vals.(*sql.NullFloat64).Float64)
			case *sql.NullInt64:
				obj.Index(i).Field(index).SetInt(vals.(*sql.NullInt64).Int64)
			default:
				str, _ := vals.(*[]byte)
				obj.Index(i).Field(index).SetBytes(*str)
				fmt.Println("this type cannot by supported now, set Type to []byte, value is: ", string(*str))
			}
		}
		i++
	}
	return nil
}
