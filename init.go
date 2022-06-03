package lightinggo

import (
	"path/filepath"
	"strings"
	"sync"
)

// PrepareOption Used to set up the lightinggo plug-in.
type PrepareOption func(wg *sync.WaitGroup)

// Registe Preparing for lightinggo.
func Registe(options ...PrepareOption) {
	var wg sync.WaitGroup

	for _, f := range options {
		wg.Add(1)

		go f(&wg)
	}

	wg.Wait()
}

// PrepareLogger register logger.
func PrepareLogger(cfg *LoggerConfig) PrepareOption {

	return func(wg *sync.WaitGroup) {
		defer wg.Done()

		if v := strings.TrimSpace(cfg.Filename); len(v) != 0 {
			cfg.Filename = filepath.Clean(v)
		}

		initLogger(cfg)
	}
}

// // PrepareLogger register customerror.
// func PrepareError(code ErrCode) PrepareOption {

// 	return func(wg *sync.WaitGroup) {
// 		defer wg.Done()

// 		_ = InitError(code)
// 	}
// }
