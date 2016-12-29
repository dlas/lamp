

package hw

import (
	"i2c"
	"sync"
)

type HW struct {
	LED* i2c.I2C
	GPIO* i2c.I2C
	Lock sync.Mutex
}


func NewHW() (*HW, error) {
	led, err:= i2c.NewI2C(65, 2);
	if (err != nil) {
		return nil, err;
	}

	gpio, err := i2c.NewI2C(38, 2);
	if (err != nil) {
		return nil, err;
	}

	
	var res HW;
	res.LED= led;
	res.GPIO = gpio
	res.INIT();
	return &res, nil
}

func (hw * HW) INIT() {
	hw.LED.WriteRegU8(0, 0);
	hw.LED.WriteRegU8(1, 4);
	
}

func (hw * HW) ReadButtons() int {
	hw.Lock.Lock();
	defer hw.Lock.Unlock();
	r, _ :=hw.GPIO.ReadRegU8(9);
	return int(r);

}
	
func (hw * HW) SetLEDs(r, b, g int) {
	hw.Lock.Lock();
	defer hw.Lock.Unlock();
			
	hw.LED.WriteRegU8(13, uint8(g));
	hw.LED.WriteRegU8(17, uint8(b));
	hw.LED.WriteRegU8(21, uint8(r));
}


