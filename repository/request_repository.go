package repository

import (
	"context"
	"errors"
	"probutu/api-platform-svc/model"

	"gorm.io/gorm"
)

type IRequestRepository interface {
	Find(context.Context, model.Request) (model.Requests, int64, error)
	FindOne(context.Context, string) (*model.Request, error)
	Create(context.Context, model.Request) (*model.Request, error)
	Update(context.Context, model.Request) error
}

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(
	db *gorm.DB,
) IRequestRepository {
	return &requestRepository{
		db,
	}
}

func (r *requestRepository) Find(ctx context.Context, f model.Request) (requests model.Requests, count int64, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Request{}).
		Count(&count).Order("name ASC").
		Preload("Headers").
		Preload("Responses").
		Where(f).
		Find(&requests).Error; err != nil {
		return
	}

	return
}

func (r *requestRepository) FindOne(ctx context.Context, requestId string) (request *model.Request, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Request{}).
		Preload("Headers").
		Preload("Responses").
		Where("request_id = ?", requestId).
		First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *requestRepository) Create(ctx context.Context, data model.Request) (collection *model.Request, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Request{}).Create(&data).Error; err != nil {
		return nil, err
	}

	return r.FindOne(ctx, data.ID)
}

func (r *requestRepository) Update(ctx context.Context, request model.Request) (err error) {
	return r.db.WithContext(ctx).Model(&model.Request{}).Save(&request).Error
}
