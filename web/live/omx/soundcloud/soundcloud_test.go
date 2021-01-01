package soundcloud

import "testing"

// 1) register to sound cloud
// 2) Get the token here:
// https://mopidy.com/ext/soundcloud/#authentication
// 3) This act like the mopidy extension
func getScHinstance() SoundCloudClient {
	return SoundCloudClient{
		ClientID:  "",
		AuthToken: "",
		UserAgent: "",
	}
}
func TestStream(t *testing.T) {

	sc := getScHinstance()

	respJs := MeInfo{}

	if err := sc.GetReq("me", &respJs); err != nil {
		t.Error("Error on connect", err)
		return
	}
	if respJs.ID != 917787916 {
		t.Error("Expected user id 917787916, but got: ", respJs)
	}
}

func TestTracks(t *testing.T) {
	sc := getScHinstance()
	respJs := make([]TrackInfo, 0)
	// query tracks for relax keyword
	if err := sc.GetReq("tracks?q=relax", &respJs); err != nil {
		t.Error("Error on connect", err)
		return
	}
	t.Error(len(respJs))
}

func TestTrackId(t *testing.T) {
	sc := getScHinstance()
	respJs := TrackInfo{}
	// Track id got from tracks query
	if err := sc.GetReq("tracks/62576046", &respJs); err != nil {
		t.Error("Error on connect", err)
		return
	}
	t.Error("OK?", respJs.URI)
}
