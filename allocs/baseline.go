//go:build !solution && !change

package allocs

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"slices"
)

type Counter interface {
	Count(r io.Reader) error
	String() string
}

type BaselineCounter struct {
	counts map[string]int
}

func NewBaselineCounter() Counter {
	return BaselineCounter{counts: map[string]int{}}
}

func (c BaselineCounter) Count(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	dataStr := string(data)
	for line := range strings.SplitSeq(dataStr, "\n") {
		for word := range strings.SplitSeq(line, " ") {
			c.counts[word]++
		}
	}
	return nil
}

func (c BaselineCounter) String() string {
	keys := make([]string, 0, 0)
	for word := range c.counts {
		keys = append(keys, word)
	}
	slices.Sort(keys)

	result := ""
	for _, key := range keys {
		line := fmt.Sprintf("word '%s' has %d occurrences\n", key, c.counts[key])
		result += line
	}
	return result
}
