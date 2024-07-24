package service

import (
	"net/http"
	"probutu/api-platform-svc/model"
	"probutu/api-platform-svc/repository"

	"github.com/gin-gonic/gin"
)

type IFolderService interface {
	HandleList(c *gin.Context)
	HandleCreate(c *gin.Context)
}

type folderService struct {
	folderRepository repository.IFolderRepository
}

func NewFolderService(
	folderRepository repository.IFolderRepository,
) IFolderService {
	return &folderService{
		folderRepository,
	}
}

func (svc *folderService) HandleList(c *gin.Context) {
	var req model.Folder
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	folders, total, err := svc.folderRepository.Find(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalData": total,
		"data":      folders,
	})
}

func (svc *folderService) HandleCreate(c *gin.Context) {
	var req model.Folder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	folder, err := svc.folderRepository.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, folder)
}
