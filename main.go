package main

import (
	"alarm_clock/browser"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	h    bool //help
	uriIndex    int //uri index
	v, V bool
	t, T string
)

const version = "1.0.0"
const timeFmt = "15:04"

var logger *log.Logger

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.StringVar(&t, "t", "06:00", "set alarm time")
}

// If later hasn't happened yet, make it happen on the day of now; if not, the day after.
func bestTime(now time.Time, later time.Time) time.Time {
	now = now.Local() // use local time to make things make sense
	nowh, nowm, nows := now.Clock()
	laterh, laterm, laters := later.Clock()
	add := false
	if nowh > laterh {
		add = true
	} else if (nowh == laterh) && (nowm > laterm) {
		add = true
	} else if (nowh == laterh) && (nowm == laterm) && (nows >= laters) {
		// >= in the case we're on the exact second; add a day because the alarm should have gone off by now otherwise!
		add = true
	}
	if add {
		now = now.AddDate(0, 0, 1)
	}
	return time.Date(now.Year(), now.Month(), now.Day(),
		laterh, laterm, laters, 0,
		now.Location())

}
func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	if v {
		fmt.Printf("version:%s", version)
		return
	}
	alarmTime, err := time.Parse(timeFmt, t)
	if err != nil {
		panic(err)
	}
	fmt.Println("alarmtime", alarmTime)
	file, err := os.OpenFile("alarm.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("fail to create alarm.log file!")
	}
	defer file.Close()
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)
	for {
		logger.Println("started !")
		now := time.Now()
		later := bestTime(now, alarmTime)
		logger.Println("later",later)
		duration := later.Sub(now)
		logger.Println("duration",duration.String())
		StartTimer(duration)
		uriIndex++
		logger.Println("uri index ",uriIndex)
	}
}

func StartTimer(t time.Duration) {
	timer := time.NewTimer(t)
	for {
		select {
		case <-timer.C:
			Fire()
			return
		}
	}
	logger.Println("unreachable !")
	panic("unreachable") // just in case
}
func Fire() {
	data, err := ioutil.ReadFile("video.conf")
	if err != nil {
		logger.Println("error", err)
		panic(err)
		return
	}
	var uri string
	str := strings.Split(string(data), "\n")
	if len(str) == 0 {
		logger.Println("no data")
		panic("no data ")
		return
	}
	logger.Println(str)
	if uriIndex > len(str) -1 {
		logger.Println("reset uriIndex to 0")
		uriIndex = 0
	}
	uri = str[uriIndex]
	err = browser.Open(uri)
	logger.Println("opened uri", uri, err)
	return

}
