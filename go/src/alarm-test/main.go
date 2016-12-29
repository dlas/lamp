

package main

import "hw"
import "time"
import "alarm"

func main() {
	h, _ := hw.NewHW();
	a:= alarm.NewAlarm(h);

	print("start");
	time.Sleep(5 * time.Second);
	print("set alarm");
	a.SetAlarm(time.Now().Add( 20 * time.Second));
	time.Sleep(25 * time.Second);
	a.AbortAlarmInProgress()
	time.Sleep(25 * time.Second);
}
