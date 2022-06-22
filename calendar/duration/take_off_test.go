package duration_test

import (
	"fmt"
	"testing"

	"github.com/RocsSun/calendar/calendar/duration"
)

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

func TestNewCountTime4(t *testing.T) {
	start := "2022-01-05 08:30"
	end := "2022-01-05 13:30"
	as := "08:30"
	ae := "11:30"
	ps := "13:00"
	pe := "17:00"
	tim := duration.EffectTimes(start, end, as, ae, ps, pe)
	fmt.Println(tim)
	if tim != 0.5 {
		t.Error("failed!")
	}
}
