package i18n

const ERR_UNSUPPORTED_METHOD = "405"
const ERR_SERVICE_UNAVAILABLE_498 = "498"
const ERR_SERVICE_UNAVAILABLE_499 = "499"
const ERR_SERVICE_UNAVAILABLE_500 = "500"

var i18NMessages = map[string]string{
	ERR_UNSUPPORTED_METHOD:      "只支持 GET 和 POST 方法",
	ERR_SERVICE_UNAVAILABLE_498: "服务异常",
	ERR_SERVICE_UNAVAILABLE_499: "服务异常",
	ERR_SERVICE_UNAVAILABLE_500: "服务异常",
}

func GetMessage(code string) string {
	if v, ok := i18NMessages[code]; ok {
		return v
	}

	return "服务正忙，请稍后重试"
}
