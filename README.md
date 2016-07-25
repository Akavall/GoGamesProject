**Build a binary file example:**

In the root directory GoGamesProject run

```
go build basic_server.go
```

**Run just one test file example:**

In GoGamesProject/dice run

```
go test
```

**Run all test files example:**

In GoGamesProject run

```
go test ./...
```

**Build local server**


```
go build
./GoGamesProject
```

or

Play against AI:

```
http://localhost:8000/zombie_dice
```

Play multi-player (running on local is mostly for testing in this case):

```
http://localhost:8000/zombie_dice_multi_player
```

Steps how to play multi-player:

1) Host player hits `Start Game`, this generates game_id, for example:
`c67c6624-9da2-486b-4726-d28a1dc3215f`

2) Host sends game_id to guest.

3) Guest pastes game_id in "Join Game" field and hits `Join Game`.
   Host will see that "Players in Game" has increased to 2, and the game
   can start.

4) Host hits `Take Turn` and the game starts.