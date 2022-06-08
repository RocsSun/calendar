package utils

import (
	"fmt"
	"strings"
)

// ModifyDateFormat 修改日期格式字符串，使符合YY(YY)-MM-DD hh:mm:ss，支持的参数格式YYYY-MM-DD hh:mm:ss
func ModifyDateFormat(t string) string {
	if t == "" {
		return ""
	}
	res := strings.Split(t, "-")

	if len(res[1]) < 2 {
		res[1] = fmt.Sprintf("0%s", res[1])
	}
	if len(res[2]) < 2 {
		res[2] = fmt.Sprintf("0%s", res[2])
	}
	return strings.Join(res, "-")
}
