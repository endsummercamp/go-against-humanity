# Go Against Humanity

A Cards Against Humanity web game written in [Go](http://www.golang.org/).
The game is intended to be played with multiple players, divided into 
"players" and "jurors". The card that gets most votes wins.

Written in 3 days during the [End Summer Camp](https://www.endsummercamp.org). Meant to be played for a good laugh together.

![Cards Against Humanity](screenshots/1.png)

## Usage

```bash
git clone --recursive https://github.com/ESCah/go-against-humanity

make ita-original # Choose your cards
make compile      # Compile them

go get -v        # Fetch deps
go run server.go # Start the web server
```

**Tell the players / jurors to go to http://<your-ip>:1323/ and enjoy the game!**

## Based on [Echo](https://echo.labstack.com/)
