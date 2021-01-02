package you

import (
	"testing"
)

func TestReadLinkDescription(t *testing.T) {
	URL := "https://www.youtube.com/watch?v=szL8QmWgCGo&ab_channel=RelaxChilloutMusic"
	res, err := readLinkDescription(URL, "/home/igors/tmp/info")
	if err != nil {
		t.Error("Somenthing wrong: ", err)
		return
	}
	//fmt.Println(res)
	if res.Duration != 14643 {
		t.Error("Duration is not as expected")
	}
}
