( variables )
variable score 	0 score !
variable gptr   0 gptr !
variable flag
variable vector

( random number generator )
variable seed
$10450405 constant generator
: rnd		seed @ generator um* drop 1+ dup seed ! ;
: random	rnd um* nip ;
utime drop seed !

( grid )
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

( display )

: .score	 ." Score: " score @ 5 u.r ;
: .row     swap 4 * + box@ ?dup 0= if 3 spaces ." __" else 5 u.r then ;
: .cols 4 0 do i .row loop ;
\ : .rows 4 0 do i .cols cr cr loop ;
: draw 		 page
  .score
  cr cr
  4 0 do ( row )
    i .cols
    \ 4 0 do ( col )
    \  j i .row
    \ loop
    cr cr
  loop
  ;

: set-gptr    vector @ 0< if 3 else 0 then gptr ! ;

: compute   ( limit index -- )
  set-gptr
  4 0 do ( row )
    2dup do ( col )
      j i row-col-box box@
      dup 0<> if
        ( if the value in b1 is non-zero, get the value of the box at gptr )
        j gptr @ row-col-box box@
        dup 0= if
          ( if the pointer box is zero, move the look box value into it )
          drop
          j gptr @ row-col-box box!
          j i row-col-box 0 swap box!
          vector @ ( put the vector on the stack for +loop )
        else
          2dup = if
            + dup score +!
            j gptr @ row-col-box box!
            j i row-col-box 0 swap box!
            vector @
            dup gptr +!
          else
            2drop
            i gptr @ vector @
            0< if swap then
            - 1 > if
              vector @ gptr +! 0
            else
              vector @
            then
          then
        then
      else
        ( if the bo we are looking at is zero, advance the loop)
        drop vector @
      then
      dup abs 1 > abort" excess increment "
    +loop
    set-gptr ( reset gptr )
  loop
  2drop
  ;

( game )
: up      1 vector ! -1 flag ! 4 1 compute ;
: left    1 vector ! 0 flag ! 4 1 compute ;
: down    -1 vector ! -1 flag ! 0 2 compute ;
: right   -1  vector ! 0 flag ! 0 2 compute ;
: ?move   ." (w,a,s,d,q) >" key ;
: play
  init-grid new
  begin
    draw ?move
         dup [char] w = if up new
    else dup [char] a = if left new
    else dup [char] s = if down new
    else dup [char] d = if right new
    then then then then
    [char] q = until
  ;
