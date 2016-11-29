#!/bin/bash -e

yell() { echo "$0: $*" >&2; }
die() { yell "$*"; exit 111; }
try() { "$@" || die "cannot $*"; }

QEMU=/usr/bin/qemu-system-x86_64

machine=$1

[[ "$machine" == "" ]] && die "Usage: run_machine.sh MACHINE_NAME"

case "$machine" in
    leibniz)
	exec "$QEMU" \
	     -enable-kvm \
	     -nographic \
	     -m 1024M \
	     -net nic,macaddr=96:03:08:82:1C:01 \
	     -net bridge,br=br0 \
	     -drive file=/disks/leibniz.qcow2,index=0,media=disk \
	     -fsdev local,security_model=passthrough,id=fsdev0,path=/data \
	     -device virtio-9p-pci,id=fs0,fsdev=fsdev0,mount_tag=hostdata \
	     -fsdev local,security_model=passthrough,id=fsdev1,path=/disks \
	     -device virtio-9p-pci,id=fs1,fsdev=fsdev1,mount_tag=hostdisks \
	     -fsdev local,security_model=passthrough,id=fsdev2,path=/misc \
	     -device virtio-9p-pci,id=fs2,fsdev=fsdev2,mount_tag=hostmisc
	;;

    coreos)
	exec "$QEMU" \
	     -enable-kvm \
	     -nographic \
	     -smp 2 \
	     -m 12G \
	     -net nic,macaddr=96:03:08:82:1C:02 \
	     -net bridge,br=br0 \
	     -drive file=/disks/coreos.qcow2,index=0,media=disk \
	     -drive file=/data/nobackup/disks/coreos-local.qcow2,index=1,media=disk \
	     -fsdev local,security_model=passthrough,id=fsdev0,path=/data \
	     -device virtio-9p-pci,id=fs0,fsdev=fsdev0,mount_tag=hostdata
	;;

    *)
	die "Unknown machine '$machine'"
	;;
esac
