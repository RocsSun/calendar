package main

import (
	"fmt"
	"gitee.com/RocsSun/calendar/calendar/holiday"
)

func main() {
	res := holiday.WorkCalendar(2021)
	//res := holiday.CurrentYearWorkCalendar()
	//holiday.CurrentYearWorkCalendarToJson("./2022.json")

	for k, v := range res {
		fmt.Println(k, v)
	}
}
