package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Containers []Container

func containerInList(container Container, list []string) bool {
	for _, listItem := range list {
		if os.ExpandEnv(listItem) == container.Name() {
			return true
		}
	}
	return false
}

func (containers Containers) filter(list []string) []Container {
	var filtered []Container
	for i := 0; i < len(containers); i++ {
		if containerInList(containers[i], list) {
			filtered = append(filtered, containers[i])
		}
	}
	return filtered
}

func (containers Containers) reversed() []Container {
	var reversed []Container
	for i := len(containers) - 1; i >= 0; i-- {
		reversed = append(reversed, containers[i])
	}
	return reversed
}

// Lift containers (provision + run).
// When forced, this will rebuild all images
// and recreate all containers.
func (containers Containers) lift(force bool, nocache bool, kill bool) {
	containers.provision(force, nocache)
	containers.runOrStart(force, kill)
}

// Provision containers.
// When forced, this will rebuild all images.
func (containers Containers) provision(force bool, nocache bool) {
	for _, container := range containers.reversed() {
		container.provision(force, nocache)
	}
}

// Run containers.
// When forced, removes existing containers first.
func (containers Containers) run(force bool, kill bool) {
	if force {
		containers.rm(force, kill)
	}
	for _, container := range containers.reversed() {
		container.run()
	}
}

// Run or start containers.
// When forced, removes existing containers first.
func (containers Containers) runOrStart(force bool, kill bool) {
	if force {
		containers.rm(force, kill)
	}
	for _, container := range containers.reversed() {
		container.runOrStart()
	}
}

// Start containers.
func (containers Containers) start() {
	for _, container := range containers.reversed() {
		container.start()
	}
}

// Kill containers.
func (containers Containers) kill() {
	for _, container := range containers {
		container.kill()
	}
}

// Stop containers.
func (containers Containers) stop() {
	for _, container := range containers {
		container.stop()
	}
}

// Remove containers.
// When forced, stops existing containers first.
func (containers Containers) rm(force bool, kill bool) {
	if force {
		if kill {
			containers.kill()
		} else {
			containers.stop()
		}
	}
	for _, container := range containers {
		container.rm()
	}
}

// Push containers.
func (containers Containers) push() {
	for _, container := range containers {
		container.push()
	}
}

// Status of containers.
func (containers Containers) status() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "Name\tRunning\tID\tIP\tPorts")
	for _, container := range containers {
		container.status(w)
	}
	w.Flush()
}
