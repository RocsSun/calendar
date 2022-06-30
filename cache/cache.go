package cache

import (
	"github.com/RocsSun/calendar/globals"
	"github.com/RocsSun/calendar/utils"
)

func InitCalendar() {
	if utils.IsFile(globals.CacheWorkCalendar) {
		utils.Decode(globals.CacheWorkCalendar, &globals.WorkCalendarMap)
	}

	if utils.IsFile(globals.CacheShareCalendar) {
		utils.Decode(globals.CacheShareCalendar, &globals.ShareCalendarMap)
	}
}

func UpdateCalendar() {
	if len(globals.WorkCalendarMap) != 0 {
		utils.Encoding(globals.CacheWorkCalendar, globals.WorkCalendarMap)
	}
	if len(globals.ShareCalendarMap) != 0 {
		utils.Encoding(globals.CacheShareCalendar, globals.ShareCalendarMap)
	}
}

func init() {
	InitCalendar()
}
