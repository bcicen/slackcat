class Slackcat < Formula
  desc "Simple command-line Utility to post snippets to Slack."
  homepage "https://github.com/vektorlab/slackcat"
  url "https://github.com/vektorlab/slackcat/archive/master.tar.gz"
  version "0.4"
  sha256 "02c23e8b7bf6a45f85c6911795324a208e1ed9b30b5e9bfb851a74217c68fbbc"

  depends_on "go"

  def install
    platform = `uname`.downcase.strip

    unless ENV['GOPATH']
      ENV['GOPATH'] = "/tmp"
    end

    system "make"
    bin.install "build/slackcat-0.4-#{platform}-amd64" => "slackcat"

    puts "Ready to go! Now just put your key in ~/.slackcat."
  end

  test do
      assert_equal(0, "/usr/local/bin/slackcat")
  end
end
