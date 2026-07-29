#!/usr/bin/env ruby
# frozen_string_literal: true

require "json"
require "open3"

root = File.expand_path("..", __dir__)
output, status = Open3.capture2e("go", "mod", "edit", "-json", chdir: root)
abort "无法读取 go.mod：#{output}" unless status.success?

manifest = JSON.parse(output)
replacements = Array(manifest["Replace"])
unless replacements.empty?
  paths = replacements.map { |item| item.dig("Old", "Path") }.compact
  abort "正式 Example 不允许本地 replace：#{paths.join(", ")}"
end

internal = Array(manifest["Require"]).select do |item|
  item.fetch("Path").start_with?("github.com/neteast-software/")
end
unstable = internal.reject { |item| item.fetch("Version", "").match?(/\Av\d+\.\d+\.\d+\z/) }
unless unstable.empty?
  versions = unstable.map { |item| "#{item.fetch("Path")}@#{item.fetch("Version", "<none>")}" }
  abort "正式 Example 只能依赖稳定语义版本：#{versions.join(", ")}"
end

puts "正式依赖通过：neteast=#{internal.length} replace=0"
