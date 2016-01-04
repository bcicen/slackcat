class Slackcat < Formula
  desc "Simple command-line Utility to post snippets to Slack."
  homepage "https://github.com/vektorlab/slackcat"
  url "https://github.com/vektorlab/slackcat/archive/master.tar.gz"
  version "0.4"
  sha256 "27fdfe083752f810cc5a8d383093a84556a355a3e1aa54a1ab03f3a183bc6b03"

  depends_on "go"

  def install
    platform = `uname`.downcase.strip

    unless ENV["GOPATH"]
      ENV["GOPATH"] = "/tmp"
    end

    system "make"
    bin.install "build/slackcat-0.4-#{platform}-amd64" => "slackcat"

    puts "Ready to go! Now just put your key in ~/.slackcat."
  end

  test do
    assert_equal(0, "/usr/local/bin/slackcat")
  end
end
