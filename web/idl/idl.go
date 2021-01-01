package idl

var (
	Appname = "live-streamer"
	Buildnr = "00.01.01.20210101-00"
)

type StreamProvider interface {
	IsUriForMe(uri string) bool
	GetStatusSleepTime() int
	GetURI() string
	GetTitle() string
	GetDescription() string
	Name() string
	GetStreamerCmd(cmdLineArr []string) string
	CheckStatus(chDbOperation chan *DbOperation) error
	CreateStopChannel() chan struct{}
	GetCmdStopChannel() chan struct{}
	CloseStopChannel()
	GetTrackDuration() (string, bool)
	GetTrackPosition() (string, bool)
	GetTrackStatus() (string, bool)
}

type DbOpType int

const (
	DbOpHistoryInsert DbOpType = iota
)

type DbOperation struct {
	DbOpType DbOpType
	Payload  interface{}
}
