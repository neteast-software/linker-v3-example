# frozen_string_literal: true

require "fileutils"
require "minitest/autorun"
require "tmpdir"
require_relative "check-go-baseline"

class ExampleGoBaselineTest < Minitest::Test
  def test_current_repository_is_consistent
    root = File.expand_path("..", __dir__)
    check = ExampleGoBaseline.new(root)

    assert_empty check.errors
    assert_equal "1.26.5", check.go_version
    assert_equal "v3.3.3", check.linker_version
  end

  def test_reports_repository_projection_drift
    Dir.mktmpdir("example-go-baseline") do |root|
      FileUtils.mkdir_p(File.join(root, ".github/workflows"))
      File.write(File.join(root, "go.mod"), <<~MOD)
        module example.com/app

        go 1.26.5

        require github.com/neteast-software/linker/v3 v3.3.0
      MOD
      File.write(File.join(root, "README.md"), "当前工具链基线为 Go `1.26.4`，framework 基线为 linker `v3.2.0`\n")
      workflow = File.join(root, ".github/workflows/check.yml")
      File.write(workflow, <<~YAML)
        go-version: 1.26.4
        repository: neteast-software/linker
        ref: 3.0.0
      YAML

      check = ExampleGoBaseline.new(root, workflows: [workflow], toolchain: "1.26.3")

      assert_equal [
        "README.md: Go 基线=\"1.26.4\"，期望 1.26.5",
        "README.md: Linker 基线=\"v3.2.0\"，期望 v3.3.0",
        ".github/workflows/check.yml: go-version=1.26.4，期望 1.26.5",
        ".github/workflows/check.yml: Linker ref=3.0.0，期望 v3.3.0",
        "当前工具链=1.26.3，期望 1.26.5"
      ], check.errors
    end
  end
end
