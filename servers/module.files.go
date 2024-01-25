package servers

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesHandler"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesUsecase"
)

type IFilesModule interface {
	Init()
	Usecase() filesUsecase.IFilesUsecase
	Handler() filesHandler.IFileHandler
}

type filesModule struct {
	*moduleFactory
	usecase filesUsecase.IFilesUsecase
	handler filesHandler.IFileHandler
}

func (m *moduleFactory) FilesModule() IFilesModule {
	usecase := filesUsecase.FilesUsecase(m.s.cfg)
	handler := filesHandler.FileHandler(m.s.cfg, usecase)

	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *filesModule) Init() {
	router := f.r.Group("/files")

	router.Post("/upload", f.mid.JwtAuth(), f.mid.Authorize(2), f.handler.UploadFiles)
	router.Patch("/delete", f.mid.JwtAuth(), f.mid.Authorize(2), f.handler.DeleteFile)
}

func (f *filesModule) Usecase() filesUsecase.IFilesUsecase { return f.usecase }
func (f *filesModule) Handler() filesHandler.IFileHandler  { return f.handler }
