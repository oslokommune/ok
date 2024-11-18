# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "4.7.5"

  depends_on "yq"

  on_macos do
    on_intel do
      url "https://github.com/oslokommune/ok/releases/download/v4.7.5/ok_4.7.5_darwin_amd64.tar.gz"
      sha256 "c1a1e5e6b74995778660553a20e6cb6cc4cf6732ff3e9517b6765e0e66e4ff23"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    on_arm do
      url "https://github.com/oslokommune/ok/releases/download/v4.7.5/ok_4.7.5_darwin_arm64.tar.gz"
      sha256 "e23e473b0d7189e83727eabeebaceb8202c710ef4992f772863baa91d4d998ef"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/oslokommune/ok/releases/download/v4.7.5/ok_4.7.5_linux_amd64.tar.gz"
        sha256 "43949e9ef8aea4f18fbfacb14aff0f335c33ddeb268a2920c2b229950aa3b0b3"

        def install
          bin.install "ok"
          bash_completion.install "completions/ok.bash" => "ok"
          zsh_completion.install "completions/ok.zsh" => "_ok"
          fish_completion.install "completions/ok.fish"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/oslokommune/ok/releases/download/v4.7.5/ok_4.7.5_linux_arm64.tar.gz"
        sha256 "eaa1e5a48447f99ca97b7a06c77f9c9eba3fe76506f992e8380b105b8a231896"

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
