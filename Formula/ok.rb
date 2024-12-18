# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "5.3.0"

  depends_on "fzf"
  depends_on "yq"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/oslokommune/ok/releases/download/v5.3.0/ok_5.3.0_darwin_amd64.tar.gz"
      sha256 "a3f070a31f98583978e94d2579ab721a74f194dbab48ca07668db01847dba52e"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/oslokommune/ok/releases/download/v5.3.0/ok_5.3.0_darwin_arm64.tar.gz"
      sha256 "10f14f7cc061ae0fe6779258a334fdf7bd868e0a57925a9aa19a751f53d5969f"

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
        url "https://github.com/oslokommune/ok/releases/download/v5.3.0/ok_5.3.0_linux_amd64.tar.gz"
        sha256 "bddc6508bd2e6de44d1225e75055f4561152f7ca2ec5a2b4958cc505d5bb14ba"

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
        url "https://github.com/oslokommune/ok/releases/download/v5.3.0/ok_5.3.0_linux_arm64.tar.gz"
        sha256 "183c26ba956feb9480a45f73f4b26794daf1b2608999ec98ceae6378d9f7ddcd"

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
