package builder

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/omarabdul3ziz/flistify/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (bl *Builder) HandleLine(line string) error {
	if isEmptyOrComment(line) {
		return nil
	}

	key, value, err := parseLine(line)
	if err != nil {
		return err
	}

	handler, ok := bl.Handlers[key]
	if !ok {
		return errors.Errorf("%v is unsupported types.command", key)
	}

	if err := handler(value); err != nil {
		return errors.Wrapf(err, "failed to handle: %v", line)
	}

	return nil
}

func (bl *Builder) handleFrom(line string) error {
	// TODO: download flist from hub
	// TODO: load flist from cache
	// TODO: check by hash the files not just the existence of some directories
	if isRootFS(bl.Paths.RootFS) {
		log.Info().Msg("there is already a rootfs in this directory")
		return nil
	}

	version, err := getBaseImageVersion(line)
	if err != nil {
		return err
	}

	if _, err := exec.LookPath("debootstrap"); err != nil {
		return errors.Wrapf(err, "debootstrap command not found. check scripts/prepare.sh")
	}

	cmd := types.Command{
		Name: "debootstrap",
		Args: []string{version, bl.Paths.RootFS, UBUNTU_ARCHIVE},
	}

	if err := executeCommand(cmd); err != nil {
		return errors.Wrapf(err, "failed executing command: %+v", cmd)
	}

	return nil
}

func (bl *Builder) handleRun(line string) error {
	name, args := parseRunLine(line)

	cmd := types.Command{
		Name: name,
		Args: args,
	}

	if err := runWithArchChroot(bl.Paths.RootFS, cmd); err != nil {
		return errors.Wrapf(err, "failed to run command '%s %s': %+v", name, strings.Join(args, " "), cmd)
	}

	return nil
}

func (bl *Builder) handleKernel(line string) error {
	// TODO: add validation on kernel version

	cmd := types.Command{
		Name: "apt",
		Args: []string{"install", "-y", fmt.Sprintf("linux-modules-extra-%v", line)},
	}
	if err := runWithArchChroot(bl.Paths.RootFS, cmd); err != nil {
		return errors.Wrapf(err, "failed installing the kernel: %v", line)
	}

	if err := editModulesFile(bl.Paths.RootFS); err != nil {
		return errors.Wrap(err, "failed to edit modules file")
	}

	if err := updateAndClean(bl.Paths.RootFS); err != nil {
		return err
	}

	if err := extractKernel(bl.Paths.RootFS, line); err != nil {
		return errors.Wrap(err, "failed extract the kernel")
	}

	return nil
}

func (bl *Builder) handleEntrypoint(line string) error {
	// TODO: needed for container vm
	return nil
}

func (bl *Builder) handleEnv(line string) error {
	// TODO: needed for container vm
	return nil
}
