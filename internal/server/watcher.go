package server

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Watcher polls the sigil directory for file changes and triggers reloads.
type Watcher struct {
	dirs     []string
	interval time.Duration
	onChange func()
	stop     chan struct{}
	wg       sync.WaitGroup
}

// NewWatcher creates a poll-based file watcher.
func NewWatcher(sigilDir string, interval time.Duration, onChange func()) *Watcher {
	return &Watcher{
		dirs: []string{
			filepath.Join(sigilDir, "pages"),
			filepath.Join(sigilDir, "themes"),
			filepath.Join(sigilDir, "datasources"),
		},
		interval: interval,
		onChange: onChange,
		stop:     make(chan struct{}),
	}
}

// Start begins polling for changes in a goroutine.
func (w *Watcher) Start() {
	w.wg.Add(1)
	go w.poll()
}

// Stop signals the watcher to stop and waits for it to finish.
func (w *Watcher) Stop() {
	close(w.stop)
	w.wg.Wait()
}

func (w *Watcher) poll() {
	defer w.wg.Done()

	snapshot := w.takeSnapshot()

	for {
		select {
		case <-w.stop:
			return
		case <-time.After(w.interval):
			current := w.takeSnapshot()
			if !snapshotsEqual(snapshot, current) {
				snapshot = current
				w.onChange()
			}
		}
	}
}

type fileSnapshot struct {
	Path    string
	ModTime time.Time
	Size    int64
}

func (w *Watcher) takeSnapshot() map[string]fileSnapshot {
	snap := map[string]fileSnapshot{}
	for _, dir := range w.dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if filepath.Ext(name) != ".yaml" && filepath.Ext(name) != ".yml" {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			path := filepath.Join(dir, name)
			snap[path] = fileSnapshot{
				Path:    path,
				ModTime: info.ModTime(),
				Size:    info.Size(),
			}
		}
	}
	return snap
}

func snapshotsEqual(a, b map[string]fileSnapshot) bool {
	if len(a) != len(b) {
		return false
	}
	for path, aSnap := range a {
		bSnap, ok := b[path]
		if !ok {
			return false
		}
		if !aSnap.ModTime.Equal(bSnap.ModTime) || aSnap.Size != bSnap.Size {
			return false
		}
	}
	return true
}
