package gov

import (
	"errors"
	"fmt"
	"github.com/RocsSun/calendar/spider/parse"
	"github.com/RocsSun/calendar/utils"
	"io"
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
		log.Fatalln("GSpider.Get, uri 不能为空。")
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
	res, err := g.SearchHolidayUri(year)
	if err != nil {
		return nil, err
	}
	if res == "" {
		return nil, errors.New(fmt.Sprintf("放假通知url为空，未找到%d年度的放假安排通知。", year))
	}
	return g.Get(res)
}

func (g GSpider) SearchHolidayUri(year int) (string, error) {
	if year < 2007 {
		return "", errors.New("search holiday must after 2007. ")
	}
	r, err := GSpider{}.SearchHoliday(year)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(r.Body)

	if r.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("response status code is %d. ", r.StatusCode))
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	res := parse.GParse{}.ParseHolidayUri(year, b)
	return res, nil
}
