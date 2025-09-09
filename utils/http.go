package utils

import (
	"time"

	"github.com/go-resty/resty/v2"
)

var Client *resty.Client

func init() {
	Client = resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)
}
