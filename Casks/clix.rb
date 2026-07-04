cask "clix" do
  version "1.1.3"

  on_arm do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-arm64"
    sha256 "e7703a5e973b8e84418740a72288a8a6d362bca1d2286688c7595ffd6057df8c"
  end
  on_intel do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-x64"
    sha256 "5b582078a512578389a448c1798cdcae2b1c3ff5138d3c5d2482712eb03210d4"
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
