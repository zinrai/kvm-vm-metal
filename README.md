# kvm-vm-metal

Defines empty KVM virtual machines that imitate the hardware shape of a physical server.

## What it does

Creates empty libvirt domains for PXE-based OS installation testing.

## Prerequisites

- sudo
- qemu-img
- virsh
- virt-install
- ovmf

## Usage

```
$ kvm-vm-metal [OPTIONS] -nic ... -disk ... VM_NAME
```

At least one `-nic` and one `-disk` are required. Both flags are repeatable.

### NIC specification

`-nic bridge=NAME[,model=MODEL]`

- `bridge` (required): host bridge name
- `model` (default `e1000e`): NIC model. Examples: `e1000e`, `e1000`, `igb`, `virtio`, `rtl8139`

### Disk specification

`-disk size=SIZE[,bus=BUS]`

- `size` (required): disk size accepted by `qemu-img create` (e.g. `20G`, `500G`)
- `bus` (default `sata`): disk bus. Examples: `sata`, `virtio`, `scsi`

Disk image files are created under `/var/lib/libvirt/images/` as `VM_NAME-disk0.qcow2`, `VM_NAME-disk1.qcow2`, and so on.

### Examples

Minimal VM with one NIC and one disk:

```
$ kvm-vm-metal -nic bridge=br0 -disk size=50G node-01
```

Two NICs on different bridges and three disks:

```
$ kvm-vm-metal \
    -nic bridge=br-prov \
    -nic bridge=br-data \
    -disk size=20G \
    -disk size=500G \
    -disk size=500G \
    node-01
```

Custom CPU and memory:

```
$ kvm-vm-metal -memory 4096 -vcpus 4 -nic bridge=br0 -disk size=100G node-03
```

## License

This project is licensed under the [MIT License](./LICENSE).
