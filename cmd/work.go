/*
Package cmd
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/RocsSun/calendar/calendar/holiday"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// workCmd represents the work command
var workCmd = &cobra.Command{
	Use:   "work",
	Short: "获取法定节假日",
	Long: `通过命令行生成一年中的节假日信息。支持生成2007年之后的年份的节假日和股市交易日历。
通过work获取正常调班的节假日
参数支持
	work:
		-y 2022 获取指定的年份的日历
		-j path/to/file.json 将节假日日历保存到指定的json文件。
	不跟参数是为当前年份的。
详情参照example。`,
	Run: func(cmd *cobra.Command, args []string) {
		workCalendar(_year, fp)
	},
	Example: `
	calendar work 生成当前年份的放假调班信息。
	calendar work -j "./2022.json" 生成当前年份的放假调班信息，并保存到"./2022.json"文件中。
	calendar work -y 2022 生成2022年份的放假调班信息。
	calendar work -y 2022 -j "./2022.json" 生成2022年份的放假调班信息，并保存到"./2022.json"文件中。
`,
}

func init() {
	rootCmd.AddCommand(workCmd)

	// Here you will define your flags and configuration settings.
	workCmd.Flags().IntVarP(&_year, "指定年份", "y", time.Now().Year(), "指定年份的节假日")
	workCmd.Flags().StringVarP(&fp, "指定文件", "j", "", "指定输出的json文件。")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func workCalendar(year int, f string) {
	if f == "" {
		c, err := holiday.WorkCalendar(year)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(c)
		}
	} else {
		holiday.WorkCalendarToJson(year, f)
	}
}
