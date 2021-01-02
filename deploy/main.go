package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aaaasmile/live-streamer/deploy/depl"
)

var (
	defOutDir = "../../Zips"
)

func main() {
	const (
		pi3 = "pi3"
		pi4 = "pi4"
	)
	var outdir = flag.String("outdir", "",
		fmt.Sprintf("Output zip directory. If empty use the hardcoded one: %s\n", defOutDir))

	tt := []string{pi3, pi4}
	var target = flag.String("target", "",
		fmt.Sprintf("Target of deployment: %v", tt))

	flag.Parse()

	rootDirRel := ".."
	pathItems := []string{"live-streamer.bin", "static", "templates"}
	switch *target {
	case pi3:
		pathItems = append(pathItems, "deploy/config_files/pi3_config.toml")
	case pi4:
		pathItems = append(pathItems, "deploy/config_files/pi4_config.toml")
	default:
		log.Fatalf("Deployment target %s is not recognized or not specified", *target)
	}
	log.Printf("Create the zip package for target %s and out dir ", *target)

	outFn := getOutFileName(*outdir, *target)
	depl.CreateDeployZip(rootDirRel, pathItems, outFn, func(pathItem string) string {
		if strings.HasPrefix(pathItem, "deploy/config_files") {
			return "config.toml"
		}
		return pathItem
	})
}

func getOutFileName(outdir string, tgt string) string {
	if outdir == "" {
		outdir = defOutDir
	}
	var err error
	outdir, err = filepath.Abs(outdir)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		log.Fatalf("Directory for zip generation not found %s", outdir)
	}
	log.Println("Using zip out directoy", outdir)
	vn := depl.GetVersionNrFromFile("../web/idl/idl.go", "")
	log.Println("Version is ", vn)

	currentTime := time.Now()
	s := fmt.Sprintf("Live-streamer_%s_%s_%s.zip", strings.Replace(vn, ".", "-", -1), currentTime.Format("02012006-150405"), tgt) // current date-time stamp using 2006 date time format template
	s = filepath.Join(outdir, s)
	return s
}

func testGetVersion() {
	buf, err := ioutil.ReadFile("../web/idl/idl.go")
	if err != nil {
		log.Fatalln("Cannot read input file", err)
	}
	s := string(buf)
	fmt.Println(s)
	vn := depl.GetBuildVersionNr(s, "")
	if vn == "" {
		log.Fatalln("Version not found")
	}
	fmt.Println("Version is ", vn)
	//depl.TestLexer()
}
