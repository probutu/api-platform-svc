package service

import (
	"net/http"
	"probutu/api-platform-svc/model"
	"probutu/api-platform-svc/repository"

	"github.com/gin-gonic/gin"
)

type IRequestService interface {
	HandleList(c *gin.Context)
	HandleCreate(c *gin.Context)
}

type requestService struct {
	requestRepository repository.IRequestRepository
}

func NewRequestService(
	requestRepository repository.IRequestRepository,
) IRequestService {
	return &requestService{
		requestRepository,
	}
}

func (svc *requestService) HandleList(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	request, total, err := svc.requestRepository.Find(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalData": total,
		"data":      request,
	})
}

func (svc *requestService) HandleCreate(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	requests, err := svc.requestRepository.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, requests)
}
