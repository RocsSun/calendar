package gov

import (
	"calendar/spider/parse"
	"calendar/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var searchUri = "http://sousuo.gov.cn/s.htm"

type GSpider struct{}

// Get http Get Method.
func (g GSpider) Get(uri string) (res *http.Response, err error) {
	if uri == "" {
		//return http.Response{}, errors.New("search info is null. ")
	}
	res, err = http.Get(uri)
	return
}

// SearchHoliday 搜索放假安排。按年份。
func (g GSpider) SearchHoliday(year int) (*http.Response, error) {
	if year == 0 {
		return nil, errors.New("search info is null. ")
	}
	info := fmt.Sprintf("国务院办公厅关于%d年部分节假日安排的通知", year)
	return g.Get(g.MakeSearchURL(info))
}

// MakeSearchURL 生成搜索的url。
func (g GSpider) MakeSearchURL(info string) string {
	if info == "" {
		//return http.Response{}, errors.New("search info is null. ")
	}
	r := utils.EncodeUri(searchUri, map[string]string{
		"t": "govall",
		"q": info,
	})
	return r
}

// HolidayDetail 获取国务院关于某年的放假安排。
func (g GSpider) HolidayDetail(year int) (*http.Response, error) {
	res := g.SearchHolidayUri(year)
	if res == "" {
		return nil, errors.New("放假通知url为空。")
	}

	return g.Get(res)
}

func (g GSpider) SearchHolidayUri(year int) string {
	if year < 2007 {
		log.Println("search holiday must after 2007.")
		return ""
	}
	r, err := GSpider{}.SearchHoliday(year)

	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		return ""
	}
	if r.StatusCode != 200 {
		log.Println("response status code is ", r.StatusCode)
		return ""
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	res := parse.GParse{}.ParseHolidayUri(year, b)
	return res
}

func SearchHolidayUri(year int) string {
	return GSpider{}.SearchHolidayUri(year)
}
