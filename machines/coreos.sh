#!/bin/bash -e

exec /usr/bin/qemu-system-x86_64 \
		-enable-kvm \
		-nographic \
		-m 8G \
		-net nic,macaddr=96:03:08:82:1C:02 \
		-net bridge,br=br0 \
		-drive file=/disks/coreos.qcow2,index=0,media=disk \
		-fsdev local,security_model=passthrough,id=fsdev0,path=/data \
		-device virtio-9p-pci,id=fs0,fsdev=fsdev0,mount_tag=hostdata
