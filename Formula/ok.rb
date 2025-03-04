# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "5.9.0"

  depends_on "fzf"
  depends_on "yq"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v5.9.0/ok_5.9.0_darwin_amd64.tar.gz"
      sha256 "b9e371f9c99c695811a15a024b43b2099f0c99c911757e45aab6780956ab8596"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v5.9.0/ok_5.9.0_darwin_arm64.tar.gz"
      sha256 "d07e1660b7db86b5f785c63012ac0624db4a04040d7a6d0d804c97bbfd078119"

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
        url "https://github.com/oslokommune/ok/releases/download/v5.9.0/ok_5.9.0_linux_amd64.tar.gz"
        sha256 "433d988f74004ceed79ab14cfc4c595804c0ba031f3e34ccbc69c50258148997"

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
        url "https://github.com/oslokommune/ok/releases/download/v5.9.0/ok_5.9.0_linux_arm64.tar.gz"
        sha256 "ce4850e439484fc3d3ab0f1b7956bbe82d2483217e5b59b87d64c9e97983c99f"

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
