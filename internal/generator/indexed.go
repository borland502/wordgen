package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// IndexedDataset stores the parsed word index in memory for repeated low-latency calls.
type IndexedDataset struct {
	entries []wordEntry
}

// LoadIndexedDataset loads and decodes the full JSON word index into memory.
func LoadIndexedDataset(path string) (*IndexedDataset, error) {
	reader, err := openWordIndexReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var entries []wordEntry
	if err := decoder.Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode indexed word dataset: %w", err)
	}

	return &IndexedDataset{entries: entries}, nil
}

// SelectWordsWithContext selects up to Count sampled words from the in-memory dataset.
func (dataset *IndexedDataset) SelectWordsWithContext(ctx context.Context, request Request) ([]string, int, error) {
	if dataset == nil {
		return nil, 0, fmt.Errorf("indexed dataset is nil")
	}
	if err := validateRequest(request); err != nil {
		return nil, 0, err
	}

	matcher := compileMatcher(request)
	seed := request.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	rng := rand.New(rand.NewSource(seed))
	selectedWords := make([]string, 0, request.Count)
	matchedCount := 0

	for _, entry := range dataset.entries {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}
		if !matcher.matches(entry) {
			continue
		}

		matchedCount++
		if len(selectedWords) < request.Count {
			selectedWords = append(selectedWords, entry.Word)
			continue
		}

		candidateIndex := rng.Intn(matchedCount)
		if candidateIndex < request.Count {
			selectedWords[candidateIndex] = entry.Word
		}
	}

	if len(selectedWords) > 1 {
		rng.Shuffle(len(selectedWords), func(left, right int) {
			selectedWords[left], selectedWords[right] = selectedWords[right], selectedWords[left]
		})
	}

	return selectedWords, matchedCount, nil
}

// StreamMatchedWordsWithContext streams matching words from the in-memory dataset.
func (dataset *IndexedDataset) StreamMatchedWordsWithContext(ctx context.Context, request Request, yield func(string) bool) (int, error) {
	if dataset == nil {
		return 0, fmt.Errorf("indexed dataset is nil")
	}
	if err := validateRequest(request); err != nil {
		return 0, err
	}

	matcher := compileMatcher(request)
	matchedCount := 0
	for _, entry := range dataset.entries {
		if err := ctx.Err(); err != nil {
			return matchedCount, err
		}
		if !matcher.matches(entry) {
			continue
		}

		matchedCount++
		if !yield(entry.Word) {
			return matchedCount, nil
		}
	}

	return matchedCount, nil
}
