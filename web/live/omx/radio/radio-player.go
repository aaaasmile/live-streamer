package radio

import (
	"fmt"
	"log"
	"strings"

	"github.com/aaaasmile/live-streamer/db"
	"github.com/aaaasmile/live-streamer/web/idl"
)

type infoFile struct {
	Title       string
	Description string
}

type RadioPlayer struct {
	URI     string
	Info    *infoFile
	chClose chan struct{}
}

func (rp *RadioPlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "http") &&
		(strings.Contains(uri, "mp3") || strings.Contains(uri, "aacp")) {
		log.Println("This is a streaming resource ", uri)
		rp.URI = uri
		return true
	}
	return false
}

func (rp *RadioPlayer) GetStatusSleepTime() int {
	return 500
}

func (rp *RadioPlayer) GetURI() string {
	return rp.URI
}
func (rp *RadioPlayer) GetTitle() string {
	if rp.Info != nil {
		return rp.Info.Title
	}
	return ""
}
func (rp *RadioPlayer) GetDescription() string {
	if rp.Info != nil {
		return rp.Info.Description
	}
	return ""
}
func (rp *RadioPlayer) Name() string {
	return "radio"
}
func (rp *RadioPlayer) GetStreamerCmd(cmdLineArr []string) string {
	args := strings.Join(cmdLineArr, " ")
	cmd := fmt.Sprintf("omxplayer %s %s", args, rp.URI)
	return cmd
}
func (rp *RadioPlayer) CheckStatus(chDbOperation chan *idl.DbOperation) error {
	if rp.Info == nil {
		info := infoFile{
			// TODO read radio info from db
		}
		hi := db.HistoryItem{
			URI:         rp.URI,
			Title:       info.Title,
			Description: info.Description,
			Type:        rp.Name(),
		}
		dop := idl.DbOperation{
			DbOpType: idl.DbOpHistoryInsert,
			Payload:  hi,
		}
		chDbOperation <- &dop
		rp.Info = &info
	}

	return nil
}

func (rp *RadioPlayer) CreateStopChannel() chan struct{} {
	if rp.chClose == nil {
		rp.chClose = make(chan struct{})
	}
	return rp.chClose
}

func (rp *RadioPlayer) GetCmdStopChannel() chan struct{} {
	return rp.chClose
}

func (rp *RadioPlayer) CloseStopChannel() {
	if rp.chClose != nil {
		close(rp.chClose)
		rp.chClose = nil
	}
}

func (rp *RadioPlayer) GetTrackDuration() (string, bool) {
	return "", false
}
func (rp *RadioPlayer) GetTrackPosition() (string, bool) {
	return "", false
}
func (rp *RadioPlayer) GetTrackStatus() (string, bool) {
	return "", false
}
