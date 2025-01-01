package journal

import (
	"bufio"
	"os"
	"testing"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name      string
		inputfile string
		wantErr   bool
	}{
		{"Test 1", "../../testdata/set1/Journal.2024-12-29T230453.01.log", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			journalFile, err := os.Open(tt.inputfile)
			if err != nil {
				t.Fatal(err)
			}
			defer journalFile.Close()

			scanner := bufio.NewScanner(journalFile)
			for scanner.Scan() {
				line := scanner.Text()
				t.Log(line)
				got, err := Parse([]byte(line))
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Log("got", got)
			}
			// t.Fail()
		})
	}

}
