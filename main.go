package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/getlantern/systray"
)

var isWebViewOpen bool
var webviewProcess *exec.Cmd

func main() {
	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printHelp()
		os.Exit(0)
	}

	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		log.Fatal("This application is designed for macOS only")
	}

	// Run systray on main thread
	systray.Run(onReady, onExit)
}

func printHelp() {
	fmt.Println("macOS Translate App")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("A simple macOS menubar application that provides quick access to Google Translate through Safari.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  macostranslate [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -h, --help    Show this help message and exit")
	fmt.Println()
	fmt.Println("FEATURES:")
	fmt.Println("  🌐 Lives in your macOS menubar")
	fmt.Println("  🚀 Quick access to Google Translate")
	fmt.Println("  📝 Text input dialog for instant translation")
	fmt.Println("  🦊 Opens Google Translate in Safari with a dedicated window")
	fmt.Println("  🎯 Simple menu controls (Open/Close/Quit)")
	fmt.Println("  📱 Automatically sized Safari window (1000x700)")
	fmt.Println("  ⚡ Auto-start with system (install/remove via menu)")
	fmt.Println("  📊 Status indicator showing current state")
	fmt.Println()
	fmt.Println("USAGE INSTRUCTIONS:")
	fmt.Println("  1. After starting the app, you'll see a 🌐 icon in your menubar")
	fmt.Println("  2. Click the icon to access the menu")
	fmt.Println("  3. Select '🚀 Open Translate' to open Google Translate in Safari")
	fmt.Println("  4. Select '📝 Translate Text' to enter text directly for translation")
	fmt.Println("  5. Use '⚡ Install Auto-Start' to make the app start automatically")
	fmt.Println("  6. Use '🗑️ Remove Auto-Start' to remove from system startup")
	fmt.Println("  7. Use '🛑 Quit' to exit the application completely")
	fmt.Println()
	fmt.Println("REQUIREMENTS:")
	fmt.Println("  - macOS (this app is designed specifically for macOS)")
	fmt.Println("  - Safari browser (pre-installed on macOS)")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/techt3/macostranslate")
}

func onReady() {
	systray.SetTitle("🌐")
	systray.SetTooltip("Google Translate (Click for instant translate)")

	// Create a minimal menu with the main action first
	openItem := systray.AddMenuItem("🚀 Open Translate", "Open Google Translate instantly")

	// Add text input option
	textInputItem := systray.AddMenuItem("📝 Translate Text", "Enter text to translate")

	systray.AddSeparator()

	systray.AddSeparator()

	// Auto-start menu items
	autoStartItem := systray.AddMenuItem("⚡ Install Auto-Start", "Install app to start with system")
	removeAutoStartItem := systray.AddMenuItem("🗑️ Remove Auto-Start", "Remove app from system startup")

	// Check if auto-start is already installed
	if isAutoStartInstalled() {
		autoStartItem.Disable()
		removeAutoStartItem.Enable()
	} else {
		autoStartItem.Enable()
		removeAutoStartItem.Disable()
	}

	systray.AddSeparator()

	quitItem := systray.AddMenuItem("🛑 Quit", "Quit the application")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-openItem.ClickedCh:
				if !isWebViewOpen {
					go openWebView()
				}
			case <-textInputItem.ClickedCh:
				go handleTextInput()

			case <-autoStartItem.ClickedCh:
				if installAutoStart() {
					autoStartItem.Disable()
					removeAutoStartItem.Enable()
				}
			case <-removeAutoStartItem.ClickedCh:
				if removeAutoStart() {
					autoStartItem.Enable()
					removeAutoStartItem.Disable()
				}
			case <-quitItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// Handle system signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		systray.Quit()
	}()
}

func openWebView() {
	if isWebViewOpen {
		// If already open, bring Safari to front
		script := `
tell application "Safari"
    activate
    if (count of windows) > 0 then
        if URL of current tab of front window contains "translate.google.com" then
            -- Window is already showing translate, just activate it
            set index of front window to 1
        else
            -- Navigate current tab to translate
            set URL of current tab of front window to "https://translate.google.com"
            set bounds of front window to {100, 100, 1100, 800}
        end if
    end if
end tell
`
		cmd := exec.Command("osascript", "-e", script)
		_ = cmd.Run()
		return
	}

	isWebViewOpen = true
	defer func() {
		isWebViewOpen = false
	}()

	// Open Safari with Google Translate
	script := `
tell application "Safari"
    activate
    if (count of windows) = 0 then
        make new document with properties {URL:"https://translate.google.com"}
    else
        set URL of current tab of front window to "https://translate.google.com"
    end if
    set bounds of front window to {100, 100, 1100, 800}
end tell
`

	cmd := exec.Command("osascript", "-e", script)
	webviewProcess = cmd

	err := cmd.Run()
	if err != nil {
		log.Printf("Error opening Safari window: %v", err)
		// Fallback to default browser
		cmd = exec.Command("open", "https://translate.google.com")
		webviewProcess = cmd
		err = cmd.Run()
		if err != nil {
			log.Printf("Error opening default browser: %v", err)
		}
	}
}

func closeWebView() {
	if webviewProcess != nil && isWebViewOpen {
		// Try to close the current Safari tab if it's showing Google Translate
		script := `
tell application "Safari"
    if URL of current tab of front window contains "translate.google.com" then
        close current tab of front window
    end if
end tell
`
		cmd := exec.Command("osascript", "-e", script)
		_ = cmd.Run() // Ignore error as window might not exist
		webviewProcess = nil
	}
}

func onExit() {
	// Cleanup
	if webviewProcess != nil && webviewProcess.Process != nil {
		_ = webviewProcess.Process.Kill() // Ignore error as process might already be dead
	}
	os.Exit(0)
}

// Auto-start functionality

func isAutoStartInstalled() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "pl.com.t3.macostranslate.plist")
	_, err = os.Stat(plistPath)
	return err == nil
}

func installAutoStart() bool {
	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("Error getting executable path: %v", err)
		return false
	}

	// Get absolute path
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
		return false
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home directory: %v", err)
		return false
	}

	// Create LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		log.Printf("Error creating LaunchAgents directory: %v", err)
		return false
	}

	// Create the plist content
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>pl.com.t3.macostranslate</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
</dict>
</plist>`, absPath)

	// Write the plist file
	plistPath := filepath.Join(launchAgentsDir, "pl.com.t3.macostranslate.plist")
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		log.Printf("Error writing plist file: %v", err)
		return false
	}

	// Load the launch agent
	cmd := exec.Command("launchctl", "load", plistPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Error loading launch agent: %v", err)
		// Try to remove the file if loading failed
		_ = os.Remove(plistPath)
		return false
	}

	log.Println("Auto-start installed successfully")
	return true
}

func removeAutoStart() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home directory: %v", err)
		return false
	}

	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "pl.com.t3.macostranslate.plist")

	// Unload the launch agent
	cmd := exec.Command("launchctl", "unload", plistPath)
	if err := cmd.Run(); err != nil {
		// If unload fails, it might not be loaded, but we should still try to remove the file
		log.Printf("Warning: Error unloading launch agent: %v", err)
	}

	// Remove the plist file
	if err := os.Remove(plistPath); err != nil {
		log.Printf("Error removing plist file: %v", err)
		return false
	}

	log.Println("Auto-start removed successfully")
	return true
}

func handleTextInput() {
	// Use AppleScript to show a text input dialog
	script := `
set userInput to text returned of (display dialog "Enter text to translate:" default answer "" with title "Google Translate")
return userInput
`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error getting text input: %v", err)
		return
	}

	// Clean up the input text
	inputText := strings.TrimSpace(string(output))
	if inputText == "" {
		return
	}

	// Open Google Translate with the input text
	openWebViewWithText(inputText)
}

func openWebViewWithText(text string) {
	// URL encode the text
	encodedText := url.QueryEscape(text)
	translateURL := fmt.Sprintf("https://translate.google.com/?text=%s", encodedText)

	isWebViewOpen = true
	defer func() {
		isWebViewOpen = false
	}()

	// Open Safari with Google Translate and the text
	script := fmt.Sprintf(`
tell application "Safari"
    activate
    if (count of windows) = 0 then
        make new document with properties {URL:"%s"}
    else
        set URL of current tab of front window to "%s"
    end if
    set bounds of front window to {100, 100, 1100, 800}
end tell
`, translateURL, translateURL)

	cmd := exec.Command("osascript", "-e", script)
	webviewProcess = cmd

	err := cmd.Run()
	if err != nil {
		log.Printf("Error opening Safari window with text: %v", err)
		// Fallback to default browser
		cmd = exec.Command("open", translateURL)
		webviewProcess = cmd
		err = cmd.Run()
		if err != nil {
			log.Printf("Error opening default browser with text: %v", err)
		}
	}
}
