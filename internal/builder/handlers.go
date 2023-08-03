package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var rootfsElements = []string{
	"/bin",
	"/dev",
	"/etc",
	"/lib",
	"/usr",
	"/var",
}

type Command struct {
	Name string
	Args []string
}

func handleFrom(line []string) error {
	// TODO: download flist from hub
	// TODO: load flist from cache

	var err error
	ubuntuVersion := line[0]
	directoryName := name

	if !directoryContainsRootFS(path) {
		cmd := exec.Command("debootstrap", ubuntuVersion, directoryName, "http://archive.ubuntu.com/ubuntu")
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}

	return err
}

func handleRun(line []string) error {
	cmd := Command{
		Name: line[0],
		Args: line[1:],
	}

	if err := runInside(path, cmd); err != nil {
		return fmt.Errorf("error handling run %+v: %+v", line, err.Error())
	}

	return nil
}

func handleKernel(line []string) error {
	// TODO: add validation on kernel version

	// install kernel
	cmd := Command{
		Name: "apt",
		Args: []string{"install", "-y", line[0]},
	}
	if err := runInside(path, cmd); err != nil {
		return fmt.Errorf("error handling run %+v: %+v", line, err.Error())
	}

	// update-initramfs
	moduleFilePath := filepath.Join(path, "/etc/initramfs-tools/modules")
	moduleFile, err := os.OpenFile(moduleFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open module file: %v", err.Error())
	}
	defer moduleFile.Close()

	_, err = moduleFile.WriteString("\nfs-virtiofs\n")
	if err != nil {
		return fmt.Errorf("couldn't add virtiofs: %v", err.Error())
	}

	if err := runInside(path, Command{Name: "update-initramfs", Args: []string{"-c", "-k", "all"}}); err != nil {
		return fmt.Errorf("couldn't update initramfs: %v", err.Error())
	}

	// // clean
	cmds := []Command{
		{
			Name: "apt-get",
			Args: []string{"clean"},
		},
		{
			Name: "cloud-init",
			Args: []string{"clean"},
		},
	}
	if err := runMultipleInside(path, cmds); err != nil {
		return fmt.Errorf("couldn't clean: %v", err.Error())
	}

	// // extract kernel

	if err := extractKernel(); err != nil {
		return fmt.Errorf("error extracting: %v", err.Error())
	}

	return nil
}

func handleEntrypoint(line []string) error {
	return nil
}

func handleEnv(line []string) error {
	return nil
}
