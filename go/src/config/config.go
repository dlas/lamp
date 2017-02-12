/* This holds a master confirmation object that is shared between
 * many modules.
 */
package config

import "encoding/json"
import "io/ioutil"
import "log"

type Config struct {
	/* A set of lamp presets */
	Presets []RGB

	AlarmEnable bool // Unused, KILL

	/* Do we try to calk to google calendar */
	GoogleEnable bool

	/* Google API secret */
	GoogleSecret []byte

	/* Google OAUTH token */
	GoogleAuth []byte

	/* What MP3 to play during the alarm */
	WakeupMP3 string

	/* Howmany minute before the first meeting to we need to wake up */
	WakeupMinsToMeeting int

	/* What is the first hour of the day topay attention to */
	MorningMeetingHourStart int

	/* What is the last hour of the day to pay attention to */
	MorningMeetingHourEnd int
}

type RGB struct {
	Red   int
	Blue  int
	Green int
}

/* Safe a configuration. */
func SaveConfig(c *Config) {
	buf, _ := json.Marshal(c)
	ioutil.WriteFile("config", buf, 0660)
	log.Printf("SAVE: %v", *c)

}

/* Try to load and parse configuration from a file.
 * If this goes wrong for any reason, make up a basic config and
 * return that instead */
func LoadConfig() *Config {
	var c Config

	buf, err := ioutil.ReadFile("config")

	if err == nil {
		err = json.Unmarshal(buf, &c)
		if err == nil {
			c.GoogleSecret, _ = ioutil.ReadFile("client_secret.json")
			log.Printf("RECOVER: %v", c)
			return &c
		}
	}

	/* That failed. Get out of here with some sensible defaults */
	c.Presets = []RGB{
		RGB{5, 2, 2},
		RGB{2, 3, 6},
		RGB{15, 15, 15},
	}

	cs, _ := ioutil.ReadFile("client_secret.json")
	c.GoogleSecret = cs
	return &c
}
