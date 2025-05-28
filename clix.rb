class Clix < Formula
  desc "Clix CLI: Install and manage Clix SDK for iOS and Android"
  homepage "https://github.com/clix-so/homebrew-clix-cli"
  url "https://github.com/clix-so/homebrew-clix-cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_REAL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", "clix"
    bin.install "clix"
  end

  test do
    system "#{bin}/clix", "--version"
  end
end
