class Domloc < Formula
  desc "Local domain routing for developers — zero config HTTPS reverse proxy"
  homepage "https://github.com/wemit/domloc"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-darwin-arm64.tar.gz"
      sha256 "REPLACE_DARWIN_ARM64_SHA256"
    else
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-darwin-amd64.tar.gz"
      sha256 "REPLACE_DARWIN_AMD64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-linux-arm64.tar.gz"
      sha256 "REPLACE_LINUX_ARM64_SHA256"
    else
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-linux-amd64.tar.gz"
      sha256 "REPLACE_LINUX_AMD64_SHA256"
    end
  end

  depends_on "caddy"
  depends_on "dnsmasq"

  def install
    bin.install "domloc"
  end

  def caveats
    <<~EOS
      Run `domloc init` to complete setup.
      This will configure dnsmasq, install a local CA, and start Caddy.
    EOS
  end

  test do
    system "#{bin}/domloc", "--version"
  end
end
