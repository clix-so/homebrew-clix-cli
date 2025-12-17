class Clix < Formula
  desc "A CLI tool for integrating and managing the Clix SDK in your mobile projects"
  homepage "https://github.com/clix-so/homebrew-clix-cli"
  url "https://registry.npmjs.org/@clix-so/clix-cli/-/clix-cli-1.0.3.tgz"
  sha256 "b0d664743d927e1f89af458dd9da5b41779bc17a6cf4bf11e716de2d2b817c14"
  license "MIT"

  depends_on "node@18"

  def install
    system "npm", "install", *std_npm_args(prefix: libexec)
    bin.install_symlink libexec/"bin/clix"
  end

  test do
    assert_match "A CLI tool for integrating and managing the Clix SDK", shell_output("#{bin}/clix --help")
  end
end
