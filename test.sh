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
    echo "🔧 Testing completed"
    
    # Kill the app
    kill $APP_PID
    echo "🛑 Application stopped"
else
    echo "❌ Application failed to start or crashed"
    exit 1
fi

echo ""
echo "🎉 All tests passed! The app is ready to use."
echo "🚀 Key improvements:"
echo "   • One-click translate access"
echo "   • Text input dialog for instant translation"
echo "   • Status display in menu"
echo "   • Better window management"
echo "   • Auto-start now handled by Homebrew"
echo ""
echo "💡 Run './macostranslate' to start the application"
echo "💡 Use 'brew install' for automatic autostart configuration"
