package accesslog

import "time"

// Record 保留 operateLog 字段协议，但只承载 Gateway 已脱敏的稳定事实。
type Record struct {
	EventID        string    `json:"eventId"`
	SendTime       time.Time `json:"sendTime"`
	Topic          string    `json:"topic"`
	RequestURL     string    `json:"requestUrl"`
	RequestMethod  string    `json:"requestMethod"`
	Success        bool      `json:"success"`
	SimpleErrorMsg string    `json:"simpleErrorMsg"`
	DetailErrorMsg string    `json:"detailErrorMsg"`
	RequestIP      string    `json:"requestIp"`
	UserAgent      string    `json:"userAgent"`
	RequestParams  string    `json:"requestParams"`
	ResponseData   string    `json:"responseData"`
	ExecutionTime  int64     `json:"executionTime"`
	UserID         uint64    `json:"userId"`
	Operator       string    `json:"operName"`
	Platform       string    `json:"operatePlatform"`
}
