package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/nightlyone/lockfile"
)

const (
	baseDir    = "/data/virt/"
	diskSuffix = ".qcow2"

	qemuBin    = "/usr/bin/qemu-system-x86_64"
	nbdBin     = "/usr/bin/qemu-nbd"
	qemuImgBin = "/usr/bin/qemu-img"
)

type Importance int

const (
	ImportanceUndefined Importance = iota
	ImportanceCritical
	ImportanceMedium
	ImportanceTrivial
)

type Disk struct {
	Name       string
	Importance Importance
}

type Machine struct {
	Name     string
	MemoryMb int
	Mac      string
	Disks    []*Disk
}

type Repository struct {
	Machines []*Machine
	Disks    []*Disk
}

var LeibnizBoot = &Disk{
	Name:       "leibniz-boot",
	Importance: ImportanceMedium,
}

var PersonalData = &Disk{
	Name:       "personal",
	Importance: ImportanceCritical,
}

var GeneralData = &Disk{
	Name:       "data",
	Importance: ImportanceCritical,
}

var BlockChainDisk = &Disk{
	Name:       "blockchain",
	Importance: ImportanceMedium,
}

var CoreosBoot = &Disk{
	Name:       "coreos-boot",
	Importance: ImportanceTrivial,
}

var Leibniz = &Machine{
	Name:     "leibniz",
	MemoryMb: 7000,
	Mac:      "96:03:08:82:1C:01",
	Disks:    []*Disk{LeibnizBoot},
}

var Coreos = &Machine{
	Name:     "coreos",
	MemoryMb: 7000,
	Mac:      "96:03:08:82:1C:02",
	Disks:    []*Disk{CoreosBoot},
	//	Disks:    []*Disk{CoreosBoot, GeneralData, BlockChainDisk, PersonalData},
}

var Repo = &Repository{
	Machines: []*Machine{Leibniz, Coreos},
}

func (r *Repository) GetMachine(name string) *Machine {
	for _, m := range r.Machines {
		if m.Name == name {
			return m
		}
	}
	return nil
}

func LockMachine(m *Machine) (func(), error) {
	lf, err := lockfile.New(path.Join(baseDir, m.Name, "lock"))
	nilfn := func() { return }
	if err != nil {
		return nilfn, err
	}
	if err := lf.TryLock(); err != nil {
		return nilfn, err
	}

	unlock := func() {
		if err := lf.Unlock(); err != nil {
			log.Print("Error unlocking file: ", err)
		}
	}

	return unlock, nil
}

func Run(m *Machine) error {
	unlock, err := LockMachine(m)
	if err != nil {
		return err
	}
	defer unlock()

	args := []string{
		"-enable-kvm",
		"-nographic",
		"-m", fmt.Sprintf("%dM", m.MemoryMb),
		"-net", "nic,macaddr=" + m.Mac,
		"-net", "bridge,br=br0",
	}

	for i, d := range m.Disks {
		path := DiskPath(d)
		if _, err := os.Stat(path); err != nil {
			return err
		}
		args = append(args, "-drive",
			fmt.Sprintf("file=%s,index=%d,media=disk", path, i))
	}

	cmd := exec.Command(qemuBin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func redirect(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func InitCoreos(tmpdir string) error {
	// disk, err := LatestDisk(Coreos, 0)
	// if err != nil {
	// 	return err
	// }
	//
	// // Install the virtual disk.
	//
	// // This could be smarter and find an unused file but whatever.
	// const devPath = "/dev/nbd0"
	// nbdCmd := exec.Command(nbdBin, "-c", devPath, disk)
	// redirect(nbdCmd)
	// if err := nbdCmd.Run(); err != nil {
	// 	return err
	// }
	// defer func() {
	// 	exec.Command(nbdBin, "-d", devPath).Run()
	// }()
	//
	// // Janky as hell.
	// time.Sleep(2 * time.Second)
	//
	// // Find the root partition.
	// blkidOut, err := exec.Command("blkid", "-t", "LABEL=ROOT", "-o", "device").Output()
	// if err != nil {
	// 	return err
	// }
	//
	// partDev := ""
	// for _, dev := range strings.Split(string(blkidOut), "\n") {
	// 	if strings.HasPrefix(dev, devPath) {
	// 		partDev = dev
	// 		break
	// 	}
	// }
	// if partDev == "" {
	// 	return errors.New("could not find root partition")
	// }
	//
	// // Mount the disk
	// mntPath := path.Join(tmpdir, "nbdmnt")
	// if err := os.Mkdir(mntPath, 0777); err != nil {
	// 	return err
	// }
	//
	// mntCmd := exec.Command("mount", string(partDev), mntPath)
	// redirect(mntCmd)
	// if err := mntCmd.Run(); err != nil {
	// 	return err
	// }
	// defer func() {
	// 	exec.Command("umount", mntPath).Run()
	// }()
	//
	// cloudConfigDest := path.Join(mntPath, "var/lib/coreos-install/user_data")
	// if err := os.MkdirAll(path.Dir(cloudConfigDest), 0777); err != nil {
	// 	return err
	// }
	//
	// if err := ioutil.WriteFile(cloudConfigDest, MustAsset(cloudConfigPath), 0666); err != nil {
	// 	return err
	// }
	return nil
}

func DiskPath(disk *Disk) string {
	return path.Join(baseDir, "disks", disk.Name+diskSuffix)
}

// func InitCoreos(m *Machine, idx int, tmpdir string) error {
// 	disk := diskPath(m, idx)

// 	if exists, err := DiskExists(m, idx); err != nil {
// 		return err
// 	} else if exists {
// 		return errors.New("cowardly refusing to overwrite disk")
// 	}

// 	createCmd := exec.Command(qemuImgBin, "create", "-f", "qcow2", disk, fmt.Sprintf("%dG", m.Disks[idx].SizeGb))

// 	if err := createCmd.Run(); err != nil {
// 		return err
// 	}

// 	// sudo qemu-img create -f qcow2 /images/leibniz.qcow2 100G

// 	devPath := "/dev/nbd0"
// 	nbdCmd := exec.Command(nbdBin, "-c", devPath, disk)
// 	nbdCmd.Stdout = os.Stdout
// 	nbdCmd.Stderr = os.Stderr
// 	if err := nbdCmd.Run(); err != nil {
// 		return err
// 	}
// 	defer func() {
// 		exec.Command(nbdBin, "-d", devPath).Run()
// 	}()

// 	if err := RestoreAsset(tmpdir, coreosInstallPath); err != nil {
// 		return err
// 	}

// 	installer := path.Join(tmpdir, coreosInstallPath)
// 	installCmd := exec.Command(installer, "-d", devPath)
// 	installCmd.Stdout = os.Stdout
// 	installCmd.Stderr = os.Stderr

// 	return installCmd.Run()
// }

// func diskPath(m *Machine, index int) string {
// 	return path.Join(baseDir, m.Name, m.Disks[index].Name+diskSuffix)
// }

// NOTE: Code stolen from ioutil/tempfile.go
//
// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
// var rand uint32
// var randmu sync.Mutex
//
// func reseed() uint32 {
// 	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
// }
//
// func nextSuffix() string {
// 	randmu.Lock()
// 	r := rand
// 	if r == 0 {
// 		r = reseed()
// 	}
// 	r = r*1664525 + 1013904223 // constants from Numerical Recipes
// 	rand = r
// 	randmu.Unlock()
// 	return strconv.Itoa(int(1e9 + r%1e9))[1:]
// }
