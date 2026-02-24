cask "clix" do
  version "1.1.2"

  on_arm do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-arm64"
    sha256 "79205eb1ec83556773ef8a050eec825892b5e134ba77e34749ee28bfaed70761"
  end
  on_intel do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-x64"
    sha256 "90633ff2bb4d7b682577d54c713406476d54f6dbd301286cb31b541f7e2af811"
  end

  name "Clix CLI"
  desc "AI-powered CLI for integrating and managing the Clix SDK in mobile projects"
  homepage "https://github.com/clix-so/clix-cli"

  binary "clix-darwin-#{Hardware::CPU.arch}", target: "clix"

  postflight do
    system_command "/usr/bin/xattr",
      args: ["-d", "com.apple.quarantine", "#{staged_path}/clix-darwin-#{Hardware::CPU.arch}"],
      sudo: false
  end

  caveats <<~EOS
    This cask installs an unsigned binary. If you encounter issues, run:
      xattr -d com.apple.quarantine $(which clix)
  EOS

  zap trash: [
    "~/.config/clix",
    "~/.local/state/clix",
  ]
end
