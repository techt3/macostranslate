# macOS Translate App

A simple macOS menubar application that provides quick access to Google Translate through Safari.

## Features

- 🌐 Lives in your macOS menubar
- 🚀 Quick access to Google Translate
- 📝 Text input dialog for instant translation
- 🦊 Opens Google Translate in Safari with a dedicated window
- 🎯 Simple menu controls (Open/Close/Quit)
- 📱 Automatically sized Safari window (1000x700)
- ⚡ Auto-start with system (install/remove via menu)
- 📊 Status indicator showing current state

## Prerequisites

- macOS (this app is designed specifically for macOS)
- Safari browser (pre-installed on macOS)

## Installation


### Option 1 Manual Homebrew Install

```bash
# Download the latest Homebrew formula
curl -L https://github.com/techt3/macostranslate/releases/latest/download/macostranslate.rb -o /tmp/macostranslate.rb

# Install with Homebrew
brew install /tmp/macostranslate.rb
```

### Option 2: Direct Download

```bash
# Download the latest release
curl -L https://github.com/techt3/macostranslate/releases/latest/download/macostranslate-*.tar.gz -o macostranslate.tar.gz

# Extract and install
tar -xzf macostranslate.tar.gz
sudo cp macostranslate /usr/local/bin/
```

### Option 3: Build from Source

If you want to build from source (requires Go 1.24 or later):

```bash
# Clone the repository
git clone https://github.com/techt3/macostranslate.git
cd macostranslate

# Install dependencies
go mod tidy

# Build the application
go build -o macostranslate

# Run the application
./macostranslate
```

## Usage

1. After starting the app, you'll see a 🌐 icon in your menubar
2. Click the icon to access the menu
3. Select "🚀 Open Translate" to open a Safari window with Google Translate
4. Select "📝 Translate Text" to enter text directly for translation
5. The Safari window will be automatically sized and positioned
6. Use "⚡ Install Auto-Start" to make the app start automatically with your system
7. Use "🗑️ Remove Auto-Start" to remove the app from system startup
8. Use "🛑 Quit" to exit the application completely

## Text Input Feature

The app now includes a convenient text input dialog:

- **📝 Translate Text**: Click to open a text input dialog
- **Instant Translation**: Enter text and it opens Google Translate with your text pre-filled
- **URL Encoding**: Properly handles special characters and spaces
- **Direct Integration**: No need to copy/paste - just type and translate!

## Auto-Start Feature

The app includes a convenient auto-start feature that allows it to launch automatically when you log into your Mac:

- **Install Auto-Start**: Creates a Launch Agent that starts the app when you log in
- **Remove Auto-Start**: Removes the Launch Agent and stops automatic startup
- **Status Display**: Shows whether auto-start is enabled or disabled
- **Safe Installation**: Uses macOS standard Launch Agents directory (`~/Library/LaunchAgents/`)

The auto-start feature:
- Creates a `pl.com.t3.macostranslate.plist` file in your Launch Agents directory
- Uses `launchctl` to manage the service
- Automatically detects if already installed
- Provides easy removal if you change your mind

## Enhanced Features

- **🚀 One-Click Access**: The first menu item opens Google Translate instantly
- **📊 Status Display**: Shows whether translate is currently open or ready
- **🔄 Smart Window Management**: If already open, brings Safari to front
- **⚡ Streamlined Menu**: Minimal, focused interface with clear actions
- **🎯 Prominent Action**: Main translate action is the first and most visible option

## How it Works

The app uses AppleScript to control Safari, creating a dedicated window for Google Translate. This approach:
- Avoids threading issues with embedded webviews
- Provides a native macOS experience
- Uses the system's default web engine (Safari)
- Maintains proper window management

## Building for Distribution

To build a standalone executable:

```bash
go build -ldflags "-s -w" -o macostranslate
```

Or use the Makefile:

```bash
make build-release
```

## Dependencies

- `github.com/getlantern/systray` - For menubar integration
- Built-in macOS Safari browser
- AppleScript (built into macOS)

## Notes

- The app requires an internet connection to load Google Translate
- Safari must be available (pre-installed on all macOS systems)
- The app is designed to be lightweight and stay in the background
- Window positioning and sizing is handled automatically

## License

This project is open source and available under the MIT License.
