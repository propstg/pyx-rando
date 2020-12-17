# pyx-rando

A simple Rando Cardrissian bot for [Pretend You're Xyzzy](https://github.com/ajanata/PretendYoureXyzzy).

Notes:
1. Only tested on a privately hosted instance. Probably doesn't work on the public server without modifications.
2. This bot will judge cards, when it is its turn. I believe simple server changes would be needed to prevent this (by skipping players named `/Rando[0-9]+/`), if this isn't desired.
3. Likely (definitely) doesn't catch all the errors it should, in general.

## Running
1. `go build .`
2. Edit .env to point to your instance and the gameId of the game you want it to join (grab the `game` number from the url)
3. `./rando.exe`