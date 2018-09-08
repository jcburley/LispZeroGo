SOURCEFILES=LispZeroGo.go zero.lisp zero-test.lisp zero-test.gold
BUILDFILES=Makefile
RESULTFILES=LispZeroGo zero-test.out

all: LispZeroGo zero-test

LispZeroGo: LispZeroGo.go
	go build

zero-test: zero-test.gold zero-test.out
	diff -u zero-test.gold zero-test.out

zero-test.out: LispZeroGo zero-test.lisp
	./LispZeroGo zero-test.lisp > zero-test.out

zero-new-gold: zero-test.out
	rm -f zero-test.gold
	cp zero-test.out zero-test.gold
	chmod a-w zero-test.gold

# 'go install' moves, rather than copies, the binary. This would force
# a re-build if 'make -k install' was run, even without changes to the
# dependencies (*.go) having been made. So, after 'go install', try to
# hardlink the executable back into the local directory, and if that
# fails, just copy it.
install: LispZeroGo zero-test
	if ! cmp -s LispZeroGo $(GOBIN)/LispZeroGo; then \
	  go install; \
	  cp -lv $(GOBIN)/LispZeroGo . || cp -pv $(GOBIN)/LispZeroGo .; \
	fi

clean:
	rm -f $(RESULTFILES)

.PHONY: all clean zero-test zero-new-gold install
