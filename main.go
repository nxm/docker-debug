package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"log"
	"os"
)

type Command struct {
	Name        string
	Description string
	Parameters  string
	Execute     func(args []string)
}

const ContainerId = "CONTAINER-ID"

func main() {
	var commands []Command
	commands = []Command{
		{
			Name:        "inspect",
			Description: "Display container details. Use --full for more information.",
			Parameters:  fmt.Sprintf("%s [--full]", ContainerId),
			Execute:     handleInspectCommand,
		},
		{
			Name:        "usage",
			Description: "Display process usage statistics for a container.",
			Parameters:  ContainerId,
			Execute:     handleUsageCommand,
		},
		{
			Name:        "json",
			Description: "Display container details in JSON format.",
			Parameters:  ContainerId,
			Execute:     handleJsonCommand,
		},
		{
			Name:        "help",
			Description: "Display this help message.",
			Parameters:  "",
			Execute:     func(args []string) { printHelp(commands) },
		},
	}

	if len(os.Args) < 2 {
		printHelp(commands)
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	for _, cmd := range commands {
		if cmd.Name == commandName {
			cmd.Execute(args)
			return
		}
	}

	fmt.Printf("Unknown command: %s\n", commandName)
	printHelp(commands)
	os.Exit(1)
}

func handleInspectCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("USAGE: inspect CONTAINER-ID [--full]")
		os.Exit(1)
	}

	id := args[0]
	fullDetails := len(args) > 1 && args[1] == "--full"

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	cli.NegotiateAPIVersion(context.Background())

	container, _, err := cli.ContainerInspectWithRaw(context.Background(), id, false)
	if err != nil {
		log.Fatalf("docker inspect for '%s' failed: %v", id, err)
	}

	printContainerDetails(container, fullDetails)
}

func handleUsageCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("USAGE: usage CONTAINER-ID")
		os.Exit(1)
	}

	id := args[0]

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	cli.NegotiateAPIVersion(context.Background())

	container, _, err := cli.ContainerInspectWithRaw(context.Background(), id, false)
	if err != nil {
		log.Fatalf("docker inspect for '%s' failed: %v", id, err)
	}

	fmt.Println("Usage:")
	getProcessUsage(int32(container.State.Pid))
}

func handleJsonCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("USAGE: json CONTAINER-ID")
		os.Exit(1)
	}

	id := args[0]

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	cli.NegotiateAPIVersion(context.Background())

	_, rawData, err := cli.ContainerInspectWithRaw(context.Background(), id, false)
	if err != nil {
		log.Fatalf("docker inspect for '%s' failed: %v", id, err)
	}

	fmt.Println(string(rawData))
}

func printHelp(commands []Command) {
	fmt.Println("USAGE:")
	for _, cmd := range commands {
		fmt.Printf("  %s %s - %s\n", cmd.Name, cmd.Parameters, cmd.Description)
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
