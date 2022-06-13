package shares

import (
	"fmt"
	"gitee.com/RocsSun/calendar/cache"
	"gitee.com/RocsSun/calendar/calendar/holiday"
	"gitee.com/RocsSun/calendar/constants"
	"gitee.com/RocsSun/calendar/utils"
	"time"
)

// ShareTradeCalendar 工作日历。false为休息。true为工作日。返回一年的所有的日期的节假日。
func ShareTradeCalendar(year int) map[string]bool {
	res := make(map[string]bool)

	if check(year) {
		readCache(year, res)
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

func readCache(year int, r map[string]bool) map[string]bool {
	st, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-01-01", year))
	ed, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-12-31", year))
	for i := st; ed.Sub(i).Hours() >= 0; i = i.Add(24 * time.Hour) {
		r[i.Format("2006-01-02")] = constants.ShareCalendarMap[i.Format("2006-01-02")]
	}
	return r
}

func updateCache(year int, r map[string]bool) {
	st, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-01-01", year))
	ed, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-12-31", year))
	for i := st; ed.Sub(i).Hours() >= 0; i = i.Add(24 * time.Hour) {
		constants.ShareCalendarMap[i.Format("2006-01-02")] = r[i.Format("2006-01-02")]
	}
	cache.UpdateCalendar()
}

func check(year int) bool {
	if _, ok := constants.ShareCalendarMap[fmt.Sprintf("%d-01-01", year)]; ok {
		return true
	}
	return false
}
