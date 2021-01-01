package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aaaasmile/live-streamer/conf"
	"github.com/aaaasmile/live-streamer/web/idl"
	"github.com/aaaasmile/live-streamer/web/live/omx"
	"github.com/aaaasmile/live-streamer/web/live/omx/fileplayer"
	"github.com/aaaasmile/live-streamer/web/live/omx/radio"
	"github.com/aaaasmile/live-streamer/web/live/omx/you"
)

func getProviderForURI(uri string, pl *omx.OmxPlayer) (idl.StreamProvider, error) {
	streamers := make([]idl.StreamProvider, 0)
	streamers = append(streamers, &you.YoutubePl{TmpInfo: conf.Current.TmpInfo})
	streamers = append(streamers, &fileplayer.FilePlayer{Dbus: pl.GetDbus()})
	streamers = append(streamers, &radio.RadioPlayer{})

	for _, prov := range streamers {
		if prov.IsUriForMe(uri) {
			return prov, nil
		}
	}
	return nil, fmt.Errorf("Unable to find a provider for the uri %s", uri)
}

func handlePlayUri(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	reqURI := struct {
		URI string `json:"uri"`
	}{}
	if err := json.Unmarshal(rawbody, &reqURI); err != nil {
		return err
	}

	if reqURI.URI == "" {
		log.Println("Ignore empty request")
		return fmt.Errorf("Ignore empty URI request")
	}
	if err := startUri(reqURI.URI, pl); err != nil {
		return err
	}

	return returnStatus(w, req, pl)
}

func startUri(uri string, pl *omx.OmxPlayer) error {
	if uri == "" {
		return fmt.Errorf("Nothing to play")
	}
	prov, err := getProviderForURI(uri, pl)
	if err != nil {
		return err
	}
	log.Println("Using provider name: ", prov.Name())
	if err := pl.StartPlay(uri, prov); err != nil {
		return err
	}
	if err := checkAfterStartPlay(prov.GetStatusSleepTime(), uri, pl); err != nil {
		return err
	}
	return nil
}

func handleNextTitle(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	uri, err := pl.NextTitle()
	if err != nil {
		return err
	}
	if uri != "" {
		if err := startUri(uri, pl); err != nil {
			return err
		}
	}

	return returnStatus(w, req, pl)
}

func handlePreviousTitle(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	uri, err := pl.PreviousTitle()
	if err != nil {
		return err
	}
	if uri != "" {
		if err := startUri(uri, pl); err != nil {
			return err
		}
	}
	return returnStatus(w, req, pl)
}

func handleSetPowerState(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	reqPower := struct {
		PowerState string `json:"power"`
	}{}

	if err := json.Unmarshal(rawbody, &reqPower); err != nil {
		return err
	}

	log.Println("Set power state request ", reqPower)

	switch reqPower.PowerState {
	case "off":
		err = pl.PowerOff()
		return nil
	case "on":
		last, err := liteDB.FetchHistory(0, 1)
		if err != nil {
			return err
		}
		if len(last) == 1 {
			log.Println("With power on try to play this uri ", last[0].URI)
			if err := startUri(last[0].URI, pl); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Toggle power state  not recognized %s", reqPower.PowerState)
	}
	if err != nil {
		return err
	}

	return returnStatusAfterCheck(w, req, pl)
}

func checkAfterStartPlay(sleepTime int, uri string, pl *omx.OmxPlayer) error {
	var err error
	log.Println("Check the status after play ", sleepTime)
	time.Sleep(200 * time.Millisecond)
	i := 0
	for i < 8 {
		err = pl.CheckStatus(uri)
		if err != nil {
			log.Println("Error and retry play ", i, err)
			i++
		} else {
			break
		}
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
	log.Println("Status player now: OK")
	return err
}

func handleChangeVolume(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	reqVol := struct {
		VolumeType string `json:"volume"`
	}{}

	if err := json.Unmarshal(rawbody, &reqVol); err != nil {
		return err
	}

	log.Println("Change volume request ", reqVol)

	switch reqVol.VolumeType {
	case "up":
		if err = pl.VolumeUp(); err != nil {
			return err
		}
		return returnStatus(w, req, pl)
	case "down":
		if err = pl.VolumeDown(); err != nil {
			return err
		}
		return returnStatus(w, req, pl)
	}

	stateMute := ""
	switch reqVol.VolumeType {
	case "mute":
		if stateMute, err = pl.VolumeMute(); err != nil {
			return err
		}
	case "unmute":
		if stateMute, err = pl.VolumeUnmute(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Change volume request not recognized %s", reqVol.VolumeType)
	}
	res := struct {
		Mute string `json:"mute"`
		Type string `json:"type"`
	}{
		Type: "mute",
		Mute: stateMute,
	}
	log.Println("Mute state ", stateMute)
	return writeResponseNoWsBroadcast(w, res)
}

func handlePauseOrResume(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer, act string) error {
	log.Println("Resume request ")
	statePlay := ""
	var err error
	switch act {
	case "Resume":
		if statePlay, err = pl.Resume(); err != nil {
			return err
		}
	case "Pause":
		if statePlay, err = pl.Pause(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Action in pause/resume %s not recognized", act)
	}

	res := struct {
		PlayState string `json:"playstate"`
		Type      string `json:"type"`
	}{
		PlayState: statePlay,
		Type:      "playsate",
	}
	log.Println("Playing state ", statePlay)
	return writeResponseNoWsBroadcast(w, res)
}

func handlePlayerState(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	return returnStatus(w, req, pl)
}

func returnStatus(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	if err := pl.CheckStatus(pl.GetCurrURI()); err != nil {
		return err
	}
	return returnStatusAfterCheck(w, req, pl)
}

func returnStatusAfterCheck(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	res := struct {
		Player        string `json:"player"`
		Mute          string `json:"mute"`
		URI           string `json:"uri"`
		TrackDuration string `json:"trackDuration"`
		TrackPosition string `json:"trackPosition"`
		TrackStatus   string `json:"trackStatus"`
		Type          string `json:"type"`
		Title         string `json:"title"`
		Description   string `json:"description"`
	}{
		Player:        pl.GetStatePlaying(),
		Mute:          pl.GetStateMute(),
		URI:           pl.GetCurrURI(),
		TrackDuration: pl.GetTrackDuration(),
		TrackPosition: pl.GetTrackPosition(),
		TrackStatus:   pl.GetTrackStatus(),
		Type:          "status",
		Title:         pl.GetStateTitle(),
		Description:   pl.GetStateDescription(),
	}

	return writeResponse(w, res)
}
