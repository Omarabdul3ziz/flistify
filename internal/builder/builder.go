package builder

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/omarabdul3ziz/flistify/pkg/types"
)

// Builder represents a builder object
type Builder struct {
	FlistName string
	Handlers  map[string]func(string) error
	Paths     types.Paths
}

// NewBuilder creates a new Builder object with the given flist name
func NewBuilder(name string) (Builder, error) {
	bl := Builder{}
	bl.SetFlistName(name)
	bl.SetPaths()
	bl.SetHandlers()

	if err := createDirectoryIfNotExist(bl.Paths.RootFS); err != nil {
		return Builder{}, errors.Wrap(err, "failed to start the flist")
	}

	return bl, nil
}

// SetFlistName sets the flist name for the builder
func (bl *Builder) SetFlistName(name string) {
	if name == "" {
		name = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	bl.FlistName = name
}

// SetHundlers sets the handlers for each key
func (bl *Builder) SetHandlers() {
	bl.Handlers = map[string]func(string) error{
		"FROM":       bl.handleFrom,
		"KERNEL":     bl.handleKernel,
		"RUN":        bl.handleRun,
		"ENV":        bl.handleEnv,
		"ENTRYPOINT": bl.handleEntrypoint,
	}
}

// SetPaths sets the paths for the builder
func (bl *Builder) SetPaths() {
	bl.Paths.ProjectDir = "/var/lib/flistify"
	bl.Paths.RootFS = filepath.Join(bl.Paths.ProjectDir, "flists", bl.FlistName)
	bl.Paths.Boot = filepath.Join(bl.Paths.RootFS, "boot")
}

// Build executes the build process by reading the commands from the file at the given path
func (bl *Builder) Build(path string) error {
	file, err := os.Open(path)
	// TODO: validate the content of the file
	if err != nil {
		return errors.Wrapf(err, "failed to open file %v", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// TODO: same command in multiple lines with `\`
		if err := bl.HandleLine(scanner.Text()); err != nil {
			return errors.Wrapf(err, "failed to handle line %v", scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "scanner encountered an error")
	}

	return nil
}
