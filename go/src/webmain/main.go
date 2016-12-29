
package main

import "os"
import "alarm"
import "config"
import "hw"
import "web"


func main() {
	h, _ := hw.NewHW();
	a := alarm.NewAlarm(h);
	c := config.LoadConfig();

	var w web.WebState;
	w.Hw = h;
	w.Alarm = a;
	w.Config = c;
	w.Webroot=os.Getenv("WEBROOT");

	web.StartServer(&w);
	x := make(chan bool);
	<- x;
}




