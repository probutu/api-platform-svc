package service

import "github.com/gin-gonic/gin"

type IUserService interface{}

type userService struct{}

func NewUserService() IUserService {
	return &userService{}
}

func (svc *userService) HandleList(c *gin.Context) {}
