![Screenshot](docs/screenshot.png)

# SSH Guess game

Guess a number from 1 to 1 Million

## Try Now!
```sh
$ ssh guess.ancoder.net -p 2220
```

## Build
```sh
$ go build .
```

## Run with docker
```sh
$ docker run -d -p 12345:12345 -v "~/.ssh:/app/.ssh" ghcr.io/ancoder2009/ssh-guess-game:[hash]
```