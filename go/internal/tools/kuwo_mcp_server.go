package tools

import "context"

func HandleSearchSongs(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	keyword, _ :=getString(args, "keyword")
	if keyword == "" {
		return err("keyword is required")
}

	page, _ :=getInt(args, "page")
	if page < 1 {
		page = 1
	}
	return ok(map[string]interface{}{
}
		"keyword": keyword,
		"page":    page,
		"message": "placeholder search result",
	})

func HandleGetSongInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ :=getString(args, "id")
	if id == "" {
		return err("id is required")
}

	return ok(map[string]interface{}{
}
		"id":   id,
		"info": "placeholder song info",
	})
}