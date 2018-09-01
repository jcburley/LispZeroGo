Translated from LispZero's lisp-zero-single.c program by c2go and then
substantially hand-edited to use native GoLang facilities (among other
things), this version passes the same tests and has roughly the same
wallclock performance (while using about twice the CPU and half the
memory!) of lisp-zero-single.c on an input file containing 10M (cons
'a 'b) forms on my MacBook Pro.

This suggests that Joker, an implementation of an interactive Clojure
written in GoLang, will have reasonably decent performance as it's
optimized (even though it's not great right now) and won't hit a wall
due to being written in GoLang (which seems to be the problem plaguing
ClojureScript and especially the original Clojure, running on the
JVM).

I'll do more perforance testing on my Ubuntu box and Ubuntu VM (my
primary Internet-facing server) when I return from vacation, to
confirm this assessment.

But it's looking like good news, as it suggests a decent-performing
system can be constructed that runs scripts written in Joker (*.joke)
even very frequently for short periods of time, which is a popular
model, historically speaking, the the Unix community -- compared to
multi-GB Java/JVM processes that take several seconds to start, and
even compared to large Node.JS/v8 programs that take upwards of a full
second to start.

It's probably the case that productization of such a system would
require optimization of Joker itself -- but not a rewrite into C,
which I'd like to avoid!
