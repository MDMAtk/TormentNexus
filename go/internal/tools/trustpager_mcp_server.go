package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func HandleSearchContacts(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	limit, _ :=getInt(args, "limit")
	if limit <= 0 {
		limit = 20
	}
	apiKey := os.Getenv("TRUSTPAG")

---
*deepseek-reasoner (deepseek)*
}