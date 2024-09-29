# Fret Files

> Sorting thruough your concerns

Simple File organizer for your Downloads folder

## Tech Stack

- Go and Wails
- "github.com/joho/godotenv"
- "github.com/sevlyar/go-daemon"

## Installation

1. Clone the repository
2. add .env file with the following content

```
WATCH_PATH=YOUR_DOWNLOADS_PATH
```

3. Run the following commands

```bash
go get -u github.com/joho/godotenv
go get -u github.com/sevlyar/go-daemon

```

4. Build the executable daemon

```bash
go build -o mydaemon main.go
```

5. Run the daemon

```bash
./mydaemon
```

6. kill the daemon

```bash
./stop_daemon.sh
```
