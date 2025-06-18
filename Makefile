


run:
	go tool templ generate -cmd="go run ./cmd/t3-clone" -watch 

css:
	./tailwindcss -i static/css/input.css -o static/css/styles.css --minify --watch



.PHONY: run, css 

