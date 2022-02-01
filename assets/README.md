## Windows images

Windows images have to be `.ico` files (with multiple resolutions)

```
convert -background none assets/sunset.svg -define icon:auto-resize asset/sunset.ico
```

## Converting images to Go binary

```
~/go/bin/2goarray sunriseIcon main < assets/sunrise.ico > sunriseIcon.go
```
