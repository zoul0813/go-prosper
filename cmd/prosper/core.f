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

\ : test branch 2 1 2 3 4 5 6 ;
\ : test? ?branch 4 1 . branch 2 0 . ;
\ : test 10 begin 1 - dup . dup until ;
: dotest 10 5 do i . loop ;
: cols
  4 0
  do
    dup ( .s .r cr )
    i ( .s .r cr )
    . . cr
    ( .s .r cr )
  loop
  ;
: rows
  4 0
  do
    i ( .s .r cr )
    cols
    drop
    cr cr
  loop
  ;

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

: VARIABLE ( "<spaces>name" -- )
  HERE
  0 ,
  CONSTANT
;

: CREATE ( -- a-addr )
  VARIABLE
;

\ TODO: implement constant in Forth?
\ : CONSTANT ( x "<spaces>name" -- )
\   >R : >R
\   POSTPONE LITERAL
\   POSTPONE ;
\ ;
