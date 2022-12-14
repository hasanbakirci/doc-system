package helpers

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type FileResponse struct {
	FileName   string
	FileRename string
	Path       string
	Extension  string
	MimeType   string
}

func AddFile(file *multipart.FileHeader, types ...string) *FileResponse {

	src, err := file.Open()
	if err != nil {
		log.Error(err)
		errorHandler.Panic(http.StatusBadRequest, err.Error())
	}

	fileByte, err := ioutil.ReadAll(src)
	if err != nil {
		log.Error(err)
		errorHandler.Panic(http.StatusBadRequest, err.Error())
	}

	mt := mimetype.Detect(fileByte)
	if !handleMimeType(types, mt.String()) {
		errorHandler.Panic(http.StatusBadRequest, "wrong file format")
	}
	fileRename := fmt.Sprintf(uuid.New().String() + "-" + time.Now().Format("2006-01-02-15-04-05"))
	path := "upload/" + fileRename + mt.Extension()

	err = ioutil.WriteFile(path, fileByte, 0777)
	if err != nil {
		log.Error(err)
		errorHandler.Panic(http.StatusBadRequest, err.Error())
	}

	defer src.Close()

	return &FileResponse{
		FileName:   file.Filename,
		FileRename: fileRename,
		Path:       path,
		Extension:  mt.Extension(),
		MimeType:   mt.String(),
	}
}

func handleMimeType(types []string, mime string) bool {
	for _, r := range types {
		if r == mime {
			return true
		}
	}
	return false
}
