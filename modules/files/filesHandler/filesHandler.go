package filesHandler

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type FileHandlerErrCode string

const (
	uploadFilesErr FileHandlerErrCode = "files-001"
	deleteFileErr  FileHandlerErrCode = "files-002"
)

type IFileHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type fileHandler struct {
	cfg         config.IConfig
	fileUsecase filesUsecase.IFilesUsecase
}

func FileHandler(cfg config.IConfig, fileUsecase filesUsecase.IFilesUsecase) IFileHandler {
	return &fileHandler{
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}

func (h *fileHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	destination := form.Value["destination"]

	// files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		// check file extension
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFilesErr),
				"invalid file extension",
			).Res()
		}
		// check file size
		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFilesErr),
				fmt.Sprintf("file size must less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: fmt.Sprintf("%s/%s", destination, filename),
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := h.fileUsecase.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadFilesErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

func (h *fileHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)

	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteFileErr),
			err.Error(),
		).Res()
	}

	if err := h.fileUsecase.DeleteFileOnGCP(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteFileErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
