package helper

import (
	redigo "github.com/gomodule/redigo/redis"
	"github.com/luoxiaojun1992/form"
	"github.com/syyongx/php2go"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// Error Handling
func PanicErr(errMsg string, errs ...error) {
	log.Panicln(errMsg, ": ", errs)
}

func CheckErr(errs ...error) bool {
	if len(errs) <= 0 {
		return false
	}

	hasError := false
	for _, err := range errs {
		if err != nil {
			hasError = true
		}
	}

	return hasError
}

func CheckErrThenPanic(errMsg string, errs ...error) {
	if !CheckErr(errs...) {
		return
	}

	PanicErr(errMsg, errs...)
}

func Abort(errMsg string, errs ...error) {
	log.Fatalln(errMsg, ": ", errs)
}

func CheckErrThenAbort(errMsg string, errs ...error) {
	if !CheckErr(errs...) {
		return
	}

	Abort(errMsg, errs...)
}

func LogErr(errMsg string, errs ...error) {
	log.Println(errMsg, ": ", errs)
}

func CheckErrThenLog(errMsg string, errs ...error) {
	if !CheckErr(errs...) {
		return
	}

	LogErr(errMsg, errs...)
}

// Expr
func ConditionalOperator(condition bool, expr1 func() interface{}, expr2 func() interface{}) interface{} {
	if condition {
		return expr1()
	} else {
		return expr2()
	}
}

// Str
func ParseStr(val interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}

	switch realVal := val.(type) {
	case bool:
		if realVal {
			return "1", nil
		} else {
			return "0", nil
		}
	case int:
		return strconv.Itoa(realVal), nil
	case int8:
		x := int(realVal)
		if int8(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case int16:
		x := int(realVal)
		if int16(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case int32:
		x := int(realVal)
		if int32(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case int64:
		x := int(realVal)
		if int64(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case uint:
		return strconv.Itoa(int(realVal)), nil
	case uint8:
		x := int(realVal)
		if uint8(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case uint16:
		x := int(realVal)
		if uint16(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case uint32:
		x := int(realVal)
		if uint32(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	case uint64:
		x := int(realVal)
		if uint64(x) != realVal {
			return "", strconv.ErrRange
		}
		return strconv.Itoa(x), nil
	}

	strVal, errStrVal := redigo.String(val, err)
	if (errStrVal != redigo.ErrNil) && CheckErr(errStrVal) {
		return strVal, errStrVal
	} else {
		return strVal, nil
	}
}

// Int
func StrIntval(str string) (int, error) {
	str = strings.TrimSpace(str)
	pre := ""

	if strings.HasPrefix(str, "-") || strings.HasPrefix(str, "+") {
		pre = str[0:1]

		if len(str) > 1 {
			str = str[1:]
		} else {
			str = ""
		}
	}

	i := 0
	if len(str) > 0 {
		i = strings.IndexFunc(str, func(r rune) bool {
			return !unicode.IsNumber(r)
		})
	}

	if i > 0 {
		str = str[0:i]
	} else if i == 0 {
		str = ""
	}

	if str == "" {
		str = "0"
	}

	return strconv.Atoi(pre + str)
}

func ParseInt(val interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	switch realVal := val.(type) {
	case bool:
		if realVal {
			return 1, nil
		} else {
			return 0, nil
		}
	case int:
		return realVal, nil
	case int8:
		x := int(realVal)
		if int8(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	case int16:
		x := int(realVal)
		if int16(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	case int32:
		x := int(realVal)
		if int32(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	case uint:
		return int(realVal), nil
	case uint8:
		x := int(realVal)
		if uint8(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	case uint16:
		x := int(realVal)
		if uint16(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	case uint32:
		x := int(realVal)
		if uint32(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	case uint64:
		x := int(realVal)
		if uint64(x) != realVal {
			return 0, strconv.ErrRange
		}
		return x, nil
	}

	strVal, errStrVal := redigo.String(val, err)
	if (errStrVal != redigo.ErrNil) && CheckErr(errStrVal) {
		return redigo.Int(val, err)
	} else {
		return StrIntval(strVal)
	}
}

// Type
func Isset(arr map[string]interface{}, key string) bool {
	if _, existed := arr[key]; !existed {
		return false
	}

	if arr[key] == nil {
		return false
	}

	return true
}

func Empty(val interface{}) bool {
	if strVal, isStrVal := val.(string); isStrVal {
		if strVal == "0" {
			return true
		}
	}

	return php2go.Empty(val)
}

// Url
func HttpBuildQuery(val interface{}) (string, error) {
	return form.EncodeToString(val)
}
