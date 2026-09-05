package cron

type CronJob struct {
	URL      string `json:"url"`
	ChatId   int64  `json:"chatId"`
	JobName  string `json:"jobName"`
	Schedule string `json:"schedule"`
}
