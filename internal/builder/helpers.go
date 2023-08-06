package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/omarabdul3ziz/flistify/internal/config"
	"github.com/omarabdul3ziz/flistify/pkg/types"
	"github.com/omarabdul3ziz/flistify/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func isEmptyOrComment(line string) bool {
	line = strings.TrimSpace(line)
	return line == "" || line[0] == '#'
}

func parseLine(line string) (string, string, error) {
	content := strings.SplitN(strings.TrimSpace(line), " ", 2)

	if len(content) < 2 {
		return "", "", errors.Errorf("invalid command format: %s", line)
	}
	return content[0], content[1], nil
}

func getBaseImageVersion(line string) (string, error) {
	content := strings.Split(strings.TrimSpace(line), ":")

	if len(content) < 2 {
		return "", errors.Errorf("invalid base format: %s", line)
	}

	image := content[0]
	version := content[1]

	if !slices.Contains(config.SUPPORTED_DISTROS, image) {
		return "", errors.Errorf("%v is an unsupported base image", image)
	}

	return version, nil
}

func parseRunLine(line string) (string, []string) {
	content := strings.Split(line, " ")
	return content[0], content[1:]
}

func runWithArchChroot(path string, cmd types.Command) error {
	cmd.Args = append([]string{path, cmd.Name}, cmd.Args...)
	return utils.ExecuteCommand(types.Command{
		Name: "arch-chroot",
		Args: cmd.Args,
	})
}

func updateAndClean(path string) error {
	cmds := []types.Command{
		{
			Name: "update-initramfs",
			Args: []string{"-c", "-k", "all"},
		},
		{
			Name: "apt",
			Args: []string{"clean"},
		},
		{
			Name: "cloud-init",
			Args: []string{"clean"},
		},
	}
	for idx := range cmds {
		if err := runWithArchChroot(path, cmds[idx]); err != nil {
			return errors.Wrapf(err, "failed running: %+v", cmds[idx])
		}
	}

	return nil
}

func editModulesFile(path string) error {
	moduleFilePath := filepath.Join(path, config.INITRAMSFS_PATH)

	moduleFile, err := os.OpenFile(moduleFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "couldn't open module file: %v", moduleFilePath)
	}

	if _, err = moduleFile.WriteString("\n" + config.VIRTIOFS + "\n"); err != nil {
		return errors.Wrapf(err, "couldn't add virtiofs")
	}

	return nil
}

func extractKernel(path string, kernelVersion string) error {
	// TODO: add info logs to track progress
	kernelName := fmt.Sprintf("vmlinuz-%v", kernelVersion)
	vmlinuzPath := filepath.Join(path, "/boot/vmlinuz")

	filename := fmt.Sprintf("%s%s", kernelName, ".elf")
	elfFilePath := filepath.Join(path, "/boot", filename)

	kerFilePath := filepath.Join(path, "/boot", kernelName)

	// Run the extract-vmlinux command and redirect its output to the output file
	extractCmd := exec.Command("/usr/bin/sudo", "extract-vmlinux", vmlinuzPath)

	outputFile, err := os.Create(elfFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed creating file: %v", elfFilePath)
	}
	defer outputFile.Close()
	extractCmd.Stdout = outputFile

	if err := extractCmd.Run(); err != nil {
		return errors.Wrapf(err, "failed extracting")
	}

	// Run the "sudo tee" command to redirect the output to /dev/null
	teeCmd := exec.Command("sudo", "tee", "/dev/null")
	teeCmd.Stdin = outputFile

	if err := teeCmd.Run(); err != nil {
		return errors.Wrapf(err, "error running tee command")
	}

	// mvCmd to rename elf file
	mvCmd := exec.Command("mv", elfFilePath, kerFilePath)
	if err := mvCmd.Run(); err != nil {
		return errors.Wrapf(err, "failed renaming the file")
	}

	return nil
}

/*
func runInsideChrootJail(chrootDir string, command types.Command) error {
	fmt.Printf("[+] Executing command \"%v\" with args \"%v\"\n", command.Name, command.Args)

	if !isExist(chrootDir) {
		return fmt.Errorf("chroot directory does not exist: %s", chrootDir)
	}

	exit, err := chroot(chrootDir)
	if err != nil {
		return fmt.Errorf("can't chroot: %v", err.Error())
	}

	if err := executeCommand(command); err != nil {
		return fmt.Errorf("can't run: %v", err.Error())
	}

	if err := exit(); err != nil {
		return fmt.Errorf("can't exit: %v", err.Error())
	}

	return nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func chroot(path string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(path); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}
*/
