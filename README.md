A simple CLI-based messaging application with (intentionally) no auth. At the point of writing this, no encryption is provided, but it is likely to include something along the lines of PGP for client-server communications.

Running Locally:
If running both parts of the application locally, you can run "docker compose up --build" in the server folder, and a simple "go run main.go" in the client folder. This will connect you to the server's websocket handler and you can begin sending messages.

Networking:
I've yet to host the server remotely, but some simple networking and environment variables changes/additions will allow you to run both application parts in their respective docker containers.
