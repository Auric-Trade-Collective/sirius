//go:build linux

package kernel

import (
	"log/slog"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

const ZFS_POOL_NAME = "sirius"

func loadZFS() error {
	modules := []string{
		"/lib/modules/spl.ko",
		"/lib/modules/zfs.ko",
	}

	for _, mod := range modules {
		fle, err := os.Open(mod)
		if err != nil {
			return err
		}

		err = unix.FinitModule(int(fle.Fd()), "", 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func MountZFS() error {
	if err := loadZFS(); err != nil {
		slog.Error("Could not load ZFS kernel modules Reason: " + err.Error())
		return err
	}

	importCmd := exec.Command("/bin/zpool", "import", "-d", "/dev", ZFS_POOL_NAME)
	importCmd.Stdout, importCmd.Stderr = os.Stdout, os.Stderr
	if err := importCmd.Run(); err != nil {
		slog.Error("Could not run zpool Reason: " + err.Error())
		return err
	}

	mountCmd := exec.Command("/bin/zfs", "mount", "-a")
	mountCmd.Stdout, mountCmd.Stderr = os.Stdout, os.Stderr
	if err := mountCmd.Run(); err != nil {
		slog.Error("Could not run zfs Reason: " + err.Error())
		return err
	}

	if err := shareMountsWithSirius(); err != nil {
		slog.Error("Failed to share important system mounts Reason: " + err.Error())
		return err
	}

	if err := unix.Chroot("/sirius/"); err != nil {
		slog.Error("Could not chroot into main FS Reason: " + err.Error())
		return err
	}

	return nil
}

func MountFs() error {
	err := unix.Mount("none", "/dev", "devtmpfs", unix.MS_NOSUID, "")
	if err != nil {
		slog.Error("Could not mount filesystem...")
		return err
	}

	err = unix.Mount("proc", "/proc", "proc", 0, "")
    if err != nil {
        return err
    }

    err = unix.Mount("sysfs", "/sys", "sysfs", 0, "")
    if err != nil {
        return err
    }

	return nil
}

func shareMountsWithSirius() error {
	mounts := []string{"/dev", "/proc", "/sys"}

	for _, mnt := range mounts {
		target := "/sirius" + mnt

		err := unix.Mount(mnt, target, "", unix.MS_BIND, "")
		if err != nil {
			slog.Error("Failed to bind device into sirius environment Reason: " + err.Error())
			return err
		}
	}
	return nil
}
