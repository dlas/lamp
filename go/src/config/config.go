
package config

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


func SaveConfig(c * Config){}

func LoadConfig() *Config{
	var c Config;
	c.Presets= []RGB{
		RGB{5, 2, 2},
		RGB{2, 3, 6},
		RGB{15, 15, 15},
	}
	return &c
}


