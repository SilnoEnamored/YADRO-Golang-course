build:
	go build -o bin/myapp.exe cmd/xkcd/main.go

clean:
	del /F /Q bin\myapp.exe
