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

// countDays 计算天数
func (c *CountTime) countDays() float64 {
	var days = 0.0
	for i := c.offsetStart; c.offsetEnd.Sub(i).Hours() >= 24; i = i.Add(24 * time.Hour) {
		if constants.WorkCalendarMap[i.Year()][i.Format("2006-01-02")] {
			days++
		}
	}
	res, _ := decimal.NewFromFloat(days).Round(2).Float64()
	return res
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
	res, _ := decimal.NewFromFloat(c.countDays() + c.countHours()).Round(2).Float64()
	return res
}

func (c *CountTime) offset() {

	if c.isOffset {
		return
	}

	c.offsetStart = c.Start
	c.offsetEnd = c.End

	// 处理请假开始时间。
	// 处理请假开始时间在下午下班之后的情况。重置时间为第二天的早上上班时间。
	if c.offsetStart.Hour()*60+c.offsetStart.Minute() >= c.PmEnd.Hour()*60+c.PmEnd.Minute() {
		c.offsetStart = c.offsetStart.Add(
			time.Duration(
				24*60-c.offsetStart.Hour()*60-c.offsetStart.Minute()+c.AmStart.Minute()+c.AmStart.Hour()*60,
			) * time.Minute,
		)
	}

	// 处理请假开始时间在早上上班之前的，重置时间为当天早上上班时间。
	if c.offsetStart.Hour()*60+c.offsetStart.Minute() < c.AmStart.Hour()*60+c.AmStart.Minute() {
		c.offsetStart = c.offsetStart.Add(
			time.Duration((c.AmStart.Hour()-c.Start.Hour())*60+c.AmStart.Minute()-c.Start.Minute()) * time.Minute,
		)
	}

	// 处理请假开始时间在在早上下班之后（含），下午上班之前的（含），重置请假开始时间为当天下午上班时间。
	if c.offsetStart.Hour()*60+c.offsetStart.Minute() >= c.AmEnd.Hour()*60+c.AmEnd.Minute() &&
		c.offsetStart.Hour()*60+c.offsetStart.Minute() <= c.PmStart.Hour()*60+c.PmStart.Minute() {
		c.offsetStart = c.offsetStart.Add(
			time.Duration(c.PmStart.Hour()*60+c.PmStart.Minute()-(c.offsetStart.Hour()*60+c.offsetStart.Minute())) * time.Minute,
		)
	}

	// 处理请假结束时间。
	// 处理请假结束时间在上班之前的。重置的前一天的下午下班时间。
	if c.offsetEnd.Hour()*60+c.offsetEnd.Minute() <= c.AmStart.Hour()*60+c.AmStart.Minute() {
		c.offsetEnd = c.offsetEnd.Add(
			time.Duration(
				(c.PmEnd.Hour()-24-c.offsetEnd.Hour())*60+c.PmEnd.Minute()-c.offsetEnd.Minute(),
			) * time.Minute,
		)

		// 更新日期。
		c.AmStart = c.AmStart.Add(-24 * time.Hour)
		c.AmEnd = c.AmEnd.Add(-24 * time.Hour)
		c.PmStart = c.PmStart.Add(-24 * time.Hour)
		c.PmEnd = c.PmEnd.Add(-24 * time.Hour)
	}
	// 处理请假结束时间在下班之后的，重置为当天的下午下班时间。
	if c.offsetEnd.Hour()*60+c.offsetEnd.Minute() > c.PmEnd.Hour()*60+c.PmEnd.Minute() {
		c.offsetEnd = c.PmEnd
	}
	// 处理请假结束时间在早上下班之后（含），下午上班之前的（含），重置为当天早上下班时间。
	if (c.offsetEnd.After(c.AmEnd) || c.offsetEnd.Equal(c.AmEnd)) && (c.offsetEnd.Before(c.PmStart) || c.offsetEnd.Equal(c.PmStart)) {
		c.offsetEnd = c.AmEnd
	}

	// 处理请假开始时间大于等于请假结束时间，跨天请假的情况。
	if c.offsetStart.Before(c.offsetEnd) && c.offsetStart.Hour() >= c.offsetEnd.Hour() {
		c.balance = float64(c.PmEnd.Hour()-c.Start.Hour()) + float64(c.PmEnd.Minute()-c.Start.Minute())/60 + c.balance
		c.offsetStart = c.offsetStart.Add(
			time.Duration((24-c.offsetStart.Hour()+c.AmStart.Hour())*60+c.AmStart.Minute()-c.offsetStart.Minute()) * time.Minute,
		)
	}

	// 处理请假结束时间在早上下班之前，早上上班之后（含），下午上班之后前，下午下班之前（含），减去中午的休息时间。
	if c.offsetStart.Hour()*60+c.offsetStart.Minute() >= c.AmStart.Hour()*60+c.AmStart.Minute() &&
		c.offsetStart.Hour()*60+c.offsetStart.Minute() < c.AmEnd.Hour()*60+c.AmEnd.Minute() &&
		c.offsetEnd.Hour()*60+c.offsetEnd.Minute() > c.PmStart.Hour()*60+c.PmStart.Minute() &&
		c.offsetEnd.Hour()*60+c.offsetEnd.Minute() <= c.PmEnd.Hour()*60+c.PmEnd.Minute() {
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
