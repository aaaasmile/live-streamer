package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaaasmile/live-streamer/util"
	_ "github.com/mattn/go-sqlite3"
)

type LiteDB struct {
	connDb       *sql.DB
	DebugSQL     bool
	SqliteDBPath string
}

type HistoryItem struct {
	ID                                int
	URI, Title, Description, Duration string
	Timestamp                         time.Time
	PlayPosition                      int
	DurationInSec                     int
	Type                              string
}

func (ld *LiteDB) OpenSqliteDatabase() error {
	var err error
	dbname := util.GetFullPath(ld.SqliteDBPath)
	if _, err := os.Stat(dbname); err != nil {
		return err
	}
	log.Println("Using the sqlite file: ", dbname)
	ld.connDb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	return nil
}

func (ld *LiteDB) FetchHistory(pageIx int, pageSize int) ([]HistoryItem, error) {
	q := `SELECT id,Timestamp,URI,Title,Description,Duration,PlayPosition,DurationInSec,Type
		  FROM History
		  ORDER BY Timestamp DESC 
		  LIMIT %d OFFSET %d;`
	offsetRows := pageIx * pageSize
	q = fmt.Sprintf(q, pageSize, offsetRows)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]HistoryItem, 0)
	var tss int64
	for rows.Next() {
		item := HistoryItem{}
		tss = 0
		if err := rows.Scan(&item.ID, &tss, &item.URI, &item.Title,
			&item.Description, &item.Duration, &item.PlayPosition,
			&item.DurationInSec, &item.Type); err != nil {
			return nil, err
		}
		item.Timestamp = time.Unix(tss, 0)
		res = append(res, item)
	}
	return res, nil
}

func (ld *LiteDB) CreateHistory(uri, title, description, duration string, durinsec int, tt string) error {
	item := HistoryItem{
		URI:           uri,
		Title:         title,
		Description:   description,
		Duration:      duration,
		DurationInSec: durinsec,
		Type:          tt,
	}
	return ld.InsertHistoryItem(&item)
}

func (ld *LiteDB) InsertHistoryItem(item *HistoryItem) error {
	q := `INSERT INTO History(Timestamp,URI,Title,Description,Duration,PlayPosition,DurationInSec,Type) VALUES(?,?,?,?,?,?,?,?);`
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	now := time.Now()
	sqlres, err := stmt.Exec(now.Local().Unix(), item.URI, item.Title, item.Description,
		item.Duration, 0, item.DurationInSec, item.Type)
	if err != nil {
		return err
	}
	log.Println("History inserted: ", sqlres)
	return nil
}
