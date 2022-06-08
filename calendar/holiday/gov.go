package holiday

import (
	"calendar/spider/gov"
	"calendar/spider/parse"
	"calendar/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

var _year = -9999

// GovHoliday 国务院放假调休安排。false为放假了。true为调班。
func GovHoliday(year int) map[string]bool {
	_year = year
	var dm = make(map[string]bool)

	r, err := gov.GSpider{}.HolidayDetail(year)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Fatalln(errors.New(fmt.Sprintf("status code is %d, exit. ", r.StatusCode)))
		return nil
	}
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	res := parse.GParse{}.ParseHolidayInfo(rb)
	for _, v := range res {
		parseHolidayDate(v, dm)
	}
	return dm
}

// WorkCalendar 工作日历。false为休息。true为工作日。返回一年的所有的日期的节假日。
func WorkCalendar(year int) map[string]bool {
	res := make(map[string]bool)

	holiday := GovHoliday(year)

	if holiday == nil {
		return nil
	}
	start := fmt.Sprintf("%d-01-01", year)
	end := fmt.Sprintf("%d-12-31", year)

	st, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil
	}

	et, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil
	}

	for i := st; et.Sub(i) >= 0; i = i.Add(24 * time.Hour) {
		if v, ok := holiday[i.Format("2006-01-02")]; ok {
			res[i.Format("2006-01-02")] = v
		} else if i.Weekday() == 0 || i.Weekday() == 6 {
			res[i.Format("2006-01-02")] = false
		} else {
			res[i.Format("2006-01-02")] = true
		}
	}
	return res
}

// CurrentYearWorkCalendar 当前年份的节假日信息。
func CurrentYearWorkCalendar() map[string]bool {
	return WorkCalendar(time.Now().Year())
}

// CurrentYearWorkCalendarToJson 将当前年份的节假日信息导出到json文件。
func CurrentYearWorkCalendarToJson(fp string) {
	WorkCalendarToJson(time.Now().Year(), fp)
}

// WorkCalendarToJson 生成指定年份的节假日信息到给定的json文件。
func WorkCalendarToJson(year int, fp string) {
	if fp == "" {
		log.Fatalln("文件名称为空。")
	}
	b, err := json.Marshal(WorkCalendar(year))

	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(fp, b, 777)
	if err != nil {
		log.Fatalln(err)
	}
}

// parseHolidayDate 解析国务院具体的放假安排，具体到每一天，false为放假。true为工作日。
func parseHolidayDate(in string, dm map[string]bool) {
	r := regexp.MustCompile(`\d{0,4}年?\d+月\d+日至?\d{0,2}月?\d{0,2}日?`)
	t := r.FindAllString(in, -1)

	for i, _ := range t {
		t[i] = strings.ReplaceAll(t[i], "月", "-")
		t[i] = strings.ReplaceAll(t[i], "日", "")
	}
	var tmp []string

	if strings.Contains(t[0], "年") {
		tm := strings.Split(t[0], "年")[1]
		tmp = strings.Split(tm, "至")

	} else {
		tmp = strings.Split(t[0], "至")
	}

	start := fmt.Sprintf("%d-%s", _year, tmp[0])
	end := ""

	if strings.Contains(tmp[1], "-") {
		end = fmt.Sprintf("%d-%s", _year, tmp[1])
	} else {
		end = fmt.Sprintf("%d-%s-%s", _year, strings.Split(tmp[0], "-")[0], tmp[1])
	}

	start = utils.ModifyDateFormat(start)
	end = utils.ModifyDateFormat(end)
	st, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Fatalln(err)
		return
	}

	et, err := time.Parse("2006-01-02", end)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for i := st; et.Sub(i) >= 0; i = i.Add(24 * time.Hour) {
		dm[i.Format("2006-01-02")] = false
	}

	if len(t) > 1 {
		for i := 1; i < len(t); i++ {
			dt := fmt.Sprintf("%d-%s", _year, t[i])
			dt = utils.ModifyDateFormat(dt)
			dm[dt] = true
		}
	}
}
