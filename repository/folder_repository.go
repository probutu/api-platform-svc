package repository

import (
	"context"
	"errors"
	"probutu/api-platform-svc/model"

	"gorm.io/gorm"
)

type IFolderRepository interface {
	Find(context.Context, model.Folder) (model.Folders, int64, error)
	FindOne(context.Context, string) (*model.Folder, error)
	Create(context.Context, model.Folder) (*model.Folder, error)
	Update(context.Context, model.Folder) error
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) IFolderRepository {
	return &folderRepository{db}
}

func (r *folderRepository) Find(ctx context.Context, f model.Folder) (requests model.Folders, count int64, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Folder{}).
		Count(&count).Order("name ASC").
		Preload("Requests").
		Where(f).
		Find(&requests).Error; err != nil {
		return
	}

	return
}

func (r *folderRepository) FindOne(ctx context.Context, requestId string) (request *model.Folder, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Folder{}).
		Preload("Requests").
		Where("folder_id = ?", requestId).
		First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *folderRepository) Create(ctx context.Context, data model.Folder) (collection *model.Folder, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Folder{}).Create(&data).Error; err != nil {
		return nil, err
	}

	return r.FindOne(ctx, data.ID)
}

func (r *folderRepository) Update(ctx context.Context, request model.Folder) (err error) {
	return r.db.WithContext(ctx).Model(&model.Folder{}).Save(&request).Error
}
