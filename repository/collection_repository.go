package repository

import (
	"context"
	"errors"
	"probutu/api-platform-svc/model"

	"gorm.io/gorm"
)

type ICollectionRepository interface {
	Find(context.Context, model.Collection) (model.Collections, int64, error)
	FindOne(context.Context, model.Collection) (*model.Collection, error)
	Create(context.Context, model.Collection) (*model.Collection, error)
	Update(context.Context, model.Collection) error
}

type collectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(
	db *gorm.DB,
) ICollectionRepository {
	return &collectionRepository{
		db,
	}
}

func (r *collectionRepository) Find(ctx context.Context, f model.Collection) (collections model.Collections, count int64, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Collection{}).
		Count(&count).Order("collection_name ASC").
		Preload("Folders", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Requests", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Headers").Preload("Responses")
			})
		}).
		Preload("Requests").
		Where(f).
		Find(&collections).Error; err != nil {
		return
	}

	return
}

func (r *collectionRepository) FindOne(ctx context.Context, f model.Collection) (collection *model.Collection, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Collection{}).
		Preload("Folders", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Requests", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Headers").Preload("Responses")
			})
		}).
		Preload("Requests").
		Where(f).
		First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *collectionRepository) Create(ctx context.Context, data model.Collection) (collection *model.Collection, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Collection{}).Create(&data).Error; err != nil {
		return nil, err
	}

	return r.FindOne(ctx, model.Collection{ID: data.ID})
}

func (r *collectionRepository) Update(ctx context.Context, collection model.Collection) (err error) {
	return r.db.WithContext(ctx).Model(&model.Collection{}).Save(&collection).Error
}
