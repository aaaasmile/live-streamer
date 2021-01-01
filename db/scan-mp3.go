package db

import (
	"log"
	"os"

	"github.com/dhowden/tag"
)

func ScanFile(fileName string) error {

	info, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	meta, err := tag.ReadFrom(info)
	if err != nil {
		return err
	}
	log.Println("Recognized ", meta.Title(), meta.FileType(), meta.Album(), meta.Artist(), meta.AlbumArtist())
	return nil
}
