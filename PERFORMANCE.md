Performance Analysis of LispZeroGo versus LispZero

Using: https://blog.golang.org/profiling-go-programs

Generating test cases via e.g.:

    $ lumo -e '(load-file "gen-big-list.clj") (gen-big-list 10000000)' > 10M.out
    
    $ lumo -e '(load-file "gen-big-list.clj") (println "(cond nil (quote (") (gen-big-list 1000000) (println ")))")' > D-1M.out

2018-08-31 jcburley: Hand-modified version of c2go-generated variant
of lisp-zero-single.c, with all Malloc() calls removed/replaced (with
new()/make()), among other "nativizations" performed to the code. (It
still doesn't pass zero-test.lisp, but does fine with the many-cons
tests.)

This version has interesting performance characteristics on my MacBook
Pro (pony):

    $ xtime LispZeroGo -q 1M.out
    3.90u 0.10s 3.19r 185700kB LispZeroGo -q 1M.out
    $ xtime LispZeroGo -q 1M.out
    4.03u 0.10s 3.25r 185948kB LispZeroGo -q 1M.out
    $ xtime LispZeroGo -q 1M.out
    3.95u 0.09s 3.19r 185680kB LispZeroGo -q 1M.out
    $ xtime ../lisp-zero-single -q 1M.out
    3.36u 0.14s 3.52r 334796kB ../lisp-zero-single -q 1M.out
    $ xtime ../lisp-zero-single -q 1M.out
    3.21u 0.13s 3.36r 334796kB ../lisp-zero-single -q 1M.out
    $ xtime ../lisp-zero-single -q 1M.out
    3.21u 0.13s 3.35r 334796kB ../lisp-zero-single -q 1M.out
    $ xtime LispZeroGo -q 10M.out
    62.62u 1.41s 30.49r 1747796kB LispZeroGo -q 10M.out
    $ xtime LispZeroGo -q 10M.out
    58.72u 1.37s 29.61r 1796368kB LispZeroGo -q 10M.out
    $ xtime LispZeroGo -q 10M.out
    65.99u 1.52s 30.47r 1847908kB LispZeroGo -q 10M.out
    $ xtime ../lisp-zero-single -q 10M.out
    33.49u 1.36s 34.89r 3437460kB ../lisp-zero-single -q 10M.out
    $ xtime ../lisp-zero-single -q 10M.out
    33.02u 1.32s 34.37r 3437432kB ../lisp-zero-single -q 10M.out
    $ xtime ../lisp-zero-single -q 10M.out
    32.66u 1.28s 33.98r 3437448kB ../lisp-zero-single -q 10M.out
    $

Takes about half the memory, and a little less wall-clock time, but
about twice the CPU time for the large (10M-cons) input! Must be
making use of more than one thread, though my (original) code doesn't
explicitly do anything to cause that -- maybe an underlying-I/O thing?

Added a -inbufsize option, tried 4MB, but no real differences
observed.

2018-09-01 jcburley: Further "nativization" and cleanup has resulted
in passing the zero-test.lisp test, suggesting that replacing '*byte'
with 'string' type (perhaps for the widely-passed 'what' variable?)
helps the garbage collector avoid freeing stuff prematurely, or some
such thing.

Performance seems about the same:

    $ xtime LispZeroGo -inbufsize $((1024 * 1024 * 4)) -q 10M.out
    56.54u 1.40s 28.38r 1916120kB LispZeroGo -inbufsize 4194304 -q 10M.out
    $ xtime LispZeroGo -inbufsize $((1024 * 1024 * 4)) -q 10M.out
    68.15u 1.37s 30.27r 1857784kB LispZeroGo -inbufsize 4194304 -q 10M.out
    $ xtime LispZeroGo -inbufsize $((1024 * 1024 * 4)) -q 10M.out
    57.58u 1.27s 31.79r 1756000kB LispZeroGo -inbufsize 4194304 -q 10M.out
    $ xtime ../lisp-zero-single -q 10M.out
    34.12u 1.46s 35.66r 3283916kB ../lisp-zero-single -q 10M.out
    $ xtime ../lisp-zero-single -q 10M.out
    33.70u 1.41s 35.15r 3437456kB ../lisp-zero-single -q 10M.out
    $ xtime ../lisp-zero-single -q 10M.out
    33.62u 1.40s 35.07r 3437472kB ../lisp-zero-single -q 10M.out
    $
