# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "5.13.0"

  depends_on "fzf"
  depends_on "yq"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v5.13.0/ok_5.13.0_darwin_amd64.tar.gz"
      sha256 "2210db23f707333c1a21edc6b960387f938a74ba8a86677c2be689c2293c108e"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v5.13.0/ok_5.13.0_darwin_arm64.tar.gz"
      sha256 "69eb0e74048ad832172a11b048f1aee209301450c77b8ec84464b5e29d847057"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel? and Hardware::CPU.is_64_bit?
      url "https://github.com/oslokommune/ok/releases/download/v5.13.0/ok_5.13.0_linux_amd64.tar.gz"
      sha256 "9e0846b72160190af1799fef6831ff9af3b1a91880bce49df2b9c5835010f2c8"
      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    if Hardware::CPU.arm? and Hardware::CPU.is_64_bit?
      url "https://github.com/oslokommune/ok/releases/download/v5.13.0/ok_5.13.0_linux_arm64.tar.gz"
      sha256 "ace8af27378fc4a05a13aba4426b4e6b57e14c6eab04e4489ca423c78b5cf02d"
      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
  end
end
