
package config

import "encoding/json"
import "io/ioutil"
import "log"

type Config struct {
	Presets []RGB
	
	AlarmEnable bool
	GoogleEnable bool
	GoogleSecret []byte
	GoogleAuth []byte

}

type RGB struct {
	Red int
	Blue int
	Green int
}


func SaveConfig(c * Config){
	buf, _ := json.Marshal(c);
	ioutil.WriteFile("config", buf, 0660);
	log.Printf("SAVE: %v", *c);

}

/* Try to load and parse configuration from a file */
func LoadConfig() *Config{
	var c Config;

	buf, err := ioutil.ReadFile("config");

	if (err == nil) {
		err = json.Unmarshal(buf, &c);
		if (err == nil) {
			log.Printf("RECOVER: %v", c);
			return &c
		}
	}		

	/* That failed. Get out of here with some sensible defaults */
	c.Presets= []RGB{
		RGB{5, 2, 2},
		RGB{2, 3, 6},
		RGB{15, 15, 15},
	}

	cs, _ := ioutil.ReadFile("client_secret.json");
	c.GoogleSecret = cs
	return &c
}


