# Go Against Humanity

A Cards Against Humanity web game written in [Go](http://www.golang.org/).
The game is intended to be played with multiple players, divided into 
"players" and "jurors". The card that gets most votes wins.

Written in 3 days during the [End Summer Camp](https://www.endsummercamp.org). Meant to be played for a good laugh together.

![Cards Against Humanity](screenshots/1.png)

## Usage

```bash
export GOPATH=~/gopath # If you have an existing gopath, use that instead
mkdir -p $GOPATH/src/github.com/ESCah
git clone --recursive https://github.com/ESCah/go-against-humanity

make compile      # Compile the cards

npm install                    # Install Node.js deps
node node_modules/.bin/webpack # Compile Web resources

go get -v        # Fetch deps
```

## Configure the application
Create a file `config.toml` and add the following:

```toml
[General]
Decks = ['ita-original-sfoltita', 'ita-espansione', 'ita-HACK']
```

where Decks is an array containing the decks you want to use
([here is the list](https://github.com/ESCah/json-against-humanity/tree/master/src))


## Running
```bash
go run server.go # Start the web server
```

**Tell the players / jurors to go to http://<your-ip>:1323/ and enjoy the game!**

## Based on [Echo](https://echo.labstack.com/)
