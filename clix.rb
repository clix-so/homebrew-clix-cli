class Clix < Formula
  desc "A CLI tool for integrating and managing the Clix SDK in your mobile projects"
  homepage "https://github.com/clix-so/homebrew-clix-cli"
  url "https://registry.npmjs.org/@clix-so/clix-cli/-/clix-cli-1.0.0.tgz"
  sha256 "6f41bee7a7f50acdee60d93b9918946f55339ae3ff338a112b72db0595740940"
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
