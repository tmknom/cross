package dir

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/tmknom/cross/internal/errlib"
)

type BaseDir struct {
	raw  string
	abs  string
	rel  string
	work string
}

func NewBaseDir(raw string) *BaseDir {
	return &BaseDir{
		raw:  raw,
		abs:  "",
		rel:  ".",
		work: "",
	}
}

func (d *BaseDir) Abs() string {
	if d.abs == "" {
		d.abs = d.generateAbs()
	}
	return d.abs
}

func (d *BaseDir) generateAbs() string {
	if filepath.IsAbs(d.raw) {
		return d.raw
	}
	return filepath.Clean(filepath.Join(d.Work(), d.raw))
}

func (d *BaseDir) Rel() string {
	return d.rel
}

func (d *BaseDir) RelByWork() string {
	result, err := filepath.Rel(d.Work(), d.Abs())
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "cannot resolve rel, work: %s, abs: %s", d.Work(), d.Abs()))
	}
	return filepath.Clean(result)
}

func (d *BaseDir) Work() string {
	if d.work == "" {
		d.work = d.generateWork()
	}
	return d.work
}

func (d *BaseDir) generateWork() string {
	result, err := os.Getwd()
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "cannot resolve work dir"))
	}
	return filepath.Clean(result)
}

func (d *BaseDir) String() string {
	return d.Abs()
}

func (d *BaseDir) GoString() string {
	return d.String()
}

func (d *BaseDir) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
