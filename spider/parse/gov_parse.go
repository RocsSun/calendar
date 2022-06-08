package parse

import (
	"bytes"
	"fmt"
	"regexp"
	"unsafe"
)

type GParse struct{}

// ParseHolidayUri 解析获取放假安排的URI。
func (g GParse) ParseHolidayUri(year int, in []byte) (res string) {

	r := `<em[^>]*>`
	re := regexp.MustCompile(r)
	in = re.ReplaceAll(in, []byte(""))

	r = `</em>`
	re = regexp.MustCompile(r)
	in = re.ReplaceAll(in, []byte(""))

	r = `<a[^>]+href="(?P<uri>[a-zA-z]+?://[\w/&%\--~]*?)"[^<]+>[\w\S\s]*?</a>`
	re = regexp.MustCompile(r)
	result := re.FindAllSubmatch(in, -1)

	if len(result) == 0 {
		return
	}

	r = fmt.Sprintf("国务院办公厅关于%d年部分节假日安排的通知", year)
	zc := regexp.MustCompile(`/zhengce/content`)
	re = regexp.MustCompile(r)

	var resMap = make(map[string]bool)

	for _, v := range result {
		if len(re.FindAll(v[0], -1)) == 1 && len(zc.FindAll(v[1], -1)) == 1 {

			if resMap[*(*string)(unsafe.Pointer(&v[1]))] {
				continue
			}

			resMap[*(*string)(unsafe.Pointer(&v[1]))] = true
			res = *(*string)(unsafe.Pointer(&v[1]))

		}
	}
	return res
}

// ParseHolidayInfo 解析放假安排的相关信息。
func (g GParse) ParseHolidayInfo(in []byte) (res []string) {
	r := regexp.MustCompile(`<span[^>]*?>`)
	in = r.ReplaceAll(in, []byte(""))

	r = regexp.MustCompile(`</span>`)
	in = r.ReplaceAll(in, []byte(""))

	r = regexp.MustCompile(`<p[^>]*>(?P<con>[\s\w\S]*?)</p>`)
	ans := r.FindAllSubmatch(in, -1)
	for _, v := range ans {
		switch {
		case bytes.Contains(v[1], []byte("一")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		case bytes.Contains(v[1], []byte("二")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		case bytes.Contains(v[1], []byte("三")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		case bytes.Contains(v[1], []byte("四")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		case bytes.Contains(v[1], []byte("五")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		case bytes.Contains(v[1], []byte("六")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		case bytes.Contains(v[1], []byte("七")):
			res = append(res, *(*string)(unsafe.Pointer(&v[1])))
		}
	}
	return res
}
