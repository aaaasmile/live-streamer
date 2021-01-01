package you

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/conf"
)

type InfoLink struct {
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	Description string `json:"description"`
}

func readLinkDescription(URI, dirout string) (*InfoLink, error) {
	cmd := fmt.Sprintf("%s --write-info-json --skip-download --youtube-skip-dash-manifest %s", getYoutubePlayer(), URI)
	log.Println("info with ", cmd, dirout)
	olddir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer os.Chdir(olddir)

	os.Chdir(dirout)
	cmdex := exec.Command("bash", "-c", cmd)
	cmdex.Dir = conf.Current.TmpInfo
	out, err := cmdex.Output()
	if err != nil {
		log.Println("Info json error: ", err)
		return nil, err
	}
	time.Sleep(200 * time.Millisecond)

	strOut := string(out)
	log.Println("Write Json out ", strOut)
	fname := ""
	arr := strings.Split(strOut, ":")
	if len(arr) > 0 {
		fname = arr[len(arr)-1]
		log.Println("Splitted: ", fname, len(arr))
	} else {
		return nil, fmt.Errorf("Unexpected output")
	}

	filelist, err := ioutil.ReadDir(dirout)
	if err != nil {
		return nil, err
	}
	fullfname := ""
	for _, f := range filelist {
		if !f.IsDir() {
			//fmt.Println("** File: ", f.Name(), filepath.Ext(f.Name()))
			ext := filepath.Ext(f.Name())
			if ext == ".json" {
				fullfname = path.Join(dirout, f.Name())
				log.Println("Use the first json file in dir ", fullfname)
			}
		}
	}
	if fullfname == "" {
		log.Println("Try to recognize filename from output")
		fname = strings.Trim(fname, "\n")
		fname = strings.Trim(fname, " ")
		fullfname = path.Join(dirout, fname)
		fullfname = fmt.Sprintf("\"%s\"", fullfname)
	}

	log.Println("Read file: ", fullfname)

	inforaw, err := ioutil.ReadFile(fullfname)
	if err != nil {
		return nil, err
	}
	infoLink := InfoLink{}
	err = json.Unmarshal(inforaw, &infoLink)

	if err != nil {
		return nil, err
	}

	log.Println("Info: ", infoLink)

	err = os.Remove(fullfname)
	if err != nil {
		return nil, err
	}
	log.Println("Removed file ", fullfname)

	return &infoLink, nil
}
