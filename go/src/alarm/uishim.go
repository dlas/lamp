package alarm

import "hw"

type UIState struct {
	Red int
	Green int
	Blue int
	ForceAlarmOff bool
	
	LightingLevel int
	LightingPreset int
}

func (self *Alarm) IncrementPreset() {
	l := self.UI.LightingPreset;
	l++;
	if (l >= len(self.Config.Presets)) {
		l = 0;
	}
	
	self.UI.LightingPreset = l
	self.Update();


}

func (self *Alarm) IncrementLevel() {
	l:= self.UI.LightingLevel;
	l+=25;
	if (l > 100) {
		l = 0;
	}

	self.UI.LightingLevel = l;
	self.Update();
}

func (self *Alarm) Update() {

	rgb := self.Config.Presets[self.UI.LightingPreset];
	i := float64(self.UI.LightingLevel);
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

	self.Lock.Lock();
	defer self.Lock.Unlock();
	self.UI.Red = red
	self.UI.Green = green
	self.UI.Blue = blue

	self.AbortAlarmInProgress()
	self.LEDs.SetLEDs(red,green,blue);
}


func (self *Alarm) ForceAlarmOff(armed bool) {
	self.UI.ForceAlarmOff=armed

}

func (self *Alarm) ButtonPress(irq int,current int) {
	self.Lock.Lock();
	defer self.Lock.Unlock()
	if (self.AlarmIsActive) {
		self.LEDs.SetLEDs(15, 15, 15);
		self.AbortAlarmInProgress();
	} else {
		if (irq&hw.BUTTON_TOGGLE_LIGHTS != 0) {
			self.IncrementLevel();
		}
		
		if (irq & hw.BUTTON_LIGHTING_MODE != 0) {
			self.IncrementPreset();
		}

		if (irq & hw.BUTTON_TOGGLE_ALARM != 0) {
			self.ForceAlarmOff(!self.UI.ForceAlarmOff)
		}
	}
}
