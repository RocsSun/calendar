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
	isOffset    bool      // 标记是否处理过
}

//var dateMap map[string]bool

// countDays 计算天数
func (c *CountTime) countDays() float64 {
	var days = 0.0
	for i := c.offsetStart; c.offsetEnd.Sub(i).Hours() >= 24; i = i.Add(24 * time.Hour) {
		if constants.WorkCalendarMap[i.Year()][i.Format("2006-01-02")] {
			days++
		}
	}

	return days
}

// countHours 计算请假的时间。
func (c *CountTime) countHours() float64 {
	effectTime := float64(c.offsetEnd.Hour()-c.offsetStart.Hour())*60 + float64(c.offsetEnd.Minute()-c.offsetStart.Minute()) + c.balance*60
	res, _ := decimal.NewFromFloat(effectTime / (c.AmEnd.Sub(c.AmStart).Minutes() + c.PmEnd.Sub(c.PmStart).Minutes())).Round(2).Float64()
	return res
}

// EffectTime 计算有效时间。
func (c *CountTime) EffectTime() float64 {
	c.offset()
	return c.countDays() + c.countHours()
}

func (c *CountTime) offset() {

	if c.isOffset {
		return
	}

	c.offsetStart = c.Start
	c.offsetEnd = c.End

	// 处理请假开始时间。
	// 处理请假开始时间在下午下班之后的情况。重置时间为第二天的早上上班时间。
	if c.Start.Hour() >= c.PmEnd.Hour() && c.Start.Minute() >= c.PmEnd.Minute() {
		c.offsetStart = c.Start.Add(
			time.Duration(24-c.Start.Hour()+c.AmStart.Hour()) * time.Hour).Add(
			time.Duration(c.AmStart.Minute()-c.Start.Minute()) * time.Minute)
	}

	// 处理请假开始时间在早上上班之前的，重置时间为当天早上上班时间。
	if c.Start.Hour() <= c.AmStart.Hour() && c.Start.Minute() < c.AmStart.Minute() {
		c.offsetStart = c.Start.Add(
			time.Duration(c.AmStart.Hour()-c.Start.Hour()) * time.Hour).Add(
			time.Duration(c.AmStart.Minute()-c.Start.Minute()) * time.Minute)
	}

	// 处理请假开始时间在在早上下班之后（含），下午上班之前的（含），重置请假开始时间为当天下午上班时间。
	if c.Start.Hour() >= c.AmEnd.Hour() && c.Start.Minute() >= c.AmEnd.Minute() && c.Start.Hour() < c.PmStart.Hour() && c.Start.Minute() < c.PmStart.Minute() {
		c.offsetStart = c.Start.Add(
			time.Duration(c.PmStart.Hour()-c.Start.Hour()) * time.Hour,
		).Add(
			-time.Duration(c.PmStart.Minute()-c.Start.Minute()) * time.Minute,
		)
	}

	// 处理请假结束时间。
	// 处理请假结束时间在上班之前的。重置的前一天的下午下班时间。
	if c.End.Hour() <= c.AmStart.Hour() && c.End.Minute() <= c.Start.Minute() {
		c.offsetEnd = c.End.Add(
			time.Duration(-24+c.PmEnd.Hour()-c.End.Hour()) * time.Hour,
		).Add(
			time.Duration(c.PmEnd.Minute()-c.End.Minute()) * time.Minute,
		)

		// 更新日期。
		c.AmStart = c.AmStart.Add(-24 * time.Hour)
		c.AmEnd = c.AmEnd.Add(-24 * time.Hour)
		c.PmStart = c.PmStart.Add(-24 * time.Hour)
		c.PmEnd = c.PmEnd.Add(-24 * time.Hour)
	}
	// 处理请假结束时间在下班之后的，重置为当天的下午下班时间。
	if c.End.Hour() >= c.PmEnd.Hour() && c.End.Minute() > c.PmEnd.Minute() {
		c.offsetEnd = c.PmEnd
	}
	// 处理请假结束时间在早上下班之后（含），下午上班之前的（含），重置为当天早上下班时间。
	if (c.End.After(c.AmEnd) || c.End.Equal(c.AmEnd)) && (c.End.Before(c.PmStart) || c.End.Equal(c.PmStart)) {
		c.offsetEnd = c.AmEnd
	}

	// 处理请假开始时间
	if c.offsetStart.Before(c.offsetEnd) && c.offsetStart.Hour() >= c.offsetEnd.Hour() {
		c.balance = float64(c.PmEnd.Hour()-c.Start.Hour()) + float64((c.PmEnd.Minute()-c.Start.Minute())/60) + c.balance
		// 前置一天。
		c.offsetStart = c.offsetStart.Add(
			time.Duration(24 - c.offsetStart.Hour() + c.AmStart.Hour()),
		).Add(
			time.Duration(c.AmStart.Minute()-c.offsetStart.Minute()) * time.Minute,
		)
	}

	// 处理请假结束时间在早上下班之前，早上上班之后（含），下午上班之后前，下午下班之前（含），减去中午的休息时间。
	if float64(c.offsetStart.Hour())+float64(c.offsetStart.Minute()/60) >= float64(c.AmStart.Hour())+float64(c.AmStart.Minute()/60) &&
		float64(c.offsetStart.Hour())+float64(c.offsetStart.Minute()/60) < float64(c.AmEnd.Hour())+float64(c.AmEnd.Minute()/60) &&
		float64(c.offsetEnd.Hour())+float64(c.offsetEnd.Minute()/60) > float64(c.PmStart.Hour())+float64(c.PmStart.Minute()/60) &&
		float64(c.offsetEnd.Hour())+float64(c.offsetEnd.Minute()/60) <= float64(c.PmEnd.Hour())+float64(c.PmEnd.Minute()/60) {
		//c.balance = c.balance - float64(c.PmStart.Hour()-c.AmEnd.Hour()) + float64(c.PmStart.Minute()-c.AmEnd.Minute())/60
		c.balance = c.balance - c.PmStart.Sub(c.AmEnd).Minutes()/60
	}

	c.isOffset = true
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
