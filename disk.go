package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const imageDir = "/var/lib/libvirt/images"

func createDisks(vmName string, specs DiskList) ([]string, func(), error) {
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		return nil, func() {}, fmt.Errorf("image directory does not exist: %s", imageDir)
	}

	var paths []string
	var created []string

	cleanup := func() {
		for _, p := range created {
			removeDisk(p)
		}
	}

	for i, d := range specs {
		path := filepath.Join(imageDir, fmt.Sprintf("%s-disk%d.qcow2", vmName, i))
		if err := runQemuImgCreate(path, d.Size); err != nil {
			return nil, cleanup, err
		}
		created = append(created, path)
		paths = append(paths, path)
	}

	return paths, cleanup, nil
}

func runQemuImgCreate(path, size string) error {
	cmd := exec.Command("sudo", "qemu-img", "create", "-f", "qcow2", path, size)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img create failed: %v: %s", err, output)
	}
	return nil
}

func removeDisk(path string) {
	cmd := exec.Command("sudo", "rm", "-f", path)
	cmd.Run()
}
