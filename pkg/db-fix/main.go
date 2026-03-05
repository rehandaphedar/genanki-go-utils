package main

import (
	"archive/zip"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func FixDb(inputPath, outputPath string) error {
	reader, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	for _, file := range reader.File {
		ioWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}

		if file.Name == "collection.anki2" {
			tempDb, _ := os.CreateTemp("", "anki-db-*.anki2")
			readCloser, _ := file.Open()
			io.Copy(tempDb, readCloser)
			readCloser.Close()
			tempDb.Close()
			defer os.Remove(tempDb.Name())

			db, err := sql.Open("sqlite3", tempDb.Name())
			if err != nil {
				return err
			}

			rows, _ := db.Query("SELECT id, tags FROM notes")
			type row struct {
				id   int64
				tags string
			}
			var updates []row
			for rows.Next() {
				var id int64
				var raw string
				rows.Scan(&id, &raw)
				if strings.HasPrefix(raw, "[") {
					var list []string
					if json.Unmarshal([]byte(raw), &list) == nil {
						updates = append(updates, row{id, " " + strings.Join(list, " ") + " "})
					}
				}
			}
			rows.Close()

			if len(updates) > 0 {
				transaction, _ := db.Begin()
				statement, _ := transaction.Prepare("UPDATE notes SET tags=? WHERE id=?")
				for _, update := range updates {
					statement.Exec(update.tags, update.id)
				}
				statement.Close()
				transaction.Commit()
			}

			cardRows, _ := db.Query("SELECT id FROM cards")
			var cardIDs []int64
			for cardRows.Next() {
				var id int64
				cardRows.Scan(&id)
				cardIDs = append(cardIDs, id)
			}
			cardRows.Close()

			if len(cardIDs) > 0 {
				transaction, _ := db.Begin()
				statement, _ := transaction.Prepare("UPDATE cards SET id=? WHERE id=?")
				for _, id := range cardIDs {
					statement.Exec(GenerateIntID(), id)
				}
				statement.Close()
				transaction.Commit()
			}

			db.Close()

			dbFile, _ := os.Open(tempDb.Name())
			io.Copy(ioWriter, dbFile)
			dbFile.Close()

		} else {
			readCloser, err := file.Open()
			if err != nil {
				return err
			}
			io.Copy(ioWriter, readCloser)
			readCloser.Close()
		}
	}

	return nil
}

func GenerateIntID() int64 {
	var b [8]byte
	rand.Read(b[:])
	return int64(binary.LittleEndian.Uint64(b[:]) & math.MaxInt64)
}
