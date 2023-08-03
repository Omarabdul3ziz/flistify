package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func startRootfs(flistName string) error {
	// TODO: check permission
	if !isExist(flistName) {
		err := os.Mkdir(flistName, 0666)
		if err != nil {
			return fmt.Errorf("couldn't create directory ")
		}
	}

	return nil
}

// func isRootFS(path string) bool {
// 	if !isExist(path) {
// 		return false
// 	}

// }

func directoryContainsRootFS(directoryPath string) bool {
	directoryContents, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return false
	}

	foundElements := make(map[string]bool)
	for _, file := range directoryContents {
		foundElements[file.Name()] = true
	}

	for _, element := range rootfsElements {
		if _, exists := foundElements[strings.TrimPrefix(element, "/")]; !exists {
			return false
		}
	}

	return true
}

func runInside(chrootDir string, command Command) error {
	fmt.Printf("[+] Executing command \"%v\" with args \"%v\"\n", command.Name, command.Args)

	if !isExist(chrootDir) {
		return fmt.Errorf("chroot directory does not exist: %s", chrootDir)
	}

	exit, err := Chroot(chrootDir)
	if err != nil {
		return fmt.Errorf("can't chroot: %v", err.Error())
	}

	// do some work
	if err := executeCommand(command); err != nil {
		return fmt.Errorf("can't run: %v", err.Error())
	}

	// exit from the chroot
	if err := exit(); err != nil {
		return fmt.Errorf("can't exit: %v", err.Error())
	}

	return nil
}

func runMultipleInside(chrootDir string, commands []Command) error {

	if !isExist(chrootDir) {
		return fmt.Errorf("chroot directory does not exist: %s", chrootDir)
	}

	exit, err := Chroot(chrootDir)
	if err != nil {
		return fmt.Errorf("can't chroot: %v", err.Error())
	}

	// do some work
	for idx := range commands {
		if err := executeCommand(commands[idx]); err != nil {
			return fmt.Errorf("can't run: %v", err.Error())
		}
	}

	// exit from the chroot
	if err := exit(); err != nil {
		return fmt.Errorf("can't exit: %v", err.Error())
	}

	return nil
}

func Chroot(path string) (func() error, error) {
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

func executeCommand(cmd Command) error {
	fmt.Println(cmd.Name, cmd.Args)
	command := exec.Command(cmd.Name, cmd.Args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func getKernelName() (string, error) {
	// TODO: this is tmp helper, get the name from the kernel in zerofile or any other way.
	directoryContents, err := ioutil.ReadDir(filepath.Join(path, "/boot"))
	if err != nil {
		return "", fmt.Errorf("error reading directory: %v", err.Error())
	}

	for _, file := range directoryContents {
		if strings.Contains(file.Name(), "vmlinuz-") {
			return file.Name(), nil
		}
	}

	return "", fmt.Errorf("couldn't fine vmlinuz file")
}

func extractKernel() error {
	kernelName, err := getKernelName()
	if err != nil {
		return fmt.Errorf("cound't get kernel name: %v", err.Error())
	}

	vmlinuzPath := filepath.Join(path, "/boot/vmlinuz")
	filename := fmt.Sprintf("%s%s", kernelName, ".elf")
	elfFilePath := filepath.Join(path, "/boot", filename)

	KerFilePath := filepath.Join(path, "/boot", kernelName)

	// Run the extract-vmlinux command and redirect its output to the output file
	extractCmd := exec.Command("/usr/bin/sudo", "extract-vmlinux", vmlinuzPath)

	outputFile, err := os.Create(elfFilePath)
	if err != nil {
		return fmt.Errorf("error creating elffile")
	}
	defer outputFile.Close()
	extractCmd.Stdout = outputFile

	if err := extractCmd.Run(); err != nil {
		return fmt.Errorf("error extracting: %v", err.Error())
	}

	// Run the "sudo tee" command to redirect the output to /dev/null
	teeCmd := exec.Command("sudo", "tee", "/dev/null")
	teeCmd.Stdin = outputFile

	if err := teeCmd.Run(); err != nil {
		return fmt.Errorf("error running tee command: %v", err.Error())
	}

	// mvCmd to rename elf file
	mvCmd := exec.Command("mv", elfFilePath, KerFilePath)
	if err := mvCmd.Run(); err != nil {
		return fmt.Errorf("couldn't move: %v", err.Error())
	}

	return nil
}
