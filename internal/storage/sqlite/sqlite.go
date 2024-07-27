package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anosovs/datalink/internal/storage"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {	
	const op = "storage.sqlite.New"
	pwd, _ := os.Getwd()
	fileAbs := pwd + strings.Replace(storagePath, ".", "", 1)
	if _, err := os.Stat(fileAbs); err!= nil {
		os.MkdirAll("./storage", os.ModePerm)
		file, err := os.Create(fileAbs)
		if err !=nil {
			return nil, err
		}
		defer file.Close()
	}
	

    
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS messages (
			uuid VARCHAR(36) PRIMARY KEY,
			message TEXT NOT NULL,
			count SMALLINT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveMsg(msgToSave string, count int) (uid string, err error) {
	const op = "storage.sqlite.SaveMsg"

	uid = uuid.New().String()
	stmt, err := s.db.Prepare("INSERT INTO messages (uuid, message, count) VALUES (?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_ , err = stmt.Exec(uid, msgToSave, count)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return uid, nil
}

func (s *Storage) GetMsg(uid string) (msg string, cnt int, err error) {
	const op = "storage.sqlite.GetMsg"

	stmt, err := s.db.Prepare("SELECT message, count FROM messages WHERE uuid=?")
	if err != nil {
		return "", 0,  fmt.Errorf("%s: %w", op, err)
	}	

	err = stmt.QueryRow(uid).Scan(&msg, &cnt)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, storage.ErrUuidNotFount
	}
	if err != nil {
		return "", 0, fmt.Errorf("%s: %w", op, err)
	}
	return msg, cnt, nil
}

func (s *Storage) DeleteMsg(uid string) error {
	const op = "storage.sqlite.DeleteMsg"
	stmt, err := s.db.Prepare("DELETE FROM messages WHERE uuid = ? ")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(uid)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrUuidNotFount
	}
	return nil
}

func (s *Storage) DecreaseCountMsg(uid string) error {
	const op = "storage.sqlite.DecreaseCount"
	stmt, err := s.db.Prepare("UPDATE messages SET count = count-1 WHERE uuid=?;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(uid)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrUuidNotFount
	}
	return nil
}

func (s *Storage) DeleteOldMessages(days int) error {
	const op = "storage.sqlite.DeleteOldMessages"
	availableUntil := time.Now().Add(-time.Hour * 24 * time.Duration(days))
	stmt, err := s.db.Prepare("DELETE FROM messages WHERE created_at < ? ")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(availableUntil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}