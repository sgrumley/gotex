module github.com/sgrumley/gotex

go 1.24

toolchain go1.24.2

require (
	github.com/gdamore/tcell/v2 v2.7.4
	github.com/rivo/tview v0.0.0-20241016194538-c5e4fb24af13
	github.com/sgrumley/gotex/pkg v0.0.0-00010101000000-000000000000
	github.com/sgrumley/hotfuzz v0.0.0-20250424093722-0c620013625e
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sgrumley/gotex/pkg => ./pkg
