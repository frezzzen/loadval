package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type WindowState struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximised bool `json:"maximised"`
}

func NewApp() *App {
	return &App{}
}

func Startup(a *App, ctx context.Context) {
	a.ctx = ctx
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}

func (a *App) MaximiseApp() {
	runtime.WindowToggleMaximise(a.ctx)
}

func (a *App) MinimiseApp() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) ReloadApp() {
	runtime.WindowReload(a.ctx)
}

func (a *App) LogDebug(message string) {
	fmt.Printf("[DEBUG] %s\n", message)
}

func (a *App) LogInfo(message string) {
	fmt.Printf("[INFO] %s\n", message)
}

func (a *App) LogWarning(message string) {
	fmt.Printf("[WARNING] %s\n", message)
}

func (a *App) LogError(message string) {
	fmt.Printf("[ERROR] %s\n", message)
}

func (a *App) GetLogFilePath() string {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = os.Getenv("APPDATA")
	}
	if appData == "" {
		return ""
	}
	return filepath.Join(appData, "LOADVAL", "logs", "loadval.log")
}

func (a *App) OpenLogFile() error {
	logPath := a.GetLogFilePath()
	if logPath == "" {
		return fmt.Errorf("could not determine log file path")
	}

	logDir := filepath.Dir(logPath)
	runtime.BrowserOpenURL(a.ctx, "file:///"+filepath.ToSlash(logDir))
	return nil
}

func (a *App) SaveWindowState() error {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = os.Getenv("APPDATA")
	}
	if appData == "" {
		return fmt.Errorf("LOCALAPPDATA/APPDATA environment variable is not set")
	}

	appDir := filepath.Join(appData, "LOADVAL")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	width, height := runtime.WindowGetSize(a.ctx)
	x, y := runtime.WindowGetPosition(a.ctx)

	state := WindowState{
		Width:     width,
		Height:    height,
		X:         x,
		Y:         y,
		Maximised: false,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(appDir, "window-state.json"), data, 0644)
}

func (a *App) LoadWindowState() (*WindowState, error) {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = os.Getenv("APPDATA")
	}
	if appData == "" {
		return nil, fmt.Errorf("LOCALAPPDATA/APPDATA environment variable is not set")
	}

	filePath := filepath.Join(appData, "LOADVAL", "window-state.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var state WindowState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func (a *App) GetWindowSize() (int, int) {
	width, height := runtime.WindowGetSize(a.ctx)
	return width, height
}

func (a *App) SetWindowSize(width, height int) {
	runtime.WindowSetSize(a.ctx, width, height)
}
