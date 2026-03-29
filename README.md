# MessageApp (TCP Chat in Go)

A simple terminal-based chat application built with Go using raw TCP sockets.

## Features

- Multi-client chat server
- Set username with `/user`
- List connected users with `/users`
- Broadcast messages with `/msg`
- Private message to a user with `/mtu`
- Group chat:
  - Create group: `/group`
  - List groups: `/groups`
  - Join group: `/join`
  - Leave group: `/getout`
  - Message group: `/mtg`
- Built-in help with `/help`

## Project Structure

- `server/`: TCP server
- `client/`: terminal client

## Requirements

- Go 1.21+ (recommended)

## Run the Server

From the `server` folder:

```bash
go run . 8080
```

The server listens on `:8080`.

## Run a Client

From the `client` folder (open a new terminal for each client):

```bash
go run . 8080
```

Then enter your username when prompted.

## Basic Usage

After connecting, you can use commands:

```text
/user alice
/users
/msg hello everyone
/mtu bob hi bob
/group backend
/groups
/join backend
/mtg backend deploy at 5pm
/getout backend
/help
/quit
```

## Notes

- This project is intentionally minimal and educational.
- Output uses ANSI colors; color rendering depends on your terminal.
- Current server logic is split into multiple files to keep code maintainable as features grow.
