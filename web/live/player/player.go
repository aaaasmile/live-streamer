package player

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/aaaasmile/live-streamer/web/idl"
	"github.com/aaaasmile/live-streamer/web/live/player/playlist"
)

type StrPlayer struct {
	mutex         *sync.Mutex
	state         StateOmx
	chDbOperation chan *idl.DbOperation
	PlayList      *playlist.LLPlayList
	Providers     map[string]idl.StreamProvider
	ChAction      chan *ActionDef
}

func NewStrPlayer(chDbop chan *idl.DbOperation) *StrPlayer {
	cha := make(chan *ActionDef)
	res := StrPlayer{
		mutex:         &sync.Mutex{},
		chDbOperation: chDbop,
		Providers:     make(map[string]idl.StreamProvider),
		ChAction:      cha,
	}

	return &res
}

func (op *StrPlayer) ListenOmxState(statusCh chan *StateOmx) {
	log.Println("start listenplayer. Waiting for status change in StrPlayer")
	for {
		st := <-statusCh
		op.mutex.Lock()
		log.Println("Set StrPlayer state ", st)
		if st.StatePlayer == SPoff {
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

func (op *StrPlayer) GetTrackDuration() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		if td, ok := prov.GetTrackDuration(); ok {
			return td
		}
	}

	return op.state.TrackDuration
}

func (op *StrPlayer) GetTrackPosition() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		if td, ok := prov.GetTrackPosition(); ok {
			return td
		}
	}
	return op.state.TrackPosition
}

func (op *StrPlayer) GetTrackStatus() string {
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

func (op *StrPlayer) GetStatePlaying() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.StatePlayer.String()
}

func (op *StrPlayer) GetStateTitle() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetTitle()
	}

	return ""
}

func (op *StrPlayer) GetStateDescription() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetDescription()
	}

	return ""
}

func (op *StrPlayer) GetCurrURI() string {
	log.Println("getCurrURI")
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.CurrURI
}

func (op *StrPlayer) StartPlay(URI string, prov idl.StreamProvider) error {
	var err error
	if op.PlayList, err = playlist.CreatePlaylistFromProvider(URI, prov); err != nil {
		return err
	}
	log.Println("StartPlay ", URI)

	return op.startPlayListCurrent(prov)
}

func (op *StrPlayer) PreviousTitle() (string, error) {
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

func (op *StrPlayer) NextTitle() (string, error) {
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

func (op *StrPlayer) CheckStatus(uri string) error {
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

func (op *StrPlayer) PowerOff() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	log.Println("Power off, terminate StrPlayer with signal kill")
	op.freeAllProviders()
	return nil
}

func (op *StrPlayer) freeAllProviders() {
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

func (op *StrPlayer) execCommand(uri, cmdText string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(cmdText string, actCh chan *ActionDef, uri string, chstop chan struct{}) {
		log.Println("Submit the command in background ", cmdText)
		cmd := exec.Command("bash", "-c", cmdText)
		actCh <- &ActionDef{
			URI:    uri,
			Action: ActPlaying,
		}

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

		if err := cmd.Start(); err == nil {
			log.Println("PID started ", cmd.Process.Pid)
			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
				log.Println("Wait ist terminated")
			}()

			select {
			case <-chstop:
				log.Println("Received stop signal, kill parent and child processes")
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
					log.Println("Error on killing the process ", err)
				}
			case err := <-done:
				log.Println("Process finished")
				if err != nil {
					log.Println("Error on process termination =>", err)
				}
				log.Println(string(stderrBuf.Bytes()))
				log.Println(string(stdoutBuf.Bytes()))
			}
			log.Println("Exit from waiting command execution")

		} else {
			log.Println("ERROR cmd.Start() failed with", err)
		}

		log.Println("Player has been terminated. Cmd was ", cmdText)
		actCh <- &ActionDef{
			URI:    uri,
			Action: ActTerminate,
		}

	}(cmdText, op.ChAction, uri, chstop)
}

func (op *StrPlayer) startPlayListCurrent(prov idl.StreamProvider) error {
	log.Println("Start current item ", op.PlayList)
	var curr *playlist.PlayItem
	var ok bool
	if curr, ok = op.PlayList.CheckCurrent(); !ok {
		return nil
	}
	log.Println("Current item is ", curr)
	op.mutex.Lock()
	defer op.mutex.Unlock()

	curURI := op.state.CurrURI
	if curURI != "" {
		log.Println("Shutting down the current player of ", curURI)
		if pp, ok := op.Providers[curURI]; ok {
			chStop := pp.GetCmdStopChannel()
			if chStop != nil {
				chStop <- struct{}{}
				pp.CloseStopChannel()
			}
			delete(op.Providers, curURI)
		}
	}
	uri := prov.GetURI()
	op.Providers[uri] = prov

	log.Println("Start player wit URI ", uri)

	cmd := prov.GetStreamerCmd()
	log.Println("Start the command: ", cmd)
	op.execCommand(uri, cmd, prov.CreateStopChannel())

	return nil
}
