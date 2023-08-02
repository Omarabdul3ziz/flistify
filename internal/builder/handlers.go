package builder

import (
	"fmt"
	"os"
	"os/exec"
)

var rootfsElements = []string{
	"/bin",
	"/dev",
	"/etc",
	"/lib",
	"/usr",
	"/var",
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
	// TODO: which shell to use

	// runCommand := strings.Join(line, " ")
	// cmd := fmt.Sprintf("/bin/bash -c \"%+v\"", runCommand)
	fmt.Printf("[+] %v\n", line)
	err := runInside(path, line[0], line[1:])
	if err != nil {
		return fmt.Errorf("couldn't chroot and run %+v: %+v", line, err)
	}
	return nil
}

func handleKernel(line []string) error {
	return nil
}

func handleEntrypoint(line []string) error {
	return nil
}

func handleEnv(line []string) error {
	return nil
}
