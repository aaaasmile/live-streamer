package dbus

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/aaaasmile/live-streamer/web/live/omx/omxstate"
	"github.com/godbus/dbus"
)

type OmxDbus struct {
	coDBus dbus.BusObject
}

func (op *OmxDbus) ClearDbus() {
	op.coDBus = nil
}

func (op *OmxDbus) connectObjectDbBus() error {
	if op.coDBus != nil {
		return nil
	}
	u, err := user.Current()
	log.Println("User ", u.Username)

	fname := fmt.Sprintf("/tmp/omxplayerdbus.%s", u.Username)
	if _, err := os.Stat(fname); err == nil {
		//busAddr := "unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d"
		raw, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(raw))
		log.Println("Env DBUS_SESSION_BUS_ADDRESS set to ", string(raw))
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	obj := conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2/omxplayer")

	address := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	log.Println("session bus addr ", address)
	log.Println("dbus connection: ", conn)

	op.coDBus = obj
	return nil
}

func (op *OmxDbus) getProperty(prop string) (*dbus.Variant, error) {
	if err := op.connectObjectDbBus(); err != nil {
		return nil, err
	}
	res, err := op.coDBus.GetProperty(prop)
	if err != nil {
		op.coDBus = nil
		return nil, err
	}
	return &res, nil
}

func (op *OmxDbus) CallSimpleAction(action string) error {
	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0)
	return nil
}

// dbus-send --print-reply=literal --session --dest=org.mpris.MediaPlayer2.omxplayer /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Action int32:18 >/dev/null
func (op *OmxDbus) CallIntAction(action string, id int) error {
	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0, id)
	return nil
}

func (op *OmxDbus) CallStrAction(action string, para string) error {
	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0, para)
	return nil
}

func (op *OmxDbus) CheckTrackStatus(st *omxstate.StateOmx) error {
	log.Println("go ahead with dbus status")
	dur, err := op.getProperty("org.mpris.MediaPlayer2.Player.Duration")
	if err != nil {
		return err
	}
	pos, err := op.getProperty("org.mpris.MediaPlayer2.Player.Position")
	if err != nil {
		return err
	}

	status, err := op.getProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		return err
	}

	st.TrackDuration = fmt.Sprint(dur)
	st.TrackPosition = fmt.Sprint(pos)
	st.TrackStatus = fmt.Sprint(status)

	log.Println("Duration, position,  status ", dur, pos, status)
	return nil
}
