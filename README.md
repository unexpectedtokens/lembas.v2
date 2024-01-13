# Lembas

To run this app you'll need.

- To install the standalone tailwind cli + add it to $PATH. The app will run without it. But any css/classname changes will have no effect. See: https://tailwindcss.com/blog/standalone-cli
- Install templ and add it to $PATH. Without it, you won't be able to generate code based on template changes
- Create a .env file and define a MONGO_URI variable, pointing to you mongodb instance
- Run the app with either `make` or `go run cmd/main.go`
