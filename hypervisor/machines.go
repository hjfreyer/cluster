package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nightlyone/lockfile"
)

const (
	baseDir    = "/data/virt/"
	diskSuffix = ".qcow2"

	qemuBin    = "/usr/bin/qemu-system-x86_64"
	nbdBin     = "/usr/bin/qemu-nbd"
	qemuImgBin = "/usr/bin/qemu-img"

	cloudConfigPath   = "data/cloud_config.yml"
	coreosInstallPath = "data/coreos-install"
)

type Machine struct {
	Name      string
	MemoryMb  int
	Mac       string
	DiskNames []string
}

type MachineSet struct {
	Machines []*Machine
}

var Leibniz = &Machine{
	Name:      "leibniz",
	MemoryMb:  7000,
	Mac:       "96:03:08:82:1C:01",
	DiskNames: []string{"main"},
}

var Coreos = &Machine{
	Name:      "coreos",
	MemoryMb:  7000,
	Mac:       "96:03:08:82:1C:02",
	DiskNames: []string{"main", "data", "blockchain"},
}

var Machines = []*Machine{Leibniz, Coreos}

func GetMachine(name string) *Machine {
	for _, m := range Machines {
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

type DiskId struct {
	ChainIdx int
	LinkIdx  int
}

var (
	diskNameRe = regexp.MustCompile(`^chain-(\d{5})-link-(\d{5}).qcow2$`)
)

func ParseDiskPath(path string) (DiskId, error) {
	group := diskNameRe.FindStringSubmatch(path)
	if group == nil {
		return DiskId{}, fmt.Errorf("invalid path %q", path)
	}
	chain, err := strconv.Atoi(group[1])
	if err != nil {
		return DiskId{}, err
	}
	link, err := strconv.Atoi(group[2])
	if err != nil {
		return DiskId{}, err
	}
	return DiskId{chain, link}, nil
}

type byChainAndLink []DiskId

func (b byChainAndLink) Len() int      { return len(b) }
func (b byChainAndLink) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byChainAndLink) Less(i, j int) bool {
	if b[i].ChainIdx != b[j].ChainIdx {
		return b[i].ChainIdx < b[j].ChainIdx
	}
	return b[i].LinkIdx < b[j].LinkIdx
}

func LatestDisk(m *Machine, didx int) (string, error) {
	diskDir := path.Join(baseDir, m.Name, m.DiskNames[didx])
	files, err := ioutil.ReadDir(diskDir)
	if err != nil {
		return "", err
	}

	parsed := make([]DiskId, len(files))
	for i, f := range files {
		parsed[i], err = ParseDiskPath(f.Name())
		if err != nil {
			return "", err
		}
	}

	sort.Sort(byChainAndLink(parsed))

	if len(parsed) == 0 {
		return "", errors.New("No files in disk directory")
	}

	prev := parsed[0]
	prev.LinkIdx--
	for _, disk := range parsed {
		if (prev.ChainIdx == disk.ChainIdx && prev.LinkIdx+1 == disk.LinkIdx) || (prev.ChainIdx+1 == disk.ChainIdx) {
			prev = disk
		} else {
			return "", errors.New("Missing disk file!")
		}
	}
	return path.Join(diskDir, fmt.Sprintf("chain-%05d-link-%05d.qcow2", prev.ChainIdx, prev.LinkIdx)), nil
}

func Ready(m *Machine) error {
	for i := 0; i < len(m.DiskNames); i++ {
		_, err := LatestDisk(m, i)
		return err
	}
	return nil
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

	for i := range m.DiskNames {
		latestPath, err := LatestDisk(m, i)
		if err != nil {
			return err
		}
		args = append(args, "-drive",
			fmt.Sprintf("file=%s,index=%d,media=disk", latestPath, i))
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
	disk, err := LatestDisk(Coreos, 0)
	if err != nil {
		return err
	}

	// Install the virtual disk.

	// This could be smarter and find an unused file but whatever.
	const devPath = "/dev/nbd0"
	nbdCmd := exec.Command(nbdBin, "-c", devPath, disk)
	redirect(nbdCmd)
	if err := nbdCmd.Run(); err != nil {
		return err
	}
	defer func() {
		exec.Command(nbdBin, "-d", devPath).Run()
	}()

	// Janky as hell.
	time.Sleep(2 * time.Second)

	// Find the root partition.
	blkidOut, err := exec.Command("blkid", "-t", "LABEL=ROOT", "-o", "device").Output()
	if err != nil {
		return err
	}

	partDev := ""
	for _, dev := range strings.Split(string(blkidOut), "\n") {
		if strings.HasPrefix(dev, devPath) {
			partDev = dev
			break
		}
	}
	if partDev == "" {
		return errors.New("could not find root partition")
	}

	// Mount the disk
	mntPath := path.Join(tmpdir, "nbdmnt")
	if err := os.Mkdir(mntPath, 0777); err != nil {
		return err
	}

	mntCmd := exec.Command("mount", string(partDev), mntPath)
	redirect(mntCmd)
	if err := mntCmd.Run(); err != nil {
		return err
	}
	defer func() {
		exec.Command("umount", mntPath).Run()
	}()

	cloudConfigDest := path.Join(mntPath, "var/lib/coreos-install/user_data")
	if err := os.MkdirAll(path.Dir(cloudConfigDest), 0777); err != nil {
		return err
	}

	if err := ioutil.WriteFile(cloudConfigDest, MustAsset(cloudConfigPath), 0666); err != nil {
		return err
	}
	return nil
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
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}
