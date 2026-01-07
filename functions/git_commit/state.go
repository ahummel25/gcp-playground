package git_commit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"time"

	"cloud.google.com/go/storage"

	"github.com/hummelgcp/go/shared/logging"
)

const stateObjectName = "state.json"

// DayState holds per-day counters and whether the day is a rest day.
type DayState struct {
	Runs    int  `json:"runs"`
	Commits int  `json:"commits"`
	RestDay bool `json:"rest_day"`
}

// State holds daily state keyed by YYYY-MM-DD. Only the last 8 days are kept.
type State struct {
	Days map[string]DayState `json:"days"`
}

// LoadState reads the state object from GCS. Missing or invalid object returns empty state.
func LoadState(ctx context.Context, bucket string) (*State, error) {
	if bucket == "" {
		return &State{Days: make(map[string]DayState)}, nil
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	o := client.Bucket(bucket).Object(stateObjectName)
	r, err := o.NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return &State{Days: make(map[string]DayState)}, nil
		}
		return nil, fmt.Errorf("NewReader: %w", err)
	}
	defer r.Close()

	var s State
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return &State{Days: make(map[string]DayState)}, nil
	}
	if s.Days == nil {
		s.Days = make(map[string]DayState)
	}
	return &s, nil
}

// SaveState writes the state to GCS and prunes to the last 8 days.
func SaveState(ctx context.Context, bucket string, s *State) error {
	if bucket == "" {
		return nil
	}
	pruneToLastNDays(s, 8)

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	o := client.Bucket(bucket).Object(stateObjectName)
	w := o.NewWriter(ctx)
	w.ContentType = "application/json"
	if err := json.NewEncoder(w).Encode(s); err != nil {
		_ = w.Close()
		return fmt.Errorf("Encode: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	return nil
}

func pruneToLastNDays(s *State, n int) {
	if len(s.Days) <= n {
		return
	}
	var keys []string
	for k := range s.Days {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	keep := make(map[string]struct{})
	for i := 0; i < n && i < len(keys); i++ {
		keep[keys[i]] = struct{}{}
	}
	for k := range s.Days {
		if _, ok := keep[k]; !ok {
			delete(s.Days, k)
		}
	}
}

// restDayProbability is the baseline chance a new day is designated a rest day (~2/7).
const restDayProbability = 2.0 / 7.0

// weekendRestBoost adds medium probability for Sat/Sun.
const weekendRestBoost = 0.63

// EnsureToday ensures today exists in state. If it was missing, it is added and
// rest_day is set randomly (~2/7) with a weekend bias. If the last 6 days already
// have 2+ rest days, do not make today rest.
func EnsureToday(s *State, today string) {
	if _, ok := s.Days[today]; ok {
		return
	}
	t := mustParseDate(today)
	restProb := restDayProbability
	if isWeekend(t) {
		restProb += weekendRestBoost
	}
	rest := rand.Float64() < restProb
	// If we already had 2+ rest days in the last 6, avoid another. Weekend days ignore this constraint.
	if rest && !isWeekend(t) {
		count := 0
		for i := 1; i <= 6; i++ {
			d := t.AddDate(0, 0, -i)
			key := d.Format(time.DateOnly)
			if s.Days[key].RestDay {
				count++
			}
		}
		if count >= 2 {
			rest = false
		}
	}
	s.Days[today] = DayState{RestDay: rest}
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func mustParseDate(s string) time.Time {
	t, _ := time.Parse(time.DateOnly, s)
	return t
}

// RecordRun updates state for today: increments runs and, if didCommit, commits.
// Call before SaveState.
func RecordRun(s *State, today string, didCommit bool) {
	d := s.Days[today]
	d.Runs++
	if didCommit {
		d.Commits++
	}
	s.Days[today] = d
}

const (
	dailyCommitCap = 8
	maxBatchSize   = 5
)

// ShouldCommit decides whether to commit using only state. No external
// probability knob: rest days, daily cap, and heavy-week cutoff are hard
// rules; the remainder uses a state-derived probability with streak bias.
//
// Hard rules (state as gate):
//   - rest day: never commit
//   - commits >= 8 today: never commit
//   - total commits in last 7 days > 18: never commit (heavy week)
//
// Otherwise: probability derived only from state (commits today, runs today,
// recent streak, total week). Roll for variance; no COMMIT_PROBABILITY.
func ShouldCommit(s *State, today string) bool {
	d := s.Days[today]
	if d.RestDay {
		return false
	}
	if d.Commits >= 8 {
		return false
	}

	t := mustParseDate(today)
	totalWeek := totalWeekCommits(s, t)
	if totalWeek > 36 {
		return false
	}

	// Prior-day streak data (last 3 days)
	y1, y2, y3 := recentCommitCounts(s, t)

	// State-derived probability (no external base)
	p := stateDerivedProbability(d.Runs, d.Commits, totalWeek, y1, y2, y3)
	roll := rand.Float64() * 100
	logging.Info("Commit roll for %s: %.2f (threshold %.1f)", today, roll, p)

	return roll < p
}

// CommitDecisionProbability returns the probability used for commit decisions,
// and a hard-rule reason when the decision is forced.
func CommitDecisionProbability(s *State, today string) (p float64, hardRule string) {
	d := s.Days[today]
	if d.RestDay {
		return 0, "rest_day"
	}
	if d.Commits >= 8 {
		return 0, "daily_cap"
	}

	t := mustParseDate(today)
	totalWeek := totalWeekCommits(s, t)
	if totalWeek > 36 {
		return 0, "heavy_week"
	}

	y1, y2, y3 := recentCommitCounts(s, t)
	return stateDerivedProbability(d.Runs, d.Commits, totalWeek, y1, y2, y3), ""
}

func totalWeekCommits(s *State, t time.Time) int {
	var totalWeek int
	for i := 0; i < 7; i++ {
		key := t.AddDate(0, 0, -i).Format(time.DateOnly)
		totalWeek += s.Days[key].Commits
	}
	return totalWeek
}

func recentCommitCounts(s *State, t time.Time) (y1, y2, y3 int) {
	y1 = s.Days[t.AddDate(0, 0, -1).Format(time.DateOnly)].Commits
	y2 = s.Days[t.AddDate(0, 0, -2).Format(time.DateOnly)].Commits
	y3 = s.Days[t.AddDate(0, 0, -3).Format(time.DateOnly)].Commits
	return y1, y2, y3
}

// CommitBatchSize decides how many commits to create in this run (1–4) when
// a commit is already approved. Heavy days get 2–4; otherwise 1. The result
// is clamped to remaining capacity for today.
func CommitBatchSize(s *State, today string) int {
	d := s.Days[today]
	remaining := dailyCommitCap - d.Commits
	if remaining <= 0 {
		return 0
	}
	t := mustParseDate(today)
	totalWeek := totalWeekCommits(s, t)
	y1, y2, y3 := recentCommitCounts(s, t)

	batch := 1
	if isHeavyDay(totalWeek, y1, y2, y3) && remaining >= 2 {
		batch = heavyDayBatchSize()
	}
	if batch > remaining {
		batch = remaining
	}
	if batch > maxBatchSize {
		batch = maxBatchSize
	}
	return batch
}

func isHeavyDay(totalWeek, y1, y2, y3 int) bool {
	if totalWeek >= 8 {
		return true
	}
	if y1 > 0 || y2 > 0 {
		return true
	}
	if y1+y2+y3 >= 3 {
		return true
	}
	return false
}

func heavyDayBatchSize() int {
	roll := rand.IntN(100)
	switch {
	case roll < 20:
		return 2
	case roll < 55:
		return 3
	case roll < 85:
		return 4
	default:
		return 5
	}
}

// stateDerivedProbability returns 0–100 from today's runs/commits, recent
// streak, and last-7-days total. Used only when no hard rule applied.
func stateDerivedProbability(runs, commits, totalWeek, y1, y2, y3 int) float64 {
	var p float64
	switch commits {
	case 0:
		p = 55
	case 1:
		p = 40
	case 2:
		p = 30
	case 3:
		p = 20
	default:
		p = 5
	}

	// Late-day boost (no hard catch-up)
	if runs >= 3 && commits == 0 {
		p += 30
	} else if runs >= 2 && commits == 0 {
		p += 15
	}

	// Streak bias: cluster commits, allow gaps
	if y1 > 0 {
		p += 10
	}
	if y1 == 0 && y2 == 0 {
		p -= 5
	}
	if y1+y2+y3 >= 6 {
		p -= 5
	}

	if totalWeek < 6 {
		p += 15
	}
	if p < 5 {
		p = 5
	}
	if p > 75 {
		p = 75
	}
	return p
}

// DecideAndRecord runs the full random-mode flow: load state, ensure today,
// decide commit or skip (state-only, no COMMIT_PROBABILITY), record the run
// (not the commit), and save. It returns true to commit and false to skip.
// On error it returns (false, err). If it returns true, the caller must
// perform the GitHub commit and then call RecordCommits.
func DecideAndRecord(ctx context.Context, bucket, today string) (doCommit bool, err error) {
	s, err := LoadState(ctx, bucket)
	if err != nil {
		return false, fmt.Errorf("LoadState: %w", err)
	}
	EnsureToday(s, today)

	doCommit = ShouldCommit(s, today)
	RecordRun(s, today, false) // only record run; commit recorded after success

	if err := SaveState(ctx, bucket, s); err != nil {
		return false, fmt.Errorf("SaveState: %w", err)
	}
	return doCommit, nil
}

// RecordCommits increments the commit count for today by count after
// successful commits. Call this only for commits that actually succeeded.
func RecordCommits(ctx context.Context, bucket, today string, count int) error {
	if count <= 0 {
		return nil
	}
	if bucket == "" {
		return nil
	}
	s, err := LoadState(ctx, bucket)
	if err != nil {
		return fmt.Errorf("LoadState: %w", err)
	}
	d := s.Days[today]
	d.Commits += count
	s.Days[today] = d
	if err := SaveState(ctx, bucket, s); err != nil {
		return fmt.Errorf("SaveState: %w", err)
	}
	return nil
}
