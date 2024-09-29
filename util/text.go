package util

import (
	"fmt"
	"strconv"
	"time"
)

func CreateOrderNo() string {
	now := time.Now()
	orderNo := fmt.Sprintf("TEDT-%s%s", now.Format("20060102"), strconv.Itoa(int(now.Unix())))
	return orderNo
}
