package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var pull = flag.Bool("pull", false, "Whether to pull new base images")
var cache = flag.Bool("cache", true, "Whether to use the docker cache for builds")
var push = flag.Bool("push", false, "Whether to push the built versions")

var defaultPackages = []string{
	"btcwallet",
	"btcd",
	"ddclient",
	"nginx",
	"transmission",
}
var bases = []string{
	"debian:jessie",
	"debian:sid",
	"golang:latest",
}

type result struct {
	key string
	err error
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
		cmd := exec.Command("docker", "push", "hjfreyer/"+name+":latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

type Task struct {
	key string
	op  func() error
}

func doAll(tasks []Task) []result {
	ch := make(chan result)
	for _, task := range tasks {
		task := task
		go func() {
			ch <- result{task.key, task.op()}
		}()
	}

	var res []result
	for _ = range tasks {
		res = append(res, <-ch)
	}
	return res
}

func logResults(res []result) error {
	var err error
	for _, r := range res {
		if r.err == nil {
			log.Print(r.key + ": SUCCESS")
		} else {
			log.Print(r.key + ": FAIL")
			log.Print(r.err)
			err = r.err
		}
	}
	return err
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
			pulls = append(pulls, Task{
				key: "Pulling " + p,
				op: func() error {
					return doPull(p)
				}})
		}
		if anyErr := logResults(doAll(pulls)); anyErr != nil {
			log.Fatal("Not all pulls succeeded.")
		}
	}

	var builds []Task
	for _, p := range packages {
		p := p
		builds = append(builds, Task{
			key: "Build and/or push " + p,
			op: func() error {
				return doBuildAndPush(p)
			}})
	}

	if anyErr := logResults(doAll(builds)); anyErr != nil {
		log.Fatal("Not all build/pulls succeeded.")
	}
}
