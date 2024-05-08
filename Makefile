build:
	go build -o bin/myapp.exe cmd/xkcd/main.go

run:
	bin/myapp.exe

run_search:
	bin/myapp.exe -s "I'm following your questions"

run_search_index:
	bin/myapp.exe -s "I'm following your questions" -i

clean:
	del /F /Q bin/myapp.exe
