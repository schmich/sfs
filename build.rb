#!env ruby

version = `git tag`.lines.last.strip
commit = `git rev-parse HEAD`.strip

system("go build -ldflags \"-w -s -X main.version=#{version} -X main.commit=#{commit}\"") || fail
