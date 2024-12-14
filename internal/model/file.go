package model

import "time"

type File struct {
	ID          int
	Name        string
	Type        string
	Size        int64
	UploadedAt  time.Time
	UserID      int
	Path        string
	ShortLink   string
	Description string
	Tags        []string
	Hash        string
	Status      string
}
