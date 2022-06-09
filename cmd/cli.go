/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"gitee.com/RocsSun/calendar/calendar/holiday"
	"gitee.com/RocsSun/calendar/calendar/shares"
	"github.com/spf13/cobra"
	"time"
)

var (
	_year = -9999
	fp    string
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "命令行操作直接生成。",
	Long: `通过命令行生成一年中的节假日信息。支持生成2007年之后的年份的节假日和股市交易日历。
	通过work获取正常调班的节假日
	通过trade 获取股市的交易日历
	参数支持
		-y 2022 获取指定的年份的日历
		-j path/to/file.json 将节假日日历保存到指定的json文件。
	不跟参数是为当前年份的。
	详情参照example。
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cli called")

		fmt.Println(_year)
		if args[0] == "work" {
			workCalendar(_year, fp)
		} else if args[0] == "share" {
			sharesCalendar(_year, fp)
		}
	},
	Args: cobra.ExactArgs(1),
	Example: `
	calendar cli work 生成当前年份的放假调班信息。
	calendar cli work -j "./2022.json" 生成当前年份的放假调班信息，并保存到"./2022.json"文件中。
	calendar cli work -y 2022 生成2022年份的放假调班信息。
	calendar cli work -y 2022 -j "./2022.json" 生成2022年份的放假调班信息，并保存到"./2022.json"文件中。
	calendar cli share 生成当前年份股市的放假调班信息
`,
}

func init() {
	rootCmd.AddCommand(cliCmd)

	// Here you will define your flags and configuration settings.

	cliCmd.Flags().IntVarP(&_year, "指定年份", "y", time.Now().Year(), "指定年份的节假日")
	cliCmd.Flags().StringVarP(&fp, "指定文件", "j", "", "指定输出的json文件。")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cliCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cliCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func workCalendar(year int, f string) {
	if f == "" {
		fmt.Println(holiday.WorkCalendar(year))
	} else {
		holiday.WorkCalendarToJson(year, f)
	}
}

func sharesCalendar(year int, f string) {
	if f == "" {
		fmt.Println(shares.ShareTradeCalendar(year))
	} else {
		shares.ShareTradeCalendarToJson(year, f)
	}
}
