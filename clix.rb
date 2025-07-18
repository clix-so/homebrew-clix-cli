# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Clix < Formula
  desc "A CLI tool for ..."
  homepage "https://github.com/clix-so/homebrew-clix-cli"
  version "0.1.10"

  on_macos do
    on_intel do
      url "https://github.com/clix-so/homebrew-clix-cli/releases/download/v0.1.10/clix_0.1.10_darwin_amd64.tar.gz"
      sha256 "a6a1dbe46a3b7a45258f397b18cbd21f8aa2d00e6c6d6127eff3e11c5477fdfd"

      def install
        bin.install "clix"
      end
    end
    on_arm do
      url "https://github.com/clix-so/homebrew-clix-cli/releases/download/v0.1.10/clix_0.1.10_darwin_arm64.tar.gz"
      sha256 "6ca764899d485563face2c05d273577b6d17ce5c7651e4d029e7e3d8a90d3c44"

      def install
        bin.install "clix"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/clix-so/homebrew-clix-cli/releases/download/v0.1.10/clix_0.1.10_linux_amd64.tar.gz"
        sha256 "73ffdd3cac29033c0f16bf25d7467a2ac55ce850c338e6c3c1415efd53fb0ec8"

        def install
          bin.install "clix"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/clix-so/homebrew-clix-cli/releases/download/v0.1.10/clix_0.1.10_linux_arm64.tar.gz"
        sha256 "f085d2d826eebcd65b1e4748146d29a16f790b21e208ca700daa30f7b962334b"

        def install
          bin.install "clix"
        end
      end
    end
  end

  def caveats
    <<~EOS
      To get started, run:
        clix --help
    EOS
  end

  test do
    system "#{bin}/clix", "--version"
  end
end
