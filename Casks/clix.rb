cask "clix" do
  version "1.1.0"

  on_arm do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-arm64"
    sha256 "57c4bd227b40e4f66aea5f8e47a8d429cbfe56e8f76cf2a5008be05a2b6aba24"
  end
  on_intel do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-x64"
    sha256 "c4397e608c7f98d2e542a1ad40546caf282d28c2951aa32cac145fc6cdf5f856"
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
