# Propser, a Forth Interpreter written in Go

Go Forth and Propser

```sh
go run ./cmd/prosper
```

```text
Go Forth 0.0.0-alpha, Copyright (C) 2024 David Higgins <zoul0813@me.com>
Compiled with go1.24.0
Core: runtime/core.f
Type 'bye' to exit

> words or .sr seed branch mod drop @ depth clear generator cells >r octal 2dup new (debug) xor 2>r showreturn vector .row um* & r@ decimal .score * do 1- and nop max 2r@ space init-grid negate variable rnd (showstack) rot < /mod - <> pick until >> row-col-box (emit) nip << draw over 2over box@ box! constant > -rot tuck -bit score / abs . 2swap create fill + ^ , 0< 2abs u/mod ?dup .r page gptr (at-xy) i 0= +! grid [char] ." hex u< test (showreturn) swap utime spaces begin (clearscreen) showstack .cols i.r dup = loop emit u* 2nip here ?branch words 1+ u .s random | ! cr bye flag min 2drop 2r> r> ms &^ allot binary  ok
>
```

Prosper is a basic Forth interpreter. Most of its runtime words are defined in [runtime/core.f](runtime/core.f).

To play 2048, run this from the project directory in a terminal:

```sh
go run ./cmd/prosper 2048.fs
```

At the `>` prompt, enter `play`. Use `w/a/s/d` to move up/left/down/right;
press `q` to return to the prompt. Moves take effect immediately without Enter.
Enter `play` again to restart with a cleared board and score. Enter `bye` or press
Ctrl+C to exit Prosper. A full board remains responsive to input and `q`.

The command accepts source filenames and loads them in order after `runtime/core.f`.
Source errors stop startup and identify the failing file.

Build and validate:

```sh
go build -o bin/prosper ./cmd/prosper
go test ./... -timeout 30s
go vet ./...
```

The game tests cover all four move directions, merge scoring, random placement,
full boards, replay, quitting, and data/return-stack balance.

[2048.fs](2048.fs) is adapted from the game published by
[Dave Bucklin](https://www.davebucklin.com/play/2019/08/10/2048.html).
[2048_raw.fs](2048_raw.fs) preserves the supplied original source and comments.
The local version corrects movement past unequal adjacent tiles, prevents tile
placement from hanging on full boards, skips spawning after unchanged moves,
and resets the score on replay. The supplied original has no explicit copyright
or license notice; its author and source attribution are retained.

This game was used to test various aspects of Prosper, and functionality was added to Prosper to get this game to run.

The `draw` word renders the 2048 display board, for example.

```text
> draw
```

```text
    Score:   112

   __   __   __   __

   __   __   __   __

   __   __   __   __

   __   __   __   __

 ok
>
```

