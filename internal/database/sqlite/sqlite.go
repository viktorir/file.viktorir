package sqlite

import (
	"database/sql"
	"file.viktorir/internal/model"
	"log"
	"os"
	path2 "path"
	"strings"
)

type Sqlite struct {
	db *sql.DB
}

func Init() (*Sqlite, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbDir := path2.Join(homeDir, "data")
	dbPath := path2.Join(dbDir, "metadata.db")

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	file, err := os.ReadFile("./internal/database/sqlite/create.sql")
	if err != nil {
		return nil, err
	}
	queries := strings.Split(string(file), ";")
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Sqlite3 connect!", dbPath)

	return &Sqlite{db: db}, nil
}

func (s Sqlite) Insert(f model.File) (id int64, err error) {
	tags := strings.Join(f.Tags, ",")
	query := "INSERT INTO files (name, type, size, uploaded_at, user_id, path, short_link, description, tags, hash, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id"
	result, err := s.db.Exec(query, f.Name, f.Type, f.Size, f.UploadedAt, f.UserID, f.Path, f.ShortLink, f.Description, tags, f.Hash, f.Status)
	if err != nil {
		return 0, err
	}
	id, err = result.LastInsertId()
	return
}

func (s Sqlite) GetFullLink(userId int, fileType, fileName string) (f model.File, err error) {
	var tags string
	query := "SELECT id, name, type, size, uploaded_at, user_id, path, short_link, description, tags, hash, status FROM files WHERE user_id = ? AND type = ? AND name = ? LIMIT 1;"
	err = s.db.QueryRow(query, userId, fileType, fileName).Scan(&f.ID, &f.Name, &f.Type, &f.Size, &f.UploadedAt, &f.UserID, &f.Path, &f.ShortLink, &f.Description, &tags, &f.Hash, &f.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.File{}, nil
		}
		return model.File{}, err
	}
	f.Tags = strings.Split(tags, ",")

	return f, nil
}

func (s Sqlite) GetShortLink(shortLink string) (f model.File, err error) {
	var tags string
	query := "SELECT id, name, type, size, uploaded_at, user_id, path, short_link, description, tags, hash, status FROM files WHERE short_link = ? LIMIT 1;"
	err = s.db.QueryRow(query, shortLink).Scan(&f.ID, &f.Name, &f.Type, &f.Size, &f.UploadedAt, &f.UserID, &f.Path, &f.ShortLink, &f.Description, &tags, &f.Hash, &f.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.File{}, nil
		}
		return model.File{}, err
	}
	f.Tags = strings.Split(tags, ",")

	return f, nil
}

func (s Sqlite) Close() error {
	return s.db.Close()
}
