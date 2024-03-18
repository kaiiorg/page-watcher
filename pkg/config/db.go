package config

var (
	DEFAULT_DB_PATH   = "./page-watcher.db"
	DEFAULT_DB_RETAIN = 3
)

type DB struct {
	// Path is the path and name of the sqlite db file
	Path string `hcl:"path"`
	// Retain is how many extra versions of the page to retain
	// At least one version will always be in the DB after the first store; otherwise we'd have no way of knowing it changed!
	Retain int `hcl:"retain"`
}

func defaultDb() *DB {
	return &DB{
		Path:   DEFAULT_DB_PATH,
		Retain: DEFAULT_DB_RETAIN,
	}
}
