package upload

import (
	"mime/multipart"
	"net/http"
)

type Uploader interface {
	UploadFile(r *http.Request, file multipart.File, handler *multipart.FileHeader) error
}
