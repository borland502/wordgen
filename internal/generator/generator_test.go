package generator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type testWordEntry struct {
	Word    string   `json:"word"`
	Sources []string `json:"sources"`
}

func TestSelectWordsWithContextAppliesFilters(t *testing.T) {
	datasetPath := writeTestDataset(t, []testWordEntry{
		{Word: "steel", Sources: []string{"fsu/wordle.txt"}},
		{Word: "stale", Sources: []string{"fsu/wordle.txt"}},
		{Word: "stone", Sources: []string{"dwyl/words.txt"}},
		{Word: "apple", Sources: []string{"fsu/wordle.txt"}},
	})

	words, matched, err := SelectWordsWithContext(context.Background(), Request{
		Dataset:   datasetPath,
		Count:     10,
		MinLength: 5,
		MaxLength: 5,
		Prefix:    "st",
		Contains:  "e",
		Sources:   []string{"FSU/WORDLE.TXT"},
		Seed:      42,
	})
	if err != nil {
		t.Fatalf("select words: %v", err)
	}

	if matched != 2 {
		t.Fatalf("expected 2 matched words, got %d", matched)
	}
	if len(words) != 2 {
		t.Fatalf("expected 2 selected words, got %d", len(words))
	}

	allowed := map[string]struct{}{
		"steel": {},
		"stale": {},
	}
	for _, word := range words {
		if _, ok := allowed[word]; !ok {
			t.Fatalf("unexpected selected word %q", word)
		}
	}
}

func TestStreamMatchedWordsWithContextStopsOnYieldFalse(t *testing.T) {
	datasetPath := writeTestDataset(t, []testWordEntry{
		{Word: "alpha", Sources: []string{"fsu/wordle.txt"}},
		{Word: "bravo", Sources: []string{"fsu/wordle.txt"}},
		{Word: "charlie", Sources: []string{"fsu/wordle.txt"}},
	})

	seen := 0
	matched, err := StreamMatchedWordsWithContext(context.Background(), Request{
		Dataset:   datasetPath,
		Count:     1,
		MinLength: 1,
		MaxLength: 0,
	}, func(word string) bool {
		seen++
		return seen < 2
	})
	if err != nil {
		t.Fatalf("stream words: %v", err)
	}

	if matched != 2 {
		t.Fatalf("expected stream to stop after 2 matches, got %d", matched)
	}
}

func TestSelectWordsWithContextCancelled(t *testing.T) {
	datasetPath := writeTestDataset(t, []testWordEntry{
		{Word: "alpha", Sources: []string{"fsu/wordle.txt"}},
		{Word: "bravo", Sources: []string{"fsu/wordle.txt"}},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := SelectWordsWithContext(ctx, Request{
		Dataset:   datasetPath,
		Count:     1,
		MinLength: 1,
	})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestIndexedDatasetSelectWordsWithContext(t *testing.T) {
	datasetPath := writeTestDataset(t, []testWordEntry{
		{Word: "steel", Sources: []string{"fsu/wordle.txt"}},
		{Word: "stale", Sources: []string{"fsu/wordle.txt"}},
		{Word: "stone", Sources: []string{"dwyl/words.txt"}},
	})

	indexed, err := LoadIndexedDataset(datasetPath)
	if err != nil {
		t.Fatalf("load indexed dataset: %v", err)
	}

	words, matched, err := indexed.SelectWordsWithContext(context.Background(), Request{
		Count:     5,
		MinLength: 5,
		Prefix:    "st",
		Sources:   []string{"fsu/wordle.txt"},
		Seed:      7,
	})
	if err != nil {
		t.Fatalf("select from indexed dataset: %v", err)
	}

	if matched != 2 {
		t.Fatalf("expected 2 matched words from indexed dataset, got %d", matched)
	}
	if len(words) != 2 {
		t.Fatalf("expected 2 selected words from indexed dataset, got %d", len(words))
	}
}

func BenchmarkSelectWordsWithContext(b *testing.B) {
	datasetPath := writeBenchmarkDataset(b)
	request := Request{
		Dataset:   datasetPath,
		Count:     10,
		MinLength: 5,
		MaxLength: 10,
		Prefix:    "wo",
		Contains:  "rd",
		Sources:   []string{"fsu/wordle.txt", "dwyl/words.txt"},
		Seed:      99,
	}

	b.ResetTimer()
	for range b.N {
		_, _, err := SelectWordsWithContext(context.Background(), request)
		if err != nil {
			b.Fatalf("select words: %v", err)
		}
	}
}

func writeTestDataset(tb testing.TB, entries []testWordEntry) string {
	tb.Helper()

	bytes, err := json.Marshal(entries)
	if err != nil {
		tb.Fatalf("marshal dataset: %v", err)
	}

	path := filepath.Join(tb.TempDir(), "all.json")
	if err := os.WriteFile(path, append(bytes, '\n'), 0o600); err != nil {
		tb.Fatalf("write dataset: %v", err)
	}

	return path
}

func writeBenchmarkDataset(tb testing.TB) string {
	tb.Helper()

	entries := make([]testWordEntry, 0, 10000)
	for index := range 10000 {
		sources := []string{"dwyl/words.txt"}
		if index%2 == 0 {
			sources = []string{"fsu/wordle.txt"}
		}
		entries = append(entries, testWordEntry{
			Word:    "word" + string(rune('a'+(index%26))) + "x",
			Sources: sources,
		})
	}

	return writeTestDataset(tb, entries)
}
