SOURCEFILES=LispZeroGo.go zero.lisp zero-test.lisp zero-test.gold
BUILDFILES=Makefile
RESULTFILES=LispZeroGo zero-test.out

all: LispZeroGo zero-test

LispZeroGo: $(GOBIN)/LispZeroGo

$(GOBIN)/LispZeroGo: LispZeroGo.go
	go build

zero-test: zero-test.gold zero-test.out
	diff -u zero-test.gold zero-test.out

zero-test.out: LispZeroGo zero-test.lisp
	$(GOBIN)/LispZeroGo zero-test.lisp > zero-test.out

zero-new-gold: zero-test.out
	rm -f zero-test.gold
	cp zero-test.out zero-test.gold
	chmod a-w zero-test.gold

install: LispZeroGo
	go install

clean:
	rm -f $(RESULTFILES)

.PHONY: all LispZeroGo clean zero-test zero-new-gold install
