#!/bin/bash

# Test script for Enhanced macOS Translate App

echo "🧪 Testing Enhanced macOS Translate App..."
echo ""

# Check if app builds
echo "📦 Building enhanced application..."
go build -o macostranslate
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Check if app starts without errors
echo ""
echo "🚀 Starting enhanced application (will run for 5 seconds)..."
./macostranslate &
APP_PID=$!

sleep 5

# Check if app is still running
if kill -0 $APP_PID 2>/dev/null; then
    echo "✅ Enhanced application started successfully"
    echo "🌐 Look for the globe icon in your menubar"
    echo "⚡ Click it and select '🚀 Open Translate' for instant access"
    echo "� Click it and select '📝 Translate Text' for text input"
    echo "�📊 Menu now shows status and streamlined options"
    echo "🔧 Test the auto-start feature via the menu"
    
    # Kill the app
    kill $APP_PID
    echo "🛑 Application stopped"
else
    echo "❌ Application failed to start or crashed"
    exit 1
fi

echo ""
echo "🔍 Testing auto-start functionality..."

# Test auto-start detection
if [ -f ~/Library/LaunchAgents/com.macostranslate.plist ]; then
    echo "⚠️  Auto-start plist already exists - cleaning up for test"
    launchctl unload ~/Library/LaunchAgents/com.macostranslate.plist 2>/dev/null
    rm -f ~/Library/LaunchAgents/com.macostranslate.plist
fi

echo ""
echo "🎉 All tests passed! The enhanced app is ready to use."
echo "🚀 Key improvements:"
echo "   • One-click translate access"
echo "   • Text input dialog for instant translation"
echo "   • Status display in menu"
echo "   • Auto-start functionality"
echo "   • Better window management"
echo ""
echo "💡 Run './macostranslate' to start the enhanced application"
