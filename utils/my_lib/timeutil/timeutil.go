package timeutil

import (
	"time"
)

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}
}

func RFC3339TOUnixMilli(value string) int64 {
	ts, _ := time.Parse(time.RFC3339, value)
	InputCST := ts.In(cst).Format(CSTLayout)
	tt, _ := time.ParseInLocation(CSTLayout, InputCST, cst) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.UnixMilli()
}
