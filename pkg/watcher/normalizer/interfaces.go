package normalizer

import "github.com/kaiiorg/page-watcher/pkg/config"

type Normalizer interface {
	Get(page *config.Page) (string, error)
}
