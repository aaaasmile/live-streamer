package db

import (
	"fmt"
	"testing"
)

func TestInsertHistory(t *testing.T) {
	lietDb := LiteDB{
		DebugSQL:     true,
		SqliteDBPath: "./test/test-data.db",
	}

	err := lietDb.OpenSqliteDatabase()
	if err != nil {
		t.Error("Error open file ", err)
		return
	}
	title := "Buddha Lounge Chillout Music ◈ Buddha Bar Chill out Music ◈ Café Bar Restaurant Background Music Mix"
	descr := "Buddha Lounge is a wonderful blend of the most beautiful musical trends of the Occident and the Orie..."
	for i := 0; i < 10; i++ {
		dur := fmt.Sprintf("1:3%d", i)
		durinsec := 12443
		err = lietDb.CreateHistory("https://youtu.be/7eFEp8b8oC4", title, descr, dur, durinsec, "youtube")
		if err != nil {
			t.Error("Error on insert history", err)
			return
		}
	}
}

func TestFetchHistory(t *testing.T) {
	lietDb := LiteDB{
		DebugSQL:     true,
		SqliteDBPath: "./test/test-data.db",
	}

	err := lietDb.OpenSqliteDatabase()
	if err != nil {
		t.Error("Error open file ", err)
		return
	}

	pageSize := 3
	list, err := lietDb.FetchHistory(2, pageSize)
	if err != nil {
		t.Error("Error on fetch history", err)
		return
	}
	if len(list) != pageSize {
		t.Errorf("Expected %d items, but fetched %d", pageSize, len(list))
		return
	}
	for _, item := range list {
		fmt.Println("ID, Duration: ", item.ID, item.Duration)
	}
	//t.Error("Force output")
}
