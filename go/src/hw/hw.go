/* This implements the high-level i2c hardware interface
 * to the custom LED control board.
 */

package hw

import (
	"fmt"
	"i2c"
	"sync"
	"time"
)

/*
 * There are two implementations of this interface.
 * The real one, "HW" and a Null implementation that
 * we use for testing on computers without appropriate
 * i2c hookups.
 */

type HWInterface interface {
	SetLEDs(r, g, b int)
	SetStatus(status int, value bool)
	RegisterButtonCallback(c func(irq int, cur int))
}

/*
 * What state do we have to hold on to?
 */
type HW struct {
	/* File object to the LED PWM chip */
	LED *i2c.I2C
	/* File object to the GPIO chip */
	GPIO *i2c.I2C
	Lock sync.Mutex
	/* Callback to call whenever buttons are pressed */
	ButtonCallback func(irq int, cur int)
}

/* Implementation of the null HWInterface */
type NullHW struct{}

func (x *NullHW) SetLEDs(r, g, b int)                             {}
func (x *NullHW) SetStatus(s int, v bool)                         {}
func (x *NullHW) RegisterButtonCallback(c func(irq int, cur int)) {}

/* Definitions! How are things hooked up */
const SOFTWARE_ALIVE_LED = 1
const CAL_SYNC_LED = 2
const CAL_ARMED_LED = 3
const ERROR_LED = 4

const BUTTON_TOGGLE_LIGHTS = 1
const BUTTON_LIGHTING_MODE = 2
const BUTTON_TOGGLE_ALARM = 4

const I2C_LED_ADDRESS = 65
const I2C_GPIO_ADDRESS = 38

/*
 * Make a new HWInterface object. On error, return a
 * Nullinterface suitable for testing.
 */
func NewHW() (HWInterface, error) {
	led, err := i2c.NewI2C(I2C_LED_ADDRESS, 2)
	if err != nil {
		return &NullHW{}, err
	}

	gpio, err := i2c.NewI2C(I2C_GPIO_ADDRESS, 2)
	if err != nil {
		return &NullHW{}, err
	}

	var res HW
	res.LED = led
	res.GPIO = gpio
	res.INIT()
	go res.ButtonPoller()
	return &res, nil
}

/* Initialize the hardware */
func (hw *HW) INIT() {

	/* Set the mode registers of the LED driver chip.
	 * It happens that 0 and 0 are the desiered values.
	 */
	hw.LED.WriteRegU8(0, 0)
	hw.LED.WriteRegU8(1, 0)

	/* Set all of the LEDs off */
	for i := 9; i <= 21; i++ {

		hw.LED.WriteRegU8(uint8(i), 0)
	}

	/* Set all of the status LEDs off. */
	/* XXX this is actually broken. its backwards and turns them on. */
	for i := 38; i <= 53; i++ {
		hw.LED.WriteRegU8(uint8(i), 0)
	}

	/* Set configuration for the GPIO chip */
	hw.GPIO.WriteRegU8(0, 0xFF)
	hw.GPIO.WriteRegU8(1, 0)
	hw.GPIO.WriteRegU8(2, 0xFF)

	hw.GPIO.WriteRegU8(3, 0xFF)
	hw.GPIO.WriteRegU8(4, 0xFF)
	hw.GPIO.WriteRegU8(6, 0xFF)
}

/* Run in a LOOP, pollin gthe GPIO chip */
func (hw *HW) ButtonPoller() {
	for {
		time.Sleep(500 * time.Millisecond)
		hw.Lock.Lock()
		/* IRQ status register is 7; current input is 9 */
		irq, _ := hw.GPIO.ReadRegU8(7)
		callback := hw.ButtonCallback
		current, _ := hw.GPIO.ReadRegU8(9)
		hw.Lock.Unlock()

		/* Call the callback if an irq happened */
		if irq != 0 {
			if callback != nil {
				callback(int(irq), int(current))
			}
			print(fmt.Sprintf("BUTTONS BUTTONS BUTTONS: %v\n", irq))
		}
	}
}

func (hw *HW) ReadButtons() int {
	hw.Lock.Lock()
	defer hw.Lock.Unlock()
	r, _ := hw.GPIO.ReadRegU8(9)
	return int(r)

}

/* Set the illumination LEDs to any brightness.
 * TODO: use the low and high register to get more fine graned
 * control.
 */
func (hw *HW) SetLEDs(r, g, b int) {
	hw.Lock.Lock()
	defer hw.Lock.Unlock()

	hw.LED.WriteRegU8(13, uint8(g))
	hw.LED.WriteRegU8(17, uint8(b))
	hw.LED.WriteRegU8(21, uint8(r))
}

/* Turn one of the status LEDs on or off */
func (hw *HW) SetStatus(status int, value bool) {
	reg := 41 + 4*status

	val := 0
	if value {
		val = 16
	}

	hw.LED.WriteRegU8(uint8(reg), uint8(val))
}

/* Set the callback to run when there is a button changes state */
func (hw *HW) RegisterButtonCallback(c func(irq int, cur int)) {
	hw.ButtonCallback = c
}
