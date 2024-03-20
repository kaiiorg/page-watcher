package watcher

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kaiiorg/page-watcher/pkg/config"
	"github.com/kaiiorg/page-watcher/pkg/models"
	"github.com/kaiiorg/page-watcher/pkg/repositories/page_repository"
	"github.com/kaiiorg/page-watcher/pkg/watcher/normalizer"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gorm.io/gorm"
)

type Watcher struct {
	config    *config.Config
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	normalizer     normalizer.Normalizer
	pageRepository page_repository.PageRepository
}

func New(config *config.Config) (*Watcher, error) {
	pageRepository, err := page_repository.NewSqlitePageRepository(config.DB)
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		config:         config,
		normalizer:     normalizer.New(),
		pageRepository: pageRepository,
	}
	w.ctx, w.ctxCancel = context.WithCancel(context.Background())

	return w, nil
}

func (w *Watcher) Close() {
	w.ctxCancel()
	w.wg.Wait()
}

func (w *Watcher) Watch() {
	for _, page := range w.config.Pages {
		w.wg.Add(1)
		go w.watchPage(page)
	}
}

func (w *Watcher) watchPage(page *config.Page) {
	defer w.wg.Done()
	every := page.EveryDuration()
	ticker := time.NewTicker(every)
	log := log.With().Str("page", page.Name).Str("every", every.String()).Logger()
	log.Info().Msg("Watching")

	for {
		select {
		case <-w.ctx.Done():
			log.Warn().Msg("Watcher exiting")
			return
		case <-ticker.C:
			log.Info().Msg("Checking")
			w.check(page, log)
		}
	}
}

func (w *Watcher) check(page *config.Page, log zerolog.Logger) {
	normalizedPage, err := w.normalizer.Get(page)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get and normalize page")
	}

	log.Info().Int("length", len(normalizedPage)).Msg("Got and normalized page")
	dbPage, err := w.pageRepository.GetLatestChange(page.Name)
	if err != nil {
		// If we didn't find a record for this page, add it
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msg("Didn't find an existing page to compare against; attempting to add it")

			dbPage := &models.Page{
				Name: page.Name,
				Text: normalizedPage,
			}
			err = w.pageRepository.SaveChange(dbPage)
			if err == nil {
				log.Info().Msg("Successfully added the new page")
			} else {
				log.Error().Err(err).Msg("Failed to add the new page")
			}
			return
		}
	}

	// Determine if there are differences and format it into a human-readable diff if soo
	diff, isDiff := w.diff(normalizedPage, dbPage.Text, page, log)

	// If the page is not different, stop now
	if !isDiff {
		log.Info().Msg("No differences found!")
		return
	}

	// If the page is different, add it to the database
	newDbPage := &models.Page{
		Name: dbPage.Name,
		Text: normalizedPage,
		Diff: diff,
	}
	err = w.pageRepository.SaveChange(newDbPage)
	if err != nil {
		log.Error().Err(err).Msg("A difference was found, but failed to save it to the database!")
		return
	}
	log.Info().Msg("Found differences; saved them to the database")

	// TODO send discord notifications informing about the change
	log.Warn().Msg("Differences were found, but notifications are currently unsupported")
}

func (w *Watcher) diff(newPage, currentPage string, pageConfig *config.Page, log zerolog.Logger) (string, bool) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(newPage, currentPage, false)
	return dmp.DiffPrettyHtml(diffs), true
}
