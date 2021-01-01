package live

import (
	"fmt"
	"log"
	"net/http"
)

func handlePost(w http.ResponseWriter, req *http.Request) error {
	var err error
	lastPath := getURLForRoute(req.RequestURI)
	log.Println("Check the last path ", lastPath)
	switch lastPath {
	case "SetPowerState":
		err = handleSetPowerState(w, req, player)
	case "GetPlayerState":
		err = handlePlayerState(w, req, player)
	case "NextTitle":
		err = handleNextTitle(w, req, player)
	case "PreviousTitle":
		err = handlePreviousTitle(w, req, player)
	case "PlayUri":
		err = handlePlayUri(w, req, player)
	case "OSRequest":
		err = handleOSRequest(w, req)
	case "FetchHistory":
		err = handleHistoryRequest(w, req)
	default:
		return fmt.Errorf("%s method is not supported", lastPath)
	}

	return err
}
