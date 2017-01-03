
package alarm

import "hw"
import "time"
import "sync"
import "config"

var AbortWakeUpKey int;

type Alarm struct {
	WakeUpAt time.Time
	AlarmIsSet bool
	Aborting bool
	AlarmIsActive bool
	LEDs hw.HWInterface
	Lock sync.Mutex
	UI UIState
	Config *config.Config
}

func NewAlarm(h hw.HWInterface, c * config.Config) *Alarm {
	var a Alarm;
	if (c == nil) {
		c = config.LoadConfig();
	}
	a.Config = c;
	a.LEDs = h
	h.RegisterButtonCallback(a.ButtonPress);
	h.SetStatus(hw.SOFTWARE_ALIVE_LED, true);
	h.SetStatus(hw.CAL_SYNC_LED, false)
	h.SetStatus(hw.CAL_ARMED_LED, false);
	h.SetStatus(hw.ERROR_LED, false);

	go a.AlarmLoop();
	go a.PanicHandler();
	return &a
}

/* Runs as a goroutine that will turn on the error light if we crash */
func (self *Alarm) PanicHandler() {
	x := make(chan bool);
	<- x
	defer func() {
		self.LEDs.SetStatus(hw.ERROR_LED, true);
	}()
}

/* This runs as a goroutine to constantly pool the current time
 * and set of the alarm */
func (self * Alarm) AlarmLoop() {
	for ;; {
		time.Sleep(1 * time.Second);
		self.Lock.Lock();
		now := time.Now();
		if (self.AlarmIsSet && now.After(self.WakeUpAt) && !self.UI.ForceAlarmOff) {
			/* Wake up time! Ring the alarm! */
			self.AlarmIsSet = false;
			self.AlarmIsActive = true;
			go self.WakeUp();
		}
		/* Set the alarm armed status LED */
		self.LEDs.SetStatus(hw.CAL_ARMED_LED, self.AlarmIsSet && !self.UI.ForceAlarmOff);
		self.Lock.Unlock();
	}
}

/*
 * This goroutine runs to actually do the alarm. It makes the
 * light go bring and (eventually) will play music!
 * it pools the "Checkin" function to give up if the alarm
 * is canceled.
 */
func (self * Alarm) WakeUp() {
	defer func() {
		/* Checkin will panic with &AbortWakeUpKey
		 * if the alarm to be canceled.
		 * If we panic with that, just terminate the goroutine,
		 * otherwise, we goofed and should rethrow the panic
		 */
		r := recover();
		if (r != &AbortWakeUpKey && r != nil)  {
			panic(r);
		}
	}()

	/* Slowly make the lights brighter */
	print("DINGDINGDING");
	for i := 0; i < 16; i++ {
		self.Checkin();
		print("BRIGHTER!")
		self.LEDs.SetLEDs(i,i,i);
		time.Sleep(1 * time.Second);
	}

}

/* Check if the alarm has been canceled and panic if it has */
func (self * Alarm) Checkin() {
	self.Lock.Lock();
	defer self.Lock.Unlock();
	if (self.Aborting == true) {
		self.Aborting = false;
		panic(&AbortWakeUpKey)
	}
}


/* Set an alarm for a given time */
func (self * Alarm) SetAlarm(wake time.Time) {
	self.Lock.Lock();
	defer self.Lock.Unlock();

	self.AlarmIsSet = true;
	self.WakeUpAt = wake
}



func (self * Alarm) AbortAlarmInProgress() {
	//self.Lock.Lock();
	//defer self.Lock.Unlock();
	if (self.AlarmIsActive) {
		self.Aborting = true;
	}
}

