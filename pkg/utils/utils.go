package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/omarabdul3ziz/flistify/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func CreateDirectoryIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.Wrapf(err, "couldn't create directory: %v", path)
		}
	}
	return nil
}

func IsRootFS(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	essentialDirsAndFiles := []string{"/bin", "/etc", "/dev", "/lib", "/usr", "/var", "/etc/passwd", "/etc/shadow", "/etc/group"}

	for _, element := range essentialDirsAndFiles {
		if _, err := os.Stat(filepath.Join(path, element)); err != nil {
			return false
		}
	}

	return true
}

func ExecuteCommand(cmd types.Command) error {
	log.Info().Msgf("[+] Executing: %v %+v", cmd.Name, strings.Join(cmd.Args, " "))

	// TODO: conceder CommandContext
	command := exec.Command(cmd.Name, cmd.Args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
