package main

import (
	"log"
	"merkaba/chromedp/sysutil"
	"time"
)

func main() {
	b := sysutil.BootTime()
	log.Printf("boot: %s", b.Format(time.RFC3339))
	log.Printf("now: %s", time.Now().Format(time.RFC3339))
}
