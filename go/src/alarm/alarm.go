
package alarm

import "hw"
import "time"
import "sync"

var AbortWakeUpKey int;

type Alarm struct {
	WakeUpAt time.Time
	AlarmIsSet bool
	Aborting bool
	LEDs *hw.HW
	Lock sync.Mutex
}

func NewAlarm(h *hw.HW) *Alarm {
	var a Alarm;
	a.LEDs = h
	go a.AlarmLoop();
	return &a
}

func (self * Alarm) AlarmLoop() {
	for ;; {
		time.Sleep(1 * time.Second);
		self.Lock.Lock();
		now := time.Now();
		if (self.AlarmIsSet && now.After(self.WakeUpAt)) {
			self.AlarmIsSet = false;
			go self.WakeUp();
		}
		self.Lock.Unlock();
	}
}

func (self * Alarm) WakeUp() {
	defer func() {
		r := recover();
		if (r != &AbortWakeUpKey && r != nil)  {
			panic(r);
		}
	}()

	print("DINGDINGDING");
	for i := 0; i < 16; i++ {
		self.Checkin();
		print("BRIGHTER!")
		self.LEDs.SetLEDs(i,i,i);
		time.Sleep(1 * time.Second);
	}

}

func (self * Alarm) Checkin() {
	self.Lock.Lock();
	defer self.Lock.Unlock();
	if (self.Aborting == true) {
		self.Aborting = false;
		panic(&AbortWakeUpKey)
	}
}


func (self * Alarm) SetAlarm(wake time.Time) {
	self.Lock.Lock();
	defer self.Lock.Unlock();

	self.AlarmIsSet = true;
	self.WakeUpAt = wake
}



func (self * Alarm) AbortAlarmInProgress() {
	self.Lock.Lock();
	defer self.Lock.Unlock();
	self.Aborting = true;
}

