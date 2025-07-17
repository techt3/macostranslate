class Macostranslate < Formula
  desc "macOS Translate App - Simple menubar application with autostart for quick access to Google Translate"
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
    require "fileutils"
    
    # Create launch agent plist for autostart
    launch_agent_dir = "#{Dir.home}/Library/LaunchAgents"
    FileUtils.mkdir_p(launch_agent_dir)
    
    plist_content = <<~PLIST
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
        <key>Label</key>
        <string>pl.com.t3.macostranslate</string>
        <key>ProgramArguments</key>
        <array>
          <string>#{bin}/macostranslate</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <false/>
      </dict>
      </plist>
    PLIST
    
    plist_path = "#{launch_agent_dir}/pl.com.t3.macostranslate.plist"
    File.write(plist_path, plist_content)
    
    # Load the launch agent
    system "launchctl", "load", plist_path
    
    puts "✅ macOS Translate installed and configured for autostart"
    puts "🌐 The app will now start automatically when you log in"
    puts "🚀 Starting the app now..."
  end
  
  def uninstall_preflight
    plist_path = "#{Dir.home}/Library/LaunchAgents/pl.com.t3.macostranslate.plist"
    if File.exist?(plist_path)
      system "launchctl", "unload", plist_path
      File.delete(plist_path)
      puts "🗑️ Removed autostart configuration"
    end
  end
  
  test do
    system "#{bin}/macostranslate", "--help"
  end
end
