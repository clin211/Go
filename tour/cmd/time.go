package cmd

import (
	"log"
	"strconv"
	"strings"
	"time"
	"tour/internal/timer"

	"github.com/spf13/cobra"
)

var calculateTime, duration string

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间格式处理",
	Long:  "时间格式处理",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("result: %s, %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "获取当前时间",
	Long:  "获取当前时间",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("result: %s, %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

var calculateTimeCmd = &cobra.Command{
	Use:   "calc",
	Short: "计算所需时间",
	Long:  "计算所需时间",
	Run: func(cmd *cobra.Command, args []string) {
		var currentTimer time.Time
		var layout = "2006-01-02 15:04:05"

		if calculateTime == "" {
			currentTimer = timer.GetNowTime()
		} else {
			var err error
			if !strings.Contains(calculateTime, "") {
				layout = "2006-01-02"
			}
			currentTimer, err = time.Parse(layout, calculateTime)
			if err != nil {
				t, _ := strconv.Atoi(calculateTime)
				currentTimer = time.Unix(int64(t), 0)
			}
		}

		calculateTime, err := timer.GetCalculateTime(currentTimer, duration)

		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v", err)
		}

		log.Printf("result: %s, %d", calculateTime.Format(layout), calculateTime.Unix())
	},
}

var formatTimestamp = &cobra.Command{
	Use:   "format",
	Short: "时间戳格式化工具",
	Long:  "将时间戳格式化为当前时区的标准时间",
	Run: func(cmd *cobra.Command, args []string) {
		timestamp := args[0]

		// 解析时间戳的位数
		timestampLen := len(timestamp)

		var t time.Time
		switch timestampLen {
		case 10:
			// 解析10位时间戳
			seconds, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			t = time.Unix(seconds, 0)
		case 13:
			// 解析13位时间戳
			milliseconds, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			seconds := milliseconds / 1000
			t = time.Unix(seconds, milliseconds%1000*int64(time.Millisecond))
		case 16:
			// 解析16位时间戳
			microseconds, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			seconds := microseconds / 1000000
			t = time.Unix(seconds, (microseconds%1000000)*int64(time.Microsecond))
		default:
			log.Fatal("无效的时间戳位数")
		}

		// 格式化为当前时区的标准时间
		localTime := t.Local()

		// 打印格式化后的时间
		log.Println("格式化后的时间:", localTime.Format("2006-01-02 15:04:05"))
	},
}

func init() {
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calculateTimeCmd)
	timeCmd.AddCommand(formatTimestamp)

	calculateTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", `需要计算的时间，有效单位为时间戳或已格式化的时间`)
	calculateTimeCmd.Flags().StringVarP(&duration, "duration", "d", "", `持续时间，有效单位为: "ns", "us" (or "μs"), "ms", "s", "m", "h"`)
}
