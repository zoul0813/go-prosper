: emit (emit) ;
: . emit ;



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
\ r@      [ ]

: 2DUP OVER OVER ;
: 2DROP DROP DROP ;
: 2SWAP ROT >R ROT R> ;
: 2OVER 2>R 2DUP 2R> 2SWAP ;
: 2NIP 2>R 2DROP 2R> ;
\ 2TUCK    [ ]
\ 2ROT     [ ]
: 2>R SWAP >R >R ;
: 2R> R> R> SWAP ;
\ 2r@      [ ]

\ clear     [x]
\ depth     [x]
: showstack (showstack) ;
: .s showstack ;