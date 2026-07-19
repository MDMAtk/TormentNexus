class Tormentnexus < Formula
  desc "AI Control Plane with Persistent Memory - 26,000+ MCP Tools"
  homepage "https://tormentnexus.site"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/MDMAtk/TormentNexus/releases/download/v1.0.0-b4/tormentnexus-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64_SHA256"
    end

    on_intel do
      url "https://github.com/MDMAtk/TormentNexus/releases/download/v1.0.0-b4/tormentnexus-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64_SHA256"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/MDMAtk/TormentNexus/releases/download/v1.0.0-b4/tormentnexus-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_LINUX_ARM64_SHA256"
    end

    on_intel do
      url "https://github.com/MDMAtk/TormentNexus/releases/download/v1.0.0-b4/tormentnexus-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_LINUX_AMD64_SHA256"
    end
  end

  def install
    bin.install "tormentnexus"
  end

  def caveats
    <<~EOS
      TormentNexus has been installed!

      To start the server:
        tormentnexus serve

      Dashboard will be available at:
        http://localhost:7778

      Configuration directory:
        ~/.tormentnexus
    EOS
  end

  test do
    system "#{bin}/tormentnexus", "--version"
  end
end
