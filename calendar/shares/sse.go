package shares

import (
	"fmt"
	"github.com/RocsSun/calendar/cache"
	"github.com/RocsSun/calendar/calendar/holiday"
	"github.com/RocsSun/calendar/constants"
	"github.com/RocsSun/calendar/utils"
	"time"
)

// ShareTradeCalendar 工作日历。false为休息。true为工作日。返回一年的所有的日期的节假日。
func ShareTradeCalendar(year int) map[string]bool {
	res := make(map[string]bool)

	if check(year) {
		return readCache(year)
	}

	hol := holiday.GovHoliday(year)

	if hol == nil {
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
		if i.Weekday() == 0 || i.Weekday() == 6 {
			res[i.Format("2006-01-02")] = false
		} else if v, ok := hol[i.Format("2006-01-02")]; ok {
			res[i.Format("2006-01-02")] = v
		} else {
			res[i.Format("2006-01-02")] = true
		}
	}

	updateCache(year, res)
	return res
}

// PreTradeDay 前一个交易日。
func PreTradeDay(t time.Time) time.Time {
	var res map[string]bool
	if v, ok := constants.ShareCalendarMap[t.Year()]; !ok {
		res = ShareTradeCalendar(time.Now().Year())
	} else {
		res = v
	}

	i := t.Add(-24 * time.Hour)

	if i.Year() < t.Year() {
		ShareTradeCalendar(time.Now().Year() - 1)
	}

	for !res[i.Format("2006-01-02")] {
		i = i.Add(-24 * time.Hour)
		if i.Year() < t.Year() {
			ShareTradeCalendar(time.Now().Year() - 1)
		}
	}

	return i
}

// CurrentYearShareTradeCalendar 当前年份的节假日信息。
func CurrentYearShareTradeCalendar() map[string]bool {
	return ShareTradeCalendar(time.Now().Year())
}

// ShareTradeCalendarToJson 指定年份的股票交易日里生成json文件。
func ShareTradeCalendarToJson(year int, fp string) {
	utils.MapToJsonFile(ShareTradeCalendar(year), fp)
}

// CurrentYearShareTradeCalendarToJson 当前年份的股票交易日里生成json文件。
func CurrentYearShareTradeCalendarToJson(fp string) {
	utils.MapToJsonFile(ShareTradeCalendar(time.Now().Year()), fp)
}

func readCache(year int) map[string]bool {
	return constants.ShareCalendarMap[year]
}

func updateCache(year int, r map[string]bool) {

	if len(r) != 0 {
		constants.ShareCalendarMap[year] = r
		cache.UpdateCalendar()
	}
}

func check(year int) bool {
	if _, ok := constants.ShareCalendarMap[year]; ok {
		return true
	}
	return false
}
