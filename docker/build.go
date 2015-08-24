package main

import (
       "fmt"
	"flag"
	"log"
	"os"
	"os/exec"
)

var pull = flag.Bool("pull", false, "Whether to pull new base images")
var cache = flag.Bool("cache", true, "Whether to use the docker cache for builds")
var push = flag.Bool("push", false, "Whether to push the built versions")

var defaultPackages = []string{
	"btc",
	"btcd",
	"ddclient",
	"nginx",
	"transmission",
}
var bases = []string{
	"debian:jessie",
	"debian:sid",
}

func doPull(name string) error {
	cmd := exec.Command("docker", "pull", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func doBuildAndPush(name string) error {
	cmd := exec.Command("docker", "build", fmt.Sprintf("--no-cache=%t", !*cache), "-t", "hjfreyer/"+name, name)
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

type Task func() error

func doAll(tasks []Task) error {
	ch := make(chan error)
	for _, task := range tasks {
		task := task
		go func() {
			ch <- task()
		}()
	}

	var result error
	for _ = range tasks {
		if err := <-ch; err != nil {
			result = err
		}
	}
	return result
}

func main() {
	flag.Parse()
	packages := defaultPackages
	if len(flag.Args()) > 0 {
		packages = flag.Args()
	}

	if *pull {
		var pulls []Task
		for _, p := range bases {
			p := p
			pulls = append(pulls, func() error {
				return doPull(p)
			})
		}
		if err := doAll(pulls); err != nil {
			log.Fatal(err)
		}
	}

	var builds []Task
	for _, p := range packages {
		p := p
		builds = append(builds, func() error {
			return doBuildAndPush(p)
		})
	}
	if err := doAll(builds); err != nil {
		log.Fatal(err)
	}
}
