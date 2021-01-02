package soundfile

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aaaasmile/live-streamer/db"
	"github.com/aaaasmile/live-streamer/web/idl"
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
	URI        string
	Info       *infoFile
	StreamDest string
	chClose    chan struct{}
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
	// %s => 5550/stream.mp3
	raw := `--play-and-exit --sout="#transcode{vcodec=none,acodec=mp3,ab=128,channels=2,samplerate=44100}:http{mux=mp3,dst=:%s}" --sout-keep`
	paraStr := fmt.Sprintf(raw, fp.StreamDest)
	cmd := fmt.Sprintf("cvlc %s %s", fp.URI, paraStr)
	return cmd
}
func (fp *FilePlayer) CheckStatus(chDbOperation chan *idl.DbOperation) error {
	if fp.Info == nil {
		info := infoFile{
			// TODO read from db
		}
		info.DurationInSec = 0 // TODO
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
