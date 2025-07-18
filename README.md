# Propser, a Forth Interpreter written in Go

Go Forth and Propser

```sh
 > go run cmd/prosper/main.go
```

```text
Go Forth 0.0.0-alpha, Copyright (C) 2024 David Higgins <zoul0813@me.com>
Compiled with go1.24.0
Core: cmd/prosper/core.f
Type 'bye' to exit

> words or .sr seed branch mod drop @ depth clear generator cells >r octal 2dup new (debug) xor 2>r showreturn vector .row um* & r@ decimal .score * do 1- and nop max 2r@ space init-grid negate variable rnd (showstack) rot < /mod - <> pick until >> row-col-box (emit) nip << draw over 2over box@ box! constant > -rot tuck -bit score / abs . 2swap create fill + ^ , 0< 2abs u/mod ?dup .r page gptr (at-xy) i 0= +! grid [char] ." hex u< test (showreturn) swap utime spaces begin (clearscreen) showstack .cols i.r dup = loop emit u* 2nip here ?branch words 1+ u .s random | ! cr bye flag min 2drop 2r> r> ms &^ allot binary  ok
>
```

Prosper is a very basic Forth interpreter, the majority of the words available for Propser are contained in [cmd/prosper/core.f](cmd/prosper/core.f)

[2048.fs](https://www.davebucklin.com/play/2019/08/10/2048.html) was obtained from the blog by [Dave Bucklin](https://www.davebucklin.com/).

This game was used to test various aspects of Prosper, and functionality was added to Prosper to get this game to run.

The `draw` word in `core.f` will render the 2048 display board, for example.

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

