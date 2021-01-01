package omx

import (
	"log"
	"sync"

	"github.com/aaaasmile/live-streamer/web/idl"
	"github.com/aaaasmile/live-streamer/web/live/omx/omxstate"
	"github.com/aaaasmile/live-streamer/web/live/omx/playlist"
)

type OmxPlayer struct {
	mutex         *sync.Mutex
	state         omxstate.StateOmx
	chDbOperation chan *idl.DbOperation
	PlayList      *playlist.LLPlayList
	Providers     map[string]idl.StreamProvider
	ChAction      chan *omxstate.ActionDef
}

func NewOmxPlayer(chDbop chan *idl.DbOperation) *OmxPlayer {
	cha := make(chan *omxstate.ActionDef)
	res := OmxPlayer{
		mutex:         &sync.Mutex{},
		chDbOperation: chDbop,
		Providers:     make(map[string]idl.StreamProvider),
		ChAction:      cha,
	}

	return &res
}

func (op *OmxPlayer) ListenOmxState(statusCh chan *omxstate.StateOmx) {
	log.Println("start listenOmxState. Waiting for status change in omxplayer")
	for {
		st := <-statusCh
		op.mutex.Lock()
		log.Println("Set OmxPlayer state ", st)
		if st.StatePlayer == omxstate.SPoff {
			k := op.state.CurrURI
			if _, ok := op.Providers[k]; ok {
				delete(op.Providers, k)
			}
			op.state.ClearTrackStatus()
		} else {
			op.state.TrackDuration = st.TrackDuration
			op.state.TrackPosition = st.TrackPosition
			op.state.TrackStatus = st.TrackStatus
			op.state.StateMute = st.StateMute
		}
		op.state.CurrURI = st.CurrURI
		op.state.StatePlayer = st.StatePlayer
		op.state.Info = st.Info
		op.mutex.Unlock()
	}
}

func (op *OmxPlayer) GetTrackDuration() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		if td, ok := prov.GetTrackDuration(); ok {
			return td
		}
	}

	return op.state.TrackDuration
}

func (op *OmxPlayer) GetTrackPosition() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		if td, ok := prov.GetTrackPosition(); ok {
			return td
		}
	}
	return op.state.TrackPosition
}

func (op *OmxPlayer) GetTrackStatus() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		log.Println("Tracking satus of ", prov)
		if td, ok := prov.GetTrackStatus(); ok {
			return td
		}
	}
	return op.state.TrackStatus
}

func (op *OmxPlayer) GetStatePlaying() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.StatePlayer.String()
}

func (op *OmxPlayer) GetStateTitle() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetTitle()
	}

	return ""
}

func (op *OmxPlayer) GetStateDescription() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetDescription()
	}

	return ""
}

func (op *OmxPlayer) GetCurrURI() string {
	log.Println("getCurrURI")
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.CurrURI
}

func (op *OmxPlayer) StartPlay(URI string, prov idl.StreamProvider) error {
	var err error
	if op.PlayList, err = playlist.CreatePlaylistFromProvider(URI, prov); err != nil {
		return err
	}
	log.Println("StartPlay ", URI)

	return op.startPlayListCurrent(prov)
}

func (op *OmxPlayer) PreviousTitle() (string, error) {
	if op.PlayList == nil {
		log.Println("Nothing to play because no playlist is provided")
		return "", nil
	}
	var curr *playlist.PlayItem
	var ok bool
	if _, ok = op.PlayList.CheckCurrent(); !ok {
		return "", nil
	}

	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI == "" {
		log.Println("Player is not active, ignore next title")
		return "", nil
	}

	if curr, ok = op.PlayList.Previous(); !ok {
		return "", nil
	}

	u := curr.URI
	log.Println("the previous title is", u)

	return u, nil
}

func (op *OmxPlayer) NextTitle() (string, error) {
	if op.PlayList == nil {
		log.Println("Nothing to play because no playlist is provided")
		return "", nil
	}
	var curr *playlist.PlayItem
	var ok bool
	if _, ok = op.PlayList.CheckCurrent(); !ok {
		return "", nil
	}

	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI == "" {
		return "", nil
	}

	if curr, ok = op.PlayList.Next(); !ok {
		return "", nil
	}

	u := curr.URI
	log.Println("the next title is", u)

	return u, nil
}

func (op *OmxPlayer) CheckStatus(uri string) error {
	if uri == "" {
		return nil
	}
	log.Println("Check state uri ", uri)
	op.mutex.Lock()
	defer op.mutex.Unlock()

	log.Println("Check status req", op.state)

	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		if err := prov.CheckStatus(op.chDbOperation); err != nil {
			return err
		}
	}

	return nil
}

func (op *OmxPlayer) PowerOff() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	log.Println("Power off, terminate omxplayer with signal kill")
	op.freeAllProviders()
	return nil
}

func (op *OmxPlayer) freeAllProviders() {
	for k, prov := range op.Providers {
		log.Println("Sending kill signal to ", k)
		ch := prov.GetCmdStopChannel()
		if ch != nil {
			log.Println("Force kill with channel")
			ch <- struct{}{}
			prov.CloseStopChannel()
		}
	}

	op.Providers = make(map[string]idl.StreamProvider)

}
