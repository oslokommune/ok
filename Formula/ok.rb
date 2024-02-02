# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "1.9.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v1.9.0/ok_1.9.0_darwin_arm64.tar.gz"
      sha256 "6691efa374406c6d76eb6acb12cbb7e3ef4f23de2f579f3ada60eed30f5f2755"

      def install
        bin.install "ok"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v1.9.0/ok_1.9.0_darwin_amd64.tar.gz"
      sha256 "8ef773ae26dd3ee55fc05235233802bc1dd345c0fda85e43bfdfc9f0f0e2d6dd"

      def install
        bin.install "ok"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/oslokommune/ok/releases/download/v1.9.0/ok_1.9.0_linux_arm64.tar.gz"
      sha256 "c6620d1128945b173624e5de9c0c828858ec8a5c308138611d33b50bc7fe60f0"

      def install
        bin.install "ok"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v1.9.0/ok_1.9.0_linux_amd64.tar.gz"
      sha256 "9dedc6af5151b9fd2c72b4a495b7a9dc2bcca00659d9bef9df512e5ca48d72ba"

      def install
        bin.install "ok"
      end
    end
  end
end
