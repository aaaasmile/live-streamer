package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func handleOSRequest(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	reqCmd := struct {
		Cmd string `json:"cmd"`
	}{}
	if err := json.Unmarshal(rawbody, &reqCmd); err != nil {
		return err
	}
	log.Println("Os Request ", reqCmd.Cmd)
	var cmdStr string
	switch reqCmd.Cmd {
	case "reboot":
		cmdStr = "sudo reboot"
		log.Println("Reboot request received")
		execCmdInBackground(cmdStr)
	case "shutdown":
		cmdStr = "sudo shutdown"
		log.Println("Shutdown request received")
		execCmdInBackground(cmdStr)
	case "service-restart":
		// NOTE: stop or start make no sense because after a stop
		// no more commands could be handled.
		// For developement use the console
		log.Println("Service restart request received")
		cmdStr = "sudo systemctl restart live-omxctrl"
		execCmdInBackground(cmdStr)
	case "kill-all-omx":
		log.Println("Kill all Omxplayer")
		cmdStr = "sudo killall omxplayer.bin"
		execCmdInBackground(cmdStr)
	default:
		return fmt.Errorf("Command not recognized %s", reqCmd.Cmd)
	}

	resObj := struct {
		Msg string `json:"msg"`
	}{
		Msg: fmt.Sprintf("Cmd %q submitted", cmdStr),
	}

	return writeResponse(w, resObj)
}

func execCmdInBackground(cmdStr string) {
	go func() {
		time.Sleep(2 * time.Second)
		cmdExec := exec.Command("bash", "-c", cmdStr)
		out, err := cmdExec.Output()
		if err != nil {
			log.Println("Error on shutdown: ", err)
		}
		res := string(out)
		log.Println("Command executed ", cmdStr, res)
	}()
}
