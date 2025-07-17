# macOS Translate App

A simple macOS menubar application that provides quick access to Google Translate through Safari.

## Features

- 🌐 Lives in your macOS menubar
- 🚀 Quick access to Google Translate
- 📝 Text input dialog for instant translation
- 🦊 Opens Google Translate in Safari with a dedicated window
- 🎯 Simple menu controls (Open/Close/Quit)
- 📱 Automatically sized Safari window (1000x700)
- ⚡ Auto-start with system (automatically configured via Homebrew)
- 📊 Status indicator showing current state

## Prerequisites

- macOS (this app is designed specifically for macOS)
- Safari browser (pre-installed on macOS)

## Installation


### Option 1: Homebrew Install (Recommended)

```bash
# Download the latest Homebrew formula
curl -L https://github.com/techt3/macostranslate/releases/latest/download/macostranslate.rb -o /tmp/macostranslate.rb 
brew install /tmp/macostranslate.rb
```

**Note:** The Homebrew installation automatically configures the app to start when you log in. To disable autostart, simply uninstall with `brew uninstall macostranslate`.

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
6. Use "🛑 Quit" to exit the application completely

**Note:** If installed via Homebrew, the app will automatically start when you log in.

## Text Input Feature

The app now includes a convenient text input dialog:

- **📝 Translate Text**: Click to open a text input dialog
- **Instant Translation**: Enter text and it opens Google Translate with your text pre-filled
- **URL Encoding**: Properly handles special characters and spaces
- **Direct Integration**: No need to copy/paste - just type and translate!

## Auto-Start Feature

The app automatically starts when you log into your Mac when installed via Homebrew:

- **Automatic Configuration**: Homebrew automatically configures autostart during installation
- **Launch Agent**: Creates a proper macOS Launch Agent that starts the app when you log in
- **Safe Installation**: Uses macOS standard Launch Agents directory (`~/Library/LaunchAgents/`)
- **Easy Management**: Install/uninstall autostart by installing/uninstalling the app via Homebrew

The auto-start feature:
- Creates a `pl.com.t3.macostranslate.plist` file in your Launch Agents directory
- Uses `launchctl` to manage the service
- Automatically configured during `brew install`
- Automatically removed during `brew uninstall`

To disable autostart, simply uninstall the app:
```bash
brew uninstall macostranslate
```

To re-enable autostart, reinstall the app:
```bash
brew install /tmp/macostranslate.rb
```

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
