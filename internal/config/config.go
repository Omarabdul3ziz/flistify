package config

const (
	FLISTIFY_ROOT   = "/var/lib/flistify"
	UBUNTU_ARCHIVE  = "http://archive.ubuntu.com/ubuntu"
	INITRAMSFS_PATH = "/etc/initramfs-tools/modules"
	VIRTIOFS        = "fs-virtiofs"
)

var (
	SUPPORTED_DISTROS = []string{"ubuntu"}
)
