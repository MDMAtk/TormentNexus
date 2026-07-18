package cloud

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Container represents a user's Docker container
type Container struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Port        int       `json:"port"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	StoppedAt   time.Time `json:"stopped_at,omitempty"`
	MemoryLimit string    `json:"memory_limit"`
	CPULimit    string    `json:"cpu_limit"`
	VolumePath  string    `json:"volume_path"`
}

// Provisioner handles Docker container provisioning
type Provisioner struct {
	mu          sync.RWMutex
	containers  map[string]*Container
	basePort    int
	networkName string
	imageName   string
}

// NewProvisioner creates a new container provisioner
func NewProvisioner() *Provisioner {
	return &Provisioner{
		containers:  make(map[string]*Container),
		basePort:    8000,
		networkName: "hypernexus-cloud",
		imageName:   "hypernexus/kernel:latest",
	}
}

// Initialize sets up the Docker network
func (p *Provisioner) Initialize() error {
	// Create network if it doesn't exist
	cmd := exec.Command("docker", "network", "inspect", p.networkName)
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("docker", "network", "create", p.networkName)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("create network: %s: %w", string(output), err)
		}
	}
	return nil
}

// Provision creates a new container for an account
func (p *Provisioner) Provision(ctx context.Context, accountID, plan string) (*Container, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if account already has a container
	for _, c := range p.containers {
		if c.AccountID == accountID {
			return c, nil
		}
	}

	// Find available port
	port := p.findAvailablePort()

	// Create container
	container := &Container{
		ID:          fmt.Sprintf("hn-%s", accountID[:8]),
		AccountID:   accountID,
		Name:        fmt.Sprintf("hypernexus-%s", accountID[:8]),
		Status:      "creating",
		Port:        port,
		CreatedAt:   time.Now(),
		MemoryLimit: p.getMemoryLimit(plan),
		CPULimit:    p.getCPULimit(plan),
		VolumePath:  fmt.Sprintf("/var/lib/hypernexus/%s", accountID),
	}

	// Run Docker container
	cmd := exec.CommandContext(ctx, "docker", "run", "-d",
		"--name", container.Name,
		"--network", p.networkName,
		"--memory", container.MemoryLimit,
		"--cpus", container.CPULimit,
		"-p", fmt.Sprintf("%d:7778", port),
		"-v", fmt.Sprintf("%s:/data", container.VolumePath),
		"-e", fmt.Sprintf("HN_ACCOUNT_ID=%s", accountID),
		"-e", "HN_MODE=cloud",
		"--restart", "unless-stopped",
		p.imageName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker run: %s: %w", string(output), err)
	}

	container.Status = "running"
	container.StartedAt = time.Now()
	p.containers[container.ID] = container

	log.Printf("[Provisioner] Created container %s for account %s on port %d", container.ID, accountID, port)
	return container, nil
}

// Stop stops a container
func (p *Provisioner) Stop(ctx context.Context, containerID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	container, ok := p.containers[containerID]
	if !ok {
		return fmt.Errorf("container not found: %s", containerID)
	}

	cmd := exec.CommandContext(ctx, "docker", "stop", container.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker stop: %s: %w", string(output), err)
	}

	container.Status = "stopped"
	container.StoppedAt = time.Now()
	return nil
}

// Start starts a stopped container
func (p *Provisioner) Start(ctx context.Context, containerID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	container, ok := p.containers[containerID]
	if !ok {
		return fmt.Errorf("container not found: %s", containerID)
	}

	cmd := exec.CommandContext(ctx, "docker", "start", container.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker start: %s: %w", string(output), err)
	}

	container.Status = "running"
	container.StartedAt = time.Now()
	container.StoppedAt = time.Time{}
	return nil
}

// Destroy removes a container and its volumes
func (p *Provisioner) Destroy(ctx context.Context, containerID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	container, ok := p.containers[containerID]
	if !ok {
		return fmt.Errorf("container not found: %s", containerID)
	}

	// Stop and remove container
	exec.CommandContext(ctx, "docker", "stop", container.Name).Run()
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", container.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker rm: %s: %w", string(output), err)
	}

	// Remove volume
	exec.CommandContext(ctx, "rm", "-rf", container.VolumePath).Run()

	delete(p.containers, containerID)
	log.Printf("[Provisioner] Destroyed container %s", containerID)
	return nil
}

// GetContainer returns a container by ID
func (p *Provisioner) GetContainer(containerID string) (*Container, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	container, ok := p.containers[containerID]
	if !ok {
		return nil, fmt.Errorf("container not found: %s", containerID)
	}
	return container, nil
}

// GetContainerByAccount returns a container by account ID
func (p *Provisioner) GetContainerByAccount(accountID string) (*Container, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, c := range p.containers {
		if c.AccountID == accountID {
			return c, nil
		}
	}
	return nil, fmt.Errorf("container not found for account: %s", accountID)
}

// ListContainers returns all containers
func (p *Provisioner) ListContainers() []*Container {
	p.mu.RLock()
	defer p.mu.RUnlock()

	containers := make([]*Container, 0, len(p.containers))
	for _, c := range p.containers {
		containers = append(containers, c)
	}
	return containers
}

// SyncContainers syncs container status with Docker
func (p *Provisioner) SyncContainers(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, container := range p.containers {
		cmd := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{.State.Status}}", container.Name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			container.Status = "unknown"
			continue
		}

		status := strings.TrimSpace(string(output))
		switch status {
		case "running":
			container.Status = "running"
		case "exited":
			container.Status = "stopped"
		case "created":
			container.Status = "created"
		default:
			container.Status = status
		}
	}
	return nil
}

// findAvailablePort finds an available port
func (p *Provisioner) findAvailablePort() int {
	usedPorts := make(map[int]bool)
	for _, c := range p.containers {
		usedPorts[c.Port] = true
	}

	port := p.basePort
	for usedPorts[port] {
		port++
	}
	return port
}

// getMemoryLimit returns memory limit based on plan
func (p *Provisioner) getMemoryLimit(plan string) string {
	switch plan {
	case "enterprise":
		return "4g"
	case "pro":
		return "2g"
	case "starter":
		return "1g"
	default:
		return "512m"
	}
}

// getCPULimit returns CPU limit based on plan
func (p *Provisioner) getCPULimit(plan string) string {
	switch plan {
	case "enterprise":
		return "4"
	case "pro":
		return "2"
	case "starter":
		return "1"
	default:
		return "0.5"
	}
}
