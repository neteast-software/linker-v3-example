package console

import "github.com/neteast-software/go-module/graph/console/theme"

func Theme() theme.Theme {
	return theme.New("default", theme.System).
		Allow(theme.Light, theme.Dark, theme.System).
		Light(theme.Style(
			theme.Colors{
				Primary: theme.Hex("#196A5B"), Info: theme.Hex("#146C94"),
				Success: theme.Hex("#2E7D32"), Warning: theme.Hex("#A65C00"),
				Danger: theme.Hex("#B42318"), Surface: theme.Hex("#FFFFFF"),
				SurfaceAlt: theme.Hex("#F4F6F8"), Text: theme.Hex("#172033"),
				TextMuted: theme.Hex("#566174"), Border: theme.Hex("#D8DEE9"),
			},
			theme.Chart(
				theme.Hex("#196A5B"), theme.Hex("#146C94"), theme.Hex("#A65C00"),
				theme.Hex("#7A5195"), theme.Hex("#B42318"), theme.Hex("#2E7D32"),
			).
				SequentialColors(theme.Hex("#DDF3EE"), theme.Hex("#196A5B")).
				DivergingColors(theme.Hex("#146C94"), theme.Hex("#F4F6F8"), theme.Hex("#B42318")),
		)).
		Dark(theme.Style(
			theme.Colors{
				Primary: theme.Hex("#5EC6B2"), Info: theme.Hex("#60A5FA"),
				Success: theme.Hex("#4ADE80"), Warning: theme.Hex("#FBBF24"),
				Danger: theme.Hex("#FB7185"), Surface: theme.Hex("#111827"),
				SurfaceAlt: theme.Hex("#1F2937"), Text: theme.Hex("#F9FAFB"),
				TextMuted: theme.Hex("#CBD5E1"), Border: theme.Hex("#374151"),
			},
			theme.Chart(
				theme.Hex("#5EC6B2"), theme.Hex("#60A5FA"), theme.Hex("#FBBF24"),
				theme.Hex("#C084FC"), theme.Hex("#FB7185"), theme.Hex("#4ADE80"),
			).
				SequentialColors(theme.Hex("#1F2937"), theme.Hex("#5EC6B2")).
				DivergingColors(theme.Hex("#60A5FA"), theme.Hex("#1F2937"), theme.Hex("#FB7185")),
		))
}
