--[ pico-8 go ]-- by harraps

--[[
goban is 19*19
each pos is 6*6px
therefore goban is 114*114px
screen is 128*128px
7px of padding around the goban
goban goes from 7 to 114
]]

-- initialize --

---[[
function _init()
 palt(0,false) --draw black
 palt(14,true) --don't for pink
 cursor=nv(9,9)
 goban={} --board of the game
 turn =false --black start
 endgm=0 --state of the game
 --0->playing
 --1->set dead stones
 --2->display territories
 --3->display scores
 --pass counter
 pass={} --count passes
 pass.c=false  --confirmation
 pass[false]=0 --passes black
 pass[true ]=0 --passes white
 terr=nil --territories list
 --scores
 score={} --score of players
 score[false]=0 --score black
 score[true ]=0 --score white
 --draw once
 d_goban () --draw goban once
 d_cursor() --draw cursor once
end
--]]

-- main loop --

---[[
function _update()
 if(not btnpress())return
 if endgm<2 then
  i_dpad() --move cursor
  if endgm<1 then
   i_goban()
   --if both have passed twice
   if  pass[false]>1
   and pass[true ]>1
   then endgm=1 end --next state
  else i_end() end
 else --score or territories
  if btnp(4) or btnp(5) then
   --switch state
   endgm=endgm>2 and 2 or 3
  end
 end
 --draw goban in most game state
 if(endgm<=2)d_goban ()
 if(endgm<=1)d_cursor()
 if(endgm==2)d_terr  ()
 if(endgm>=3)d_score ()
end
--]]

-- interaction --

---[[
--handle cursor movement
function i_dpad()
 local x,y=cursor.x,cursor.y
 --move the cursor
 if(btnp(0)) x-=1 --left
 if(btnp(1)) x+=1 --right
 if(btnp(2)) y+=1 --up
 if(btnp(3)) y-=1 --down
 --cap the value
 if     x<0  then x=0
 elseif x>18 then x=18 end
 if     y<0  then y=0
 elseif y>18 then y=18 end
 cursor.x=x
 cursor.y=y
end
--]]

---[[
--manage player inputs
function i_goban()
 if     btnp(5) then --button x
  --cancel passing
  if pass.c then pass.c=false
  else           pass.c=true end
 elseif btnp(4) then --button o
  if pass.c then --confirm pass
   pass.c=false  --reset passing
   pass[turn]+=1 --pass turn
   turn=not turn --switch turn
   return --stop
  end
  --place a goishi at cursor
  local g=n_goishi(turn,cursor)
  local ll=g_lib(g)
  --no more liberties
  if(#ll<=0)return --stop
  --adding stone failed
  if(not g_add(g))return --stop
  --check surrounding stone
  for d in all(_dir) do
   --try capture them
   g_capture(v_add(g,d))
  end
  pass[turn]=0 --reset pass turn
  turn=not turn --switch turn
 end
end
--]]

---[[
--clean groups of stones
function i_end()
 if     btnp(5) then --button x
  endgame() --end the game
 elseif btnp(4) then --button o
  --get goishi at cursor
  local o=g_goishi(cursor)
  --no goishi ?
  if(o==nil)return --stop
  --we need the second list
  local s=not(o.die or false)
  local ll,lg=g_lib(o)
  for g in all(lg) do
   g.die=s --switch goishi state
  end
 end
end
--]]

-- draw --

---[[
--draw the goban
function d_goban()
 --draw the background
 for x=0,3 do for y=0,3 do
  spr(8,x*32,y*32,4,4)
 end end
 --draw the goban
 --[[ heavy draw calls
 for i=0,18 do for j=0,18 do
  local c,x,y=19,7+i*6,115-j*6
  if     i==0  then c-=1
  elseif i==18 then c+=1  end
  if     j==0  then c+=16
  elseif j==18 then c-=16 end
  if(i%6==3 and j%6==3) c=33
  spr(c,x,y)
 end end
 --]]
 ---[[ light draw calls
 --draw grid
 for x=0,1 do for y=0,1 do
  spr(73,7+x*54,7+y*54,7,7)
 end end
 --draw borders
 for i=0,1 do
  spr(64,7+i*54,  7,7,1) --top
  spr(80,7+i*54,115,7,1) --bott
  spr(71,  7,7+i*54,1,7) --left
  spr(72,115,7+i*54,1,7) --right
 end
 --draw corners
 spr( 2,  7,  7) --top left
 spr( 4,115,  7) --top right
 spr(34,  7,115) --bottom left
 spr(36,115,115) --bottom right
 --draw main positions
 for x=3,15,6 do for y=3,15,6 do
  spr(33,7+x*6,7+y*6)
 end end
 --]]
 --draw the stones
 for g in all(goban) do
  local c,x,y=1,d_coord(g)
  if(g.col)c+=16 --white
  if(g.die)c+=4  --dead
  spr(c,x,y)
 end
end
--]]

---[[
--draw the cursor
function d_cursor()
 --if passing display message
 if pass.c then
  local c=turn and 7 or 0
  rectfill(32,50,95,76,3)
  rect    (32,50,95,76,c)
  print("pass turn ?",42,56,c)
  print("\142yes  \151no",42,66,c)
  return
 end
 --draw the cursor
 local c,x,y=
  endgm>0 and 32 or --cleaner
  turn    and 16 or 0,
  d_coord(cursor)
 spr(c,x,y)
end
--]]

---[[
--draw the score
function d_score()
 for x=0,3 do for y=0,3 do
  spr(12,x*32,y*32,4,4)
 end end
 rectfill(42,32,86,96, 3)
 rect    (42,32,86,96,11)
 print("score" ,54,36,11)
 --draw black score
 spr( 6,52,54,2,2)
 print(form(score[false]),
  66,58,0)
 --draw white score
 spr(38,52,74,2,2)
 print(form(score[true ]),
  66,78,7)
end
--]]

---[[
--draw territories on the goban
function d_terr()
 for t in all(terr) do
  local c=52 --none
  if false==t.col then c=37
  elseif    t.col then c=53 end
  for u in all(t) do
   local x,y=d_coord(u)
   spr(c,x,y)
  end
 end
end
--]]

---[[
--return x,y for draw functions
function d_coord(v)
 return 7+v.x*6,115-v.y*6
end
--]]

-- end game --

---[[
--generate the territories
function endgame()
 endgm=2 --territories state
 terr={}
 local v=nv(0,0)
 --for each square of the goban
 for x=0,18 do v.x=x
 for y=0,18 do v.y=y
  if not t_present(v) then
   local n=n_terr(v)
   --add new territory
   if(#n>0)add(terr,n)
  end
 end end
end
--]]

-- territories --

---[[
--create a new territory at pos
function n_terr(v)
 local t={} --territory
 t.cnt=0 --count positions
 --counter for goishis
 t[false]=0 --black goishis
 t[true ]=0 --white goishis
 --generate if valid pos
 if(valid(v))t_gen(t,v)
 --recover color of territory
 local b,w,c=t[false],t[true]
 if b>0 or w>0 then
  if     b<1 then c=true
  elseif w<1 then c=false end
 end
 t.col=c --color of territory
 --add the score
 if(c~=nil)score[c]+=t.cnt
 --we've specified the owner
 return t
end
--]]

---[[
--generate territory recursively
--(territory, position)
function t_gen(t,v)
 local g=g_goishi(v)
 local c=g==nil --no goishi
 if(not c)c=g.die==true --dead
 if c then --no goishi at pos
  add(t,v_clone(v)) --add pos
  t.cnt+=1 --count up
  for d in all(_dir) do
   _t_gen(t,v_add(v,d))
  end
 else --goishi at position
  t[g.col]+=1 --count goishi
  --we don't care if a goishi
  --is counted multiple times
 end
end
--]]

---[[
--used by g_terr
function _t_gen(t,v)
 --position invalid ?
 if(not valid(v))return --stop
 --pos already in territory ?
 for u in all(t) do
  if(v_equ(v,u))return --stop
 end
 t_gen(t,v) --add new position
end
--]]

---[[
--is position in territories
function t_present(v)
 for t in all(terr) do
  for u in all(t) do
   --position present
   if(v_equ(v,u))return true
  end
 end
 return false --not present
end
--]]

-- goban management --

---[[
--limits of the goban
function valid(v)
 return -1<v.x and v.x<19
    and -1<v.y and v.y<19 
end
--]]

-- stone/goishi functions --

---[[
--create a new goishi
--(color, position)
function n_goishi(c,v)
 return {col=c,x=v.x,y=v.y}
end
--]]

---[[
--get goishi at position
--(position, [list])
function g_goishi(v,l)
 l=l or goban
 for g in all(l) do
  if(v_equ(v,g))return g
 end
end
--]]

---[[
--capture goishis from position
function g_capture(v)
 local o=g_goishi(v)
 if(o==nil)return --stop
 local ll,lg=g_lib(o)
 --still has liberties
 if(#ll>0)return --stop
 --no more liberties
 for g in all(lg) do
  del(goban,g) --remove goishis
 end
end
--]]

---[[
--add goishi if position is free
function g_add(g)
 local o=g_goishi(g)
 --goishi already at position
 if(o~=nil)return false
 add(goban,g) --add new goishi
 return true
end
--]]

---[[
--get liberties for group
--(current goishi,check next,
-- liberties,goishis checked)
function g_lib(g,n,ll,lg)
 n =n  or n==nil --check next
 ll=ll or {} --liberties
 lg=lg or {} --goishis
 --add stone to list
 add(lg,g)
 for d in all(_dir) do
  --get stone at position
  _g_lib(g,n,ll,lg,v_add(g,d))
 end
 --return lists for recursive
 return ll,lg
end
--]]

---[[
function _g_lib(g,n,ll,lg,v)
 local a=g_goishi(v,lg)
 local o=a
 --goishi not already checked
 if(a==nil) o=g_goishi(v)
 --goishi not on the board
 if o~=nil then
  --if stone are enemies
  if g.col~=o.col then
   --doesn't need check next
   if(not n)return --stop
   local le=g_lib(o,false)
   --more than one liberty
   if(#le~=1)return --stop
   --goishi close liberty
   if v_equ(g,le[1]) then
    l_v_add(ll,v) --valid move
   end
  elseif a==nil then --new
   --get liberties from goishi
   g_lib(o,n,ll,lg)
  end
 --goishi already on the board
 elseif valid(v) then
  l_v_add(ll,v) --valid move
 end
end
--]]

-- liberties list --

---[[
--add vector to list
function l_v_add(l,v)
 for u in all(l) do
  if(v_equ(v,u))return --already
 end
 add(l,v)
end
--]]

-- vector --

---[[
--create a new vector
function nv(x,y)
 return {x=x,y=y}
end
--]]

---[[
--true if vec are equals
function v_equ(v1,v2)
 return v1.x==v2.x
    and v1.y==v2.y
end
--]]

---[[
--add v1 and v2 to make v3
function v_add(v1,v2)
 return {x=v1.x+v2.x,
         y=v1.y+v2.y}
end
--]]

---[[
--clone vector
function v_clone(v)
 return {x=v.x,y=v.y}
end
--]]

-- utils --

---[[
--true if btn has been pressed
function btnpress()
 return btnp(0) or btnp(1)
     or btnp(2) or btnp(3)
     or btnp(4) or btnp(5)
end
--]]

---[[
--create string based on number
function form(n)
 local s=""
 --add space
 if     n< 10 then s=s.."  "
 elseif n<100 then s=s.." " end
 return s..n
end
--]]

-- liberties direction --

---[[
_dir={
 nv( 0, 1),nv( 0,-1),
 nv( 1, 0),nv(-1, 0)
}
--]]
