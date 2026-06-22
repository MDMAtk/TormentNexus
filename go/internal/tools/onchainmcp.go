package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

const banklessBaseURL = "https://api.bankless.com"
const banklessTokenEnv = "BANKLESS_API_TOKEN"

func getAPIKey() string {
    return os.Getenv(banklessTokenEnv)
}

func apiRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
    apiKey := getAPIKey()
    if apiKey == "" {
        return nil, fmt.Errorf("BANKLESS_API_TOKEN environment variable is not set")
}

    var reqBody io.Reader
    if body != nil {
        jsonBody, jsonErr := json.Marshal(body)
        if jsonErr != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", jsonErr)
}

        reqBody = strings.NewReader(string(jsonBody))

    req, reqErr := http.NewRequestWithContext(ctx, method, banklessBaseURL+path, reqBody)
    if reqErr != nil {
        return nil, fmt.Errorf("failed to create request: %w", reqErr)
}

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-BANKLESS-TOKEN", apiKey)
    
    client := http.DefaultClient
    resp, fetchErr := client.Do(req)
    if fetchErr != nil {
        return nil, fmt.Errorf("request failed: %w", fetchErr)
}

    defer resp.Body.Close()
    
    if resp.StatusCode == http.StatusUnauthorized {
        return nil, fmt.Errorf("authentication failed: invalid API token")
}

    if resp.StatusCode == http.StatusTooManyRequests {
        return nil, fmt.Errorf("rate limit exceeded")
}

    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("resource not found")
}

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
}

    result, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
        return nil, fmt.Errorf("failed to read response: %w", readErr)
}

    return result, nil
}
val, _ :=getString(args, "key")
package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

const banklessBaseURL = "https://api.bankless.com"

}

func getAPIKey() string {
    return os.Getenv("BANKLESS_API_TOKEN")
}

func banklessAPIRequest(ctx context.Context, method, path string, query url.Values, body interface{}) ([]byte, error) {
    apiKey := getAPIKey()
    if apiKey == "" {
        return nil, fmt.Errorf("BANKLESS_API_TOKEN environment variable is not set")
}

    var reqBody io.Reader
    if body != nil {
        jsonBody, jsonMarshalErr := json.Marshal(body)
        if jsonMarshalErr != nil {
            return nil, fmt.Errorf("failed to marshal request: %w", jsonMarshalErr)
}

        reqBody = strings.NewReader(string(jsonBody))

    fullURL := banklessBaseURL + path
    if query != nil {
        fullURL = fullURL + "?" + query.Encode()

    req, reqErr := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
    if reqErr != nil {
        return nil, fmt.Errorf("failed to create request: %w", reqErr)
}

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-BANKLESS-TOKEN", apiKey)
    
    client := http.DefaultClient
    resp, fetchErr := client.Do(req)
    if fetchErr != nil {
        return nil, fmt.Errorf("request failed: %w", fetchErr)
}

    defer resp.Body.Close()
    
    respBody, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
        return nil, fmt.Errorf("failed to read response: %w", readErr)
}

    if resp.StatusCode == http.StatusUnauthorized {
        return nil, fmt.Errorf("authentication failed: invalid API token")
}

    if resp.StatusCode == http.StatusTooManyRequests {
        return nil, fmt.Errorf("rate limit exceeded")
}

    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("resource not found")
}

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(respBody))
}

    return respBody, nil
}

}
}

// HandleReadContract reads contract state from a blockchain
func HandleReadContract(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    network, _ :=getString(args, "network")
    contract, _ :=getString(args, "contract")
    method, _ :=getString(args, "method")
    
    inputs, found := args["inputs"].([]interface{})
    if !found {
        return err("inputs must be an array")
}

    outputs, found := args["outputs"].([]interface{})
    if !found {
        return err("outputs must be an array")
}

    payload := map[string]interface{}{
        "network": network,
        "contract": contract,
        "method": method,
        "inputs": inputs,
        "outputs": outputs,
    }
    
    respBody, apiErr := banklessAPIRequest(ctx, "POST", "/v1/contracts/read", nil, payload)
    if apiErr != nil {
        return err(apiErr.Error())
}

    return ok(string(respBody))
}

// HandleGetProxy gets the proxy address for a given network and contract
func HandleGetProxy(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    network, _ :=getString(args, "network")
    contract, _ :=getString(args, "contract")
    
    path := fmt.Sprintf("/v1/contracts/%s/%s/proxy", url.PathEscape(network), url.PathEscape(contract))
    
    respBody, apiErr := banklessAPIRequest(ctx, "GET", path, nil, nil)
    if apiErr != nil {
        return err(apiErr.Error())
}

    return ok(string(respBody))
}

// HandleGetAbi gets the ABI for a given contract on a specific network
func HandleGetAbi(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    network, _ :=getString(args, "network")
    contract, _ :=getString(args, "contract")
    
    path := fmt.Sprintf("/v1/contracts/%s/%s/abi", url.PathEscape(network), url.PathEscape(contract))
    
    respBody, apiErr := banklessAPIRequest(ctx, "GET", path, nil, nil)
    if apiErr != nil {
        return err(apiErr.Error())
}

    return ok(string(respBody))
}

// HandleGetEvents fetches event logs for a given network and filter criteria
func HandleGetEvents(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    network, _ :=getString(args, "network")
    
    addresses, found := args["addresses"].([]interface{})
    if !found {
        return err("addresses must be an array")
}

    topic, _ :=getString(args, "topic")
    
    var optionalTopics []interface{}
    if ot, hasOT := args["optionalTopics"]; hasOT {
        optionalTopics, ok = ot.([]interface{})
        if !found {
            return err("optionalTopics must be an array")

    }
    
    payload := map[string]interface{}{
        "network": network,
        "addresses": addresses,
        "topic": topic,
    }
    if optionalTopics != nil {
        payload["optionalTopics"] = optionalTopics
    }
    
    respBody, apiErr := banklessAPIRequest(ctx, "POST", "/v1/events", nil, payload)
    if apiErr != nil {
        return err(apiErr.Error())
}

    return ok(string(respBody))
}

}

// HandleBuildEventTopic builds an event topic signature based on event name and arguments
func HandleBuildEventTopic(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    network, _ :=getString(args, "network")
    name, _ :=getString(args, "name")
    
    arguments, found := args["arguments"].([]interface{})
    if !found {
        return err("arguments must be an array")
}

    payload := map[string]interface{}{
        "network": network,
        "name": name,
        "arguments": arguments,
    }
    
    respBody, apiErr := banklessAPIRequest(ctx, "POST", "/v1/events/topic", nil, payload)
    if apiErr != nil {
        return err(apiErr.Error())
}

    return ok(string(respBody))
}