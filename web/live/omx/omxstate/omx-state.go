package omxstate

import "log"

type SPstateplaying int

func (s *SPstateplaying) String() string {
	switch *s {
	case SPundef:
		return "undef"
	case SPoff:
		return "off"
	case SPplaying:
		return "playing"
	case SPpause:
		return "pause"
	case SPrestart:
		return "restart"
	}
	return ""
}

const (
	SPundef = iota
	SPoff
	SPplaying
	SPpause
	SPrestart
)

type SMstatemute int

func (s *SMstatemute) String() string {
	switch *s {
	case SMnormal:
		return "normal"
	case SMmuted:
		return "muted"
	}
	return ""
}

const (
	SMnormal = iota
	SMmuted
	SMundef
)

type StateOmx struct {
	CurrURI       string
	StatePlayer   SPstateplaying
	StateMute     SMstatemute
	Info          string
	ItemType      string
	NextItem      string
	PrevItem      string
	TrackDuration string
	TrackPosition string
	TrackStatus   string
}

func (so *StateOmx) ClearTrackStatus() {
	so.TrackDuration = ""
	so.TrackPosition = ""
	so.TrackStatus = ""
	so.StateMute = SMundef
}

type ActionTD int

///home/igors/projects/go/bin/stringer -type=actionTD

const (
	ActTerminate ActionTD = iota
	ActPlaying
	ActPause
	ActMute
	ActUnmute
)

type ActionDef struct {
	URI        string
	Action     ActionTD
	ChStateRsp chan *StateOmx
}

type WorkerState struct {
	ChStatus chan *StateOmx
}

func ListenStateAction(actCh chan *ActionDef, workers []WorkerState) {
	log.Println("Waiting for action to change the state")
	var stateCurrent SPstateplaying
	var muteStateCurrent SMstatemute
	stateCurrent = SPoff
	muteStateCurrent = SMnormal
	uriPlaying := ""
	for {
		st := <-actCh
		olduriPlaying := uriPlaying
		log.Println("New action in state: ", st.Action.String(), stateCurrent.String())
		stateNext := StateOmx{CurrURI: olduriPlaying, StatePlayer: SPundef, StateMute: SMundef}
		switch stateCurrent {
		case SPoff:
			switch st.Action {
			case ActPlaying:
				if st.URI != "" {
					stateNext.StatePlayer = SPplaying
					uriPlaying = st.URI
				}
			case ActMute:
				stateNext.StateMute = SMmuted
			case ActUnmute:
				stateNext.StateMute = SMnormal
			}
		case SPplaying:
			switch st.Action {
			case ActPlaying:
				if st.URI != "" {
					stateNext.StatePlayer = SPrestart
					uriPlaying = st.URI
				}
			case ActPause:
				stateNext.StatePlayer = SPpause
			case ActTerminate:
				stateNext.StatePlayer = SPoff
				stateNext.StateMute = SMnormal
				uriPlaying = ""
			case ActMute:
				stateNext.StateMute = SMmuted
			case ActUnmute:
				stateNext.StateMute = SMnormal
			}
		case SPpause:
			switch st.Action {
			case ActPlaying:
				stateNext.StatePlayer = SPplaying
			case ActMute:
				stateNext.StateMute = SMmuted
			case ActUnmute:
				stateNext.StateMute = SMnormal
			case ActTerminate:
				stateNext.StatePlayer = SPoff
				stateNext.StateMute = SMnormal
				uriPlaying = ""
			}
		case SPrestart:
			switch st.Action {
			case ActTerminate:
				stateNext.StatePlayer = SPplaying
			}
		}

		log.Println("Calculated next state ", stateNext.StatePlayer.String(), stateNext.StateMute.String())
		ntfyChange := false
		if stateNext.StatePlayer != SPundef {
			log.Println("State trigger a change")
			stateCurrent = stateNext.StatePlayer
			stateNext.CurrURI = uriPlaying
			if stateNext.StateMute == SMundef {
				stateNext.StateMute = muteStateCurrent
			}
			ntfyChange = true
		} else if stateNext.StateMute != SMundef {
			log.Println("State trigger a mute change")
			stateNext.StatePlayer = stateCurrent
			ntfyChange = true
		}
		if stateNext.StateMute != SMundef {
			muteStateCurrent = stateNext.StateMute
		}
		if olduriPlaying != uriPlaying {
			stateNext.CurrURI = uriPlaying
			ntfyChange = true
		}
		if ntfyChange {
			for _, worker := range workers {
				worker.ChStatus <- &stateNext
			}
		} else {
			log.Println("Ignored action ", st.Action.String())
		}
		if st.ChStateRsp != nil {
			st.ChStateRsp <- &stateNext
		}
	}
}
