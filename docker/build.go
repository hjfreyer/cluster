package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"sync"
)

var pull = flag.Bool("pull", false, "Whether to pull new base images")
var push = flag.Bool("push", false, "Whether to push the built versions")

var defaultPackages = []string{
	"btcd",
	"btcd2",
	"transmission",
	"nginx",
}
var bases = []string{
	"debian:jessie",
	"golang",
}

func doPull(name string) error {
	cmd := exec.Command("docker", "pull", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func doBuildAndPush(name string) error {
	cmd := exec.Command("docker", "build", "-t", "hjfreyer/"+name, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if *push {
		cmd := exec.Command("docker", "push", "hjfreyer/"+name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

func main() {
	flag.Parse()
	packages := defaultPackages
	if len(flag.Args()) > 0 {
		packages = flag.Args()
	}

	if *pull {
		var pullGroup sync.WaitGroup
		pullGroup.Add(len(bases))
		for _, p := range bases {
			p := p
			go func() {
				if err := doPull(p); err != nil {
					log.Fatal(err)
				}
				pullGroup.Done()
			}()
		}
		pullGroup.Wait()
	}

	var buildGroup sync.WaitGroup
	buildGroup.Add(len(packages))
	for _, p := range packages {
		p := p
		go func() {
			if err := doBuildAndPush(p); err != nil {
				log.Fatal(err)
			}
			buildGroup.Done()
		}()
	}
	buildGroup.Wait()
}
