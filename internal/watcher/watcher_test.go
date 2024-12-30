package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func touch(filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create("temp.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(filename, currentTime, currentTime)
		if err != nil {
			panic(err)
		}
	}
}

func Test_WatchForAWhile(t *testing.T) {

	tests := []struct {
		name    string
		basedir string
		touch   []string
	}{
		{"Test 1", "../../testdata/set1/", []string{"Status.json", "Journal.2024-12-29T230453.01.log"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basedir := tt.basedir
			w, err := New(basedir)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()
			count := 0
			matchCount := 0
			go w.Watch()

			go func() {
				for event := range w.Events {
					t.Log("something happened", "event", event)
					count++
					for _, f := range tt.touch {
						if f == event.Name {
							matchCount++
						}
					}
				}
				t.Log("watcher event channel closed")
			}()
			time.Sleep(time.Second)
			// touch some files
			for _, f := range tt.touch {
				touch(filepath.FromSlash(basedir + "/" + f))
				time.Sleep(100 * time.Millisecond)
			}
			// wait for events
			time.Sleep(2 * time.Second)
			assert.Equal(t, len(tt.touch), count, "number of events")
			assert.Equal(t, len(tt.touch), matchCount, "number of matched events")
			// t.Fail()
		})
	}
}
