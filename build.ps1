$Name="origins-gardening-sounds"
$env:GOOS="windows"; go build -o bin/${Name}-win.exe -ldflags="-H=windowsgui" .
echo "Build complete."