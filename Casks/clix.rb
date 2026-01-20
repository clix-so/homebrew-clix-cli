cask "clix" do
  version "1.1.1"

  on_arm do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-arm64"
    sha256 "b4ceea0a8701255368d4e87b513b166984c11eab196580948382dcf43c0002b5"
  end
  on_intel do
    url "https://github.com/clix-so/clix-cli/releases/download/v#{version}/clix-darwin-x64"
    sha256 "24849313cf8a3ee6a0b550993787ea41d812b887e47c47396538b24ceb5740da"
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
