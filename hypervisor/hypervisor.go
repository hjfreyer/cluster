package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type disk struct {
	Path string
}

type machine struct {
	MemoryMb int
	Mac      string
	Disks    []*disk
}

type context struct {
	Machines []*machine
	Disks    []*disk
}

func (c *context) InitLeibniz() {
	d := &disk{
		"/virt/images/leibniz.qcow2",
	}
	m := &machine{
		MemoryMb: 1024,
		Mac:      "96:03:08:82:1C:01",
		Disks:    []*disk{d},
	}
	c.Machines = append(c.Machines, m)
	c.Disks = append(c.Disks, d)
}

func (c *context) RunAll() {
	var cmds []*exec.Cmd
	for _, m := range c.Machines {
		cmd := getQemuCmd(m)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			log.Print(err)
		}
	}
}

func getQemuCmd(m *machine) *exec.Cmd {
	const qemuBin = "/usr/bin/qemu-system-x86_64"
	args := []string{
		"-enable-kvm",
		"-nographic",
		"-m", fmt.Sprintf("%dM", m.MemoryMb),
		"-net", "nic,macaddr=" + m.Mac,
		"-net", "bridge,br=br0",
	}

	for i, d := range m.Disks {
		args = append(args, "-drive",
			fmt.Sprintf("file=%s,index=%d,media=disk", d.Path, i))
	}

	return exec.Command(qemuBin, args...)
}

func main() {
	var c context
	c.InitLeibniz()
	c.RunAll()
}

// exec  \
// -enable-kvm \
// -m 1G -nographic \
// -hda /virt/images/leibniz.qcow2 \
// -net nic,macaddr=96:03:08:82:1C:01 -net bridge,br=br0
