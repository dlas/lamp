
package config

type Config struct {
	Presets [3]RGB
	
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
	return &Config{}
}


