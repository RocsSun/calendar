package duration_test

import (
	"testing"

	"github.com/RocsSun/calendar/calendar/duration"
)

// TestEffectTimes 测试不跨天的正常请假
func TestEffectTimes(t *testing.T) {
	start := "2022-01-05 08:30"
	end := "2022-01-05 17:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 1 {
		t.Error("failed!")
	}
}

// TestNewCountTime1 测试请假某几个小时。
func TestNewCountTime1(t *testing.T) {
	start := "2022-01-05 15:30"
	end := "2022-01-05 17:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)
	if tim != 0.21 {
		t.Error("failed!")
	}
}

// TestNewCountTime3 测试请假半天
func TestNewCountTime3(t *testing.T) {
	start := "2022-01-05 13:30"
	end := "2022-01-05 17:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 0.5 {
		t.Error("failed!")
	}
}

// TestNewCountTime4 测试请假半天，除去中午休息时间。
func TestNewCountTime4(t *testing.T) {
	start := "2022-01-05 08:30"
	end := "2022-01-05 13:30"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 0.5 {
		t.Error("failed!")
	}
}

// TestNewCountTime5 测试跨天。
func TestNewCountTime5(t *testing.T) {
	start := "2022-01-04 13:00"
	end := "2022-01-05 13:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 1 {
		t.Error("failed!")
	}
}

// TestNewCountTime6 测试跨天请假开始时间的小时数大于请假结束时间的小时数
func TestNewCountTime6(t *testing.T) {
	start := "2022-01-04 14:00"
	end := "2022-01-06 13:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 1.86 {
		t.Error("failed!")
	}
}

// TestNewCountTime7 测试请假结束时间是第二天早上上班之前。
func TestNewCountTime7(t *testing.T) {
	start := "2022-01-04 17:00"
	end := "2022-01-06 08:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)
	if tim != 1 {
		t.Error("failed!")
	}
}

// TestNewCountTime8 测试请假开始在早上上班之前，请假结束时间是第二天早上上班之前。
func TestNewCountTime8(t *testing.T) {
	start := "2022-01-04 07:40"
	end := "2022-01-06 08:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 2 {
		t.Error("failed!")
	}
}

// TestNewCountTime9 测试请假开始时间在中午休息时间段的。
func TestNewCountTime9(t *testing.T) {
	start := "2022-01-04 11:30"
	end := "2022-01-06 18:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 2.57 {
		t.Error("failed!")
	}
}

// TestNewCountTime10 测试重复修正时间。
func TestNewCountTime10(t *testing.T) {
	start := "2022-01-04 11:30"
	end := "2022-01-06 18:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	cli := duration.NewCountTime(start, end, as, ae, ps, pe)
	if cli == nil {
		t.Error("init cli failed! ")
	}

	tim := cli.EffectTime()
	if tim != 2.57 {
		t.Error("failed!")
	}

	tim = cli.EffectTime()
	if tim != 2.57 {
		t.Error("failed!")
	}
}

// TestNewCountTime11 测试跨年。
func TestNewCountTime11(t *testing.T) {
	start := "2021-12-31 08:30"
	end := "2022-01-05 13:30"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)

	if tim != 2.5 {
		t.Error("failed!")
	}
}
