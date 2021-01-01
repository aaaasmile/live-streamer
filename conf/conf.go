package conf

import (
	"encoding/json"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceURL      string
	RootURLPattern  string
	UseRelativeRoot bool
	DebugVerbose    bool
	OmxCmdParams    string
	DBPath          string
	TmpInfo         string
	VueLibName      string
	SoundCloud      SoundCloud
}

type SoundCloud struct {
	CfgFile     string
	IsAvailable bool
	ClientID    string
	AuthToken   string
	UserAgent   string
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
	if err := readSoundCloudFile(&Current.SoundCloud); err != nil {
		log.Println("Soundcloud configuration error", err)
		Current.SoundCloud.IsAvailable = false
	}
	return Current
}

func readSoundCloudFile(sc *SoundCloud) error {
	log.Println("Read configuration file for sound cloud ", sc.CfgFile)
	f, err := os.Open(sc.CfgFile)
	if err != nil {
		return err
	}

	defer f.Close()
	info := struct {
		ClientID  string
		AuthToken string
		UserAgent string
	}{}

	err = json.NewDecoder(f).Decode(&info)
	if err != nil {
		return err
	}
	sc.ClientID = info.ClientID
	sc.AuthToken = info.AuthToken
	sc.UserAgent = info.UserAgent

	return nil
}
