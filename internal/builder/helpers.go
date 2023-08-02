package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func startRootfs(flistName string) error {
	// TODO: check permission
	if _, err := os.Stat(flistName); os.IsNotExist(err) {
		err := os.Mkdir(flistName, 0666)
		if err != nil {
			return fmt.Errorf("couldn't create directory ")
		}
	}

	return nil
}

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

// mountDirectories mounts essential directories into the chroot environment
func mountDirectories(chrootDir string) error {
	mounts := []struct {
		source string
		target string
		fstype string
		flags  uintptr
	}{
		{"/dev", chrootDir + "/dev", "bind", syscall.MS_BIND | syscall.MS_REC},
		{"/proc", chrootDir + "/proc", "proc", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC},
		{"/sys", chrootDir + "/sys", "sysfs", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC},
	}

	for _, m := range mounts {
		if err := syscall.Mount(m.source, m.target, m.fstype, m.flags, ""); err != nil {
			return fmt.Errorf("failed to mount %s: %v", m.target, err)
		}
	}

	return nil
}

// unmountDirectories unmounts the directories mounted for the chroot environment
func unmountDirectories(chrootDir string) {
	mounts := []string{"dev", "proc", "sys"}

	for _, m := range mounts {
		target := chrootDir + "/" + m
		if err := syscall.Unmount(target, syscall.MNT_DETACH); err != nil {
			fmt.Printf("failed to unmount %s: %v\n", target, err)
		}
	}
}

func exitChroot() error {
	// Change working directory to a location outside the chroot
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change working directory outside the chroot: %v", err)
	}
	return nil
}

func runInside(chrootDir, command string, args []string) error {
	// commands run with this function act as arch-chroot

	if _, err := os.Stat(chrootDir); os.IsNotExist(err) {
		return fmt.Errorf("chroot directory does not exist: %s", chrootDir)
	}

	// Mount necessary directories
	// if err := mountDirectories(chrootDir); err != nil {
	// 	return err
	// }
	// defer unmountDirectories(chrootDir)
	defer exitChroot()

	// Change root directory
	if err := syscall.Chroot(chrootDir); err != nil {
		return fmt.Errorf("failed to chroot: %v", err)
	}

	// Change working directory to the new root
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change working directory to new root: %v", err)
	}

	// Execute the provided command inside the chroot environment
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}
