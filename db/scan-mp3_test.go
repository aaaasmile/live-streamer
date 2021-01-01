package db

import (
	"testing"

	"github.com/aaaasmile/live-streamer/conf"
	"github.com/aaaasmile/live-streamer/util"
)

func initConf() error {
	conf.ReadConfig(util.GetFullPath("../config.toml"))
	return nil
}

func TestParseMeta(t *testing.T) {
	//path := "C:\\local\\Music\\Colonna Sonora - Il Padrino [godfather]\\12 - Marcia Religiosa.mp3"
	path := "C:\\local\\Music\\youtube\\Siberia-Vento.mp3"
	if err := ScanFile(path); err != nil {
		t.Error(err)
	}
	t.Error("Try the scan")
}
