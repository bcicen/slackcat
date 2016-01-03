class Slackcat < Formula
  desc "Simple command-line Utility to post snippets to Slack."
  homepage "https://github.com/vektorlab/slackcat"
  url "https://github.com/vektorlab/slackcat/archive/master.tar.gz"
  version "0.4"
  sha256 "02c23e8b7bf6a45f85c6911795324a208e1ed9b30b5e9bfb851a74217c68fbbc"

  depends_on "go"

  def install
    unless ENV['GOPATH']
      ENV['GOPATH'] = "/tmp"
    end
    system "make"
    cp "build/slackcat-0.4-darwin-amd64", "#{HOMEBREW_PREFIX}/bin/"
    mv "#{HOMEBREW_PREFIX}/bin/slackcat-0.4-darwin-amd64", "#{HOMEBREW_PREFIX}/bin/slackcat"
    chmod "u+x", "#{HOMEBREW_PREFIX}/bin/slackcat"
    puts "Ready to go! Now just put your key in ~/.slackcat."
  end

  test do
    File.exist? "#{HOMEBREW_PREFIX}/bin/slackcat"
  end
end
