class Macostranslate < Formula
  desc "macostranslate App - Simple menubar application with autostart for quick access to Google Translate"
  homepage "https://github.com/techt3/macostranslate"
  version "PLACEHOLDER_VERSION"
  
  depends_on :macos
  
  if Hardware::CPU.arm?
    url "PLACEHOLDER_URL"
    sha256 "PLACEHOLDER_SHA256"
  elsif Hardware::CPU.intel?
    url "PLACEHOLDER_URL"
    sha256 "PLACEHOLDER_SHA256"
  end
  
  def install
    bin.install "macostranslate"
  end
  
  def post_install
    puts "✅ macostranslate installed successfully!"
    puts "🌐 The app is now available in your PATH"
    puts "🚀 Run 'macostranslate' to start the app"
    puts "⚙️ Use the app menu to install autostart and keyboard shortcut services"
    puts ""
    puts "To set up services:"
    puts "1. Start the app: macostranslate"
    puts "2. Click the 🌐 icon in your menubar"
    puts "3. Select '⚙️ Install Service' to enable autostart and keyboard shortcuts"
    puts "4. Follow the instructions to set up global keyboard shortcut (Cmd+Shift+T)"
  end
  
  def uninstall_preflight
    puts "🗑️ Cleaning up macostranslate services..."
    puts "Note: Services are now managed by the app itself."
    puts "If you have installed services via the app, they will remain."
    puts "To remove them, run the app and use '🗑️ Uninstall Service' before uninstalling."
  end
  
  test do
    system "#{bin}/macostranslate", "--help"
  end
end
