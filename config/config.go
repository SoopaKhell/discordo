package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

const USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0"

type GeneralConfig struct {
	UserAgent          string `toml:"user_agent"`
	Mouse              bool   `toml:"mouse"`
	Notifications      bool   `toml:"desktop_notifications"`
	FetchMessagesLimit int    `toml:"fetch_messages_limit"`
}

type ThemeConfig struct{}
type KeybindingsConfig struct{}

type Config struct {
	General     GeneralConfig     `toml:"general"`
	Theme       ThemeConfig       `toml:"theme"`
	Keybindings KeybindingsConfig `toml:"keybindings"`
}

func NewConfig() Config {
	return Config{
		General: GeneralConfig{
			UserAgent:          USER_AGENT,
			Mouse:              true,
			Notifications:      true,
			FetchMessagesLimit: 50,
		},
		Theme:       ThemeConfig{},
		Keybindings: KeybindingsConfig{},
	}
}

// LoadConfig creates default configuration (if the configuration file does not exist already) or reads the configuration file from user configuration directory.
func LoadConfig() Config {
	configPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(configPath+"/discordo.toml", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var c Config
	if fileInfo, _ := f.Stat(); fileInfo.Size() == 0 {
		e := toml.NewEncoder(f)
		e.Indent = "\t"

		c = NewConfig()
		err = e.Encode(&c)
		if err != nil {
			panic(err)
		}
	} else {
		_, err = toml.NewDecoder(f).Decode(&c)
		if err != nil {
			panic(err)
		}
	}

	return c
}
