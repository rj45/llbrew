# TODO

This is mainly a parking lot for ideas / optimizations to do at a later date.
The code also has many `todo` comments / panics for things that need to be implemented.

- [ ] Allow const globals to have an offset that's assembled to global+offset
- [ ] Figure out how to eliminate adds to zero after the function logue is generated
    - Complicating factor is the add is in the target's instruction set
- [ ] Regalloc idea: when scanning connected moves, try to find a colour that doesn't interfere with any connected move (saved & not saved) and keep track of whether any of those values crosses a call
    - if it fails to find a move that's already coloured, try to pick a valid non-interfering colour that doesn't interferer with any connected moves either
    - Hopefully this will result in more moves being coalesced