package holiday

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/RocsSun/calendar/cache"
	"github.com/RocsSun/calendar/globals"
	"github.com/RocsSun/calendar/spider/gov"
	"github.com/RocsSun/calendar/spider/parse"
	"github.com/RocsSun/calendar/utils"
)

var _year = -9999

// GovHoliday 国务院放假调休安排。false为放假了。true为调班。
func GovHoliday(year int) (map[string]bool, error) {
	_year = year
	var dm = make(map[string]bool)

	r, err := gov.GSpider{}.HolidayDetail(year)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(r.Body)

	if r.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code is %d, exit. ", r.StatusCode))
	}
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	res := parse.GParse{}.ParseHolidayInfo(rb)
	for _, v := range res {
		if err = parseHolidayDate(v, dm); err != nil {
			return nil, err
		}
	}
	return dm, nil
}

// WorkCalendar 工作日历。false为休息。true为工作日。返回一年的所有的日期的节假日。
func WorkCalendar(year int) (map[string]bool, error) {
	res := make(map[string]bool)

	if check(year) {
		return readCache(year), nil
	}

	holiday, err := GovHoliday(year)
	if err != nil {
		return nil, err
	}
	if holiday == nil {
		return nil, errors.New("放假日期为空，获取失败。")
	}

	start := fmt.Sprintf("%d-01-01", year)
	end := fmt.Sprintf("%d-12-31", year)

	st, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, err
	}

	et, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, err
	}

	for i := st; et.Sub(i).Hours() >= 0; i = i.Add(24 * time.Hour) {
		if v, ok := holiday[i.Format("2006-01-02")]; ok {
			res[i.Format("2006-01-02")] = v
		} else if i.Weekday() == 0 || i.Weekday() == 6 {
			res[i.Format("2006-01-02")] = false
		} else {
			res[i.Format("2006-01-02")] = true
		}
	}

	updateCache(year, res)
	return res, nil
}

// WorkCalendarToJson 生成指定年份的节假日信息到给定的json文件。
func WorkCalendarToJson(year int, fp string) {
	if b, err := WorkCalendar(year); err == nil {
		utils.MapToJsonFile(b, fp)
	} else {
		log.Println(err)
		return
	}

}

// parseHolidayDate 解析国务院具体的放假安排，具体到每一天，false为放假。true为工作日。
func parseHolidayDate(in string, dm map[string]bool) error {
	r := regexp.MustCompile(`\d{0,4}年?\d+月\d+日至?\d{0,2}月?\d{0,2}日?`)
	t := r.FindAllString(in, -1)

	for i := range t {
		t[i] = strings.ReplaceAll(t[i], "月", "-")
		t[i] = strings.ReplaceAll(t[i], "日", "")
	}
	var tmp []string

	if strings.Contains(t[0], "年") {
		tmp = strings.Split(strings.Split(t[0], "年")[1], "至")
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
		return err
	}

	et, err := time.Parse("2006-01-02", end)
	if err != nil {
		return err
	}

	for i := st; et.Sub(i).Hours() >= 0; i = i.Add(24 * time.Hour) {
		dm[i.Format("2006-01-02")] = false
	}

	if len(t) > 1 {
		for i := 1; i < len(t); i++ {
			dt := fmt.Sprintf("%d-%s", _year, t[i])
			dt = utils.ModifyDateFormat(dt)
			dm[dt] = true
		}
	}
	return nil
}

func check(year int) bool {
	if _, ok := globals.WorkCalendarMap[year]; ok {
		return true
	}
	return false
}

func readCache(year int) map[string]bool {
	return globals.WorkCalendarMap[year]
}

func updateCache(year int, r map[string]bool) {

	if len(r) != 0 {
		globals.WorkCalendarMap[year] = r
		cache.UpdateCalendar()
	}
}
