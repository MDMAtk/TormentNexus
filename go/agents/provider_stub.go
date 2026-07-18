package agents

import (
	"context"
	"fmt"

	"github.com/MDMAtk/HyperNexus/foundation/adapters"
)

// DefaultProvider simulates the LLM locally.
type DefaultProvider struct{}

func (p *DefaultProvider) Chat(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
	prompt := ""
	if len(messages) > 0 {
		prompt = messages[len(messages)-1].Content
	}
	execution := adapters.PrepareProviderExecution(adapters.ProviderExecutionRequest{Prompt: prompt, CostPreference: "budget"})
	return Message{
		Role:    RoleAssistant,
		Content: fmt.Sprintf("I am the new Native Go HyperNexus Director. %s", execution.ExecutionHint),
	}, nil
}

func (p *DefaultProvider) Stream(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
	chunkChan <- "I am the "
	chunkChan <- "Native Go HyperNexus Director."
	close(chunkChan)
	return nil
}

func (p *DefaultProvider) GetModelName() string {
	return "hypernexus-native-stub-1.0"
}
