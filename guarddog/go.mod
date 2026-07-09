module github.com/YendisFish/sirius/guarddog

go 1.26.2

require github.com/YendisFish/sirius/toybox v0.0.0

require (
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)

replace github.com/YendisFish/sirius/toybox => ../toybox
