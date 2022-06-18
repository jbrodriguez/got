package recorder

import (
	"fmt"
	"path/filepath"
	"time"

	"got/cli/lib"
)

type Recorder struct {
	DataDir string
}

func CreateRecorder(dataDir string) *Recorder {
	return &Recorder{
		DataDir: dataDir,
	}
}

func (r *Recorder) AddActivity(activity string) error {
	return r.add(activity, "")
}

func (r *Recorder) AddOn() error {
	return r.add("on", "\n")
}

func (r *Recorder) add(activity, prepend string) error {
	now := time.Now()

	utc := now.UTC().Format(time.RFC3339)

	line := fmt.Sprintf("%s\n%s %s", prepend, utc, activity)

	f := filepath.Join(r.DataDir, fmt.Sprintf("%d.log", now.Year()))

	err := lib.AppendToFile(f, line)
	if err != nil {
		return err
	}

	return nil
}
