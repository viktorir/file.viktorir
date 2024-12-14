package database

import "file.viktorir/internal/model"

type FileAdapter interface {
	Insert(file model.File) (id int64, err error)
	GetFullLink(userId int, fileType, fileName string) (file model.File, err error)
	GetShortLink(shortLink string) (model.File, error)
	Close() error
}
