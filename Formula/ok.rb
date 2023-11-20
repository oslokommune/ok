# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
require_relative "lib/private_strategy"
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "1.5.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v1.5.0/ok_1.5.0_darwin_arm64.tar.gz", using: GitHubPrivateRepositoryReleaseDownloadStrategy
      sha256 "6638b834283f47b10203cd707f5fe0141b045aa5c74e487a401c90fc9518a4ab"

      def install
        bin.install "ok"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/oslokommune/ok/releases/download/v1.5.0/ok_1.5.0_linux_arm64.tar.gz", using: GitHubPrivateRepositoryReleaseDownloadStrategy
      sha256 "ac2ca89acf9fd4f98d989bedf8f83f4ce5019fc5fd19bd9661a8273f6b5ed8c9"

      def install
        bin.install "ok"
      end
    end
  end
end
