/*
Package cmd
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/RocsSun/calendar/calendar/duration"
	"github.com/spf13/cobra"
)

var (
	_start   string
	_end     string
	_amStart string
	_amEnd   string
	_pmStart string
	_pmEnd   string
)

// effectTimeCmd represents the effectTime command
var effectTimeCmd = &cobra.Command{
	Use:   "time",
	Short: "计算有效的调休/请假时长。",
	Long: `计算两个时间日期中间的工作日请假或者是调休时间。
	time:
		--start YYYY-MM-DD hh:mm 开始时间
		--end YYYY-MM-DD hh:mm 结束时间
		--amStart 早上上班时间
		--amEnd 早上下班时间
		--pmStart 下午上班时间
		--pmEnd 下午下班时间
`,
	Run: func(cmd *cobra.Command, args []string) {
		if _start != "" && _end != "" && _amStart != "" && _amEnd != "" && _pmStart != "" && _pmEnd != "" {
			fmt.Println(duration.EffectTimes(_start, _end, _amStart, _amEnd, _pmStart, _pmEnd))
		}
	},
	Example: `calendar time --start YYYY-MM-DD hh:mm --end YYYY-MM-DD hh:mm --amStart hh:mm --amEnd hh:mm --pmStart hh:mm --pmEnd hh:mm 计算有效的请假或者调休时长。`,
}

func init() {
	rootCmd.AddCommand(effectTimeCmd)

	// Here you will define your flags and configuration settings.
	effectTimeCmd.Flags().StringVar(&_start, "start", "", "开始时间")
	effectTimeCmd.Flags().StringVar(&_end, "end", "", "结束时间")
	effectTimeCmd.Flags().StringVar(&_amStart, "amStart", "", "早上上班时间")
	effectTimeCmd.Flags().StringVar(&_amEnd, "amEnd", "", "早上下班时间")
	effectTimeCmd.Flags().StringVar(&_pmStart, "pmStart", "", "下午上班时间")
	effectTimeCmd.Flags().StringVar(&_pmEnd, "pmEnd", "", "下午下班时间")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// effectTimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// effectTimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
