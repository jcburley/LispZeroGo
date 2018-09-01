Performance Analysis of LispZeroGo versus LispZero

2018-08-31 jcburley: Hand-modified version of c2go-generated variant
of lisp-zero-single.c, with all Malloc() calls removed/replaced (with
new()/make()), among other "nativizations" performed to the code. (It
still doesn't pass zero-test.lisp, but does fine with the many-cons
tests.)

This version has interesting performance characteristics on my MacBook
Pro (pony):

    craig@pony:~/github/LispZero/perftests (master u=)$ xtime LispZeroGo -q 1M.out
    3.90u 0.10s 3.19r 185700kB LispZeroGo -q 1M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime LispZeroGo -q 1M.out
    4.03u 0.10s 3.25r 185948kB LispZeroGo -q 1M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime LispZeroGo -q 1M.out
    3.95u 0.09s 3.19r 185680kB LispZeroGo -q 1M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime ../lisp-zero-single -q 1M.out
    3.36u 0.14s 3.52r 334796kB ../lisp-zero-single -q 1M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime ../lisp-zero-single -q 1M.out
    3.21u 0.13s 3.36r 334796kB ../lisp-zero-single -q 1M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime ../lisp-zero-single -q 1M.out
    3.21u 0.13s 3.35r 334796kB ../lisp-zero-single -q 1M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime LispZeroGo -q 10M.out
    62.62u 1.41s 30.49r 1747796kB LispZeroGo -q 10M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime LispZeroGo -q 10M.out
    58.72u 1.37s 29.61r 1796368kB LispZeroGo -q 10M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime LispZeroGo -q 10M.out
    65.99u 1.52s 30.47r 1847908kB LispZeroGo -q 10M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime ../lisp-zero-single -q 10M.out
    33.49u 1.36s 34.89r 3437460kB ../lisp-zero-single -q 10M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime ../lisp-zero-single -q 10M.out
    33.02u 1.32s 34.37r 3437432kB ../lisp-zero-single -q 10M.out
    craig@pony:~/github/LispZero/perftests (master u=)$ xtime ../lisp-zero-single -q 10M.out
    32.66u 1.28s 33.98r 3437448kB ../lisp-zero-single -q 10M.out
    craig@pony:~/github/LispZero/perftests (master u=)$

Takes about half the memory, and a little less wall-clock time, but
about twice the CPU time for the large (10M-cons) input! Must be
making use of more than one thread, though my (original) code doesn't
explicitly do anything to cause that -- maybe an underlying-I/O thing?
