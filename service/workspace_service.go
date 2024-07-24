package service

import (
	"net/http"
	"probutu/api-platform-svc/model"
	"probutu/api-platform-svc/repository"

	"github.com/gin-gonic/gin"
)

type IWorkspaceService interface {
	HandleList(c *gin.Context)
	HandleCreate(c *gin.Context)
	HandleUpdate(c *gin.Context)
}

type workspaceService struct {
	workspaceRepository repository.IWorkspaceRepository
}

func NewWorkspaceService(
	workspaceRepository repository.IWorkspaceRepository,
) IWorkspaceService {
	return &workspaceService{
		workspaceRepository,
	}
}

func (svc *workspaceService) HandleList(c *gin.Context) {
	var req model.Workspace
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	workspaces, total, err := svc.workspaceRepository.Find(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalData": total,
		"data":      workspaces,
	})
}

func (svc *workspaceService) HandleCreate(c *gin.Context) {
	var req model.Workspace
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	project, err := svc.workspaceRepository.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (svc *workspaceService) HandleUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
