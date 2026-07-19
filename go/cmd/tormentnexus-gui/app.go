package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx        context.Context
	startupURL string
}

// NewApp creates a new App struct instance
func NewApp(url string) *App {
	return &App{startupURL: url}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetStartupURL returns the deep link URL the app was opened with, if any
func (a *App) GetStartupURL() string {
	return a.startupURL
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, welcome to TormentNexus!", name)
}
