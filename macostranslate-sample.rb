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
    require "fileutils"
    
    puts "🚀 Setting up macostranslate..."
    
    # Create launch agent plist for autostart
    begin
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
      puts "✅ Created autostart configuration"
      
      # Try to load the launch agent, but don't fail if it doesn't work
      if system("launchctl", "load", plist_path, out: "/dev/null", err: "/dev/null")
        puts "✅ Loaded autostart agent"
      else
        puts "⚠️  Autostart agent created but not loaded (you may need to log out/in)"
      end
    rescue => e
      puts "⚠️  Could not set up autostart: #{e.message}"
    end
    
    # Create a macOS Service for global keyboard shortcut
    begin
      services_dir = "#{Dir.home}/Library/Services"
      FileUtils.mkdir_p(services_dir)
      
      service_content = <<~SERVICE
        <?xml version="1.0" encoding="UTF-8"?>
        <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
        <plist version="1.0">
        <dict>
          <key>NSMenuItem</key>
          <dict>
            <key>default</key>
            <string>Open macostranslate</string>
          </dict>
          <key>NSMessage</key>
          <string>runWorkflowAsService</string>
          <key>NSPortName</key>
          <string>NSPerformService</string>
          <key>NSRequiredContext</key>
          <array/>
          <key>NSReturnTypes</key>
          <array/>
          <key>NSSendTypes</key>
          <array/>
        </dict>
        </plist>
      SERVICE
      
      # Create the service workflow
      service_path = "#{services_dir}/macostranslate.workflow"
      FileUtils.mkdir_p(service_path)
      FileUtils.mkdir_p("#{service_path}/Contents")
      
      # Create Info.plist for the service
      File.write("#{service_path}/Contents/Info.plist", service_content)
      
      # Create the document.wflow file
      workflow_content = <<~WORKFLOW
        <?xml version="1.0" encoding="UTF-8"?>
        <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
        <plist version="1.0">
        <dict>
          <key>AMApplicationBuild</key>
          <string>521</string>
          <key>AMApplicationVersion</key>
          <string>2.10</string>
          <key>AMDocumentVersion</key>
          <string>2</string>
          <key>actions</key>
          <array>
            <dict>
              <key>action</key>
              <dict>
                <key>AMAccepts</key>
                <dict>
                  <key>Container</key>
                  <string>List</string>
                  <key>Optional</key>
                  <true/>
                  <key>Types</key>
                  <array>
                    <string>com.apple.applescript.object</string>
                  </array>
                </dict>
                <key>AMActionVersion</key>
                <string>1.0.2</string>
                <key>AMApplication</key>
                <array>
                  <string>Automator</string>
                </array>
                <key>AMParameterProperties</key>
                <dict>
                  <key>source</key>
                  <dict/>
                </dict>
                <key>AMProvides</key>
                <dict>
                  <key>Container</key>
                  <string>List</string>
                  <key>Types</key>
                  <array>
                    <string>com.apple.applescript.object</string>
                  </array>
                </dict>
                <key>ActionBundlePath</key>
                <string>/System/Library/Automator/Run AppleScript.action</string>
                <key>ActionName</key>
                <string>Run AppleScript</string>
                <key>ActionParameters</key>
                <dict>
                  <key>source</key>
                  <string>do shell script "#{bin}/macostranslate &"</string>
                </dict>
              </dict>
            </dict>
          </array>
          <key>connectors</key>
          <dict/>
          <key>workflowMetaData</key>
          <dict>
            <key>serviceInputTypeIdentifier</key>
            <string>com.apple.Automator.nothing</string>
            <key>serviceOutputTypeIdentifier</key>
            <string>com.apple.Automator.nothing</string>
            <key>serviceApplicationBundleIdentifier</key>
            <string>com.apple.finder</string>
            <key>workflowTypeIdentifier</key>
            <string>com.apple.Automator.servicesMenu</string>
          </dict>
        </dict>
        </plist>
      WORKFLOW
      
      File.write("#{service_path}/Contents/document.wflow", workflow_content)
      puts "✅ Created keyboard shortcut service"
    rescue => e
      puts "⚠️  Could not create keyboard shortcut service: #{e.message}"
    end
    
    puts ""
    puts "✅ macostranslate installed and configured!"
    puts "🌐 The app will start automatically when you log in"
    puts "⌨️ Global keyboard shortcut available:"
    puts "   1. Go to System Preferences > Keyboard > Shortcuts"
    puts "   2. Select 'Services' in the left panel"
    puts "   3. Find 'Open macostranslate' service"
    puts "   4. Assign your preferred shortcut (recommended: Cmd+Shift+T)"
    puts "🚀 Starting the app now..."
    
    # Start the app
    system "#{bin}/macostranslate &"
  end
  
  def uninstall_preflight
    begin
      plist_path = "#{Dir.home}/Library/LaunchAgents/pl.com.t3.macostranslate.plist"
      if File.exist?(plist_path)
        system("launchctl", "unload", plist_path, out: "/dev/null", err: "/dev/null")
        File.delete(plist_path)
        puts "🗑️ Removed autostart configuration"
      end
      
      # Remove the macOS Service
      service_path = "#{Dir.home}/Library/Services/macostranslate.workflow"
      if File.exist?(service_path)
        require "fileutils"
        FileUtils.rm_rf(service_path)
        puts "🗑️ Removed keyboard shortcut service"
      end
    rescue => e
      puts "Warning during uninstall: #{e.message}"
    end
  end
  
  test do
    system "#{bin}/macostranslate", "--help"
  end
end
