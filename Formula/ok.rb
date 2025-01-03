# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "5.3.2"

  depends_on "fzf"
  depends_on "yq"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v5.3.2/ok_5.3.2_darwin_amd64.tar.gz"
      sha256 "32143b5182acbaefd6d3d988f1db6494afbaf31ee53f31fcf1f04a6bde6affb4"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v5.3.2/ok_5.3.2_darwin_arm64.tar.gz"
      sha256 "34831a472be15be276052b60e8b3fd7d0b80771e4d60a88a8d817d7ca286c678"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/oslokommune/ok/releases/download/v5.3.2/ok_5.3.2_linux_amd64.tar.gz"
        sha256 "92a0e1a25d05157b46a5d02d5f8fff16c25b4af94cb379dbe884406c5f2f85ae"

        def install
          bin.install "ok"
          bash_completion.install "completions/ok.bash" => "ok"
          zsh_completion.install "completions/ok.zsh" => "_ok"
          fish_completion.install "completions/ok.fish"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/oslokommune/ok/releases/download/v5.3.2/ok_5.3.2_linux_arm64.tar.gz"
        sha256 "972e8eb502a9ecac2db810a19341c2facaf0f52f4136f5bad18ec35d61d634fc"

        def install
          bin.install "ok"
          bash_completion.install "completions/ok.bash" => "ok"
          zsh_completion.install "completions/ok.zsh" => "_ok"
          fish_completion.install "completions/ok.fish"
        end
      end
    end
  end
end
