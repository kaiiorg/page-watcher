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
	"github.com/kaiiorg/page-watcher/pkg/web"

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
	web            web.Web
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
	w.web, err = web.NewWebDiffPreviewer(config.Web, w.pageRepository)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Watcher) Close() {
	w.ctxCancel()
	w.wg.Wait()
}

func (w *Watcher) Watch() {
	w.web.Run()

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
	log.Info().Msg("Initial check")
	w.check(page, log)

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

func (w *Watcher) check(pageConfig *config.Page, log zerolog.Logger) {
	rawPage, normalizedPage, err := w.normalizer.Get(pageConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get and normalize pageConfig")
	}

	log.Info().Int("length", len(normalizedPage)).Msg("Got and normalized pageConfig")
	dbPage, err := w.pageRepository.GetLatestChange(pageConfig.Name)
	if err != nil {
		// If we didn't find a record for this pageConfig, add it
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msg("Didn't find an existing pageConfig to compare against; attempting to add it")

			dbPage := &models.Page{
				Name:           pageConfig.Name,
				NormalizedText: normalizedPage,
				RawText:        rawPage,
			}
			err = w.pageRepository.SaveChange(dbPage)
			if err == nil {
				log.Info().Msg("Successfully added the new pageConfig")
			} else {
				log.Error().Err(err).Msg("Failed to add the new pageConfig")
			}
			return
		}
	}

	// Determine if there are differences and format it into a human-readable diff if soo
	diff, isDiff := w.diff(normalizedPage, dbPage.NormalizedText, pageConfig, log)

	// If the pageConfig is not different, stop now
	if !isDiff {
		log.Info().Msg("No differences found!")
		return
	}

	// If the pageConfig is different, add it to the database
	newDbPage := &models.Page{
		Name:           dbPage.Name,
		NormalizedText: normalizedPage,
		RawText:        rawPage,
	}
	err = newDbPage.EncodeDiff(diff)
	if err != nil {
		log.Error().Err(err).Msg("A difference was found, but failed to encode the diffs before trying to save it to the database")
		return
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

func (w *Watcher) diff(newPage, currentPage string, pageConfig *config.Page, log zerolog.Logger) ([]diffmatchpatch.Diff, bool) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(newPage, currentPage, false)
	log.Debug().Str("page", pageConfig.Name).Int("diffCount", len(diffs)).Send()
	if len(diffs) == 1 {
		return nil, false
	}
	return diffs, true
}
