module github.com/YendisFish/sirius/guarddog

go 1.26.2

require (
	github.com/YendisFish/sirius/toybox v0.0.0
	golang.org/x/term v0.45.0
)

require golang.org/x/sys v0.47.0 // indirect

replace github.com/YendisFish/sirius/toybox => ../toybox
