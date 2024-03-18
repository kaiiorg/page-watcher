package page_repository

import "github.com/kaiiorg/page-watcher/pkg/models"

type PageRepository interface {
	GetLatestChange(name string) (*models.Page, error)
	SaveChange(page *models.Page) error
}
