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

	var w web.WebState
	w.Hw = h
	w.Config = c
	w.Webroot = os.Getenv("WEBROOT")
	w.GCal = google.NewCS(c)
	a := alarm.NewAlarm(h, w.GCal, c)
	w.Alarm = a

	web.StartServer(&w)
	x := make(chan bool)
	<-x
}
