package page_repository

import (
	"errors"

	"github.com/kaiiorg/page-watcher/pkg/config"
	"github.com/kaiiorg/page-watcher/pkg/models"
	"github.com/kaiiorg/page-watcher/pkg/repositories/gorm_logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlitePageRepository struct {
	config *config.DB
	db     *gorm.DB
}

func NewSqlitePageRepository(dbConfig *config.DB) (*SqlitePageRepository, error) {
	db, err := gorm.Open(
		sqlite.Open(dbConfig.Path),
		&gorm.Config{
			Logger: gorm_logger.NewGormLogger(log.With().Str("component", "db-gorm").Logger(), zerolog.GlobalLevel()),
		},
	)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Page{})
	if err != nil {
		return nil, err
	}

	spr := &SqlitePageRepository{
		config: dbConfig,
		db:     db,
	}

	return spr, nil
}

func (spr *SqlitePageRepository) GetLatestChange(name string) (*models.Page, error) {
	dbPage := &models.Page{}
	err := spr.db.Where("name = ?", name).Last(dbPage).Error
	if err != nil {
		return nil, err
	}
	return dbPage, nil
}

func (spr *SqlitePageRepository) SaveChange(page *models.Page) error {
	err := spr.db.Save(page).Error
	if err != nil {
		return err
	}

	// Find the oldest record that we want to keep
	oldestToKeep := &models.Page{}
	err = spr.db.
		Select("created_at").
		Where("name = ?", page.Name).
		Order("created_at DESC").
		Limit(1).
		Offset(spr.config.Retain). // newest record must always be kept + how many records the user wants us to keep
		Row().
		Scan(oldestToKeep)
	if err != nil {
		// If we didn't find one, then we're below the number of records we're configured to keep
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}

	// Delete any record older than this record
	err = spr.db.
		Where("name = ? AND created_at < ?", page.Name, oldestToKeep.CreatedAt).
		Delete(&models.Page{}).Error

	return err
}
