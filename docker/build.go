package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

var pull = flag.Bool("pull", false, "Whether to pull new base images")

var packages = []string{
	"btcd",
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
func build(name string) error {
	cmd := exec.Command("docker", "build", "-t", "hjfreyer/"+name, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	flag.Parse()
	if *pull {
		for _, p := range bases {
			if err := doPull(p); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, p := range packages {
		if err := build(p); err != nil {
			log.Fatal(err)
		}
	}
}
