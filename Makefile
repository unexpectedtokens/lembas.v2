run:
	@templ generate
	@tailwindcss -i ./assets/css/input.css -o ./assets/public/css/output.css --minify
	@go run cmd/main.go

tailwind:
	@tailwindcss -i ./assets/css/input.css -o ./assets/public/css/output.css --minify --watch
