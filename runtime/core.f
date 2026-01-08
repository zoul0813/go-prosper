: emit (emit) ;
: . emit ;
: hex (hex) ;
: decimal (decimal) ;
: octal (octal) ;
: binary (binary) ;

: 0< 0 < ;
: 0= 0 = ;
: 1+ 1 + ;
: 1- 1 - ;
: +! dup >r @ + r> ! ;
: 2ABS ABS SWAP ABS SWAP ;
: U< 2ABS < ;
: U 2ABS > ;
: U* 2ABS * ;
: U/MOD 2ABS /MOD ;


\ DUP     [x]
\ ?DUP    [ ]
: ?DUP DUP 0 <> ?BRANCH 1 DUP ;
\ DROP    [x]
\ SWAP    [x]
\ OVER    [x]
\ ROT     [x]
: -ROT ROT ROT ;
: NIP SWAP DROP ;
: TUCK SWAP OVER ;

\ >R      [x]
\ R>      [x]
: R@ ( -- N ) R> DUP >R ;

: AND ( N1 N2 -- N3 ) & ;
: OR ( N1 N2 -- N3 ) | ;
: XOR ( N1 N2 -- N3 ) ^ ;
: -BIT ( N1 N2 -- N3 ) &^ ;
: NOP ( -- ) ;

: 2DUP OVER OVER ;
: 2DROP DROP DROP ;
: 2SWAP ROT >R ROT R> ;
: 2OVER 2>R 2DUP 2R> 2SWAP ;
: 2NIP 2>R 2DROP 2R> ;
\ 2TUCK    [ ]
\ 2ROT     [ ]
: 2>R SWAP >R >R ;
: 2R> R> R> SWAP ;
: 2R@ 2R> 2DUP 2>R ;

\ clear     [x]
\ depth     [x]
: SHOWSTACK (showstack) ;
: .S SHOWSTACK ;
: SHOWRETURN (showreturn) ;
: .R SHOWRETURN ;
: .SR ." S: " SHOWSTACK ." R: " SHOWRETURN ;

: SPACE ."  " ;
: SPACES 0 DO SPACE LOOP ;
\ : PAGE CR ; \ TODO: temporary, should clear the screen?
\ : PAGES 0 DO PAGE LOOP ;
: PAGE (clearscreen) 0 0 (at-xy) ;

: VARIABLE ( "<spaces>name" -- )
  HERE
  0 ,
  CONSTANT
;

: CREATE ( -- a-addr )
  VARIABLE
;


\\\\\\\ TESTS

\ : test branch 2 1 2 3 4 5 6 ;
\ : test? ?branch 4 1 . branch 2 0 . ;
\ : test 10 begin 1 - dup . dup until ;
\ : dotest 10 5 do i . loop ;
\ : cols 4 0 do dup i . . cr loop ;
\ : rows 4 0 do i cols drop cr cr loop ;

: test if 10 10 * . ." hello" else 5 5 * . ." goodbye" then ;

\ 10 5
\ 2>R
\ 2R@ SWAP DROP .
\ 2R> 1 + 2DUP <> 2>R
\ ?branch -12
\ 2R> 2DROP

\: back
\  branch 2
\  ?branch 5
\  20 .
\  0
\  branch -7
\  1 + .s
\  ;
\15 back

\ TODO: implement constant in Forth?
\ : CONSTANT ( x "<spaces>name" -- )
\   >R : >R
\   POSTPONE LITERAL
\   POSTPONE ;
\ ;



\\\\\ 2048 tests


\ 2048
variable score 	112 score !
variable gptr   0 gptr !
variable flag
variable vector

\ ( random number generator )
variable seed
$10450405 constant generator
: rnd		                seed @ generator um* drop 1+ dup seed ! ;
: random	( cell -- )   rnd um* nip ;
utime drop seed !

\ ( grid )
create grid 16 cells allot
: init-grid	grid 16 cells 0 fill ;
: row-col-box   ( row col -- box )    flag @ if swap then swap 4 * + ;
: box@          ( u-box -- u-value )  cells grid + @ ;
: box!          ( u-val u-box -- )    cells grid + ! ;
: new           ( --  )
  begin
    16 random dup
    box@ 0= if 2 swap box! -1 else drop 0 then
  until ;

: .score    ."     Score: " score @ 5 i.r ;
: .row      swap 4 * + box@ ?dup 0= if 3 spaces ." __" else 5 i.r then ;
: .cols
  4 0 do
    dup
    i .row
  loop
  ;


: draw
  page
  .score
  cr cr
  4 0 do ( row )
    i .cols
    cr cr
  loop
  ;



\ : comma ." ," ;
\ : test 5 0 do dup . comma i . space loop ;
\ : test2 5 0 do i test cr loop ;
\ : test      ." \"hello world\"" ;

\ : test 16 0 do i box@ ?dup 0= if 3 spaces ." __" else . then loop ;
\ : test if ." hello world" cr cr else ." goodbye" cr cr then ;

\ : .row swap 4 * + box@ ?dup 0= if 3 spaces ." __" else 5 i.r then ;