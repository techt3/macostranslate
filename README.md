# macostranslate App

A simple macOS menubar application that provides quick access to Google Translate through Safari.

## Features

- 🌐 Lives in your macOS menubar
- 🚀 Quick access to Google Translate
- 📝 Text input dialog for instant translation
- 🦊 Opens Google Translate in Safari with a dedicated window
- 🎯 Simple menu controls (Open/Close/Quit)
- 📱 Automatically sized Safari window (1000x700)
- ⚙️ **Service Management** - Install/uninstall autostart and keyboard shortcuts directly from the app
- ⚡ Auto-start with system (optional, managed via app menu)
- ⌨️ Global keyboard shortcut (Cmd+Shift+T) - optional, managed via app menu
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

### Option 2: Build from Source

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
5. Use "⚙️ Install Service" to set up autostart and keyboard shortcuts
6. The Safari window will be automatically sized and positioned
7. Use "🛑 Quit" to exit the application completely

## Service Management

The app now includes built-in service management:

- **⚙️ Install Service**: Sets up autostart and keyboard shortcut services
- **🗑️ Uninstall Service**: Removes autostart and keyboard shortcut services
- **Dynamic Menu**: The menu updates based on current service installation status
- **User-Friendly**: No need to manually configure system settings

### What the Service Installation Does:

1. **Autostart Service**: Creates a Launch Agent that starts the app when you log in
2. **Keyboard Shortcut Service**: Creates a macOS Service for global keyboard shortcuts
3. **Proper Integration**: Uses macOS standard directories and methods
4. **Safe Installation**: All files are created in user directories, no system modifications

## Text Input Feature

The app includes a convenient text input dialog:

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


## Keyboard Shortcut Feature

The Homebrew installation automatically sets up a macOS service that allows you to assign a global keyboard shortcut:

### Setting Up the Keyboard Shortcut

1. After installing via Homebrew, go to **System Preferences** → **Keyboard** → **Shortcuts**
2. Select **Services** in the left panel
3. Scroll down to find **"Open macostranslate"** service
4. Click on it and assign your preferred shortcut (we recommend **Cmd+Shift+T**)
5. The shortcut will now work system-wide to instantly open Google Translate

### Benefits

- **⌨️ Global Access**: Works from any application, anywhere in macOS
- **🚀 Instant Launch**: No need to find the menubar icon
- **🎯 One-Key Translation**: Direct access to translate functionality
- **🔧 Customizable**: Assign any shortcut you prefer via System Preferences

## License

This project is open source and available under the MIT License.
