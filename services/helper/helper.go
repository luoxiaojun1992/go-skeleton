package helper

import (
	"log"
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

// TODO parseInt parseStr strIntval (redigo)
