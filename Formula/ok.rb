# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "1.15.0"

  on_macos do
    on_intel do
      url "https://github.com/oslokommune/ok/releases/download/v1.15.0/ok_1.15.0_darwin_amd64.tar.gz"
      sha256 "dcc8a8578c3269dddf23c7b306ac91439b9f5dfc73741a6052d6bd11ed775916"

      def install
        bin.install "ok"
      end
    end
    on_arm do
      url "https://github.com/oslokommune/ok/releases/download/v1.15.0/ok_1.15.0_darwin_arm64.tar.gz"
      sha256 "88d0676553bb70cb18e95882c3f268fb2fd1ca98139f7ccba301b4616c0cbb91"

      def install
        bin.install "ok"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/oslokommune/ok/releases/download/v1.15.0/ok_1.15.0_linux_amd64.tar.gz"
        sha256 "b6af5d62fb7b7a80f1ce589a8132cc0daccebee866c6601c5fa9ea26042f074b"

        def install
          bin.install "ok"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/oslokommune/ok/releases/download/v1.15.0/ok_1.15.0_linux_arm64.tar.gz"
        sha256 "cfe7a786a775baf9fe43e49dadc48adcda8697bd047429af6ec6caf48a6861e7"

        def install
          bin.install "ok"
        end
      end
    end
  end
end
