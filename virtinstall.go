package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func defineVM(name string, nics NICList, disks DiskList, diskPaths []string, memory, vcpus int) error {
	args := buildVirtInstallArgs(name, nics, disks, diskPaths, memory, vcpus)

	xml, err := runVirtInstallPrintXML(args)
	if err != nil {
		return fmt.Errorf("virt-install --print-xml failed: %w", err)
	}

	if err := virshDefine(xml); err != nil {
		return fmt.Errorf("virsh define failed: %w", err)
	}

	return nil
}

func buildVirtInstallArgs(name string, nics NICList, disks DiskList, diskPaths []string, memory, vcpus int) []string {
	args := []string{
		"virt-install",
		"--name", name,
		"--memory", fmt.Sprintf("%d", memory),
		"--vcpus", fmt.Sprintf("%d", vcpus),
		"--os-variant", "generic",
		"--boot", "uefi,hd,network",
		"--print-xml",
	}

	for i, d := range disks {
		spec := fmt.Sprintf("path=%s,format=qcow2,bus=%s", diskPaths[i], d.Bus)
		args = append(args, "--disk", spec)
	}

	for _, n := range nics {
		spec := fmt.Sprintf("bridge=%s,model=%s", n.Bridge, n.Model)
		args = append(args, "--network", spec)
	}

	return args
}

func runVirtInstallPrintXML(args []string) (string, error) {
	cmd := exec.Command("sudo", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func virshDefine(xml string) error {
	cmd := exec.Command("sudo", "virsh", "define", "/dev/stdin")
	cmd.Stdin = strings.NewReader(xml)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}
	return nil
}
