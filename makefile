.PHONY:clean
clean:

.PHONY:run
run: main.go
	go run main.go

.PHONY:debug
debug: main.go
	dlv debug