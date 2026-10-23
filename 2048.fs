\ Original game by Dave Bucklin.
\ Source: https://www.davebucklin.com/play/2019/08/10/2048.html
\ Local adaptation for Prosper; unmodified reference: 2048_raw.fs.
\ Changes: guard full-grid placement, spawn only after a move, reset score,
\ and advance
\ the store pointer past adjacent unequal tiles to prevent invalid merges.
\
\ 2048. To play, load this file and then run `play`.  In gforth, you can
\ invoke gforth with the name of the file to load the source. i.e. `gforth
\ 2048.fs`.

\ random number generator
Variable seed
$10450405 Constant generator
: rnd  ( -- n )  seed @ generator um* drop 1+ dup seed ! ;
: random ( n -- 0..n1 )  rnd um* nip ;
utime drop seed !
\ debug flag
variable debug 0 debug !
\ grid pointer
variable gptr 0 gptr !
\ score
variable score 0 score !
\ vector
variable vector
\ flag
variable flag
\ Track whether a move changed the board before adding another tile.
variable moved

\ init 4x4 grid
create grid 16 cells allot
: init-grid grid 16 cells 0 fill ;

\ grid access
\ row and col starting index is zero
: row-col-box ( row col -- box ) flag @ if swap then swap 4 * + ;
: box@ ( u-box -- u-value ) debug @ if dup cr ." getting box " . then cells grid + @ ;
: box! ( u-val u-box -- ) debug @ if 2dup swap cr ." storing " . ." in box " . then cells grid + ! ;
\ Count empty cells so a full board cannot hang random tile placement.
: empty-cells ( -- n )
  0 16 0 do i box@ 0= if 1+ then loop ;
\ place tile in random open spot on the edge
: new ( -- )
  empty-cells if
  begin
    16 random dup
    box@ 0= if 2 swap box! -1 else drop 0 then
  until
  then ;
\ draw grid
: draw ( -- ) page
  ." Score: " score @ 5 u.r cr cr
  4 0 do ( row )
    4 0 do ( col )
      j i swap 4 * + box@ ?dup 0= if 3 spaces ." __" else 5 u.r then
    loop
    cr cr
  loop ;

: compute ( limit index -- )
  0 moved !
  \ if the vector is less than 0, set gptr to 3, else 0
  vector @ 0< if 3 else 0 then
  debug @ if cr ." setting gptr to " dup . then
  gptr !
  \ Four rows
  4 0 do
    2dup do
      \ get the value of the box at the loop index
      j i row-col-box box@
      dup 0<> if
        \ if the value in b1 is non-zero, get the value of the box at the position of gptr
        j gptr @ row-col-box box@
        dup 0= if
          \ if the pointer box is zero, move the look box value into it
          drop
          -1 moved !
          debug @ if cr ." moving " dup . then
          j gptr @ row-col-box box!
          j i row-col-box 0 swap box!
          \ put the vector on the stack for +loop
          vector @
        else
          2dup = if
            \ if the values are non-zero and equal, merge them
            \ add the sum to the score
            -1 moved !
            + dup score +!
            debug @ if dup cr ." writing sum " . then
            j gptr @ row-col-box box!
            j i row-col-box 0 swap box!
            \ Increment both the loop and the pointer by the vector.
            vector @
            debug @ if dup cr ." increment gptr by " . then
            dup gptr +!
          else
            2drop
            \ if two nonzero values differ,
            \ either advance the loop or advance the pointer
            \ and put the vector on the stack for +loop
            i gptr @ vector @
            0< if swap then
            - 1 > if
              vector @ gptr +! 0
            else
              \ Adjacent unequal tiles occupy both slots; advance both pointers.
              vector @ dup gptr +!
            then
          then
        then
      else
        \ if the box we are looking at is zero, advance the loop
        drop vector @
      then
      dup abs 1 > abort" excess increment "
    +loop
    \ reset store pointer to 3 when vector is -1
    vector @ 0< if 3 else 0 then gptr !
  loop 2drop ;

\ gameplay
: up 1 vector ! -1 flag ! 4 1 compute ;
: left 1 vector ! 0 flag ! 4 1 compute ;
: down -1 vector ! -1 flag ! 0 2 compute ;
: right -1 vector ! 0 flag ! 0 2 compute ;
\ get move
: ?move ." (w,a,s,d,q) >" key ;
: play
  0 score !
  init-grid new
  begin
    draw ?move
         dup [char] w = if up moved @ if new then
    else dup [char] a = if left moved @ if new then
    else dup [char] s = if down moved @ if new then
    else dup [char] d = if right moved @ if new then
    then then then then
  [char] q = until ;
