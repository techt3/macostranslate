#!/bin/bash

# Test script for macostranslate service installation
echo "=== macostranslate Service Installation Test ==="
echo ""

# Check if binary exists
if [ ! -f "./macostranslate" ]; then
    echo "❌ Error: macostranslate binary not found. Please build first with 'go build'"
    exit 1
fi

echo "✅ Binary found: ./macostranslate"
echo ""

# Check current service status
echo "📋 Current service status:"
echo "  User LaunchAgent: $([ -f "$HOME/Library/LaunchAgents/pl.com.t3.macostranslate.plist" ] && echo "✅ Installed" || echo "❌ Not installed")"
echo "  System LaunchDaemon: $([ -f "/Library/LaunchDaemons/pl.com.t3.macostranslate.plist" ] && echo "✅ Installed" || echo "❌ Not installed")"
echo "  Keyboard Service: $([ -d "$HOME/Library/Services/macostranslate.workflow" ] && echo "✅ Installed" || echo "❌ Not installed")"
echo ""

echo "🚀 To test the application:"
echo "  1. Run: ./macostranslate"
echo "  2. Look for the 🌐 icon in your menu bar"
echo "  3. Click the icon and select '⚙️ Install Service'"
echo "  4. Choose your preferred installation type:"
echo "     - System Install: Better integration, requires admin password"
echo "     - User Install: No password required, user-level only"
echo ""

echo "🔧 Installation features:"
echo "  • System-level: Installs to /Library/LaunchDaemons (requires sudo)"
echo "  • User-level: Installs to ~/Library/LaunchAgents (no sudo required)"
echo "  • Keyboard shortcuts: Always installed to ~/Library/Services"
echo "  • Automatic fallback from system to user if admin access fails"
echo ""

echo "🧹 To uninstall:"
echo "  • Use the '🗑️ Uninstall Service' menu option in the app"
echo "  • Or manually remove the files listed above"
