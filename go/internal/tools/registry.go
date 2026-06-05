package tools

/**
 * @file registry.go
 * @module go/internal/tools
 *
 * WHAT: Go-native registry for standard library and parity tools.
 * Maps tool names to their native Go implementations.
 */

import (
	"context"
	"fmt"
)

type ToolHandler func(ctx context.Context, args map[string]interface{}) (ToolResponse, error)

type Registry struct {
	handlers map[string]ToolHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		handlers: make(map[string]ToolHandler),
	}
	r.registerAll()
	return r
}

func (r *Registry) registerAll() {
	// Native Handlers
	r.handlers["read_file"] = HandleRead
	r.handlers["write_file"] = HandleWrite
	r.handlers["edit_file"] = HandleEdit
	r.handlers["str_replace_editor"] = HandleEdit
	r.handlers["grep_search"] = HandleGrep
	r.handlers["search_files"] = HandleGrep
	r.handlers["glob"] = HandleGlob
	r.handlers["find_files"] = HandleGlob
	r.handlers["apply_patch"] = HandleApplyPatch
	r.handlers["multi_edit"] = HandleMultiEdit
	r.handlers["bash"] = HandleBash
	r.handlers["ls"] = HandleLS
	r.handlers["list_directory"] = HandleLS
	r.handlers["web_fetch"] = HandleWebFetch
	r.handlers["ingest_git"] = HandleGitIngest
	r.handlers["sqlite_get_catalog"] = HandleSqliteGetCatalog
	r.handlers["sqlite_execute"] = HandleSqliteExecute
	r.handlers["search"] = HandleDDGSearch
	r.handlers["fetch_content"] = HandleDDGFetchContent
	r.handlers["slack_list_channels"] = HandleSlackListChannels
	r.handlers["slack_post_message"] = HandleSlackPostMessage
	r.handlers["slack_reply_to_thread"] = HandleSlackReplyToThread
	r.handlers["slack_add_reaction"] = HandleSlackAddReaction
	r.handlers["slack_get_channel_history"] = HandleSlackGetChannelHistory
	r.handlers["slack_get_thread_replies"] = HandleSlackGetThreadReplies
	r.handlers["slack_get_users"] = HandleSlackGetUsers
	r.handlers["slack_get_user_profile"] = HandleSlackGetUserProfile

	// Filesystem MCP Tools
	r.handlers["read_text_file"] = HandleReadTextFile
	r.handlers["create_directory"] = HandleCreateDirectory
	r.handlers["list_directory"] = HandleListDirectory
	r.handlers["list_directory_with_sizes"] = HandleListDirectoryWithSizes
	r.handlers["directory_tree"] = HandleDirectoryTree
	r.handlers["move_file"] = HandleMoveFile
	r.handlers["get_file_info"] = HandleGetFileInfo
	r.handlers["search_files"] = HandleSearchFiles

	// Ollama MCP Tools (AI & LLM Integration)
	r.handlers["list_local_models"] = HandleListLocalModels
	r.handlers["local_llm_chat"] = HandleLocalLLMChat
	r.handlers["ollama_health_check"] = HandleOllamaHealthCheck
	r.handlers["system_resource_check"] = HandleSystemResourceCheck

	// TTS MCP Tools (Media & Design)
	r.handlers["say_tts"] = HandleSayTTS
	r.handlers["openai_tts"] = HandleOpenAITTS

	// Vercel MCP Tools (Cloud & DevOps)
	r.handlers["vercel_list_projects"] = HandleVercelListProjects
	r.handlers["vercel_get_project"] = HandleVercelGetProject
	r.handlers["vercel_list_deployments"] = HandleVercelListDeployments
	r.handlers["vercel_get_deployment"] = HandleVercelGetDeployment
	r.handlers["vercel_cancel_deployment"] = HandleVercelCancelDeployment
	r.handlers["vercel_list_env_vars"] = HandleVercelListEnvVars
	r.handlers["vercel_create_env_var"] = HandleVercelCreateEnvVar
	r.handlers["vercel_delete_env_var"] = HandleVercelDeleteEnvVar

	// DexPaprika MCP Tools (Finance & Crypto)
	r.handlers["getCapabilities"] = HandleDexPaprikaGetCapabilities
	r.handlers["getNetworks"] = HandleDexPaprikaGetNetworks
	r.handlers["getStats"] = HandleDexPaprikaGetStats
	r.handlers["search"] = HandleDexPaprikaSearch
	r.handlers["getNetworkDexes"] = HandleDexPaprikaGetNetworkDexes
	r.handlers["getNetworkPools"] = HandleDexPaprikaGetNetworkPools
	r.handlers["getDexPools"] = HandleDexPaprikaGetDexPools
	r.handlers["getNetworkPoolsFilter"] = HandleDexPaprikaGetNetworkPoolsFilter
	r.handlers["getPoolDetails"] = HandleDexPaprikaGetPoolDetails
	r.handlers["getPoolOHLCV"] = HandleDexPaprikaGetPoolOHLCV
	r.handlers["getPoolTransactions"] = HandleDexPaprikaGetPoolTransactions
	r.handlers["getTokenDetails"] = HandleDexPaprikaGetTokenDetails
	r.handlers["getTokenPools"] = HandleDexPaprikaGetTokenPools
	r.handlers["getTokenMultiPrices"] = HandleDexPaprikaGetTokenMultiPrices
	r.handlers["filterNetworkTokens"] = HandleDexPaprikaFilterNetworkTokens
	r.handlers["getTopTokens"] = HandleDexPaprikaGetTopTokens
	r.handlers["submitFeedback"] = HandleDexPaprikaSubmitFeedback

	// National Weather Service (NWS) MCP Tools (Weather & Location)
	r.handlers["nws_get_forecast"] = HandleNWSGetForecast
	r.handlers["nws_search_alerts"] = HandleNWSSearchAlerts
	r.handlers["nws_get_observations"] = HandleNWSGetObservations
	r.handlers["nws_find_stations"] = HandleNWSFindStations
	r.handlers["nws_list_alert_types"] = HandleNWSListAlertTypes
	r.handlers["nws_get_office_discussion"] = HandleNWSGetOfficeDiscussion
	r.handlers["nws_get_zone_forecast"] = HandleNWSGetZoneForecast

	// ast-grep-mcp Tools (Category 11)
	r.handlers["ast_grep_dump_syntax_tree"] = HandleDumpSyntaxTree
	r.handlers["ast_grep_test_match_code_rule"] = HandleTestMatchCodeRule
	r.handlers["ast_grep_find_code"] = HandleFindCode
	r.handlers["ast_grep_find_code_by_rule"] = HandleFindCodeByRule

	// PAL Tools (Category 12)
	r.handlers["pal_chat"] = HandlePalChat
	r.handlers["pal_thinkdeep"] = HandlePalThinkDeep
	r.handlers["pal_planner"] = HandlePalPlanner
	r.handlers["pal_consensus"] = HandlePalConsensus
	r.handlers["pal_codereview"] = HandlePalCodeReview
	r.handlers["pal_precommit"] = HandlePalPrecommit
	r.handlers["pal_debug"] = HandlePalDebug
	r.handlers["pal_challenge"] = HandlePalChallenge

	// Short/alias mappings for PAL tools without prefix
	r.handlers["chat"] = HandlePalChat
	r.handlers["thinkdeep"] = HandlePalThinkDeep
	r.handlers["planner"] = HandlePalPlanner
	r.handlers["consensus"] = HandlePalConsensus
	r.handlers["codereview"] = HandlePalCodeReview
	r.handlers["precommit"] = HandlePalPrecommit
	r.handlers["debug"] = HandlePalDebug
	r.handlers["challenge"] = HandlePalChallenge

	// Serena Tools (Category 13)
	r.handlers["get_symbols_overview"] = HandleGetSymbolsOverview
	r.handlers["find_symbol"] = HandleFindSymbol
	r.handlers["find_referencing_symbols"] = HandleFindReferencingSymbols
	r.handlers["find_implementations"] = HandleFindImplementations
	r.handlers["find_declaration"] = HandleFindDeclaration
	r.handlers["rename_symbol"] = HandleRenameSymbol
	r.handlers["onboarding"] = HandleOnboarding

	// Claude Code Aliases
	r.handlers["Read"] = HandleRead
	r.handlers["Write"] = HandleWrite
	r.handlers["Edit"] = HandleEdit
	r.handlers["Bash"] = HandleBash
	r.handlers["LS"] = HandleLS
	r.handlers["WebFetch"] = HandleWebFetch
	r.handlers["Glob"] = HandleGlob
	r.handlers["Grep"] = HandleGrep
	r.handlers["MultiEdit"] = HandleMultiEdit

	// Codex Aliases
	r.handlers["shell"] = HandleBash
	r.handlers["create_file"] = HandleWrite
	r.handlers["view_file"] = HandleRead
	r.handlers["apply_diff"] = HandleApplyPatch
	r.handlers["search_files_codex"] = HandleGrep

	// OpenCode / Pi Aliases
	r.handlers["read"] = HandleRead
	r.handlers["write"] = HandleWrite
	r.handlers["edit"] = HandleEdit
	r.handlers["grep"] = HandleGrep
	r.handlers["ls"] = HandleLS
	r.handlers["glob_pi"] = HandleGlob

	// Thoughtbox Tools (Category 14)
	r.handlers["thoughtbox_search"] = HandleThoughtboxSearch
	r.handlers["thoughtbox_execute"] = HandleThoughtboxExecute
	r.handlers["thoughtbox_peer_notebook"] = HandleThoughtboxPeerNotebook

	// Fetch Tool (Assimilated)
	r.handlers["fetch"] = HandleFetch

	// Tavily Tools (Assimilated)
	r.handlers["tavily-search"] = HandleTavilySearch

	// Chrome DevTools Tools (Assimilated)
	r.handlers["chrome-devtools"] = HandleChromeDevTools

	// Firecrawl Tools (Assimilated from firecrawl-mcp)
	r.handlers["firecrawl_scrape"] = HandleFirecrawl
	r.handlers["firecrawl_crawl"] = HandleFirecrawl
	r.handlers["firecrawl"] = HandleFirecrawl

	// Exa Search Tools (Assimilated from SSE exa)
	r.handlers["exa_search"] = HandleExaSearch
	r.handlers["exa_find_similar"] = HandleExaFindSimilar
	r.handlers["exa_get_contents"] = HandleExaGetContents

	// arXiv Tools (Assimilated from arxiv-mcp-server)
	r.handlers["arxiv_search"] = HandleArxivSearch
	r.handlers["arxiv_get_paper"] = HandleArxivGetPaper
	r.handlers["arxiv_list_recent"] = HandleArxivListRecent

	// Semantic Scholar Tools (Assimilated from paper_search_server)
	r.handlers["paper_search"] = HandleSemanticScholarSearch
	r.handlers["paper_details"] = HandleSemanticScholarGetPaper
	r.handlers["paper_citations"] = HandleSemanticScholarGetCitations
	r.handlers["semantic_scholar_search"] = HandleSemanticScholarSearch
	r.handlers["semantic_scholar_paper"] = HandleSemanticScholarGetPaper

	// mem0 Memory Tools (Assimilated from @mem0/mcp-server)
	r.handlers["mem0_add_memory"] = HandleMem0AddMemory
	r.handlers["mem0_search_memory"] = HandleMem0SearchMemory
	r.handlers["mem0_get_memories"] = HandleMem0GetMemories
	r.handlers["mem0_delete_memory"] = HandleMem0DeleteMemory
	r.handlers["mem0_update_memory"] = HandleMem0UpdateMemory
	r.handlers["add_memory"] = HandleMem0AddMemory
	r.handlers["search_memory"] = HandleMem0SearchMemory

	// Alpaca Trading Tools (Assimilated from alpaca-mcp-server)
	r.handlers["alpaca_get_account"] = HandleAlpacaGetAccount
	r.handlers["alpaca_get_positions"] = HandleAlpacaGetPositions
	r.handlers["alpaca_get_orders"] = HandleAlpacaGetOrders
	r.handlers["alpaca_place_order"] = HandleAlpacaPlaceOrder
	r.handlers["alpaca_cancel_order"] = HandleAlpacaCancelOrder
	r.handlers["alpaca_get_bars"] = HandleAlpacaGetBars
	r.handlers["alpaca_get_latest_quote"] = HandleAlpacaGetLatestQuote

	// Alpha Vantage Financial Tools (Assimilated from av-mcp)
	r.handlers["av_quote"] = HandleAVGlobalQuote
	r.handlers["av_time_series"] = HandleAVTimeSeries
	r.handlers["av_forex_rate"] = HandleAVForexRate
	r.handlers["av_crypto_rate"] = HandleAVCryptoRate
	r.handlers["av_symbol_search"] = HandleAVSearch
	r.handlers["av_economic_indicator"] = HandleAVEconomicIndicator
	r.handlers["alpha_vantage_quote"] = HandleAVGlobalQuote

	// Hugging Face Hub Tools (Assimilated from SSE huggingface)
	r.handlers["hf_search_models"] = HandleHFSearchModels
	r.handlers["hf_get_model"] = HandleHFGetModel
	r.handlers["hf_search_datasets"] = HandleHFSearchDatasets
	r.handlers["hf_text_generation"] = HandleHFTextGeneration
	r.handlers["hf_classify_text"] = HandleHFClassification
	r.handlers["hf_embeddings"] = HandleHFEmbeddings
	r.handlers["hf_search_spaces"] = HandleHFSearchSpaces

	// Semgrep Security Tools (Assimilated from semgrep + semgrepstream)
	r.handlers["semgrep_scan"] = HandleSemgrepScan
	r.handlers["semgrep_cloud_scan"] = HandleSemgrepCloudScan
	r.handlers["semgrep_search_rules"] = HandleSemgrepRuleSearch

	// Octagon Financial Intelligence (Assimilated from octagon + octagon-deep-research)
	r.handlers["octagon_research"] = HandleOctagonResearch
	r.handlers["octagon_company_search"] = HandleOctagonCompanySearch
	r.handlers["octagon_financials"] = HandleOctagonFinancials
	r.handlers["octagon_news"] = HandleOctagonNews

	// Browser Automation Tools (Assimilated from playwright/browser-use/browsermcp/puppeteer/browserbase)
	r.handlers["browser_navigate"] = HandleBrowserNavigate
	r.handlers["browser_screenshot"] = HandleBrowserScreenshot
	r.handlers["browser_get_html"] = HandleBrowserGetHTML
	r.handlers["browser_evaluate"] = HandleBrowserEvaluate
	r.handlers["browser_click"] = HandleBrowserClick
	r.handlers["browser_fill_form"] = HandleBrowserFillForm

	// ChromaDB Vector Store Tools (Assimilated from chroma-mcp)
	r.handlers["chroma_list_collections"] = HandleChromaListCollections
	r.handlers["chroma_create_collection"] = HandleChromaCreateCollection
	r.handlers["chroma_add_documents"] = HandleChromaAddDocuments
	r.handlers["chroma_query"] = HandleChromaQuery
	r.handlers["chroma_delete_collection"] = HandleChromaDeleteCollection
	r.handlers["chroma_get_documents"] = HandleChromaGetCollection

	// Basic Memory Tools (Assimilated from basic-memory)
	r.handlers["basic_memory_write"] = HandleBasicMemoryWrite
	r.handlers["basic_memory_read"] = HandleBasicMemoryRead
	r.handlers["basic_memory_search"] = HandleBasicMemorySearch
	r.handlers["basic_memory_list"] = HandleBasicMemoryList
	r.handlers["basic_memory_delete"] = HandleBasicMemoryDelete
	r.handlers["memory_write"] = HandleBasicMemoryWrite
	r.handlers["memory_read"] = HandleBasicMemoryRead
	r.handlers["memory_search"] = HandleBasicMemorySearch

	// MindsDB ML Database Tools (Assimilated from SSE mindsdb)
	r.handlers["mindsdb_query"] = HandleMindsDBQuery
	r.handlers["mindsdb_list_models"] = HandleMindsDBListModels
	r.handlers["mindsdb_predict"] = HandleMindsDBPredict
}

func (r *Registry) Execute(ctx context.Context, name string, args map[string]interface{}) (ToolResponse, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return ToolResponse{}, fmt.Errorf("tool handler not found for: %s", name)
	}
	return handler(ctx, args)
}

func (r *Registry) HasTool(name string) bool {
	_, ok := r.handlers[name]
	return ok
}
