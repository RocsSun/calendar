package cache

import (
	"gitee.com/RocsSun/calendar/constants"
	"gitee.com/RocsSun/calendar/utils"
)

func InitCalendar() {
	if utils.IsFile(constants.CacheWorkCalendar) {
		utils.Decode(constants.CacheWorkCalendar, &constants.WorkCalendarMap)
	}

	if utils.IsFile(constants.CacheShareCalendar) {
		utils.Decode(constants.CacheShareCalendar, &constants.ShareCalendarMap)
	}
}

func UpdateCalendar() {
	if len(constants.WorkCalendarMap) != 0 {
		utils.Encoding(constants.CacheWorkCalendar, constants.WorkCalendarMap)
	}
	if len(constants.ShareCalendarMap) != 0 {
		utils.Encoding(constants.CacheShareCalendar, constants.ShareCalendarMap)
	}
}

func init() {
	InitCalendar()
}
