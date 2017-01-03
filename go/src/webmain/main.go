package main

import "os"
import "alarm"
import "config"
import "hw"
import "web"
import "google"

func main() {
	h, _ := hw.NewHW()
	c := config.LoadConfig()
	a := alarm.NewAlarm(h, c)

	var w web.WebState
	w.Hw = h
	w.Alarm = a
	w.Config = c
	w.Webroot = os.Getenv("WEBROOT")
	w.GCal = google.NewCS(c)

	web.StartServer(&w)
	x := make(chan bool)
	<-x
}
