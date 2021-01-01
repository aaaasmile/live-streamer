package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceURL      string
	RootURLPattern  string
	UseRelativeRoot bool
	DebugVerbose    bool
	DBPath          string
	TmpInfo         string
	VueLibName      string
}


var Current = &Config{}

func ReadConfig(configfile string) *Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		log.Fatal(err)
	}
	
	return Current
}
