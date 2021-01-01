package fileplayer

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aaaasmile/live-streamer/db"
	"github.com/aaaasmile/live-streamer/web/idl"
	"github.com/aaaasmile/live-streamer/web/live/omx/omxstate"
)

type infoFile struct {
	Title         string
	Description   string
	DurationInSec int
	TrackDuration string
	TrackPosition string
	TrackStatus   string
}

type FilePlayer struct {
	URI     string
	Info    *infoFile
	chClose chan struct{}
}

func (fp *FilePlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "/home") &&
		(strings.Contains(uri, ".mp3") || strings.Contains(uri, ".ogg") || strings.Contains(uri, ".wav")) {
		log.Println("this is a music file ", uri)
		fp.URI = uri
		return true
	}
	return false
}

func (fp *FilePlayer) GetStatusSleepTime() int {
	return 300
}

func (fp *FilePlayer) GetURI() string {
	return fp.URI
}
func (fp *FilePlayer) GetTitle() string {
	if fp.Info != nil {
		return fp.Info.Title
	}
	return ""
}
func (fp *FilePlayer) GetDescription() string {
	if fp.Info != nil {
		return fp.Info.Description
	}
	return ""
}
func (fp *FilePlayer) Name() string {
	return "file"
}
func (fp *FilePlayer) GetStreamerCmd() string {
	//args := strings.Join(cmdLineArr, " ")
	cmd := fmt.Sprintf("cvlc %s %s", fp.URI, `--sout="#transcode{vcodec=none,acodec=mp3,ab=128,channels=2,samplerate=44100}:http{mux=mp3,dst=:5550/stream.mp3}" --sout-keep`)
	return cmd
}
func (fp *FilePlayer) CheckStatus(chDbOperation chan *idl.DbOperation) error {
	st := &omxstate.StateOmx{}

	if fp.Info == nil {
		info := infoFile{
			// TODO read from db
		}
		info.DurationInSec, _ = strconv.Atoi(st.TrackDuration)
		info.TrackDuration = time.Duration(int64(info.DurationInSec) * int64(time.Second)).String()
		hi := db.HistoryItem{
			URI:           fp.URI,
			Title:         info.Title,
			Description:   info.Description,
			DurationInSec: info.DurationInSec,
			Type:          fp.Name(),
			Duration:      info.TrackDuration,
		}
		dop := idl.DbOperation{
			DbOpType: idl.DbOpHistoryInsert,
			Payload:  hi,
		}
		chDbOperation <- &dop
		fp.Info = &info
		log.Println("file-player info status set")
	}

	fp.Info.TrackPosition = st.TrackPosition
	fp.Info.TrackStatus = st.TrackStatus
	log.Println("Status set to ", fp.Info)
	return nil
}

func (fp *FilePlayer) CreateStopChannel() chan struct{} {
	if fp.chClose == nil {
		fp.chClose = make(chan struct{})
	}
	return fp.chClose
}

func (fp *FilePlayer) GetCmdStopChannel() chan struct{} {
	return fp.chClose
}

func (fp *FilePlayer) CloseStopChannel() {
	if fp.chClose != nil {
		close(fp.chClose)
		fp.chClose = nil
	}
}

func (fp *FilePlayer) GetTrackDuration() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackDuration, true
	}
	return "", false

}
func (fp *FilePlayer) GetTrackPosition() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackPosition, true
	}
	return "", false

}
func (fp *FilePlayer) GetTrackStatus() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackStatus, true
	}
	return "", false
}
