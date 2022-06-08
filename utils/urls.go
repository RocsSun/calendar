package utils

import (
	"fmt"
	"net/url"
)

// EncodeUri 编码URL，添加query。
func EncodeUri(u string, query map[string]string) string {
	uu, err := url.Parse(u)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	q := uu.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	uu.RawQuery = q.Encode()
	return uu.String()
}
