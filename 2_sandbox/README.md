~ Things that nee to run ~
Check clock.
Get absolute timing position according to the beacon spacing.
update main.position state
Check clock again.
If the time slot != last check, update 
else; wait

~~VVVV~~  the main program only alters with msg input
type model struct {
    keys
    help
    inputStyle
    lastKey
    position
}