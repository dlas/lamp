package alarm

import "hw"

/* Main object to track user interface state. This is part of the Alarm object
 * all all the UI methods are defined on the Alarm object but also touch alarm.UI
 */

type UIState struct {
	/* Current intensities */
	Red   int
	Green int
	Blue  int

	/* Force the alarm off */
	ForceAlarmOff bool

	/* Cycle through preset lighting and light levels */
	LightingLevel  int /* 0..100*/
	LightingPreset int
}

/* Run when the user presses the preset button. Cycles to the next defined
 * preset */

func (self *Alarm) IncrementPreset() {
	l := self.UI.LightingPreset
	l++
	if l >= len(self.Config.Presets) {
		l = 0
	}

	self.UI.LightingPreset = l
	self.Update()

}

/* Runs when the user hits the brightness button.  Increment lamp brightness,
 * or go back to zero.*/
func (self *Alarm) IncrementLevel() {
	l := self.UI.LightingLevel
	l += 25
	if l > 100 {
		l = 0
	}

	self.UI.LightingLevel = l
	self.Update()
}

/* We run this to switch the lights to the a new level/preset combination
 */
func (self *Alarm) Update() {
	rgb := self.Config.Presets[self.UI.LightingPreset]
	i := float64(self.UI.LightingLevel)
	r := int(float64(rgb.Red) * i / 100.0)
	g := int(float64(rgb.Green) * i / 100.0)
	b := int(float64(rgb.Blue) * i / 100.0)
	self.LEDs.SetLEDs(r, g, b)
}

/*
 * This expects to be called when an EXTERNAL interface (like the web)
 * wants to change the color.
 */
func (self *Alarm) UIChangeLights(red int, green int, blue int) {

	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.UI.Red = red
	self.UI.Green = green
	self.UI.Blue = blue

	/* An externally motivated change cancels any alarm in progress */
	self.AbortAlarmInProgress()
	self.LEDs.SetLEDs(red, green, blue)
}

/*
 * Set the ForceAlarmOff flag
 */
func (self *Alarm) ForceAlarmOff(armed bool) {
	self.UI.ForceAlarmOff = armed
}

/* What to do when the user presses a button? */
func (self *Alarm) ButtonPress(irq int, current int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if self.AlarmIsActive {
		/* If the alarm is active, cancel it and turn the light to
		 * full.
		 */
		/* FIXME: Don't hardcode brightness here!!! Use preset 0*/
		self.LEDs.SetLEDs(15, 15, 15)
		self.AbortAlarmInProgress()
	} else {
		/* Otherwise, do something particular to that button */
		if irq&hw.BUTTON_TOGGLE_LIGHTS != 0 {
			self.IncrementLevel()
		}

		if irq&hw.BUTTON_LIGHTING_MODE != 0 {
			self.IncrementPreset()
		}

		if irq&hw.BUTTON_TOGGLE_ALARM != 0 {
			//self.ForceAlarmOff(!self.UI.ForceAlarmOff)
		}
	}
}
