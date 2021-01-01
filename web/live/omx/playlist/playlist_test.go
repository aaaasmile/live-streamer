package playlist

import (
	"path/filepath"
	"testing"
	"time"
)

func setup() {
	dirPlaylistData = "../../../../playlist-data"
}

func TestCreateDefault(t *testing.T) {
	setup()
	pli := PlayItem{
		URI:      "http://stream.srg-ssr.ch/m/rsc_de/aacp_96",
		Info:     "Radio Swiss Classic",
		ItemType: PITRadio,
	}
	list := make([]*PlayItem, 0)
	list = append(list, &pli)

	strl := PlayList{
		Name:    "RadioCH",
		List:    list,
		Created: time.Now().Format("02.01.2006 15:04:05"),
	}
	playlistName := "default"
	strl.SavePlaylist(playlistName)

	if err := CheckIfPlaylistExist("default"); err != nil {

		t.Error("Play list not created", playlistName, err)
	}
}

func TestCreateChello(t *testing.T) {
	setup()
	list := make([]*PlayItem, 0)

	titles := []string{"Serotonin.mp3",
		"Beautiful Relaxing Music for Stress Relief • Meditation Music, Sleep Music, Ambient Study Music-2OEL4P1Rz04.mp",
		"MentalEnergy.mp3", "Happiness Frequency-Serotonin.mp3", "Serene-Assemby-chello.mp3", "chello_escape.mp3"}

	dirroot := "/home/igors/Music/youtube/relax"
	for _, item := range titles {
		name := item
		if len(item) > 30 {
			name = name[0:30] + "..."
		}
		pli := PlayItem{
			URI:      filepath.Join(dirroot, item),
			Info:     name,
			ItemType: PITMp3,
		}

		list = append(list, &pli)
	}

	strl := PlayList{
		Name:    "Chello",
		List:    list,
		Created: time.Now().Format("02.01.2006 15:04:05"),
	}
	playlistName := strl.Name + ".json"
	strl.SavePlaylist(playlistName)

	if err := CheckIfPlaylistExist(playlistName); err != nil {

		t.Error("Play list not created", playlistName, err)
	}
}
