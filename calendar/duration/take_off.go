package duration

import (
	"fmt"
	"strings"
	"time"

	"github.com/RocsSun/calendar/calendar/holiday"
	"github.com/RocsSun/calendar/constants"
	"github.com/shopspring/decimal"
)

// CountTime 计算有效的请假时间，调休时间。
type CountTime struct {
	Start       time.Time // 请假开始时间
	End         time.Time // 请假结束时间
	AmStart     time.Time // 早上上班时间
	AmEnd       time.Time // 早上下班时间
	PmStart     time.Time // 下午上班时间
	PmEnd       time.Time // 下午下班时间
	offsetStart time.Time // 校验后的请假开始时间
	offsetEnd   time.Time // 校验后的请假结束时间
	balance     float64   // 中午休息时间
}

//var dateMap map[string]bool

// countDays 计算天数
func (c *CountTime) countDays() float64 {
	var days = 0.0
	for i := c.Start; c.End.Sub(i).Hours() >= 24; i = i.Add(24 * time.Hour) {
		if constants.WorkCalendarMap[i.Year()][i.Format("2006-01-02")] {
			days++
		}
	}

	return days
}

// countHours 计算请假的时间。
func (c *CountTime) countHours() float64 {
	var effectTime = 0.0
	var balance = 0.0

	// 请假时间早于上班时间
	if c.AmStart.Sub(c.Start) > 0 {
		c.Start = c.AmStart
	}

	// 请假时间晚于下班时间
	if c.PmEnd.Sub(c.End) < 0 {
		c.End = c.PmEnd
	}

	// 请假结束时间到第二天早上的
	if c.AmStart.Sub(c.End).Minutes() > 0 {
		c.End = c.AmStart
	}

	// 请假开始时间在第二天早上上班之前的。
	if c.PmEnd.Sub(c.Start) < 0 {
		c.Start = c.AmStart
	}

	if c.Start.Sub(c.End) > 0 {
		panic("时间有误。")
	}

	// 请假时间开始时间中午上班之前，减去中午休息时间
	if c.PmStart.Sub(c.Start) > 0 {
		balance = -c.PmStart.Sub(c.AmEnd).Minutes()
	}

	// 请假时间开始时间中午上班之前，加上中午休息时间
	if c.PmStart.Sub(c.End) > 0 {
		balance += c.PmStart.Sub(c.AmEnd).Minutes()
	}

	effectTime = c.End.Sub(c.Start).Minutes() + balance
	res, _ := decimal.NewFromFloat(effectTime / (c.AmEnd.Sub(c.AmStart).Minutes() + c.PmEnd.Sub(c.PmStart).Minutes())).Round(2).Float64()
	return res
}

// EffectTime 计算有效时间。
func (c *CountTime) EffectTime() float64 {
	return c.countDays() + c.countHours()
}

func (c *CountTime) offset() {

	// 处理请假开始时间。
	// 处理请假开始时间在下午下班之后的情况。重置时间为第二天的早上上班时间。
	if c.Start.Hour() >= c.PmEnd.Hour() && c.Start.Minute() >= c.PmEnd.Minute() {
		c.offsetStart = c.Start.Add(
			time.Duration(24-c.Start.Hour()+c.AmStart.Hour()) * time.Hour).Add(
			time.Duration(c.Start.Minute()+c.AmStart.Minute()) * time.Minute)
	}

	// 处理请假开始时间在早上上班之前的，重置时间为当天早上上班时间。
	if c.Start.Hour() <= c.AmStart.Hour() && c.Start.Minute() < c.AmStart.Minute() {
		c.offsetStart = c.Start.Add(
			time.Duration(c.AmStart.Hour()-c.Start.Hour()) * time.Hour).Add(
			-time.Duration(c.AmStart.Minute()-c.Start.Minute()) * time.Minute)
	}

	// 处理请假开始时间在在早上下班之后（含），下午上班之前的（含），重置请假开始时间为当天下午上班时间。

	// 处理请假结束时间。
	// 处理请假结束时间在上班之前的。重置的前一天的下午下班时间。

	// 处理请假结束时间在下班之后的，重置为当天的下午下班时间。

	// 处理请假结束时间在早上下班之后（含），下午上班之前的（含），重置为当天早上下班时间。
}

// NewCountTime 生成实例。
func NewCountTime(start, end, amStart, amEnd, pmStart, pmEnd string) *CountTime {

	if len(start) != 16 && len(end) != 16 && len(amStart) != 5 && len(amEnd) != 6 && len(pmStart) != 5 && len(pmEnd) != 6 {
		return nil
	}

	var parse = func(s string) time.Time {
		t, err := time.Parse("2006-01-02 15:04", s)
		if err != nil {
			fmt.Println(err)
		}
		return t
	}

	st := parse(start)
	ed := parse(end)
	dateTmp := strings.Split(end, " ")[0]

	if _, ok := constants.WorkCalendarMap[st.Year()]; !ok {
		holiday.WorkCalendar(st.Year())
	}
	if _, ok := constants.WorkCalendarMap[ed.Year()]; !ok {
		holiday.WorkCalendar(ed.Year())
	}

	if _, ok := constants.WorkCalendarMap[st.Year()]; !ok {
		fmt.Println("未生成该年份的节假日日历。")
		return nil
	}

	return &CountTime{
		Start:   st,
		End:     ed,
		AmStart: parse(fmt.Sprintf("%s %s", dateTmp, amStart)),
		AmEnd:   parse(fmt.Sprintf("%s %s", dateTmp, amEnd)),
		PmStart: parse(fmt.Sprintf("%s %s", dateTmp, pmStart)),
		PmEnd:   parse(fmt.Sprintf("%s %s", dateTmp, pmEnd)),
	}
}

func EffectTimes(start, end, amStart, amEnd, pmStart, pmEnd string) float64 {
	cli := NewCountTime(start, end, amStart, amEnd, pmStart, pmEnd)
	return cli.EffectTime()
}
