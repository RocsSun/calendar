package duration_test

import (
	"fmt"
	"gitee.com/RocsSun/calendar/calendar/duration"
	"testing"
)

func TestNewCountTime(t *testing.T) {
	start := "2021-12-31 08:00"
	end := "2022-01-05 08:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"

	cli := duration.NewCountTime(start, end, as, ae, ps, pe)
	if cli == nil {
		t.Error("NewCountTime nil.")
	}
	if cli.EffectTime() != 2.0 {

		fmt.Println(cli.EffectTime())
		t.Error("count filed! ")
	}
}

func TestNewCountTime1(t *testing.T) {
	start := "2021-12-31 08:00"
	end := "2022-01-05 13:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"

	cli := duration.NewCountTime(start, end, as, ae, ps, pe)
	if cli == nil {
		t.Error("NewCountTime nil.")
	}

	if cli.EffectTime() == float64(3.43) {
		t.Error("count filed! ")
	}
}

func TestNewCountTime3(t *testing.T) {
	start := "2021-12-31 08:00"
	end := "2022-01-06 13:00"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"

	cli := duration.NewCountTime(start, end, as, ae, ps, pe)
	if cli == nil {
		t.Error("NewCountTime nil.")
	}

	fmt.Println(cli.EffectTime())
	if cli.EffectTime() != 3.43 {
		t.Error("count filed! ")
	}
}
