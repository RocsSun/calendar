/*
Package cmd
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"gitee.com/RocsSun/calendar/calendar/shares"
	"time"

	"github.com/spf13/cobra"
)

var (
	_year int
	fp    string
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "获取股市的开闭市日历",
	Long: `通过命令行生成一年中的节假日信息。支持生成2007年之后的年份的节假日和股市交易日历。
通过share获取股市的开闭市日历
参数支持
	share:
		-y 2022 获取指定的年份的日历
		-j path/to/file.json 将节假日日历保存到指定的json文件。
	不跟参数是为当前年份的。
详情参照example。`,
	Run: func(cmd *cobra.Command, args []string) {
		sharesCalendar(_year, fp)
	},
	Example: `
	calendar share 生成当前年份股市的开闭市信息。
	calendar share -j "./2022.json" 生成当前年份股市的开闭市信息，并保存到"./2022.json"文件中。
	calendar share -y 2022 生成2022年份股市的开闭市信息。
	calendar share -y 2022 -j "./2022.json" 生成2022年份股市的开闭市信息，并保存到"./2022.json"文件中。
`,
}

func init() {
	rootCmd.AddCommand(shareCmd)

	// Here you will define your flags and configuration settings.
	shareCmd.Flags().IntVarP(&_year, "指定年份", "y", time.Now().Year(), "指定年份的节假日")
	shareCmd.Flags().StringVarP(&fp, "指定文件", "j", "", "指定输出的json文件。")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sharesCalendar(year int, f string) {
	if f == "" {
		fmt.Println(shares.ShareTradeCalendar(year))
	} else {
		shares.ShareTradeCalendarToJson(year, f)
	}
}
