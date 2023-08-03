package mounter

import (
	"fmt"
	"syscall"
)

func Run() error {
	return nil
}

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
