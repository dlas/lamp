

package hw

import (
	"i2c"
	"sync"
	"time"
	"fmt"
)

type HW struct {
	LED* i2c.I2C
	GPIO* i2c.I2C
	Lock sync.Mutex
	ButtonCallback func(irq int, cur int)
}

type NullHW struct {

}

type HWInterface interface {
	SetLEDs(r, g, b int)
	SetStatus(status int, value bool)
	RegisterButtonCallback(c func(irq int, cur int))
}

func (x * NullHW) SetLEDs(r,g,b int){}
func (x * NullHW) SetStatus(s int, v bool) {}
func (x * NullHW) RegisterButtonCallback(c func(irq int, cur int)) {}


const SOFTWARE_ALIVE_LED = 1
const CAL_SYNC_LED=2
const CAL_ARMED_LED=3
const ERROR_LED=4

const BUTTON_TOGGLE_LIGHTS= 1
const BUTTON_LIGHTING_MODE = 2
const BUTTON_TOGGLE_ALARM = 4

func NewHW() (HWInterface, error) {
	led, err:= i2c.NewI2C(65, 2);
	if (err != nil) {
		return &NullHW{}, err;
	}

	gpio, err := i2c.NewI2C(38, 2);
	if (err != nil) {
		return &NullHW{}, err;
	}

	
	var res HW;
	res.LED= led;
	res.GPIO = gpio
	res.INIT();
	go res.ButtonPoller()
	return &res, nil
}

func (hw * HW) INIT() {
	hw.LED.WriteRegU8(0, 0);
	hw.LED.WriteRegU8(1, 0);
	for i := 9; i <= 21; i++ {

		hw.LED.WriteRegU8(uint8(i), 0);
	}

	for i := 38; i <= 53; i++ {
		hw.LED.WriteRegU8(uint8(i), 0)
	}
	hw.GPIO.WriteRegU8(0, 0xFF);
	hw.GPIO.WriteRegU8(1, 0);
	hw.GPIO.WriteRegU8(2, 0xFF);

	hw.GPIO.WriteRegU8(3, 0xFF)
	hw.GPIO.WriteRegU8(4, 0xFF)
	hw.GPIO.WriteRegU8(6, 0xFF)
}
func (hw * HW) ButtonPoller() {
	for ;; {
		time.Sleep(2000 * time.Millisecond);
		hw.Lock.Lock();
		irq, _ := hw.GPIO.ReadRegU8(7);
		//irq = 0xFF ^ irq;
		callback := hw.ButtonCallback;
		current, _ := hw.GPIO.ReadRegU8(9)
		hw.Lock.Unlock()

		if (irq != 0) {
			if (callback != nil) {
				callback(int(irq), int(current));
			}
			print(fmt.Sprintf("BUTTONS BUTTONS BUTTONS: %v\n", irq));
		}
	}
}


func (hw * HW) ReadButtons() int {
	hw.Lock.Lock();
	defer hw.Lock.Unlock();
	r, _ :=hw.GPIO.ReadRegU8(9);
	return int(r);

}
	
func (hw * HW) SetLEDs(r, g, b int) {
	hw.Lock.Lock();
	defer hw.Lock.Unlock();
			
	hw.LED.WriteRegU8(13, uint8(g));
	hw.LED.WriteRegU8(17, uint8(b));
	hw.LED.WriteRegU8(21, uint8(r));
}

func (hw * HW) SetStatus(status int, value bool) {
	reg := 41 + 4 * status;
	
	val := 0;
	if value {
		val = 16
	}

	hw.LED.WriteRegU8(uint8(reg), uint8(val));
}

func (hw *HW) RegisterButtonCallback(c func(irq int, cur int)) {
	hw.ButtonCallback = c;
}

