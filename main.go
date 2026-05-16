package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	checkRequiredCommands()

	var nics NICList
	var disks DiskList

	flag.Var(&nics, "nic", "NIC specification (repeatable): bridge=NAME[,model=MODEL]")
	flag.Var(&disks, "disk", "Disk specification (repeatable): size=SIZE[,bus=BUS]")
	memory := flag.Int("memory", 1024, "Memory in MB")
	vcpus := flag.Int("vcpus", 1, "Number of virtual CPUs")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] -nic ... -disk ... VM_NAME\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: VM name is required as the last argument")
		flag.Usage()
		os.Exit(1)
	}
	vmName := args[0]

	if err := validate(nics, disks); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	diskPaths, cleanup, err := createDisks(vmName, disks)
	if err != nil {
		cleanup()
		log.Fatalf("Failed to create disks: %v", err)
	}

	if err := defineVM(vmName, nics, disks, diskPaths, *memory, *vcpus); err != nil {
		cleanup()
		log.Fatalf("Failed to define VM: %v", err)
	}

	fmt.Printf("VM '%s' has been defined.\n", vmName)
}

func validate(nics NICList, disks DiskList) error {
	if len(nics) == 0 {
		return fmt.Errorf("at least one -nic is required")
	}
	if len(disks) == 0 {
		return fmt.Errorf("at least one -disk is required")
	}
	return nil
}

func checkRequiredCommands() {
	for _, cmd := range []string{"qemu-img", "virsh", "virt-install"} {
		if _, err := exec.LookPath(cmd); err != nil {
			log.Fatalf("Required command not found: %s", cmd)
		}
	}
}
