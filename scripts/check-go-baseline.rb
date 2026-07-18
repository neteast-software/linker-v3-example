#!/usr/bin/env ruby
# frozen_string_literal: true

require "open3"

class ExampleGoBaseline
  GO_VERSION_PATTERN = /\A\d+\.\d+\.\d+\z/
  LINKER_VERSION_PATTERN = /\Av3\.\d+\.\d+\z/

  def initialize(root, workflows: nil, toolchain: nil)
    @root = File.expand_path(root)
    @workflows = workflows
    @toolchain = toolchain
  end

  def go_version
    module_content[/^go\s+(\S+)\s*$/, 1].to_s
  end

  def linker_version
    module_content[/^\s*(?:require\s+)?github\.com\/neteast-software\/linker\/v3\s+(\S+)\s*$/, 1].to_s
  end

  def workflow_paths
    @workflows || Dir.glob(path(".github/workflows/*.{yaml,yml}")).sort
  end

  def toolchain
    return @toolchain unless @toolchain.nil?

    output, error, status = Open3.capture3("go", "env", "GOVERSION", chdir: @root)
    raise "go env GOVERSION 失败: #{error.strip}" unless status.success?

    output.strip.delete_prefix("go")
  end

  def errors
    failures = []
    unless go_version.match?(GO_VERSION_PATTERN)
      failures << "go.mod: go=#{go_version.inspect}，期望 MAJOR.MINOR.PATCH"
    end
    unless linker_version.match?(LINKER_VERSION_PATTERN)
      failures << "go.mod: linker=#{linker_version.inspect}，期望 v3.MINOR.PATCH"
    end
    if module_content.match?(/^replace\s+github\.com\/neteast-software\/linker\/v3\b/)
      failures << "go.mod: 不应使用 replace 覆盖已发布的 Linker"
    end

    readme = File.read(path("README.md"))
    documented_go = readme[/当前工具链基线为 Go `([^`]+)`/, 1]
    documented_linker = readme[/framework 基线为 linker `([^`]+)`/, 1]
    failures << "README.md: Go 基线=#{documented_go.inspect}，期望 #{go_version}" unless documented_go == go_version
    unless documented_linker == linker_version
      failures << "README.md: Linker 基线=#{documented_linker.inspect}，期望 #{linker_version}"
    end

    workflow_versions = workflow_paths.flat_map do |workflow|
      content = File.read(workflow)
      content.scan(/^\s*go-version:\s*["']?([^\s"']+)["']?\s*$/).flatten.map do |version|
        [display(workflow), version]
      end
    end
    failures << ".github/workflows: 缺少 setup-go 的 go-version" if workflow_versions.empty?
    workflow_versions.each do |workflow, version|
      failures << "#{workflow}: go-version=#{version}，期望 #{go_version}" unless version == go_version
    end

    linker_refs = workflow_paths.each_with_object([]) do |workflow, refs|
      content = File.read(workflow)
      match = content.match(/repository:\s*neteast-software\/linker\s*\n\s*ref:\s*["']?([^\s"']+)["']?/)
      refs << [display(workflow), match[1]] unless match.nil?
    end
    failures << ".github/workflows: 缺少 Linker source-ready checkout" if linker_refs.empty?
    linker_refs.each do |workflow, ref|
      failures << "#{workflow}: Linker ref=#{ref}，期望 #{linker_version}" unless ref == linker_version
    end

    current = toolchain
    failures << "当前工具链=#{current}，期望 #{go_version}" unless current == go_version
    failures
  rescue Errno::ENOENT => e
    ["缺少基线文件: #{e.message}"]
  rescue RuntimeError => e
    [e.message]
  end

  private

  def module_content
    @module_content ||= File.read(path("go.mod"))
  end

  def path(value)
    File.join(@root, value)
  end

  def display(value)
    value.delete_prefix(@root + File::SEPARATOR)
  end
end

if $PROGRAM_NAME == __FILE__
  root = File.expand_path("..", __dir__)
  check = ExampleGoBaseline.new(root)
  failures = check.errors
  unless failures.empty?
    failures.each { |failure| warn "go-baseline: #{failure}" }
    warn "go-baseline: 请同步 go.mod、README、workflow 和已验证工具链后重试"
    exit 1
  end

  puts "Example 基线通过: go=#{check.go_version} linker=#{check.linker_version} toolchain=#{check.toolchain}"
end
