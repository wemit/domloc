class Domloc < Formula
  desc "Local domain routing for developers — zero config HTTPS reverse proxy"
  homepage "https://github.com/wemit/domloc"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-darwin-arm64.tar.gz"
      sha256 "4d1695a5ae6c94a648042085c8264d5146915028621cf9db7c174709d856f780"
    else
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-darwin-amd64.tar.gz"
      sha256 "2f5c0b001cf5de0ff7372eed5860ee87383c8f61b0593a38371ea4525915389e"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-linux-arm64.tar.gz"
      sha256 "1a86c0a6f64767247436fcac108b9b078a726999f5fcd0ecfd05ee0fb21b58ad"
    else
      url "https://github.com/wemit/domloc/releases/download/v#{version}/domloc-linux-amd64.tar.gz"
      sha256 "d039dd16ad2f97e728ceef5fc239e1ba1bbb485b5a09b4067f6fc6786f66a983"
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
