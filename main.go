package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

func main() {
	fullDetails := flag.Bool("full", false, "Display full container details including GraphDriver, Node, SizeRw, and SizeRootFs")
	json := flag.Bool("json", false, "Display json container details")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf("USAGE: %s [--full] CONTAINER-ID\n", os.Args[0])
		os.Exit(1)
	}

	id := flag.Arg(0)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	cli.NegotiateAPIVersion(context.Background())

	container, rawData, err := cli.ContainerInspectWithRaw(context.Background(), id, false)
	if err != nil {
		log.Fatalf("Docker inspect for '%s' failed: %v", id, err)
	}

	if *json {
		fmt.Println(string(rawData))
	} else {
		printContainerDetails(container, *fullDetails)
	}

}

func printContainerDetails(container types.ContainerJSON, fullDetails bool) {
	state := container.State
	hostConfig := container.HostConfig

	fmt.Printf("Container ID: %s\n", container.ID)
	fmt.Printf("Name: %s\n", container.Name)
	fmt.Printf("Image: %s\n", container.Image)
	fmt.Printf("Created: %s\n", container.Created)
	fmt.Printf("Path: %s\n", container.Path)
	fmt.Printf("Args: %v\n", container.Args)
	fmt.Printf("Driver: %s\n", container.Driver)
	fmt.Printf("Platform: %s\n", container.Platform)
	fmt.Printf("MountLabel: %s\n", container.MountLabel)
	fmt.Printf("ProcessLabel: %s\n", container.ProcessLabel)
	fmt.Printf("AppArmorProfile: %s\n", container.AppArmorProfile)
	fmt.Printf("ExecIDs: %v\n", container.ExecIDs)

	if state != nil {
		fmt.Printf("Status: %s\n", state.Status)
		fmt.Printf("Running: %v\n", state.Running)
		fmt.Printf("Paused: %v\n", state.Paused)
		fmt.Printf("Restarting: %v\n", state.Restarting)
		fmt.Printf("OOMKilled: %v\n", state.OOMKilled)
		fmt.Printf("Dead: %v\n", state.Dead)
		fmt.Printf("PID: %d\n", state.Pid)
		fmt.Printf("ExitCode: %d\n", state.ExitCode)
		fmt.Printf("Error: %s\n", state.Error)
		fmt.Printf("StartedAt: %s\n", state.StartedAt)
		fmt.Printf("FinishedAt: %s\n", state.FinishedAt)
		if state.Health != nil {
			fmt.Printf("Health: %v\n", state.Health.Status)
		}
	}

	fmt.Printf("ResolvConfPath: %s\n", container.ResolvConfPath)
	fmt.Printf("HostnamePath: %s\n", container.HostnamePath)
	fmt.Printf("HostsPath: %s\n", container.HostsPath)
	fmt.Printf("LogPath: %s\n", container.LogPath)
	fmt.Printf("RestartCount: %d\n", container.RestartCount)

	if hostConfig != nil && fullDetails {
		fmt.Printf("HostConfig: %+v\n", hostConfig)
	}

	if fullDetails {
		if container.SizeRw != nil {
			fmt.Printf("SizeRw: %d\n", *container.SizeRw)
		}

		if container.SizeRootFs != nil {
			fmt.Printf("SizeRootFs: %d\n", *container.SizeRootFs)
		}

		if container.Node != nil {
			fmt.Printf("Node: %+v\n", container.Node)
		}

		fmt.Printf("GraphDriver: %+v\n", container.GraphDriver)
	}
}
