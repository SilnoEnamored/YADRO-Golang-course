build:
	go build -o bin/myapp.exe myapp/main.go

clean:
	del /F /Q bin\myapp.exe
