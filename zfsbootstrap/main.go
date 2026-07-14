//go:build linux

package main

import (
	"log/slog"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	slog.Info("Building ZFS drive in Sirius-like environment")

	if err := MountFs(); err != nil {
		slog.Error("Could not mount necessary devices/fs elements Reason: " + err.Error())
	}

	if err := loadZFS(); err != nil {
		slog.Error("Could not load ZFS kernel modules Reason: " + err.Error())
	}

	if err := createZFS(); err != nil {
		slog.Error("Could not correctly format zfs Reason: " + err.Error())
	}

	unix.Sync()
	time.Sleep(1 * time.Second)

	err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		slog.Error("Failed to shutdown qemu environment... Plese manually intervene")
		os.Exit(0)
	}
}

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

func createZFS() error {
	createCmd := exec.Command("/bin/zpool", "create", "sirius", "/dev/vda")
	createCmd.Stdout, createCmd.Stderr = os.Stdout, os.Stderr
	if err := createCmd.Run(); err != nil {
		slog.Error("Could not run zpool Reason: " + err.Error())
		return err
	}

	if err := os.CopyFS("/sirius/", os.DirFS("/sirius-fs/")); err != nil {
		slog.Error("Couldn't copy payload FS into drive Reason: " + err.Error())
		return err
	}

	finalizeCmd := exec.Command("/bin/zpool", "export", "sirius")
	finalizeCmd.Stdout, finalizeCmd.Stderr = os.Stdout, os.Stderr
	if err := finalizeCmd.Run(); err != nil {
		slog.Error("Failed to finalize FS Reason: " + err.Error())
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
