package live

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func handleHistoryRequest(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	paraReq := struct {
		PageIx   int `json:"pageix"`
		PageSize int `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("history Request", paraReq)
	list, err := liteDB.FetchHistory(paraReq.PageIx, paraReq.PageSize)
	if err != nil {
		return err
	}

	type HistoryItemRes struct {
		ID          int    `json:"id"`
		Type        string `json:"type"`
		PlayedAt    string `json:"playedAt"`
		Title       string `json:"title"`
		URI         string `json:"uri"`
		DurationStr string `json:"durationstr"`
	}

	res := struct {
		History []HistoryItemRes `json:"history"`
		PageIx  int              `json:"pageix"`
	}{
		History: make([]HistoryItemRes, 0),
		PageIx:  paraReq.PageIx,
	}
	for _, item := range list {
		pp := HistoryItemRes{
			ID:          item.ID,
			Type:        item.Type,
			PlayedAt:    item.Timestamp.Format("January 02, 2006 15:04:05"),
			Title:       item.Title,
			URI:         item.URI,
			DurationStr: item.Duration,
		}
		res.History = append(res.History, pp)
	}

	return writeResponseNoWsBroadcast(w, res)
}
