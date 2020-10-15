package timeline

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type TimeEntry struct {
	ID    string
	Start time.Time
	End   *time.Time
	Notes []string
}

type EntryCollection struct {
	Owner   string       `yaml:"owner"`
	Entries []*TimeEntry `yaml:"entries"`
}

type Timeline struct {
	fPath   string
	mu      sync.RWMutex
	entries *EntryCollection
}

func New(fPath string, username string) (*Timeline, error) {
	b, err := ioutil.ReadFile(fPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Timeline{
				fPath: fPath,
				entries: &EntryCollection{
					Owner:   username,
					Entries: []*TimeEntry{},
				},
			}, nil
		}
		return nil, fmt.Errorf("failed to read file at path %q: %v", fPath, err)
	}

	ents := &EntryCollection{}
	if err = yaml.Unmarshal(b, ents); err != nil {
		return nil, fmt.Errorf("failed to read entry data: %v", err)
	}

	sort.SliceStable(ents.Entries, func(i, j int) bool {
		return ents.Entries[i].Start.Before(ents.Entries[j].Start)
	})

	return &Timeline{
		fPath:   fPath,
		entries: ents,
	}, nil
}

var ErrNotExist = errors.New("no running time events exist")

func (l *Timeline) RunningEvent() (*TimeEntry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, ent := range l.entries.Entries {
		if ent.End == nil {
			return ent, nil
		}
	}

	return nil, ErrNotExist
}

func (l *Timeline) Start(when time.Time, notes []string) *TimeEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	if notes == nil {
		notes = []string{}
	}

	id := uuid.New().String()[:8]

	ent := &TimeEntry{
		ID:    id,
		Start: when,
		Notes: notes,
	}

	l.entries.Entries = append(l.entries.Entries, ent)

	return ent
}

func (l *Timeline) AddNote(note string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	cur, err := l.RunningEvent()
	if err != nil {
		return err
	}

	cur.Notes = append(cur.Notes, note)

	return nil
}

func (l *Timeline) ReportFormat() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("Time for %s\n", l.entries.Owner))

	grps := map[string][]*TimeEntry{}
	for _, ent := range l.entries.Entries {
		header := ent.Start.Format("January 2, 2006")

		if grp, ok := grps[header]; ok {
			grps[header] = append(grp, ent)
			continue
		}

		grps[header] = []*TimeEntry{ent}
	}

	var groups = []struct {
		header    string
		totalTime time.Duration
		entries   []*TimeEntry
	}{}

	now := time.Now()
	for header, ents := range grps {
		var dur time.Duration
		for _, ent := range ents {
			end := now
			if ent.End != nil {
				end = *ent.End
			}

			d := end.Sub(ent.Start)
			dur = dur + d
		}

		groups = append(groups, struct {
			header    string
			totalTime time.Duration
			entries   []*TimeEntry
		}{
			header:    header,
			totalTime: dur,
			entries:   ents,
		})
	}

	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].header > groups[j].header
	})

	for _, group := range groups {
		b.WriteString("\n")
		h := fmt.Sprintf("%s - Total: %v\n", group.header, group.totalTime.Round(time.Second))
		hr := strings.Repeat("-", len(h))
		b.WriteString(h)
		b.WriteString(hr)
		b.WriteString("\n")

		for _, ent := range group.entries {
			ip := " (in progress)"
			end := now
			if ent.End != nil {
				end = *ent.End
				ip = ""
			}
			d := end.Sub(ent.Start)
			b.WriteString(fmt.Sprintf("%v - %v%s\n", ent.Start.Format("3:04:05 PM"), d.Round(time.Second), ip))
		}

		b.WriteString(hr)
		b.WriteString("\n")
	}

	return b.String()
}

func (l *Timeline) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, err := yaml.Marshal(l.entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %v", err)
	}

	os.MkdirAll(path.Dir(l.fPath), os.ModePerm)

	f, err := os.Create(l.fPath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}

	defer func() {
		f.Sync()
		f.Close()
	}()

	_, err = bytes.NewBuffer(b).WriteTo(f)
	if err != nil {
		return fmt.Errorf("failed to write entries to file: %v", err)
	}

	return nil
}
