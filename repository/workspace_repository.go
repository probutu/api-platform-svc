package repository

import (
	"context"
	"errors"
	"probutu/api-platform-svc/model"

	"gorm.io/gorm"
)

type IWorkspaceRepository interface {
	Find(context.Context, model.Workspace) ([]model.Workspace, int64, error)
	FindOne(context.Context, model.Workspace) (*model.Workspace, error)
	Create(context.Context, model.Workspace) (*model.Workspace, error)
	Update(context.Context, model.Workspace) error
}

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(
	db *gorm.DB,
) IWorkspaceRepository {
	return &workspaceRepository{
		db,
	}
}

func (r *workspaceRepository) Find(ctx context.Context, f model.Workspace) (projects []model.Workspace, count int64, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Workspace{}).
		Count(&count).Order("workspace_name ASC").
		Preload("Collections", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Requests", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Headers").Preload("Responses")
			}).Preload("Folders", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Requests", func(tx *gorm.DB) *gorm.DB {
					return tx.Preload("Headers").Preload("Responses")
				})
			})
		}).
		Preload("Environments").
		Where(f).
		Find(&projects).Error; err != nil {
		return
	}

	return
}

func (r *workspaceRepository) FindOne(ctx context.Context, f model.Workspace) (project *model.Workspace, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Workspace{}).
		Preload("Collections", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Requests", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Headers").Preload("Responses")
			}).Preload("Folders", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Requests", func(tx *gorm.DB) *gorm.DB {
					return tx.Preload("Headers").Preload("Responses")
				})
			})
		}).
		Preload("Environments").
		Where(f).
		First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *workspaceRepository) Create(ctx context.Context, data model.Workspace) (project *model.Workspace, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Workspace{}).Create(&data).Error; err != nil {
		return nil, err
	}

	return r.FindOne(ctx, model.Workspace{ID: data.ID})
}

func (r *workspaceRepository) Update(ctx context.Context, project model.Workspace) (err error) {
	return r.db.WithContext(ctx).Model(&model.Workspace{}).Save(&project).Error
}
