package timer

import "time"

func GetNowTime() time.Time {
	// 设置当前时区位 "Asia/shanghai"
	location, _ := time.LoadLocation("Asia/shanghai")
	return time.Now().In(location)
}

func GetCalculateTime(currentTimer time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}
	return currentTimer.Add(duration), nil
}
