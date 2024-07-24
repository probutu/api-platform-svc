package service

import (
	"net/http"
	"probutu/api-platform-svc/model"
	"probutu/api-platform-svc/repository"

	"github.com/gin-gonic/gin"
)

type ICollectionService interface {
	HandleList(c *gin.Context)
	HandleCreate(c *gin.Context)
	HandleUpdate(c *gin.Context)
}

type collectionService struct {
	collectionRepository repository.ICollectionRepository
}

func NewCollectionService(
	collectionRepository repository.ICollectionRepository,
) ICollectionService {
	return &collectionService{
		collectionRepository,
	}
}

func (svc *collectionService) HandleList(c *gin.Context) {
	var req model.Collection
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	collections, total, err := svc.collectionRepository.Find(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalData": total,
		"data":      collections,
	})
}

func (svc *collectionService) HandleCreate(c *gin.Context) {
	var req model.Collection
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	project, err := svc.collectionRepository.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (svc *collectionService) HandleUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
