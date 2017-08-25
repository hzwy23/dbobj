package utils

import (
	"errors"
	"regexp"
	"strings"
)

// 将字符串类型的时间转换成日期
func DateFormat(str string, model string) (string, error) {
	switch model {
	case "YYYY-MM-DD":
		pattern := `^[1-2]{1}[0-9]{3}[-|/][0-9]{1,2}[-|/][0-9]{1,2}`
		re, err := regexp.Compile(pattern)
		if err != nil {
			return str, err
		}

		rst := re.FindString(str)
		if rst == "" {
			return str, errors.New("no match")
		}
		return rst, nil
	case "YYYY-MM-DD HH24:MM:SS":
		pattern := `^[1-2]{1}[0-9]{3}(-|/)[0-9]{2}(-|/)[0-9]{2}(T)[0-9]{2}:[0-9]{2}:[0-9]{2}`
		re, err := regexp.Compile(pattern)
		if err != nil {
			return str, err
		}
		rst := re.FindString(str)
		if rst == "" {
			return str, errors.New("no match")
		}
		return strings.Replace(rst, "T", " ", 1), nil
	}
	return str, errors.New("model is unsupported.")
}
