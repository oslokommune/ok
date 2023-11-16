# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
require_relative "lib/private_strategy"
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "1.4.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v1.4.0/ok_1.4.0_darwin_arm64.tar.gz", using: GitHubPrivateRepositoryReleaseDownloadStrategy
      sha256 "48cd50fdde96c32aab5dec170aa8e0e05ea38f581b2783dde6dc0c2e6dbb5d7a"

      def install
        bin.install "ok"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/oslokommune/ok/releases/download/v1.4.0/ok_1.4.0_linux_arm64.tar.gz", using: GitHubPrivateRepositoryReleaseDownloadStrategy
      sha256 "36745758b0f0304d1bc012753db95ab3e3707ca5d8719a5c3c03c5f6da529c8e"

      def install
        bin.install "ok"
      end
    end
  end
end
