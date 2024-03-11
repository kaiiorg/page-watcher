package watcher

import (
	"context"
	"github.com/kaiiorg/page-watcher/pkg/watcher/normalizer"
	"github.com/rs/zerolog"
	"sync"
	"time"

	"github.com/kaiiorg/page-watcher/pkg/config"

	"github.com/rs/zerolog/log"
)

type Watcher struct {
	config    *config.Config
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	normalizer normalizer.Normalizer
}

func New(config *config.Config) *Watcher {
	w := &Watcher{
		config:     config,
		normalizer: normalizer.New(),
	}
	w.ctx, w.ctxCancel = context.WithCancel(context.Background())

	return w
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

	select {
	case <-w.ctx.Done():
		log.Warn().Msg("Watcher exiting")
		return
	case <-ticker.C:
		log.Info().Msg("Checking")
		w.check(page, log)
	}
}

func (w *Watcher) check(page *config.Page, log zerolog.Logger) {
	normalizedPage, err := w.normalizer.Get(page)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get and normalize page")
	}

	log.Info().Int("length", len(normalizedPage)).Msg("Got and normalized page")
}
