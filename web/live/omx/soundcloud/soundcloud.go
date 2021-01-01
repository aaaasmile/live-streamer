package soundcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type SoundCloudClient struct {
	ClientID  string
	AuthToken string
	UserAgent string
}

type MeInfo struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

type TrackInfo struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	TagList string `json:"tag_list"`
	Title   string `json:"title"`
	URI     string `json:"uri"`
}

func (sc *SoundCloudClient) GetReq(part string, respJs interface{}) error {
	url := "https://api.soundcloud.com/" + part
	client := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("user-agent", sc.UserAgent)
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", sc.AuthToken))
	req.Header.Set("client_id", sc.ClientID)
	req.Header.Set("Limit", "25")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 401 {
		return fmt.Errorf("Host returned: 401 Unauthorized")
	}
	if res.StatusCode == 402 {
		return fmt.Errorf("Host returned: 402 Not found")
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//fmt.Println("*** Body is ", string(body))
	//fmt.Println("**  res is ", res.StatusCode, res.Header)

	if err := json.Unmarshal(rawbody, respJs); err != nil {
		fmt.Println("*** Body is ", string(rawbody))
		return err
	}
	log.Println("Res info: ", respJs)
	return nil
}
