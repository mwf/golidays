package backuper

import (
	"fmt"
	"github.com/mwf/golidays/model"
	"io/ioutil"
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
	minPeriod = time.Minute
)

type Backuper struct {
	storage store.Store
	period  time.Duration
	logger  logger.Logger

	basePath string
	files    *stringStack

	runOnce sync.Once
	done    chan struct{}
}

// New returns new backuper instance
func New(storage store.Store, period time.Duration, basePath string, maxBackups int, log logger.Logger) (*Backuper, error) {
	if period < minPeriod {
		return nil, fmt.Errorf("period is too low: %s < %s", period, minPeriod)
	}
	if maxBackups < 0 {
		return nil, fmt.Errorf("too few backups: %d", maxBackups)
	}

	info, err := os.Stat(basePath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", basePath)
	}

	b := &Backuper{
		storage:  storage,
		period:   period,
		logger:   log,
		basePath: basePath,
		files:    newStringStack(maxBackups),
		done:     make(chan struct{}),
	}

	b.restoreList()
	return b, nil
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

// RestoreStorage restores storage from the last backup
func (b *Backuper) RestoreStorage() error {
	lastBackupPath := b.files.Head()
	if lastBackupPath == "" {
		return fmt.Errorf("no backup files")
	}
	b.logger.Debugf("restoring last backup from '%s'", lastBackupPath)

	bytes, err := ioutil.ReadFile(lastBackupPath)
	if err != nil {
		return fmt.Errorf("error reading backup '%s': %s", lastBackupPath, err)
	}

	if err := b.restoreData(bytes); err != nil {
		return err
	}

	b.logger.Infof("backup '%s' restored OK", lastBackupPath)
	return nil
}

func (b *Backuper) restoreData(bytes []byte) error {
	holidays := make(model.Holidays, 0)
	if err := yaml.Unmarshal(bytes, &holidays); err != nil {
		return fmt.Errorf("error unmarshaling data: %s", err)
	}

	if err := b.storage.Set(holidays); err != nil {
		return fmt.Errorf("error restoring storage: %s", err)
	}
	return nil
}

func (b *Backuper) restoreList() {
	// try to restore backup list
	pattern := filepath.Join(b.basePath, "holidays.*\\.yml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		b.logger.Warningf("backups list restore failed: %s", err)
	}

	b.logger.Debugf("found backups: %s", matches)
	sort.Strings(matches)

	for _, fpath := range matches {
		b.preserveFile(fpath)
	}
}

func (b *Backuper) loop() {
	b.logger.Infof("%s started", b)
	defer b.logger.Infof("%s stopped", b)

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
