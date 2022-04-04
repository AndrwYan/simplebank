package api

import (
	"github.com/AndrwYan/simplebank/db/util"
	"github.com/go-playground/validator/v10"
)

//自定义货币校验器

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {

	//反射值，必须获取它的值作为一个空接口
	if currency, ok := fl.Field().Interface().(string); ok {
		//校验是否支持这种货币
		return util.IsSupportedCurrency(currency)
	} else {
		return false
	}
}
