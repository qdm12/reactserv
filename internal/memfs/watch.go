package memfs

import (
	"context"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/qdm12/golibs/logging"
)

func (m *memFS) Watch(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		m.logger.Error(err)
		return
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			m.logger.Error(err)
		}
	}()

	directoriesWatched := make(map[string]struct{})

	m.mu.RLock()
	directories := copyDirectories(m.directories)
	m.mu.RUnlock()
	adaptWatcher(watcher, m.logger, directories, directoriesWatched)

	for {
		select {
		case <-ctx.Done():
			m.logger.Warn("exiting file system watch (%s)", ctx.Err())
			return

		case _, ok := <-watcher.Events:
			if !ok {
				m.logger.Warn("exiting file system watch (watcher events channel closed)")
				return
			}
			// TODO load file affected only
			if err := m.loadAll(); err != nil {
				m.logger.Error(err)
			}
			m.mu.RLock()
			directories := copyDirectories(m.directories)
			m.mu.RUnlock()
			adaptWatcher(watcher, m.logger, directories, directoriesWatched)

		case err, ok := <-watcher.Errors:
			if !ok {
				m.logger.Warn("exiting file system watch (watcher errors channel closed)")
				return
			}
			m.logger.Error("watcher: %s", err)
		}
	}
}

func copyDirectories(directories map[string]struct{}) (directoriesCopy map[string]struct{}) {
	directoriesCopy = make(map[string]struct{}, len(directories))
	for directory := range directories {
		directoriesCopy[directory] = struct{}{}
	}
	return directoriesCopy
}

func adaptWatcher(watcher *fsnotify.Watcher, logger logging.Logger,
	directories, directoriesWatched map[string]struct{}) {
	directoriesToAdd := make(map[string]struct{})
	for directory := range directories {
		if _, ok := directoriesWatched[directory]; !ok {
			directoriesToAdd[directory] = struct{}{}
		}
	}
	directoriesToRemove := make(map[string]struct{})
	for directoryWatched := range directoriesWatched {
		if _, ok := directories[directoryWatched]; !ok {
			directoriesToRemove[directoryWatched] = struct{}{}
		}
	}
	for directoryToRemove := range directoriesToRemove {
		delete(directoriesWatched, directoryToRemove)
		// No need to watcher.Remove() as it gets deleted
		// when the directory is deleted
	}
	for directoryToAdd := range directoriesToAdd {
		err := watcher.Add(directoryToAdd)
		if err != nil {
			logger.Error(err)
			continue
		}
		directoriesWatched[directoryToAdd] = struct{}{}
		logger.Info("Watching directory %s", directoryToAdd)
	}
}
