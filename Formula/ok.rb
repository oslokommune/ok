# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ok < Formula
  desc "A CLI called ok"
  homepage "https://github.com/oslokommune/ok"
  version "4.3.0"

  on_macos do
    on_intel do
      url "https://github.com/oslokommune/ok/releases/download/v4.3.0/ok_4.3.0_darwin_amd64.tar.gz"
      sha256 "861808f19c2ad11ab9c1694aa57d26264b52f0529bbf2b77ff7e726c0cdede74"

      def install
        bin.install "ok"
        bash_completion.install "completions/ok.bash" => "ok"
        zsh_completion.install "completions/ok.zsh" => "_ok"
        fish_completion.install "completions/ok.fish"
      end
    end
    on_arm do
      url "https://github.com/oslokommune/ok/releases/download/v4.3.0/ok_4.3.0_darwin_arm64.tar.gz"
      sha256 "facf7e83c8f0bdccdc3b8f715f0da92d7090472196b18e4ad0f5407e5da86bd9"

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
        url "https://github.com/oslokommune/ok/releases/download/v4.3.0/ok_4.3.0_linux_amd64.tar.gz"
        sha256 "21e5937b02c2d1224758c80225d6153b62342a4c83e9f491b98dd9331a714d5f"

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
        url "https://github.com/oslokommune/ok/releases/download/v4.3.0/ok_4.3.0_linux_arm64.tar.gz"
        sha256 "c71546304423ca1b40bd8ca8df58c965f270c4fc25a92d7dd5b3a9da8b5c44e2"

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
