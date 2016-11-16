package service

import (
	"container/list"
	"fmt"
	"github.com/mwf/golidays/model"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/mwf/golidays/service/logger"
	"github.com/mwf/golidays/service/store"
	"gopkg.in/yaml.v2"
)

const (
	// no reason to keep it higher
	minBackupPeriod     = time.Minute
	defaultBackupPeriod = 24 * time.Hour
)

type stringStack struct {
	list   *list.List
	maxLen int
}

func newStringStack(maxLen int) *stringStack {
	return &stringStack{
		list:   list.New(),
		maxLen: maxLen,
	}
}

// Put adds string to stack. Returns popped item if maxLen is exceeded
func (s *stringStack) Put(str string) (string, bool) {
	s.list.PushBack(str)
	if s.list.Len() > s.maxLen {
		popped := s.list.Remove(s.list.Front())
		return popped.(string), true
	}

	return "", false
}

// Head returns last-added string
func (s *stringStack) Head() string {
	last := s.list.Back()
	if last == nil {
		return ""
	}
	return last.Value.(string)
}

type Backuper struct {
	storage store.Store
	period  time.Duration
	logger  logger.Logger

	basePath string
	files    *stringStack

	runOnce sync.Once
	done    chan struct{}
}

// NewBackuper returns new backuper instance
func NewBackuper(storage store.Store, period time.Duration, basePath string, maxBackups int, log logger.Logger) (*Backuper, error) {
	if period < minBackupPeriod {
		return nil, fmt.Errorf("period is too low: %s < %s", period, minBackupPeriod)
	}

	info, err := os.Stat(basePath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", basePath)
	}

	return &Backuper{
		storage:  storage,
		period:   period,
		logger:   log,
		basePath: basePath,
		files:    newStringStack(maxBackups),
		done:     make(chan struct{}),
	}, nil
}

func (b *Backuper) String() string {
	return fmt.Sprintf("Backuper {period: %s}", b.period)
}

// Run runs asynchronous update loop. Multiple calls do nothing - the loop started
// exactly once.
func (b *Backuper) Run() {
	b.runOnce.Do(func() {
		go b.loop()
	})
}

// Stop stops update loop.
func (b *Backuper) Stop() {
	b.runOnce.Do(func() {
		close(b.done)
	})
}

func (b *Backuper) loop() {
	b.logger.Debugf("started %s", b)

	for {
		select {
		case <-time.After(b.period):
			if err := b.perform(); err != nil {
				b.logger.Error(err.Error())
			}
		case <-b.done:
			return
		}
	}
}

func (b *Backuper) perform() error {
	startedAt := time.Now()

	b.logger.Debugf("perform %s", b)
	defer func() {
		b.logger.Infof("perform finished in %s", time.Now().Sub(startedAt))
	}()

	filename := b.generateBackupName()
	fpath := filepath.Join(b.basePath, filename)
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("os.Create error: %s", err)
	}
	defer f.Close()

	if err := b.collectAndWrite(f); err != nil {
		// remove file on error
		os.Remove(fpath)
		return err
	}

	b.preserveFile(fpath)
	return nil
}

func (b *Backuper) generateBackupName() string {
	dt := time.Now().Format("2006-01-02T03:04")
	return fmt.Sprintf("holidays.%s.yml", dt)
}

func (b *Backuper) collectAndWrite(f *os.File) error {
	holidays := b.storage.Dump()
	sort.Sort(model.HolidaysByDate(holidays))

	bytes, err := yaml.Marshal(holidays)
	if err != nil {
		return fmt.Errorf("error marshaling data: %s", err)
	}
	if _, err := f.Write(bytes); err != nil {
		return fmt.Errorf("error writing data: %s", err)
	}
	return nil
}

func (b *Backuper) preserveFile(fpath string) {
	// put new file in list and remove purged, if any
	name, purged := b.files.Put(fpath)
	if purged {
		if err := os.Remove(name); err != nil {
			b.logger.Warningf("failed to purge '%s': %s", name, err)
		}
	}
}
