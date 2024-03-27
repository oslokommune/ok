# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "1.13.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v1.13.0/ok_1.13.0_darwin_arm64.tar.gz"
      sha256 "719d7be58622c88be0a55f7ed3745c194352b3a73d8d60b98a79f0701d008ec9"

      def install
        bin.install "ok"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v1.13.0/ok_1.13.0_darwin_amd64.tar.gz"
      sha256 "33752a0412ba0001ca69d47a131256d32db0f7b10e3b24b03f0c7c4fa1b286af"

      def install
        bin.install "ok"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/oslokommune/ok/releases/download/v1.13.0/ok_1.13.0_linux_arm64.tar.gz"
      sha256 "e2fa43750b28c4d4392639525d55850995746deb5060a55f346d1cf686185763"

      def install
        bin.install "ok"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v1.13.0/ok_1.13.0_linux_amd64.tar.gz"
      sha256 "39f99a9ea494db712cb5c22f14139e799e113a758a060ef33ad5efa68bb02802"

      def install
        bin.install "ok"
      end
    end
  end
end
