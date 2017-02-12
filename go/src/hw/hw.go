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

	LastButtons uint8
}

/* Implementation of the null HWInterface */
type NullHW struct{}

func (x *NullHW) SetLEDs(r, g, b int)                             {}
func (x *NullHW) SetStatus(s int, v bool)                         {}
func (x *NullHW) RegisterButtonCallback(c func(irq int, cur int)) {}

/* Definitions! How are things hooked up */
const SOFTWARE_ALIVE_LED = 0
const CAL_SYNC_LED = 1
const CAL_ARMED_LED = 2
const ERROR_LED = 3

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

	/* Set the prescaler to ~ 200 hz */
	hw.LED.WriteRegU8(0, 16)
	time.Sleep(10 * time.Millisecond)
	hw.LED.WriteRegU8(254, 3)
	time.Sleep(10 * time.Millisecond)

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

		/* What's going on here?
		 * We use the IRQ feature of the GPIO chip to find buttons
		 * who'se state has changed since the last loop. That
		 * prevents us from missing a trigger if the user doesn't
		 * hold the button down for 500ms.
		 *
		 * However, we only consdier interrupts for buttons that
		 * were not pressed last time we did this loop. Otherwise, we'd
		 * keep retriggering the same event.
		 */
		/* IRQ status register is 7; current input is 9 */
		irq, _ := hw.GPIO.ReadRegU8(7)
		callback := hw.ButtonCallback
		current, _ := hw.GPIO.ReadRegU8(9)

		/* Ignore interrupts for buttons that were pressed last time
		 * we did this loop */
		irq = irq & hw.LastButtons
		hw.LastButtons = current
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
 * r,g,b are intensities from 0..100 and we scale them to
 * the PWM chips dynamic range.
 */
func (hw *HW) SetLEDs(r, g, b int) {
	hw.Lock.Lock()
	defer hw.Lock.Unlock()

	scaled_red := r * 4095 / 100
	scaled_green := g * 4095 / 100
	scaled_blue := b * 4095 / 100

	r_high := scaled_red >> 8
	g_high := scaled_green >> 8
	b_high := scaled_blue >> 8

	r_low := scaled_red & 255
	g_low := scaled_green & 255
	b_low := scaled_blue & 255

	hw.LED.WriteRegU8(13, uint8(r_high))
	hw.LED.WriteRegU8(12, uint8(r_low))
	hw.LED.WriteRegU8(17, uint8(g_high))
	hw.LED.WriteRegU8(16, uint8(g_low))
	hw.LED.WriteRegU8(21, uint8(b_high))
	hw.LED.WriteRegU8(20, uint8(b_low))

	/*
		hw.LED.WriteRegU8(13, uint8(g))
		hw.LED.WriteRegU8(17, uint8(b))
		hw.LED.WriteRegU8(21, uint8(r))
	*/
}

/* Turn one of the status LEDs on or off */
func (hw *HW) SetStatus(status int, value bool) {
	reg := 39 + 4*status

	val := 0
	if !value {
		val = 16
	}

	hw.LED.WriteRegU8(uint8(reg), uint8(val))
}

/* Set the callback to run when there is a button changes state */
func (hw *HW) RegisterButtonCallback(c func(irq int, cur int)) {
	hw.ButtonCallback = c
}
