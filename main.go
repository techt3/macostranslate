package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
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
	fmt.Println("macostranslate App")
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
	fmt.Println("  ⚡ Auto-start with system (system-level installation with admin privileges)")
	fmt.Println("  ⚙️ Service installation choice (system-level or user-level)")
	fmt.Println("  🔒 Enhanced security and reliability with sudo integration")
	fmt.Println("  ⌨️ Global keyboard shortcut (Cmd+Shift+T)")
	fmt.Println("  📊 Status indicator showing current state")
	fmt.Println()
	fmt.Println("USAGE INSTRUCTIONS:")
	fmt.Println("  1. After starting the app, you'll see a 🌐 icon in your menubar")
	fmt.Println("  2. Click the icon to access the menu")
	fmt.Println("  3. Select '🚀 Open Translate' to open Google Translate in Safari")
	fmt.Println("  4. Select '📝 Translate Text' to enter text directly for translation")
	fmt.Println("  5. Select '⚙️ Install Service' to set up autostart and keyboard shortcuts")
	fmt.Println("      - Choose 'System Install' for better integration (requires admin password)")
	fmt.Println("      - Choose 'User Install' for no password requirement")
	fmt.Println("  6. Use '🛑 Quit' to exit the application completely")
	fmt.Println("  7. To uninstall services, use '🗑️ Uninstall Service' menu option")
	fmt.Println("  7. Set up global hotkey Cmd+Shift+T in System Preferences > Keyboard > Shortcuts")
	fmt.Println()
	fmt.Println("REQUIREMENTS:")
	fmt.Println("  - macOS (this app is designed specifically for macOS)")
	fmt.Println("  - Safari browser (pre-installed on macOS)")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/techt3/macostranslate")
}

func onReady() {
	systray.SetTitle("🌐")
	systray.SetTooltip("Google Translate (Click for instant translate or use Cmd+Shift+T)")

	// Create a minimal menu with the main action first
	openItem := systray.AddMenuItem("🚀 Open Translate", "Open Google Translate instantly (Cmd+Shift+T)")

	// Add text input option
	textInputItem := systray.AddMenuItem("📝 Translate Text", "Enter text to translate")

	systray.AddSeparator()

	// Service management items
	var serviceItem *systray.MenuItem
	if isServiceInstalled() {
		serviceItem = systray.AddMenuItem("🗑️ Uninstall Service", "Remove autostart and keyboard shortcut services")
	} else {
		serviceItem = systray.AddMenuItem("⚙️ Install Service", "Install autostart and keyboard shortcut services")
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
			case <-serviceItem.ClickedCh:
				go handleServiceToggle(serviceItem)
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

func onExit() {
	// Cleanup
	if webviewProcess != nil && webviewProcess.Process != nil {
		_ = webviewProcess.Process.Kill() // Ignore error as process might already be dead
	}
	os.Exit(0)
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

func handleServiceToggle(serviceItem *systray.MenuItem) {
	if isServiceInstalled() {
		// Uninstall service
		if err := uninstallService(); err != nil {
			showErrorDialog("Failed to uninstall service", err.Error())
			return
		}
		showInfoDialog("Service Uninstalled", "Autostart and keyboard shortcut services have been removed successfully.")
		serviceItem.SetTitle("⚙️ Install Service")

	} else {
		// Install service with elevated privileges for better integration
		showInstallationChoiceDialog(serviceItem)
	}
}

func showInstallationChoiceDialog(serviceItem *systray.MenuItem) {
	// Show dialog asking for installation type
	script := `
set userChoice to button returned of (display dialog "Choose installation type:

🔒 System Installation (Recommended)
- Requires administrator password
- Better system integration
- Runs at system startup
- More reliable autostart

👤 User Installation
- No password required
- User-level only
- May require manual setup" buttons {"Cancel", "User Install", "System Install"} default button "System Install" with title "macostranslate Service Installation")

return userChoice
`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		// User cancelled or error occurred
		return
	}

	choice := strings.TrimSpace(string(output))

	switch choice {
	case "System Install":
		installServiceWithFeedback(serviceItem, true)
	case "User Install":
		installServiceWithFeedback(serviceItem, false)
	default:
		// User cancelled
		return
	}
}

func installServiceWithFeedback(serviceItem *systray.MenuItem, useSystem bool) {
	var err error

	if useSystem {
		// Use system installation
		err = installServiceWithChoice(true)
	} else {
		// Use user installation
		err = installServiceWithChoice(false)
	}

	if err != nil {
		showErrorDialog("Failed to install service", err.Error())
		return
	}

	var message string
	if useSystem {
		message = "System-level autostart and keyboard shortcut services have been installed successfully.\n\nThe application will now start automatically when your Mac boots up.\n\nTo set up the global keyboard shortcut:\n1. Go to System Preferences > Keyboard > Shortcuts\n2. Select 'Services' in the left panel\n3. Find 'Open macostranslate' service\n4. Assign your preferred shortcut (recommended: Cmd+Shift+T)"
	} else {
		message = "User-level autostart and keyboard shortcut services have been installed successfully.\n\nTo set up the global keyboard shortcut:\n1. Go to System Preferences > Keyboard > Shortcuts\n2. Select 'Services' in the left panel\n3. Find 'Open macostranslate' service\n4. Assign your preferred shortcut (recommended: Cmd+Shift+T)"
	}

	showInfoDialog("Service Installed", message)
	serviceItem.SetTitle("🗑️ Uninstall Service")
}

func showErrorDialog(title, message string) {
	script := fmt.Sprintf(`display dialog "%s" with title "%s" buttons {"OK"} default button "OK" with icon stop`,
		strings.ReplaceAll(message, "\"", "\\\""),
		strings.ReplaceAll(title, "\"", "\\\""))
	cmd := exec.Command("osascript", "-e", script)
	_ = cmd.Run()
}

func showInfoDialog(title, message string) {
	script := fmt.Sprintf(`display dialog "%s" with title "%s" buttons {"OK"} default button "OK"`,
		strings.ReplaceAll(message, "\"", "\\\""),
		strings.ReplaceAll(title, "\"", "\\\""))
	cmd := exec.Command("osascript", "-e", script)
	_ = cmd.Run()
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
