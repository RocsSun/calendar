package holiday

import (
	"calendar/spider/gov"
	"calendar/spider/parse"
	"calendar/utils"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

var DateMap map[string]bool
var Year = -9999

func init() {
	DateMap = make(map[string]bool)
}

func Holiday(year int) {
	Year = year
	r, err := gov.GSpider{}.HolidayDetail(year)
	defer r.Body.Close()

	if err != nil {
		return
	}
	if r.StatusCode != 200 {
		return
	}
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	res := parse.GParse{}.ParseHolidayInfo(rb)
	for _, v := range res {
		fmt.Println(v)
		parseHolidayDate(v)
	}
}

// parseHolidayDate 解析国务院具体的放假安排，具体到每一天，false为放假。true为工作日。
func parseHolidayDate(in string) {
	r := regexp.MustCompile(`\d{0,4}年?\d+月\d+日至?\d{0,2}月?\d{0,2}日?`)
	t := r.FindAllString(in, -1)

	for i, _ := range t {
		t[i] = strings.ReplaceAll(t[i], "月", "-")
		t[i] = strings.ReplaceAll(t[i], "日", "")
	}
	tmp := strings.Split(t[0], "至")
	start := fmt.Sprintf("%d-%s", Year, tmp[0])
	end := ""

	if strings.Contains(tmp[1], "-") {
		end = fmt.Sprintf("%d-%s", Year, tmp[1])
	} else {
		end = fmt.Sprintf("%d-%s-%s", Year, strings.Split(tmp[0], "-")[0], tmp[1])
	}

	start = utils.ModifyDateFormat(start)
	end = utils.ModifyDateFormat(end)

	st, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Println(err)
		return
	}

	et, err := time.Parse("2006-01-02", end)
	if err != nil {
		log.Println(err)
		return
	}

	for i := st; et.Sub(i) >= 0; i = i.Add(24 * time.Hour) {
		DateMap[i.Format("2006-01-02")] = false
	}

	if len(t) > 1 {
		for i := 1; i < len(t); i++ {
			dt := fmt.Sprintf("%d-%s", Year, t[i])
			dt = utils.ModifyDateFormat(dt)
			DateMap[dt] = true
		}
	}
}

//
//func parseInfo(in string, year int) map[string]bool {
//	r := regexp.MustCompile(`\d{0,4}年?\d+月\d+日至?\d{0,2}月?\d{0,2}日?`)
//	t := r.FindAllString(in, -1)
//	r = regexp.MustCompile(`月`)
//	r1 := regexp.MustCompile(`日`)
//	for i, _ := range t {
//		t[i] = r.ReplaceAllString(t[i], "-")
//		t[i] = r1.ReplaceAllString(t[i], "")
//	}
//
//	start := fmt.Sprintf("%d-%s", year, *(*string)(unsafe.Pointer(&tmp[0])))
//	end := ""
//
//	start = utils.ModifyDateFormat(start)
//	end = utils.ModifyDateFormat(end)
//
//	st, err := time.Parse("2006-01-02", start)
//	if err != nil {
//		log.Panicln(err)
//	}
//
//	et, err := time.Parse("2006-01-02", end)
//	if err != nil {
//		log.Panicln(err)
//	}
//	res := make(map[string]bool)
//	for i := st; et.Sub(i) >= 0; i = i.Add(24 * time.Hour) {
//		res[i.Format("2006-01-02")] = false
//	}
//
//	return res
//}
