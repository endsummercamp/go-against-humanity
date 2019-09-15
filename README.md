# Go Against Humanity

A Cards Against Humanity web game written in [Go](http://www.golang.org/).
The game is intended to be played with multiple players, divided into 
"players" and "jurors". The card that gets most votes wins.

Written in 3 days during the [End Summer Camp](https://www.endsummercamp.org). Meant to be played for a good laugh together.

![Cards Against Humanity](screenshots/1.png)

### Get deps
```bash
go get -v
```

### Start the web server:
```
go run server.go
```

### Tell the players / jurors to go to http://<your-ip>:1323/ and enjoy the game!


## Based on [Echo](https://echo.labstack.com/)
