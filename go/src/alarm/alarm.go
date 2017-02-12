/* This manages the scheduling and "sounding" of alarms.
 * By historical accident, this is also the "main" coordination point
 * for just about everything
 */

package alarm

import "hw"
import "time"
import "sync"
import "config"
import "google"
import "log"
import "math/rand"
import "io/ioutil"

//import "os"
import "os/exec"

var AbortWakeUpKey int

/* What do we need to manage alarms, and lights in general? */
type Alarm struct {
	WakeUpAt      time.Time      // When is alarm set for
	AlarmIsSet    bool           // Is alarm set?
	Aborting      bool           // Are we in the process of aborting an active alarm
	AlarmIsActive bool           // Is the alarm going off right now?
	LEDs          hw.HWInterface // Pointer to HW implementation
	Lock          sync.Mutex
	UI            UIState        // see uishim.go
	Config        *config.Config // All our config info
	CS            *google.CalendarState
	Player        *exec.Cmd
}

// THis is the main initialization proceedure */
func NewAlarm(h hw.HWInterface, cs *google.CalendarState, c *config.Config) *Alarm {
	var a Alarm
	if c == nil {
		c = config.LoadConfig()
	}
	a.Config = c
	a.LEDs = h
	a.CS = cs
	/* Callback every time a physical button is pressed */
	h.RegisterButtonCallback(a.ButtonPress)

	/* Put the LEDs in their default state */

	h.SetStatus(hw.SOFTWARE_ALIVE_LED, true)
	h.SetStatus(hw.CAL_SYNC_LED, false)
	h.SetStatus(hw.CAL_ARMED_LED, false)
	h.SetStatus(hw.ERROR_LED, false)

	go a.AlarmLoop()
	go a.PanicHandler()
	if cs != nil {
		go a.SyncCalendarLoop()
	}
	return &a
}

/*
 * Loop forever, synchronizing our alarm to google calendar
 */
func (a *Alarm) SyncCalendarLoop() {
	a.LEDs.SetStatus(hw.CAL_SYNC_LED, false)
	for ; ; time.Sleep(30 * time.Minute) {
		/* Scan the next 24 hours for a "wakeup" time */
		now := time.Now().Local()
		start_scan := now.Add(90 * time.Minute)
		end_scan := start_scan.Add(24 * time.Hour)

		wakeup, err := a.CS.GetNextWakeup(start_scan, end_scan, 0, 23)

		if err != nil {
			log.Printf("ERR::: %v", err)
			a.LEDs.SetStatus(hw.CAL_SYNC_LED, false)
			continue
		}

		a.LEDs.SetStatus(hw.CAL_SYNC_LED, true)

		if wakeup.IsZero() {
			log.Printf("NO EVENT")
		} else {
			/* Got it! Schedule an alarm. */

			wakeup_at := wakeup.Add(-90 * time.Minute)
			log.Printf("SET ALARM TO: %v (now is: %v)", wakeup_at, now)
			a.SetAlarm(wakeup_at)
		}

	}
}

/* Runs as a goroutine that will turn on the error light if we crash */
func (self *Alarm) PanicHandler() {
	x := make(chan bool)
	<-x
	defer func() {
		self.LEDs.SetStatus(hw.ERROR_LED, true)
	}()
}

/* This runs as a goroutine to constantly pool the current time
 * and set of the alarm */
func (self *Alarm) AlarmLoop() {
	for {
		time.Sleep(1 * time.Second)
		self.Lock.Lock()
		now := time.Now()
		if self.AlarmIsSet && now.After(self.WakeUpAt) && !self.UI.ForceAlarmOff {
			/* Wake up time! Ring the alarm! */
			self.AlarmIsSet = false
			self.AlarmIsActive = true
			// XXX I think we need to say self.Aborting = false here, too.
			go self.WakeUp()
		}
		/* Set the alarm armed status LED */
		self.LEDs.SetStatus(hw.CAL_ARMED_LED, self.AlarmIsSet && !self.UI.ForceAlarmOff)
		self.Lock.Unlock()
	}
}

/*
 * This goroutine runs to actually do the alarm. It makes the
 * light go bring and (eventually) will play music!
 * it pools the "Checkin" function to give up if the alarm
 * is canceled.
 */
func (self *Alarm) WakeUp() {
	defer func() {
		/* Checkin will panic with &AbortWakeUpKey
		 * if the alarm to be canceled.
		 * If we panic with that, just terminate the goroutine,
		 * otherwise, we goofed and should rethrow the panic
		 */
		r := recover()
		if r != &AbortWakeUpKey && r != nil {
			panic(r)
		}
		if self.Player != nil {
			self.Player.Wait()
			self.Player = nil
		}
		self.AlarmIsActive = false
	}()

	/* Slowly make the lights brighter */
	print("DINGDINGDING")
	for i := 0; i <= 100; i++ {
		if i == 50 {
			self.StartMP3Player()
		}
		self.Checkin()
		print("BRIGHTER!")
		self.LEDs.SetLEDs(i, i, i)
		time.Sleep(2 * time.Second)
	}

}

func (self *Alarm) GetMP3() string {

	path := self.Config.WakeupMP3
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Printf("Can't read directory: %v", err)
		return ""
	}

	mp3_files := make([]string, 0, 0)

	for _, fi := range files {
		f := fi.Name()

		if len(f) > 4 &&
			f[0] != '.' &&
			f[len(f)-4:len(f)] == ".mp3" {

			mp3_files = append(mp3_files, f)
		}
	}

	if len(mp3_files) == 0 {
		return ""
	} else {
		return path + "/" + mp3_files[rand.Int()%len(mp3_files)]
	}
}

func (self *Alarm) StartMP3Player() {
	self.Lock.Lock()
	file_to_play := self.GetMP3()
	self.Player = exec.Command("mplayer", "-srate", "48000", file_to_play)
	self.Player.Start()
	self.Lock.Unlock()
}

/* Check if the alarm has been canceled and panic if it has */
func (self *Alarm) Checkin() {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if self.Aborting == true {
		self.Aborting = false
		panic(&AbortWakeUpKey)
	}
}

/* Set an alarm for a given time */
func (self *Alarm) SetAlarm(wake time.Time) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	log.Printf("Set alarm for: %v", wake)
	self.AlarmIsSet = true
	self.WakeUpAt = wake
}

func (self *Alarm) TestAlarm() {
	self.SetAlarm(time.Now())
}

func (self *Alarm) AbortAlarmInProgress() {
	//self.Lock.Lock();
	//defer self.Lock.Unlock();
	log.Printf("Abort alarm!")
	if self.AlarmIsActive {
		self.Aborting = true
	}
	if self.Player != nil && self.Player.Process != nil {
		self.Player.Process.Kill()
	}
}
