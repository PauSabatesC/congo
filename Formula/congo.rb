# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Congo < Formula
  desc "Easy and unified way to connect from your terminal to AWS EC2 and ECS"
  homepage "https://github.com/PauSabatesC/congo"
  version "0.1.0"
  license "Apache 2.0 license"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/PauSabatesC/congo/releases/download/v0.1.0/congo_0.1.0_Darwin_x86_64.tar.gz"
      sha256 "d502ae1bbaa80e67f6ad1bd3398fcd3f4a2dc0690d22beba7e888daf11a6a5d2"

      def install
        bin.install "congo"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/PauSabatesC/congo/releases/download/v0.1.0/congo_0.1.0_Darwin_arm64.tar.gz"
      sha256 "04ec8f730d9ea9f4aece15c2a56800537cc03c9832b432168323431a39062908"

      def install
        bin.install "congo"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && !Hardware::CPU.is_64_bit?
      url "https://github.com/PauSabatesC/congo/releases/download/v0.1.0/congo_0.1.0_Linux_armv6.tar.gz"
      sha256 "a742d824348acbe156fe0ebfb47d71d44d24ff8da809bfbbf74fda1fe348a011"

      def install
        bin.install "congo"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/PauSabatesC/congo/releases/download/v0.1.0/congo_0.1.0_Linux_x86_64.tar.gz"
      sha256 "cabc6f35e230af4c0917e2c10f364961bbaad0b187bedea0cd637f0c52cf453f"

      def install
        bin.install "congo"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/PauSabatesC/congo/releases/download/v0.1.0/congo_0.1.0_Linux_arm64.tar.gz"
      sha256 "987acdcb9bc81a67b823da930f9f3eb9dfd8e0e78293cc5ce7b13475e7238bf0"

      def install
        bin.install "congo"
      end
    end
  end
end
