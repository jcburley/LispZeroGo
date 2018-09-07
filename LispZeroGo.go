package main

import "bufio"
import "bytes"
import "flag"
import "fmt"
import "github.com/elliotchance/c2go/linux"
import "github.com/elliotchance/c2go/noarch"
import "github.com/pkg/profile"
import "io"
import "os"
import "runtime"
import "runtime/pprof"
import "strings"
import "unsafe"

const _ISspace uint16 = ((1 << 5) << 8)

var stdin *os.File
var stdout *bufio.Writer
var stderr *bufio.Writer

func my_assert(t *byte, w *byte, l uint32, x *byte) {
	if *t == 0 {
		return
	}
	stderr.Flush()
	stdout.Flush()
	linux.AssertFail(t, w, l, x)
}

var profiler string
var cpuprofile string
var quiet bool
var tracing bool

var allocations uint64 = 0

func nl(s *bufio.Writer) {
	s.WriteString("\n")
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type compiled_fn func(string, *Object_s, *Object_s) *Object_s
type Symbol_s struct {
	n string
}
type Cdr_u struct{ memory unsafe.Pointer }

func (unionVar *Cdr_u) copy() Cdr_u {
	var buffer [8]byte
	for i := range buffer {
		buffer[i] = (*((*[8]byte)(unionVar.memory)))[i]
	}
	var newUnion Cdr_u
	newUnion.memory = unsafe.Pointer(&buffer)
	return newUnion
}
func (unionVar *Cdr_u) obj() **Object_s {
	if unionVar.memory == nil {
		var buffer [8]byte
		unionVar.memory = unsafe.Pointer(&buffer)
	}
	return (**Object_s)(unionVar.memory)
}
func (unionVar *Cdr_u) p_obj() **Object_s {
	return (**Object_s)(unsafe.Pointer(&unionVar.memory))
}
func (unionVar *Cdr_u) sym() **Symbol_s {
	if unionVar.memory == nil {
		var buffer [8]byte
		unionVar.memory = unsafe.Pointer(&buffer)
	}
	return (**Symbol_s)(unionVar.memory)
}
func (unionVar *Cdr_u) fn() *compiled_fn {
	if unionVar.memory == nil {
		var buffer [8]byte
		unionVar.memory = unsafe.Pointer(&buffer)
	}
	return (*compiled_fn)(unionVar.memory)
}

type Object_s struct {
	car *Object_s
	cdr Cdr_u
}
type Object = Object_s

var atomic Object
var p_atomic *Object_s = (*Object_s)(unsafe.Pointer(&atomic))
var compiled Object
var p_compiled *Object_s = (*Object_s)(unsafe.Pointer(&compiled))
var p_environment *Object_s = nil
var p_nil *Object_s = nil
var p_quote *Object_s = nil
var p_atom *Object_s = nil
var p_eq *Object_s = nil
var p_cons *Object_s = nil
var p_car *Object_s = nil
var p_cdr *Object_s = nil
var p_cond *Object_s = nil
var p_eval *Object_s = nil
var p_apply *Object_s = nil
var p_defglobal *Object_s = nil
var p_dot_symbol_dump *Object_s = nil

// nilp - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:401
/* Forward references. */ //
/* Objects (lists, atoms, etc.). */ //
/* Head item in this list, unless a BUILTIN_CAR node. */ //
/* Tail list, unless car is a BUILTIN_CAR node. */ //
/* If CAR of an Object is not one of these, it is a List. */ //
/* Object.cdr is a struct Symbol_s *. */ //
/* Object.cdr is a compiled_fn. */ //
/* Current evaluation environment (list of bindings). */ //
/* Input/output as (). */ //
/* Also input as '. */ //
//
func nilp(list *Object_s) int8 {
	return int8((int8(map[bool]int32{false: 0, true: 1}[int64(uintptr(unsafe.Pointer(list))) == int64(uintptr(unsafe.Pointer(p_nil)))])))
}

// atomp - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:407
/* Whether object is an atom in the traditional Lisp sense */ //
//
func atomp(list *Object_s) int8 {
	return int8((int8(map[bool]int32{false: 0, true: 1}[int32(int8((nilp(list)))) != 0 || int64(uintptr(unsafe.Pointer((*list).car))) == int64(uintptr(unsafe.Pointer(p_atomic)))])))
}

// atomicp - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:413
/* Whether object is an atom in this implementation */ //
//
func atomicp(list *Object_s) int8 {
	return int8((int8(map[bool]int32{false: 0, true: 1}[int8((noarch.NotInt8(nilp(list)))) != 0 && int64(uintptr(unsafe.Pointer((*list).car))) == int64(uintptr(unsafe.Pointer(p_atomic)))])))
}

// compiledp - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:418
func compiledp(list *Object_s) int8 {
	return int8((int8(map[bool]int32{false: 0, true: 1}[int8((noarch.NotInt8(atomp(list)))) != 0 && int64(uintptr(unsafe.Pointer((*list).car))) == int64(uintptr(unsafe.Pointer(p_compiled)))])))
}

// listp - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:423
func listp(list *Object_s) int8 {
	return int8((int8(map[bool]int32{false: 0, true: 1}[int8((noarch.NotInt8(atomp(list)))) != 0 && int8((noarch.NotInt8(compiledp(list)))) != 0])))
}

// finalp - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:428
func finalp(list *Object_s) int8 {
	return int8((int8(map[bool]int32{false: 0, true: 1}[int32(int8((listp(list)))) != 0 && int32(int8((nilp((*(*list).cdr.obj()))))) != 0])))
}

var filename string
var lineno uint32 = uint32(int32(1))
var max_object_write int32 = -int32(1)

func assert_or_dump(srcline uint32, ok int8, obj *Object_s, what string) {
	if ok != 0 || max_object_write != -1 {
		return
	}
	fmt.Fprintf(stderr, "ERROR at %d: %s, but got:\n", lineno, what)
	max_object_write = 10
	object_write(stderr, obj)
	fmt.Fprintf(stderr, "\nEnvironment:\n")
	object_write(stderr, p_environment)
	fmt.Fprintf(stderr, "\n/home/craig/github/LispZero/lisp-zero-single.c:%d: aborting\n", srcline)
	stderr.Flush()
	stdout.Flush()
	panic("assertion failure")
}

// list_car - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:456
func list_car(list *Object_s) *Object_s {
	assert_or_dump(uint32(int32(458)), (listp(list)), (list), "expected list")
	return (*list).car
}

// list_cdr - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:462
func list_cdr(list *Object_s) *Object_s {
	assert_or_dump(uint32(int32(464)), (listp(list)), (list), "expected list")
	return (*(*list).cdr.obj())
}

// object_symbol - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:468
func object_symbol(atom *Object_s) *Symbol_s {
	assert_or_dump(uint32(int32(470)), (atomicp(atom)), (atom), "expected implementation atom")
	return (*(*atom).cdr.sym())
}

// object_compiled - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:474
func object_compiled(compiled *Object_s) compiled_fn {
	assert_or_dump(uint32(int32(476)), (compiledp(compiled)), (compiled), "expected compiled function")
	return compiled_fn((*(*compiled).cdr.fn()))
}

// object_new - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:480
func object_new(car *Object_s, cdr *Object_s) *Object_s {
	var obj *Object_s
	new_obj := new(Object_s)
	obj = (*Object_s)(unsafe.Pointer(new_obj))
	allocations += 1
	(*obj).car = car
	(*(*obj).cdr.obj()) = cdr
	return obj
}

// object_new_compiled - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:497
func object_new_compiled(fn compiled_fn) *Object_s {
	var obj *Object_s
	new_obj := new(Object_s)
	obj = (*Object_s)(unsafe.Pointer(new_obj))
	allocations += 1
	(*obj).car = p_compiled
	(*(*obj).cdr.fn()) = fn
	return obj
}

/* Change to key on a 'string' type. This necessitated changing all
/* callers of symbol_sym() to pass a 'string' rather than '*byte'
/* type. That fixed some things but things still don't really work. */

type Symbol_MAP =
map[string]*Symbol_s

var map_sym Symbol_MAP

/* Map of symbols (keys) to values. */
func symbol_lookup(name string) *Symbol_s {
	sym, found := map_sym[name]
	if found {
		return sym
	} else {
		return (nil)
	}
}

func symbol_name(s *Symbol_s) string {
	return s.n
}

func symbol_name_as_bytes(s *Symbol_s) *byte {
	return (*byte)(unsafe.Pointer(&s.n))
}

var symbol_strdup int8 = int8((int8(int32(1))))

// symbol_sym - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:551
func symbol_sym(name string) *Symbol_s {
	var sym *Symbol_s = symbol_lookup(name)
	if sym != nil {
		return sym
	}
	sym = new(Symbol_s)
	sym.n = name
	map_sym[name] = sym
	return sym
}

func symbol_dump() {
	for key, value := range map_sym {
		fmt.Printf("%s -> %p\n", key, value)
	}
}

var p_sym_t *Symbol_s = nil
var p_sym_quote *Symbol_s = nil

// binding_new - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:581
/* Environment (bindings). */ //
//
func binding_new(sym *Object_s, val *Object_s) *Object_s {
	assert_or_dump(uint32(int32(583)), (atomicp(sym)), (sym), "expected implementation atom")
	return object_new(sym, val)
}

// binding_lookup - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:594
/* Bindings; each binding is either an atom (meaning its symbol is
   explicitly unbound) or a key/value cons (the symbol is in the car,
   its binding is in the cdr). */ //
/* Originally this used a recursive algorithm, but tail-recursion
   optimization wasn't being done by gcc -g -O0, and it was annoying
   to find oneself inside such deep stack traces during debugging.
*/ //
//
func binding_lookup(what string, key *Symbol_s, bindings *Object_s) *Object_s {
	if int8((nilp(bindings))) != 0 {
		return p_nil
	}
	if tracing {
		fmt.Fprintf(stderr, "%s:%d: Searching for `%s' in:\n", filename, lineno, (*key).n)
		max_object_write = int32(10)
		object_write(stderr, bindings)
		max_object_write = -int32(1)
		io.WriteString(stderr,"\n\n")
	}
	for ; int8((noarch.NotInt8(nilp(bindings)))) != 0; bindings = list_cdr(bindings) {
		assert_or_dump(uint32(int32(616)), (listp(bindings)), (bindings), "expected list")
		var binding *Object_s = list_car(bindings)
		if int32(int8((atomicp(binding)))) != 0 && int64(uintptr(unsafe.Pointer(object_symbol(binding)))) == int64(uintptr(unsafe.Pointer(key))) {
			return p_nil
		}
		{
			var symbol *Object_s = list_car(binding)
			if int32(int8((atomicp(symbol)))) != 0 && int64(uintptr(unsafe.Pointer(object_symbol(symbol)))) == int64(uintptr(unsafe.Pointer(key))) {
				return binding
			}
		}
	}
	return p_nil
}

var token_lookahead string
var lookahead_valid bool = false

func buffer_append(buf *bytes.Buffer, ch byte) {
	err := buf.WriteByte(ch)
	check(err)
}

func buffer_to_string(buf *bytes.Buffer) string {
	return string(buf.Bytes())
}

// token_putback - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:675
func token_putback(token string) {
	if lookahead_valid {
		my_assert((&[]byte("lookahead_valid is true!\x00")[0]), (&[]byte("/home/craig/github/LispZero/lisp-zero-single.c\x00")[0]), uint32(int32(677)), (&[]byte("void print_number(int *)\x00")[0]))
	}
	token_lookahead = token
	lookahead_valid = true
}

func my_getc(input *bufio.Reader) int32 {
	b, err := input.ReadByte()
	if err == io.EOF {
		return -1
	}
	check(err)
	return (int32)(b)
}

func my_ungetc(ch int32, input *bufio.Reader) {
	input.UnreadByte()
}

// token_get - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:707
func token_get(input *bufio.Reader, buf *bytes.Buffer) string {
	var ch int32

	buf.Reset()

	if lookahead_valid {
		lookahead_valid = false
		return token_lookahead
	}
	for {
		if (func() int32 {
			tempVar := my_getc(input)
			ch = tempVar
			return tempVar
		}()) == -int32(1) {
			my_exit(0)
		}
		if ch == int32(';') {
			for (func() int32 {
				tempVar := my_getc(input)
				ch = tempVar
				return tempVar
			}()) != -int32(1) && ch != int32('\n') {
			}
		}
		if ch == int32('\n') {
			lineno += 1
		}
		if noarch.NotInt32((int32(*((*uint16)(func() unsafe.Pointer {
			tempVar := (*linux.CtypeLoc())
			return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)((ch))*unsafe.Sizeof(*tempVar))
		}()))) & int32(uint16(_ISspace)))) != 0 {
			break
		}
	}
	buffer_append(buf, byte(ch))
	if noarch.Strchr((&[]byte("()'\x00")[0]), ch) != nil {
		return buffer_to_string(buf)
	}
	for {
		if (func() int32 {
			tempVar := my_getc(input)
			ch = tempVar
			return tempVar
		}()) == -int32(1) {
			my_exit(0)
		}
		if noarch.Strchr((&[]byte("()'\x00")[0]), ch) != nil || int32(*((*uint16)(func() unsafe.Pointer {
			tempVar := (*linux.CtypeLoc())
			return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)((ch))*unsafe.Sizeof(*tempVar))
		}())))&int32(uint16(_ISspace)) != 0 {
			my_ungetc(ch, input)
			return buffer_to_string(buf)
		}
		buffer_append(buf, byte(ch))
	}
}

var latest_lineno uint32

func cstr_to_string(token *byte) string {
	var tokbuf strings.Builder
	var p unsafe.Pointer = unsafe.Pointer(token)
	for {
		if (*((*byte) (p))) == 0 {
			break
		}
		tokbuf.WriteByte(*((*byte) (p)))
		p = unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 1)
	}
	return tokbuf.String()
}

// object_read - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:751
func object_read(input *bufio.Reader, buf *bytes.Buffer) *Object_s {
	token := token_get(input, buf)
	if token == "(" {
		return list_read(input, buf)
	}
	if token == "'" {
		var tmp *Object_s = object_read(input, buf)
		return object_new(object_new(p_atomic, (*Object_s)(unsafe.Pointer((p_sym_quote)))), object_new(tmp, p_nil))
	}
	if token == ")" {
		fmt.Fprintf(stderr, "unbalanced close paren\n")
		my_exit(1)
	}
	if tracing && lineno != latest_lineno {
		latest_lineno = lineno
		fmt.Fprintf(stderr, "%s:%d: Seen `%s'.\n", filename, lineno, token)
		stderr.Flush()
	}
	return object_new(p_atomic, (*Object_s)(unsafe.Pointer((symbol_sym(token)))))
}

// list_read - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:781
/* Make sure we first read the object before going on to read the rest of the list. */ //
//
func list_read(input *bufio.Reader, buf *bytes.Buffer) *Object_s {
	var first *Object_s = p_nil
	var next **Object_s = &first

	var cur *Object_s
	for {
		var token string = token_get(input, buf)
		if token == ")" {
			cur = p_nil
			break
		}
		if token == "." {
			cur = object_read(input, buf)
			if token_get(input, buf) != ")" {
				if token_get(input, buf) != ")" {
					fmt.Fprintf(stderr, "missing close parenthese for simple list\n")
					my_exit(3)
				}
			}
			break
		}

		token_putback(token)

		cur = object_new(object_read(input, buf), nil)
		*next = cur
		next = (cur.cdr.obj())
	}

	return first
}

// quotep - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:806
/* Output. */ //
/* true if object is (quote arg) */ //
/* TODO: Decide whether this look up quote in the current env to do the check */ //
//
func quotep(obj *Object_s) (c2goDefaultReturn int8) {
	if int8((noarch.NotInt8(listp(obj)))) != 0 || int8((noarch.NotInt8(finalp(list_cdr(obj))))) != 0 {
		return int8((int8(int32(0))))
	}
	{
		var car *Object_s = list_car(obj)
		return int8((int8(map[bool]int32{false: 0, true: 1}[int32(int8((compiledp(car)))) != 0 && int64(uintptr(noarch.CastInterfaceToPointer((func(string, *Object_s, *Object_s) *Object_s)((object_compiled(car)))))) == int64(uintptr(noarch.CastInterfaceToPointer(f_quote)))])))
	}
	return
}

/* TODO: Print name of function. */
func object_write(output *bufio.Writer, obj *Object_s) {
	stderr.Flush()
	if int8((nilp(obj))) != 0 {
		fmt.Fprintf(output, "()")
		return
	}
	if int8((atomicp(obj))) != 0 {
		fmt.Fprintf(output, "%s", symbol_name(object_symbol(obj)))
		return
	}
	if int8((compiledp(obj))) != 0 {
		fmt.Fprintf(output, "*COMPILED*")
		return
	}
	if int8((quotep(obj))) != 0 {
		fmt.Fprintf(output, "'")
		object_write(output, list_car(list_cdr(obj)))
		return
	}
	if max_object_write == int32(0) {
		fmt.Fprintf(output, "(...)")
		return
	}
	if max_object_write > int32(0) {
		max_object_write -= 1
	}
	fmt.Fprintf(output, "(")
	for {
		object_write(output, list_car(obj))
		if int8((finalp(obj))) != 0 {
			fmt.Fprintf(output, ")")
			return
		}
		obj = list_cdr(obj)
		if int8((noarch.NotInt8(listp(obj)))) != 0 {
			fmt.Fprintf(output, " . ")
			object_write(output, obj)
			fmt.Fprintf(output, ")")
			return
		}
		fmt.Fprintf(output, " ")
	}
}

// binding_for - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:882
/* Evaluation */ //
/* TODO: Throw an exception etc. */ //
//
func binding_for(what string, sym *Symbol_s, env *Object_s) *Object_s {
	var tmp *Object_s
	tmp = binding_lookup(what, sym, env)
	if int8((nilp(tmp))) != 0 {
		assert_or_dump(908, ^nilp(tmp), env, fmt.Sprintf("Unbound symbol `%s'", symbol_name(sym)))
	}
	return list_cdr(tmp)
}

// eval - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:901
/* Does not support traditional lambdas or labels; just the built-ins and
   our "unique" apply.  */ //
//
func eval(what string, exp *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	if int32(int8((nilp(exp)))) != 0 || int32(int8((compiledp(exp)))) != 0 {
		return exp
	}
	if int8((atomicp(exp))) != 0 {
		return binding_for(what, object_symbol(exp), env)
	}
	assert_or_dump(uint32(int32(909)), (listp(exp)), (exp), "expected list")
	{
		var func_ *Object_s = eval(what, list_car(exp), env)
		var forms *Object_s = list_cdr(exp)
		if int8((atomp(list_car(exp)))) != 0 {
			what = symbol_name(object_symbol(list_car(exp)))
		}
		if int8((compiledp(func_))) != 0 {
			var fn compiled_fn
			fn = object_compiled(func_)
			return fn(what, forms, env)
		}
		return apply(what, func_, func_, forms, env)
	}
	return
}

// assert_zedbap - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:963
/* Need to solve these problems:

   (defun x (a) (progn (setq a 6) (a)))
   (setq a 5)
   (x a)

   It should return 5, but if defun is implemented naively, it
   might return 6.

   (defun y (a b c) (cons (a b c)))
   (setq a 5)
   (setq b 6)
   (setq c 7)
   (y c a b)

   That should return (7 5 6).

   Consider: Make the built-in form work as much like a compiled form
   as possible.  The compiled form is passed in an arglist and the
   current environment, allowing it to manipulate both to a high
   degree.  The compiler lexically binds those arguments to parameter
   names, so an instance of the compiled function doesn't put those
   parameter names into the environment.  So, either implement some
   form of lexical binding of names here, or consider a nameless
   solution such as having a particular form eval to the argument list
   and another to the environment.  This means explicitly supporting
   an eval() function that, unlike traditional Lisp, takes an env
   argument.
*/ //
//
func assert_zedbap(zedba *Object_s) {
	assert_or_dump(uint32(int32(965)), (listp(zedba)), (zedba), "expected list")
	assert_or_dump(uint32(int32(967)), (listp(list_car(zedba))), (zedba), "expected list with car being arglist")
	assert_or_dump(uint32(int32(968)), (atomicp(list_car(list_car(zedba)))), (zedba), "expected zedba with 1st arg being mename")
	assert_or_dump(uint32(int32(969)), (atomicp(list_car(list_cdr(list_car(zedba))))), (zedba), "expected zedba with 2nd arg being formlistparamname")
	assert_or_dump(uint32(int32(970)), (atomicp(list_car(list_cdr(list_cdr(list_car(zedba)))))), (zedba), "expected zedba with 3rd arg being envparamname")
	assert_or_dump(uint32(int32(971)), (finalp(list_cdr(list_cdr(list_car(zedba))))), (zedba), "expected zedba with only 3 args")
	assert_or_dump(uint32(int32(973)), (finalp(list_cdr(zedba))), (zedba), "expected zedba body to be last element of zedba as list")
}

// apply - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:989
/* Apply a zedba, which is an self/arglist/env macro version of the
   classic lambda.  The form of a zedba is

   ((mename formlistparamname envparamname) body)

   where body is to be evaluated with formlistparamname bound to the
   list of forms of the invocation, envparamname to the environment
   for the evaluation of the containing expression, and mename bound
   to the zedba itself (for easy recursive references).  Note that the
   caller might want to pass something else to be bound to mename, in
   case this proves useful (e.g. a limit on the # of recursive
   invocations could be implemented this way), so this is allowed.  */ //
//
func apply(what string, func_ *Object_s, me *Object_s, forms *Object_s, env *Object_s) *Object_s {
	var meparamname *Object_s
	var formlistparamname *Object_s
	var envparamname *Object_s
	assert_zedbap(func_)
	assert_zedbap(me)
	{
		var params *Object_s = list_car(func_)
		assert_or_dump(uint32(int32(1001)), (listp(params)), (params), "expected list")
		meparamname = list_car(params)
		assert_or_dump(uint32(int32(1005)), (listp(list_cdr(params))), (params), "expected 2-element list")
		formlistparamname = list_car(list_cdr(params))
		assert_or_dump(uint32(int32(1009)), (finalp(list_cdr(list_cdr(params)))), (params), "expected 2-element list")
		envparamname = list_car(list_cdr(list_cdr(params)))
	}
	return eval(what, list_car(list_cdr(func_)), (object_new(binding_new((meparamname), (func_)), (object_new(binding_new((formlistparamname), (forms)), (object_new(binding_new((envparamname), (env)), (env))))))))
}

// f_quote - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1022
/* (quote form) => form */ //
//
func f_quote(what string, args *Object_s, env *Object_s) *Object_s {
	assert_or_dump(uint32(int32(1024)), (finalp(args)), (args), "expected 1-element list")
	return list_car(args)
}

// f_atom - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1030
/* (atom atom) => t if atom is an atom (including nil), nil otherwise */ //
//
func f_atom(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	assert_or_dump(uint32(int32(1032)), (finalp(args)), (args), "expected 1-element list")
	{
		var arg *Object_s = eval(what, list_car(args), env)
		return func() *Object_s {
			if int32(int8((atomp(arg)))) != 0 {
				return object_new(p_atomic, (*Object_s)(unsafe.Pointer((p_sym_t))))
			} else {
				return p_nil
			}
		}()
	}
	return
}

// f_eq - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1042
/* (eq left-atom right-atom) => t if args are equal, nil otherwise */ //
/* All nils are equal to each other in this implementation */ //
//
func f_eq(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	assert_or_dump(uint32(int32(1044)), (listp(args)), (args), "expected 1-element list")
	assert_or_dump(uint32(int32(1045)), (finalp(list_cdr(args))), (args), "expected 1-element list")
	{
		var left *Object_s = eval(what, list_car(args), env)
		var right *Object_s = eval(what, list_car(list_cdr(args)), env)
		assert_or_dump(uint32(int32(1051)), (atomp(left)), (left), "expected Lisp atom")
		assert_or_dump(uint32(int32(1052)), (atomp(right)), (right), "expected Lisp atom")
		if int64(uintptr(unsafe.Pointer(left))) == int64(uintptr(unsafe.Pointer(right))) {
			return object_new(p_atomic, (*Object_s)(unsafe.Pointer((p_sym_t))))
		}
		if int32(int8((nilp(left)))) != 0 || int32(int8((nilp(right)))) != 0 {
			return p_nil
		}
		return func() *Object_s {
			if int64(uintptr(unsafe.Pointer(object_symbol(left)))) == int64(uintptr(unsafe.Pointer(object_symbol(right)))) {
				return object_new(p_atomic, (*Object_s)(unsafe.Pointer((p_sym_t))))
			} else {
				return p_nil
			}
		}()
	}
	return
}

// f_cons - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1065
/* (cons car-arg cdr-arg) => (car-arg cdr-arg) */ //
//
func f_cons(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	assert_or_dump(uint32(int32(1067)), (listp(args)), (args), "expected arglist")
	assert_or_dump(uint32(int32(1068)), (finalp(list_cdr(args))), (args), "expected only two args")
	{
		var car *Object_s = eval(what, list_car(args), env)
		var cdr *Object_s = eval(what, list_car(list_cdr(args)), env)
		return object_new(car, cdr)
	}
	return
}

// f_car - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1079
/* (car cons-arg) : cons-arg is a list => car of cons-arg */ //
//
func f_car(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	assert_or_dump(uint32(int32(1081)), (finalp(args)), (args), "expected WHAT??")
	{
		var arg *Object_s = eval(what, list_car(args), env)
		assert_or_dump(uint32(int32(1086)), (listp(arg)), (arg), "expected WHAT??")
		return list_car(arg)
	}
	return
}

// f_cdr - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1093
/* (cdr cons-arg) : cons-arg is a list => cdr of cons-arg */ //
//
func f_cdr(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	assert_or_dump(uint32(int32(1095)), (finalp(args)), (args), "expected WHAT??")
	{
		var arg *Object_s = eval(what, list_car(args), env)
		assert_or_dump(uint32(int32(1100)), (listp(arg)), (arg), "expected WHAT??")
		return list_cdr(arg)
	}
	return
}

// f_cond - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1110
/* (cond ifthen-args ...) : each ifthen-args is an ifthen-pair; each
   ifthen-pair is a list of form (if-arg then-form) => eval(then-form)
   for the first if-arg in the list that is not nil (true), otherwise
   nil. */ //
//
func f_cond(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	if int8((nilp(args))) != 0 {
		return p_nil
	}
	assert_or_dump(uint32(int32(1115)), (listp(args)), (args), "expected WHAT??")
	{
		var pair *Object_s = list_car(args)
		assert_or_dump(uint32(int32(1120)), (listp(pair)), (pair), "expected WHAT??")
		assert_or_dump(uint32(int32(1121)), (finalp(list_cdr(pair))), (pair), "expected WHAT??")
		{
			var if_arg *Object_s = list_car(pair)
			var then_form *Object_s = list_car(list_cdr(pair))
			if int8((noarch.NotInt8(nilp(eval(what, if_arg, env))))) != 0 {
				return eval(what, then_form, env)
			}
			return f_cond(what, list_cdr(args), env)
		}
	}
	return
}

// f_defglobal - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1140
/* (defglobal) => global environment
   (defglobal newenv) => '() with newenv as the new global environment (SIDE EFFECT)
   (defglobal key value) => '() with new global environment prepended (via cons)
   			    with (key . value) (SIDE EFFECT)
*/ //
/* This form allows direct replacement of global environment.
   E.g. (defglobal (cdr (defglobal))) pops off the top binding. */ //
//
func f_defglobal(what string, args *Object_s, env *Object_s) *Object_s {
	if int8((nilp(args))) != 0 {
		return p_environment
	}
	assert_or_dump(uint32(int32(1145)), (listp(args)), (args), "expected WHAT??")
	if int8((nilp(list_cdr(args)))) != 0 {
		p_environment = eval(what, list_car(args), env)
	} else {
		assert_or_dump(uint32(int32(1153)), (finalp(list_cdr(args))), (args), "expected WHAT??")
		{
			var sym *Object_s = eval(what, list_car(args), env)
			var form *Object_s = eval(what, list_car(list_cdr(args)), env)
			assert_or_dump(uint32(int32(1159)), (atomicp(sym)), (sym), "expected WHAT??")
			p_environment = object_new(binding_new((object_new(p_atomic, (*Object_s)(unsafe.Pointer((object_symbol(sym)))))), (form)), (p_environment))
		}
	}
	return p_nil
}

// f_eval - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1171
/* (eval arg [env]) => arg evaluated with respect to environment env (default is current env) */ //
/* Eval this early, rather than in the final eval() below, so .c
   version .go compilers don't choose different order of evaluations
   and so mess up the tracefiles. */ //
//
func f_eval(what string, args *Object_s, env *Object_s) *Object_s {
	assert_or_dump(uint32(int32(1173)), (listp(args)), (args), "expected WHAT??")
	assert_or_dump(uint32(int32(1174)), int8((int8((map[bool]int32{false: 0, true: 1}[int32(int8((nilp(list_cdr(args))))) != 0 || int32(int8((finalp(list_cdr(args))))) != 0])))), (args), "expected WHAT??")
	var n_env *Object_s = func() *Object_s {
		if int32(int8((nilp(list_cdr(args))))) != 0 {
			return env
		} else {
			return eval(what, list_car(list_cdr(args)), env)
		}
	}()
	return eval(what, eval(what, list_car(args), env), n_env)
}

// f_apply - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1187
/* (apply zedba me forms [env]) => zedba invoked with reference to
   (presumably) itself, forms to be bound to zedba's arguments, and
   environment for such bindings (default is current env) */ //
//
func f_apply(what string, args *Object_s, env *Object_s) (c2goDefaultReturn *Object_s) {
	assert_or_dump(uint32(int32(1189)), (listp(args)), (args), "expected WHAT??")
	{
		var func_ *Object_s = eval(what, list_car(args), env)
		var rest *Object_s = list_cdr(args)
		if int8((atomp(list_car(args)))) != 0 {
			what = symbol_name(object_symbol(list_car(args)))
		}
		assert_or_dump(uint32(int32(1200)), (listp(rest)), (rest), "expected WHAT??")
		{
			var me *Object_s = eval(what, list_car(rest), env)
			var new_rest *Object_s = list_cdr(rest)
			rest = new_rest
			assert_or_dump(uint32(int32(1207)), (listp(rest)), (rest), "expected WHAT??")
			{
				var forms *Object_s = eval(what, list_car(rest), env)
				var new_rest *Object_s = list_cdr(rest)
				rest = new_rest
				assert_or_dump(uint32(int32(1214)), int8((int8((map[bool]int32{false: 0, true: 1}[int32(int8((nilp(rest)))) != 0 || int32(int8((finalp(rest)))) != 0])))), (rest), "expected WHAT??")
				return apply(what, func_, me, forms, func() *Object_s {
					if int32(int8((nilp(rest)))) != 0 {
						return p_nil
					} else {
						return eval(what, list_car(rest), env)
					}
				}())
			}
		}
	}
	return
}

// f_dot_symbol_dump - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1223
/* (.symbol_dump) : dump symbol names along with their struct Symbol_s * objects */ //
//
func f_dot_symbol_dump(what string, args *Object_s, env *Object_s) *Object_s {
	assert_or_dump(uint32(int32(1225)), int8((int8((map[bool]int32{false: 0, true: 1}[args == nil])))), (args), "expected WHAT??")
	symbol_dump()
	return p_nil
}

// initialize_builtin - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1232
func initialize_builtin(sym string, fn compiled_fn) *Object_s {
	var tmp *Object_s
	p_environment = object_new(binding_new((object_new(p_atomic, (*Object_s)(unsafe.Pointer((symbol_sym(sym)))))), (func() *Object_s {
		tmp = object_new_compiled(compiled_fn(fn))
		return tmp
	}())), (p_environment))
	return tmp
}

// initialize - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1243
/* TODO: Decide on a better name. */ //
//
func initialize() {
	p_sym_t = symbol_sym("t")
	p_sym_quote = symbol_sym("quote")
	symbol_strdup = int8((int8(int32(0))))
	p_quote = initialize_builtin("quote", f_quote)
	p_atom = initialize_builtin("atom", f_atom)
	p_eq = initialize_builtin("eq", f_eq)
	p_cons = initialize_builtin("cons", f_cons)
	p_car = initialize_builtin("car", f_car)
	p_cdr = initialize_builtin("cdr", f_cdr)
	p_cond = initialize_builtin("cond", f_cond)
	p_eval = initialize_builtin("eval", f_eval)
	p_apply = initialize_builtin("apply", f_apply)
	p_defglobal = initialize_builtin("defglobal", f_defglobal)
	p_dot_symbol_dump = initialize_builtin(".symbol_dump", f_dot_symbol_dump)
	symbol_strdup = int8((int8(int32(1))))
}

var prof interface { Stop() }
var inBufSize int = 4096  // Default on my MacBook Pro 2018-08-31

func main() {
	flag.Parse()
	if cpuprofile != "" {
		switch profiler {
		case "pkg/profile":
			prof = profile.Start(profile.ProfilePath(cpuprofile))
			defer finish()
		case "runtime/pprof":
			f, err := os.Create(cpuprofile)
			check(err)
			runtime.SetCPUProfileRate(500)
			pprof.StartCPUProfile(f)
			fmt.Fprintf(stderr, "Profiling started. See file `%s'.\n", cpuprofile);
			stderr.Flush()
			defer finish()
		default:
			fmt.Fprintf(stderr, "Unrecognized profiler: %s\n  Use 'pkg/profile' or 'runtime/pprof'.\n", profiler);
			os.Exit(96)
		}
	}
	var in *bufio.Reader
	if len(flag.Args()) == 1 {
		filename = flag.Arg(0)
		unbuf_in, err := os.Open(filename)
		check(err)
		in = bufio.NewReaderSize(unbuf_in, inBufSize)  // Get a buffered Reader
	} else if len(flag.Args()) == 0 {
		filename = "<stdin>"
		in = bufio.NewReaderSize(stdin, inBufSize)
	} else {
		fmt.Fprintf(stderr, "Excess command-line arguments starting with: %s\n", flag.Arg(1))
		my_exit(97)
	}

	if !quiet {
		fmt.Fprintf(stderr, "Underlying input buffer size: %d\n", in.Size())
		stderr.Flush()
	}
	
	map_sym = make(Symbol_MAP)
	var buf *bytes.Buffer = new(bytes.Buffer)

	initialize()

	for {
		var obj *Object_s = object_read(in, buf)
		if !quiet {
			object_write(stdout, eval(filename, obj, p_environment)); nl(stdout)
			stdout.Flush()
		}
	}

	my_exit(0)
}

func finish() {
	if prof != nil {
		prof.Stop()
	} else if cpuprofile != "" {
		pprof.StopCPUProfile()
		fmt.Fprintf(stderr, "Profiling stopped. See file `%s'.\n", cpuprofile);
	}
	if !quiet {
		fmt.Fprintf(stderr, "allocations: %d\n", allocations)
	}
	stderr.Flush()
	stdout.Flush()
}

func my_exit(rc int) {
	os.Exit(rc)
}

// debug_output - transpiled function from  /home/craig/github/LispZero/lisp-zero-single.c:1328
func debug_output(obj *Object_s) {
	object_write(stdout, obj); nl(stdout)
	stdout.Flush()
}

func init() {
	flag.StringVar(&profiler, "profiler", "runtime/pprof", "Specify type of profiler to use")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.IntVar(&inBufSize, "inbufsize", 4096, "Input buffer size to use")
	flag.BoolVar(&quiet, "q", false, "quiet; do not print top-level eval results")
	flag.BoolVar(&tracing, "t", false, "print diagnostic trace during evaluation")

	stdin = os.Stdin
	stdout = bufio.NewWriter(os.Stdout)
	stderr = bufio.NewWriter(os.Stderr)
}
