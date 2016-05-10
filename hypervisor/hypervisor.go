package main

//go:generate go-bindata data/

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

// type disk struct {
// 	Path string
// }

// type machine struct {
// 	MemoryMb int
// 	Mac      string
// 	Disks    []*disk
// }

// type context struct {
// 	TempDir  string
// 	Machines []*machine
// 	Disks    []*disk
// }

// func (c *context) InitLeibniz() {
// 	c.Machines = append(c.Machines, m)
// 	c.Disks = append(c.Disks, d)
// }

// const (
// 	coreImgPath     = "/virt/images/coreos/"
// 	cloudConfigPath = "data/cloud_config.yml"
// )

// func (c *context) InitCore() {
// 	hash := cloudConfigHash()
// 	path := path.Join(coreImgPath, hash+".qcow2")
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		if err := initCoreOsDrive(path); err != nil {
// 			log.Fatal(err)
// 		}
// 	} else if err != nil {
// 		log.Fatal(err)
// 	}

// 	d := &disk{path}
// 	m := &machine{
// 		MemoryMb: 1024,
// 		Mac:      "96:03:08:82:1C:02",
// 		Disks:    []*disk{d},
// 	}
// 	c.Machines = append(c.Machines, m)
// 	c.Disks = append(c.Disks, d)
// }

// func initCoreOsDrive(path string) error {
// 	return nil
// }

// func cloudConfigHash() string {
// 	hash := sha256.Sum256(MustAsset(cloudConfigPath))
// 	return hex.Dump(hash[:])
// }

// func (c *context) getCloudConfig() (string, error) {
// 	if err := RestoreAsset(c.TempDir, cloudConfigPath); err != nil {
// 		return "", err
// 	}
// 	return path.Join(c.TempDir, cloudConfigPath), nil
// }

// func (c *context) RunAll() {
// 	var cmds []*exec.Cmd
// 	for _, m := range c.Machines {
// 		cmd := getQemuCmd(m)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		if err := cmd.Start(); err != nil {
// 			log.Fatal(err)
// 		}
// 		cmds = append(cmds, cmd)
// 	}

// 	for _, cmd := range cmds {
// 		if err := cmd.Wait(); err != nil {
// 			log.Print(err)
// 		}
// 	}
// }

// const nbdPath = "/usr/bin/qemu-nbd"

// func installDrive(d *disk) (string, error) {
// 	path := "/dev/nbd" + nextSuffix()
// 	cmd := exec.Command(nbdPath, "-c", path, d.Path)
// 	if err := cmd.Run(); err != nil {
// 		return "", err
// 	}
// 	return path, nil
// }

// func uninstallDrive(path string) error {
// 	cmd := exec.Command(nbdPath, "-d", path)
// 	return cmd.Run()
// }

// func getQemuCmd(m *machine) *exec.Cmd {
// 	const qemuBin = "/usr/bin/qemu-system-x86_64"
// 	args := []string{
// 		"-enable-kvm",
// 		"-nographic",
// 		"-m", fmt.Sprintf("%dM", m.MemoryMb),
// 		"-net", "nic,macaddr=" + m.Mac,
// 		"-net", "bridge,br=br0",
// 	}

// 	for i, d := range m.Disks {
// 		args = append(args, "-drive",
// 			fmt.Sprintf("file=%s,index=%d,media=disk", d.Path, i))
// 	}

// 	return exec.Command(qemuBin, args...)
// }

type command struct {
	tempDir string
}

func (c *command) parseMachine(ctx *cli.Context) *Machine {
	name := ctx.Args().Get(0)
	if name == "" {
		log.Fatal("Must specify a <machine>")
	}
	m := Repo.GetMachine(ctx.Args().Get(0))
	if m == nil {
		log.Fatalf("No machine named %q", name)
	}
	return m
}

func (c *command) init(ctx *cli.Context) {
	machine := c.parseMachine(ctx)
	if machine != Coreos {
		log.Fatal("I only know how to setup the coreos box")
	}
	if err := InitCoreos(c.tempDir); err != nil {
		log.Fatal(err)
	}
}

func (c *command) run(ctx *cli.Context) {
	machine := c.parseMachine(ctx)
	log.Fatal(Run(machine))
}

func main() {
	td, err := ioutil.TempDir("", "hypervisor")
	if err != nil {
		log.Fatal(err)
	}

	var c command
	c.tempDir = td

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Action: c.init,
		},
		{
			Name:   "run",
			Action: c.run,
		},
	}

	app.Run(os.Args)

	// var c context

	// c.InitLeibniz()
	// c.RunAll()

	os.RemoveAll(c.tempDir)
}

// exec  \
// -enable-kvm \
// -m 1G -nographic \
// -hda /virt/images/leibniz.qcow2 \
// -net nic,macaddr=96:03:08:82:1C:01 -net bridge,br=br0
