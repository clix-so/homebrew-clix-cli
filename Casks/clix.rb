cask "clix" do
  version "1.0.0"

  on_arm do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-arm64"
    sha256 ""
  end
  on_intel do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-x64"
    sha256 ""
  end

  name "Clix CLI"
  desc "AI-powered CLI for integrating and managing the Clix SDK in mobile projects"
  homepage "https://github.com/clix-so/clix-cli"

  binary "clix-darwin-#{Hardware::CPU.arch}", target: "clix"

  zap trash: [
    "~/.config/clix",
    "~/.local/state/clix",
  ]
end
