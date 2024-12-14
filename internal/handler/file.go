package handler

import (
	"database/sql"
	"file.viktorir/internal/database"
	"file.viktorir/internal/model"
	"file.viktorir/pkg/hash"
	"file.viktorir/pkg/link"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileHandler struct {
	DB database.FileAdapter
}

func Init(db database.FileAdapter) *FileHandler {
	return &FileHandler{DB: db}
}

func (f FileHandler) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	userID, err := strconv.Atoi(ctx.FormValue("user_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}
	description := ctx.FormValue("description")
	tags := ctx.FormValue("tags")

	fileContent, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer fileContent.Close()

	buffer := make([]byte, 512)
	_, err = fileContent.Read(buffer)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	contentType := http.DetectContentType(buffer)

	homeDir, err := os.UserHomeDir()
	fileDir := path.Join(homeDir, "data", "files", strconv.Itoa(userID), contentType)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	hash, err := hash.GenerateToFile(fileContent)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	findFile, err := f.DB.GetFullLink(userID, contentType, file.Filename)
	if err != nil {
		if err == sql.ErrNoRows {
			findFile = model.File{}
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if findFile.Name != "" {
		return fiber.NewError(fiber.StatusConflict, "File already exists with the same name")
	}

	fileName := file.Filename
	filePath := path.Join(fileDir, fileName)
	if err := ctx.SaveFile(file, filePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	shortLink := link.GenerateShort(8)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileData := model.File{
		ID:          0,
		Name:        fileName,
		Type:        contentType,
		Size:        file.Size,
		UploadedAt:  time.Now(),
		UserID:      userID,
		Path:        filePath,
		ShortLink:   shortLink,
		Description: description,
		Tags:        strings.Split(tags, ","),
		Hash:        string(hash),
		Status:      "active",
	}

	id, err := f.DB.Insert(fileData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": id, "short_link": "/" + shortLink, "full_link": "/file/" + strconv.FormatInt(id, 10) + "/" + contentType + "/" + fileName})
}

func (f FileHandler) GetByShort(ctx *fiber.Ctx) error {
	shortLink := ctx.Params("short_link")
	if shortLink == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Short link is required")
	}

	fileData, err := f.DB.GetShortLink(shortLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "File not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if _, err := os.Stat(fileData.Path); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "File not found on server")
	}

	return ctx.SendFile(fileData.Path)
}

func (f FileHandler) GetByFull(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Params("user_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}
	fileType := ctx.Params("type")
	if fileType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File type is required")
	}
	fileSubtype := ctx.Params("subtype")
	if fileSubtype == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File subtype is required")
	}
	fileName := ctx.Params("name")
	if fileName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File name is required")
	}

	fileData, err := f.DB.GetFullLink(userID, fileType+"/"+fileSubtype, fileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "File not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if _, err := os.Stat(fileData.Path); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "File not found on server")
	}

	return ctx.SendFile(fileData.Path)
}
