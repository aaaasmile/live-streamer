package playlist

import (
	"log"

	"github.com/aaaasmile/live-streamer/web/idl"
)

type PlayinfoLast struct {
	Playlist string
	URI      string
}

type PlayItem struct {
	URI            string
	StreamProvider idl.StreamProvider
}

type PlayList struct {
	Name    string
	Created string
	List    []*PlayItem
}

type LLPlayItem struct {
	PlayItem *PlayItem
	Next     *LLPlayItem
	Previous *LLPlayItem
}

func NewLLPlayItem(nx, pr *LLPlayItem, plit *PlayItem) *LLPlayItem {
	res := LLPlayItem{
		PlayItem: plit,
		Next:     nx,
		Previous: pr,
	}
	return &res
}

type LLPlayList struct {
	Name      string
	Count     int
	FirstItem *LLPlayItem
	LastItem  *LLPlayItem
	CurrItem  *LLPlayItem
}

func (ll *LLPlayList) First() {
	ll.CurrItem = ll.FirstItem
}

func (ll *LLPlayList) Last() {
	ll.CurrItem = ll.LastItem
}

func (ll *LLPlayList) Next() (*PlayItem, bool) {
	if ll.CurrItem == nil {
		return nil, false
	}
	ll.CurrItem = ll.CurrItem.Next
	if ll.CurrItem == nil {
		ll.CurrItem = ll.FirstItem
	}
	if ll.CurrItem != nil {
		return ll.CurrItem.PlayItem, ll.CurrItem.PlayItem != nil
	} else {
		return nil, false
	}
}

func (ll *LLPlayList) Previous() (*PlayItem, bool) {
	if ll.CurrItem == nil {
		return nil, false
	}
	ll.CurrItem = ll.CurrItem.Previous
	if ll.CurrItem == nil {
		ll.CurrItem = ll.LastItem
	}
	if ll.CurrItem != nil {
		return ll.CurrItem.PlayItem, ll.CurrItem.PlayItem != nil
	} else {
		return nil, false
	}
}

func (ll *LLPlayList) CheckCurrent() (*PlayItem, bool) {
	if ll.FirstItem == nil ||
		ll.LastItem == nil ||
		ll.CurrItem == nil {
		log.Println("Invalid current item.")
		return nil, false
	}
	if ll.CurrItem.PlayItem == nil {
		return nil, false
	}
	return ll.CurrItem.PlayItem, true
}
func (ll *LLPlayList) IsEmpty() bool {
	return ll.FirstItem == nil ||
		ll.LastItem == nil
}

func CreatePlaylistFromProvider(URI string, prov idl.StreamProvider) (*LLPlayList, error) {
	res := &LLPlayList{}

	item := PlayItem{
		URI:            URI,
		StreamProvider: prov,
		//ItemType: it,
	}
	res.Name = URI
	res.CurrItem = NewLLPlayItem(nil, nil, &item)
	res.LastItem = res.CurrItem
	res.FirstItem = res.CurrItem

	return res, nil
}
