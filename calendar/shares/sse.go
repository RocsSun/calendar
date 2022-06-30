package shares

import (
	"errors"
	"fmt"
	"github.com/RocsSun/calendar/globals"
	"log"
	"time"

	"github.com/RocsSun/calendar/cache"
	"github.com/RocsSun/calendar/calendar/holiday"
	"github.com/RocsSun/calendar/utils"
)

// ShareTradeCalendar 工作日历。false为休息。true为工作日。返回一年的所有的日期的节假日。
func ShareTradeCalendar(year int) (map[string]bool, error) {
	res := make(map[string]bool)

	if check(year) {
		return readCache(year), nil
	}

	hol, err := holiday.GovHoliday(year)
	if err != nil {
		return nil, err
	}
	if hol == nil {
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

	for i := st; et.Sub(i) >= 0; i = i.Add(24 * time.Hour) {
		if i.Weekday() == 0 || i.Weekday() == 6 {
			res[i.Format("2006-01-02")] = false
		} else if v, ok := hol[i.Format("2006-01-02")]; ok {
			res[i.Format("2006-01-02")] = v
		} else {
			res[i.Format("2006-01-02")] = true
		}
	}

	updateCache(year, res)
	return res, nil
}

// PreTradeDay 前一个交易日。
func PreTradeDay(t time.Time) (time.Time, error) {
	var res map[string]bool
	var err error
	if v, ok := globals.ShareCalendarMap[t.Year()]; !ok {
		res, err = ShareTradeCalendar(time.Now().Year())
		if err != nil {
			return time.Time{}, err
		}
	} else {
		res = v
	}

	i := t.Add(-24 * time.Hour)
	j := t

	if i.Year() < t.Year() {
		res, err = ShareTradeCalendar(time.Now().Year() - 1)
		if err != nil {
			return time.Time{}, err
		}
	}

	for !res[i.Format("2006-01-02")] {
		if i.Year() < j.Year() {
			j = i

			res, err = ShareTradeCalendar(time.Now().Year() - 1)
			if err != nil {
				return time.Time{}, err
			}
		}
		i = i.Add(-24 * time.Hour)
	}

	return i, nil
}

// ShareTradeCalendarToJson 指定年份的股票交易日里生成json文件。
func ShareTradeCalendarToJson(year int, fp string) {
	if b, err := ShareTradeCalendar(year); err != nil {
		log.Println(err)
		return
	} else {
		utils.MapToJsonFile(b, fp)
	}
}

func readCache(year int) map[string]bool {
	return globals.ShareCalendarMap[year]
}

func updateCache(year int, r map[string]bool) {

	if len(r) != 0 {
		globals.ShareCalendarMap[year] = r
		cache.UpdateCalendar()
	}
}

func check(year int) bool {
	if _, ok := globals.ShareCalendarMap[year]; ok {
		return true
	}
	return false
}
