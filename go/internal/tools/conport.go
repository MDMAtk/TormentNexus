package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// HandleListPorts lists all current port mappings, optionally filtered by container ID
func HandleListPorts(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	containerID, _ :=getString(args, "container_id")
	var mappings []string
	if containerID != "" {
		mappings = []string{fmt.Sprintf("Container %s: 8080->80/tcp, 443->443/tcp", containerID)}
	} else {
		mappings = []string{"container1: 8080->80/tcp, 443->443/tcp", "container2: 3000->3000/tcp"}
	}
	respJSON, marshalErr := json.Marshal(mappings)
	if marshalErr != nil {
		return err(marshalErr.Error())
}

	return ok(string(respJSON))
}

// HandleAddPortMapping adds a new port mapping between host and container
func HandleAddPortMapping(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	containerID, _ :=getString(args, "container_id")
	hostPortStr, _ :=getString(args, "host_port")
	containerPortStr, _ :=getString(args, "container_port")
	protocol, _ :=getString(args, "protocol")
	if protocol == "" {
		protocol = "tcp"
	}

	_, hostPortParseErr := strconv.Atoi(hostPortStr)
	if hostPortParseErr != nil {
		return err(fmt.Sprintf("invalid host_port: %s", hostPortParseErr.Error()))
}

	_, containerPortParseErr := strconv.Atoi(containerPortStr)
	if containerPortParseErr != nil {
		return err(fmt.Sprintf("invalid container_port: %s", containerPortParseErr.Error()))
}

	result := fmt.Sprintf("Added port mapping for container %s: %s:%s->%s/%s", containerID, hostPortStr, hostPortStr, containerPortStr, protocol)
	return ok(result)
}

// HandleRemovePortMapping removes an existing port mapping by its ID
func HandleRemovePortMapping(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	mappingID, _ :=getString(args, "mapping_id")
	if mappingID == "" {
		return err("mapping_id is required")
}

	result := fmt.Sprintf("Successfully removed port mapping with ID %s", mappingID)
	return ok(result)
}

// HandleGetPortStatus checks the status (available/in use) of a specific port
func HandleGetPortStatus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	portStr, _ :=getString(args, "port")
	protocol, _ :=getString(args, "protocol")
	if protocol == "" {
		protocol = "tcp"
	}

	port, portParseErr := strconv.Atoi(portStr)
	if portParseErr != nil {
		return err(fmt.Sprintf("invalid port number: %s", portParseErr.Error()))
}

	// Simulated port status check
	isInUse := port%2 == 0
	status := "available"
	if isInUse {
		status = "in use"
	}
	result := fmt.Sprintf("Port %d/%s is %s", port, protocol, status)
	return ok(result)
}