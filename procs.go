package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"time"
)

func ensureArrayMap(args []Object, index int) *ArrayMap {
	switch obj := args[index].(type) {
	case *ArrayMap:
		return obj
	default:
		panic(RT.newArgTypeError(index, obj, "Map"))
	}
}

var procMeta Proc = func(args []Object) Object {
	switch obj := args[0].(type) {
	case Meta:
		meta := obj.GetMeta()
		if meta != nil {
			return meta
		}
	}
	return NIL
}

var procWithMeta Proc = func(args []Object) Object {
	checkArity(args, 2, 2)
	m := ensureMeta(args, 0)
	if args[1].Equals(NIL) {
		return args[0]
	}
	return m.WithMeta(ensureArrayMap(args, 1))
}

var procIsZero Proc = func(args []Object) Object {
	n := ensureNumber(args, 0)
	ops := GetOps(n)
	return Bool{b: ops.IsZero(n)}
}

var procIsPos Proc = func(args []Object) Object {
	n := ensureNumber(args, 0)
	ops := GetOps(n)
	return Bool{b: ops.Gt(n, Int{i: 0})}
}

var procIsNeg Proc = func(args []Object) Object {
	n := ensureNumber(args, 0)
	ops := GetOps(n)
	return Bool{b: ops.Lt(n, Int{i: 0})}
}

var procAdd Proc = func(args []Object) Object {
	x := assertNumber(args[0], "")
	y := assertNumber(args[1], "")
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Add(x, y)
}

var procAddEx Proc = func(args []Object) Object {
	x := assertNumber(args[0], "")
	y := assertNumber(args[1], "")
	ops := GetOps(x).Combine(GetOps(y)).Combine(BIGINT_OPS)
	return ops.Add(x, y)
}

var procMultiply Proc = func(args []Object) Object {
	x := assertNumber(args[0], "")
	y := assertNumber(args[1], "")
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Multiply(x, y)
}

var procMultiplyEx Proc = func(args []Object) Object {
	x := assertNumber(args[0], "")
	y := assertNumber(args[1], "")
	ops := GetOps(x).Combine(GetOps(y)).Combine(BIGINT_OPS)
	return ops.Multiply(x, y)
}

var procSubtract Proc = func(args []Object) Object {
	var a, b Object
	if len(args) == 1 {
		a = Int{i: 0}
		b = args[0]
	} else {
		a = args[0]
		b = args[1]
	}
	ops := GetOps(a).Combine(GetOps(b))
	return ops.Subtract(assertNumber(a, ""), assertNumber(b, ""))
}

var procSubtractEx Proc = func(args []Object) Object {
	var a, b Object
	if len(args) == 1 {
		a = Int{i: 0}
		b = args[0]
	} else {
		a = args[0]
		b = args[1]
	}
	ops := GetOps(a).Combine(GetOps(b)).Combine(BIGINT_OPS)
	return ops.Subtract(assertNumber(a, ""), assertNumber(b, ""))
}

var procDivide Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	y := ensureNumber(args, 1)
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Divide(x, y)
}

var procQuot Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	y := ensureNumber(args, 1)
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Quotient(x, y)
}

var procRem Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	y := ensureNumber(args, 1)
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Rem(x, y)
}

var procBitNot Proc = func(args []Object) Object {
	x := assertInt(args[0], "Bit operation not supported for "+args[0].GetType().ToString(false))
	return Int{i: ^x.i}
}

func assertInts(args []Object) (Int, Int) {
	x := assertInt(args[0], "Bit operation not supported for "+args[0].GetType().ToString(false))
	y := assertInt(args[1], "Bit operation not supported for "+args[1].GetType().ToString(false))
	return x, y
}

var procBitAnd Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i & y.i}
}

var procBitOr Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i | y.i}
}

var procBitXor Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i ^ y.i}
}

var procBitAndNot Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i &^ y.i}
}

var procBitClear Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i &^ (1 << uint(y.i))}
}

var procBitSet Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i | (1 << uint(y.i))}
}

var procBitFlip Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i ^ (1 << uint(y.i))}
}

var procBitTest Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Bool{b: x.i&(1<<uint(y.i)) != 0}
}

var procBitShiftLeft Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i << uint(y.i)}
}

var procBitShiftRight Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: x.i >> uint(y.i)}
}

var procUnsignedBitShiftRight Proc = func(args []Object) Object {
	x, y := assertInts(args)
	return Int{i: int(uint(x.i) >> uint(y.i))}
}

var procExInfo Proc = func(args []Object) Object {
	checkArity(args, 2, 2)
	return &ExInfo{
		msg:  ensureString(args, 0),
		data: ensureArrayMap(args, 1),
		rt:   RT.clone(),
	}
}

var procSetMacro Proc = func(args []Object) Object {
	vr := args[0].(*Var)
	vr.isMacro = true
	return vr
}

var procList Proc = func(args []Object) Object {
	return NewListFrom(args...)
}

var procCons Proc = func(args []Object) Object {
	checkArity(args, 2, 2)
	s := ensureSeqable(args, 1).Seq()
	return s.Cons(args[0])
}

var procFirst Proc = func(args []Object) Object {
	checkArity(args, 1, 1)
	s := ensureSeqable(args, 0).Seq()
	return s.First()
}

var procNext Proc = func(args []Object) Object {
	checkArity(args, 1, 1)
	s := ensureSeqable(args, 0).Seq()
	res := s.Rest()
	if res.IsEmpty() {
		return NIL
	}
	return res
}

var procRest Proc = func(args []Object) Object {
	checkArity(args, 1, 1)
	s := ensureSeqable(args, 0).Seq()
	return s.Rest()
}

var procConj Proc = func(args []Object) Object {
	switch c := args[0].(type) {
	case Conjable:
		return c.Conj(args[1])
	case Seq:
		return c.Cons(args[1])
	default:
		panic(RT.newError("conj's first argument must be a collection"))
	}
}

var procSeq Proc = func(args []Object) Object {
	checkArity(args, 1, 1)
	s := ensureSeqable(args, 0).Seq()
	if s.IsEmpty() {
		return NIL
	}
	return s
}

var procIsInstance Proc = func(args []Object) Object {
	checkArity(args, 2, 2)
	switch t := args[0].(type) {
	case *Type:
		if args[1].Equals(NIL) {
			return Bool{b: false}
		}
		if t.reflectType.Kind() == reflect.Interface {
			return Bool{b: args[1].GetType().reflectType.Implements(t.reflectType)}
		} else {
			return Bool{b: args[1].GetType().reflectType == t.reflectType}
		}
	default:
		panic(RT.newError("First argument to instance? must be a type"))
	}
}

var procAssoc Proc = func(args []Object) Object {
	return ensureAssociative(args, 0).Assoc(args[1], args[2])
}

var procEquals Proc = func(args []Object) Object {
	return Bool{b: args[0].Equals(args[1])}
}

var procCount Proc = func(args []Object) Object {
	switch obj := args[0].(type) {
	case Counted:
		return Int{i: obj.Count()}
	default:
		s := assertSeqable(obj, "count not supported on this type: "+obj.GetType().ToString(false))
		return Int{i: SeqCount(s.Seq())}
	}
}

var procSubvec Proc = func(args []Object) Object {
	// TODO: implement proper Subvector structure
	v := args[0].(*Vector)
	start := args[1].(Int).i
	end := args[2].(Int).i
	subv := make([]Object, 0, end-start)
	for i := start; i < end; i++ {
		subv = append(subv, v.at(i))
	}
	return NewVectorFrom(subv...)
}

var procCast Proc = func(args []Object) Object {
	t := ensureType(args, 0)
	if t.reflectType.Kind() == reflect.Interface &&
		args[1].GetType().reflectType.Implements(t.reflectType) ||
		args[1].GetType().reflectType == t.reflectType {
		return args[1]
	}
	panic(RT.newError("Cannot cast " + args[1].GetType().ToString(false) + " to " + t.ToString(false)))
}

var procVec Proc = func(args []Object) Object {
	return NewVectorFromSeq(ensureSeqable(args, 0).Seq())
}

var procHashMap Proc = func(args []Object) Object {
	if len(args)%2 != 0 {
		panic(RT.newError("No value supplied for key " + args[len(args)-1].ToString(false)))
	}
	res := EmptyArrayMap()
	for i := 0; i < len(args); i += 2 {
		res.Set(args[i], args[i+1])
	}
	return res
}

var procHashSet Proc = func(args []Object) Object {
	res := EmptySet()
	for i := 0; i < len(args); i++ {
		res.Add(args[i])
	}
	return res
}

var procStr Proc = func(args []Object) Object {
	var buffer bytes.Buffer
	for _, obj := range args {
		if !obj.Equals(NIL) {
			buffer.WriteString(obj.ToString(false))
		}
	}
	return String{s: buffer.String()}
}

var procSymbol Proc = func(args []Object) Object {
	if len(args) == 1 {
		return MakeSymbol(ensureString(args, 0).s)
	}
	return Symbol{
		ns:   STRINGS.Intern(ensureString(args, 0).s),
		name: STRINGS.Intern(ensureString(args, 1).s),
	}
}

var procKeyword Proc = func(args []Object) Object {
	if len(args) == 1 {
		switch obj := args[0].(type) {
		case String:
			return MakeKeyword(obj.s)
		case Symbol:
			return Keyword{
				ns:   obj.ns,
				name: obj.name,
			}
		default:
			return NIL
		}
	}
	return Keyword{
		ns:   STRINGS.Intern(ensureString(args, 0).s),
		name: STRINGS.Intern(ensureString(args, 1).s),
	}
}

var procGensym Proc = func(args []Object) Object {
	return genSym(ensureString(args, 0).s, "")
}

var procApply Proc = func(args []Object) Object {
	// TODO:
	// Stacktrace is broken. Need to somehow know
	// the name of the function passed ...
	f := ensureCallable(args, 0)
	return f.Call(ToSlice(ensureSeqable(args, 1).Seq()))
}

var procLazySeq Proc = func(args []Object) Object {
	return &LazySeq{
		fn: args[0].(*Fn),
	}
}

var procDelay Proc = func(args []Object) Object {
	return &Delay{
		fn: args[0].(*Fn),
	}
}

var procForce Proc = func(args []Object) Object {
	switch d := args[0].(type) {
	case *Delay:
		return d.Force()
	default:
		return d
	}
}

var procIdentical Proc = func(args []Object) Object {
	return Bool{b: args[0] == args[1]}
}

var procCompare Proc = func(args []Object) Object {
	k1, k2 := args[0], args[1]
	if k1.Equals(k2) {
		return Int{i: 0}
	}
	switch k2.(type) {
	case Nil:
		return Int{i: 1}
	}
	switch k1 := k1.(type) {
	case Nil:
		return Int{i: -1}
	case Comparable:
		return Int{i: k1.Compare(k2)}
	}
	panic(RT.newError(fmt.Sprintf("%s (type: %s) is not a Comparable", k1.ToString(true), k1.GetType().ToString(false))))
}

var procInt Proc = func(args []Object) Object {
	switch obj := args[0].(type) {
	case Char:
		return Int{i: int(obj.ch)}
	case Number:
		return obj.Int()
	default:
		panic(RT.newError(fmt.Sprintf("Cannot cast %s (type: %s) to Int", obj.ToString(true), obj.GetType().ToString(false))))
	}
}

var procNumber Proc = func(args []Object) Object {
	return assertNumber(args[0], fmt.Sprintf("Cannot cast %s (type: %s) to Number", args[0].ToString(true), args[0].GetType().ToString(false)))
}

var procDouble Proc = func(args []Object) Object {
	n := assertNumber(args[0], fmt.Sprintf("Cannot cast %s (type: %s) to Double", args[0].ToString(true), args[0].GetType().ToString(false)))
	return n.Double()
}

var procChar Proc = func(args []Object) Object {
	switch c := args[0].(type) {
	case Char:
		return c
	case Number:
		i := c.Int().i
		if i < MIN_RUNE || i > MAX_RUNE {
			panic(RT.newError(fmt.Sprintf("Value out of range for char: %d", i)))
		}
		return Char{ch: rune(i)}
	default:
		panic(RT.newError(fmt.Sprintf("Cannot cast %s (type: %s) to Char", c.ToString(true), c.GetType().ToString(false))))
	}
}

var procBool Proc = func(args []Object) Object {
	return Bool{b: toBool(args[0])}
}

var procNumerator Proc = func(args []Object) Object {
	bi := ensureRatio(args, 0).r.Num()
	return &BigInt{b: *bi}
}

var procDenominator Proc = func(args []Object) Object {
	bi := ensureRatio(args, 0).r.Denom()
	return &BigInt{b: *bi}
}

var procBigInt Proc = func(args []Object) Object {
	switch n := args[0].(type) {
	case Number:
		return &BigInt{b: *n.BigInt()}
	case String:
		bi := big.Int{}
		if _, ok := bi.SetString(n.s, 10); ok {
			return &BigInt{b: bi}
		}
		panic(RT.newError("Invalid number format " + n.s))
	default:
		panic(RT.newError(fmt.Sprintf("Cannot cast %s (type: %s) to BigInt", n.ToString(true), n.GetType().ToString(false))))
	}
}

var procBigFloat Proc = func(args []Object) Object {
	switch n := args[0].(type) {
	case Number:
		return &BigFloat{b: *n.BigFloat()}
	case String:
		b := big.Float{}
		if _, ok := b.SetString(n.s); ok {
			return &BigFloat{b: b}
		}
		panic(RT.newError("Invalid number format " + n.s))
	default:
		panic(RT.newError(fmt.Sprintf("Cannot cast %s (type: %s) to BigFloat", n.ToString(true), n.GetType().ToString(false))))
	}
}

var procNth Proc = func(args []Object) Object {
	n := ensureNumber(args, 1).Int().i
	switch coll := args[0].(type) {
	case Indexed:
		if len(args) == 3 {
			return coll.TryNth(n, args[2])
		}
		return coll.Nth(n)
	case Seqable:
		if len(args) == 3 {
			return SeqTryNth(coll.Seq(), n, args[2])
		}
		return SeqNth(coll.Seq(), n)
	default:
		panic(RT.newError("nth not supported on this type: " + coll.GetType().ToString(false)))
	}
}

var procLt Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Bool{b: GetOps(a).Combine(GetOps(b)).Lt(a, b)}
}

var procLte Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Bool{b: GetOps(a).Combine(GetOps(b)).Lte(a, b)}
}

var procGt Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Bool{b: GetOps(a).Combine(GetOps(b)).Gt(a, b)}
}

var procGte Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Bool{b: GetOps(a).Combine(GetOps(b)).Gte(a, b)}
}

var procEq Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Bool{b: GetOps(a).Combine(GetOps(b)).Eq(a, b)}
}

var procMax Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Max(a, b)
}

var procMin Proc = func(args []Object) Object {
	a := assertNumber(args[0], "")
	b := assertNumber(args[1], "")
	return Min(a, b)
}

var procIncEx Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	ops := GetOps(x).Combine(BIGINT_OPS)
	return ops.Add(x, Int{i: 1})
}

var procDecEx Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	ops := GetOps(x).Combine(BIGINT_OPS)
	return ops.Subtract(x, Int{i: 1})
}

var procInc Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	ops := GetOps(x).Combine(INT_OPS)
	return ops.Add(x, Int{i: 1})
}

var procDec Proc = func(args []Object) Object {
	x := ensureNumber(args, 0)
	ops := GetOps(x).Combine(INT_OPS)
	return ops.Subtract(x, Int{i: 1})
}

var procPeek Proc = func(args []Object) Object {
	s := assertStack(args[0], "")
	return s.Peek()
}

var procPop Proc = func(args []Object) Object {
	s := assertStack(args[0], "")
	return s.Pop().(Object)
}

var procContains Proc = func(args []Object) Object {
	switch c := args[0].(type) {
	case Gettable:
		ok, _ := c.Get(args[1])
		if ok {
			return Bool{b: true}
		}
		return Bool{b: false}
	}
	panic(RT.newError("contains? not supported on type " + args[0].GetType().ToString(false)))
}

var procGet Proc = func(args []Object) Object {
	switch c := args[0].(type) {
	case Gettable:
		ok, v := c.Get(args[1])
		if ok {
			return v
		}
	}
	if len(args) == 3 {
		return args[2]
	}
	return NIL
}

var procDissoc Proc = func(args []Object) Object {
	return ensureMap(args, 0).Without(args[1])
}

var procDisj Proc = func(args []Object) Object {
	return ensureSet(args, 0).Disjoin(args[1])
}

var procFind Proc = func(args []Object) Object {
	res := ensureAssociative(args, 0).EntryAt(args[1])
	if res == nil {
		return NIL
	}
	return res
}

var procKeys Proc = func(args []Object) Object {
	return ensureMap(args, 0).Keys()
}

var procVals Proc = func(args []Object) Object {
	return ensureMap(args, 0).Vals()
}

var procRseq Proc = func(args []Object) Object {
	return ensureReversible(args, 0).Rseq()
}

var procName Proc = func(args []Object) Object {
	return String{s: ensureNamed(args, 0).Name()}
}

var procNamespace Proc = func(args []Object) Object {
	ns := ensureNamed(args, 0).Namespace()
	if ns == "" {
		return NIL
	}
	return String{s: ns}
}

var procFindVar Proc = func(args []Object) Object {
	sym := ensureSymbol(args, 0)
	if sym.ns == nil {
		panic(RT.newError("find-var argument must be namespace-qualified symbol"))
	}
	if v, ok := GLOBAL_ENV.Resolve(sym); ok {
		return v
	}
	return NIL
}

var procSort Proc = func(args []Object) Object {
	cmp := ensureComparator(args, 0)
	coll := ensureSeqable(args, 1)
	s := SortableSlice{
		s:   ToSlice(coll.Seq()),
		cmp: cmp,
	}
	sort.Sort(s)
	return &ArraySeq{arr: s.s}
}

var procEval Proc = func(args []Object) Object {
	parseContext := &ParseContext{globalEnv: GLOBAL_ENV}
	expr := parse(args[0], parseContext)
	return eval(expr, nil)
}

var procType Proc = func(args []Object) Object {
	return args[0].GetType()
}

func pr(args []Object, escape bool) Object {
	n := len(args)
	if n > 0 {
		for _, arg := range args[:n-1] {
			print(arg.ToString(escape))
			print(" ")
		}
		print(args[n-1].ToString(escape))
	}
	return NIL
}

var procPr Proc = func(args []Object) Object {
	return pr(args, true)
}

var procPrint Proc = func(args []Object) Object {
	return pr(args, false)
}

var procNewline Proc = func(args []Object) Object {
	println()
	return NIL
}

func readFromReader(reader io.RuneReader) Object {
	r := NewReader(reader)
	obj, err := TryRead(r)
	if err != nil {
		panic(RT.newError(err.Error()))
	}
	return obj
}

var procRead Proc = func(args []Object) Object {
	return readFromReader(bufio.NewReader(os.Stdin))
}

var procReadString Proc = func(args []Object) Object {
	return readFromReader(strings.NewReader(ensureString(args, 0).s))
}

var procReadLine Proc = func(args []Object) Object {
	var line string
	fmt.Scanln(&line)
	return &String{s: line}
}

var procNanoTime Proc = func(args []Object) Object {
	return &BigInt{b: *big.NewInt(time.Now().UnixNano())}
}

var procMacroexpand1 Proc = func(args []Object) Object {
	switch s := args[0].(type) {
	case Seq:
		parseContext := &ParseContext{globalEnv: GLOBAL_ENV}
		return macroexpand1(s, parseContext)
	default:
		return s
	}
}

func loadReader(reader *Reader) (Object, error) {
	parseContext := &ParseContext{globalEnv: GLOBAL_ENV}
	var lastObj Object = NIL
	for {
		obj, err := TryRead(reader)
		if err == io.EOF {
			return lastObj, nil
		}
		if err != nil {
			return nil, err
		}
		expr, err := TryParse(obj, parseContext)
		if err != nil {
			return nil, err
		}
		lastObj, err = TryEval(expr)
		if err != nil {
			return nil, err
		}
	}
}

var procLoadString Proc = func(args []Object) Object {
	s := ensureString(args, 0)
	obj, err := loadReader(NewReader(strings.NewReader(s.s)))
	if err != nil {
		panic(err)
	}
	return obj
}

var procFindNamespace Proc = func(args []Object) Object {
	ns := GLOBAL_ENV.FindNamespace(ensureSymbol(args, 0))
	if ns == nil {
		return NIL
	}
	return ns
}

var procCreateNamespace Proc = func(args []Object) Object {
	return GLOBAL_ENV.EnsureNamespace(ensureSymbol(args, 0))
}

var procRemoveNamespace Proc = func(args []Object) Object {
	ns := GLOBAL_ENV.RemoveNamespace(ensureSymbol(args, 0))
	if ns == nil {
		return NIL
	}
	return ns
}

var procAllNamespaces Proc = func(args []Object) Object {
	s := make([]Object, 0, len(GLOBAL_ENV.namespaces))
	for _, ns := range GLOBAL_ENV.namespaces {
		s = append(s, ns)
	}
	return &ArraySeq{arr: s}
}

var procNamespaceName Proc = func(args []Object) Object {
	return ensureNamespace(args, 0).name
}

var procNamespaceMap Proc = func(args []Object) Object {
	r := &ArrayMap{}
	for k, v := range ensureNamespace(args, 0).mappings {
		r.Add(MakeSymbol(*k), v)
	}
	return r
}

var procNamespaceUnmap Proc = func(args []Object) Object {
	ns := ensureNamespace(args, 0)
	sym := ensureSymbol(args, 1)
	if sym.ns != nil {
		panic(RT.newError("Can't unintern namespace-qualified symbol"))
	}
	delete(ns.mappings, sym.name)
	return NIL
}

var procVarNamespace Proc = func(args []Object) Object {
	v := ensureVar(args, 0)
	return v.ns
}

var procRefer Proc = func(args []Object) Object {
	ns := ensureNamespace(args, 0)
	sym := ensureSymbol(args, 1)
	v := ensureVar(args, 2)
	return ns.Refer(sym, v)
}

var procAlias Proc = func(args []Object) Object {
	ensureNamespace(args, 0).AddAlias(ensureSymbol(args, 1), ensureNamespace(args, 2))
	return NIL
}

var procNamespaceAliases Proc = func(args []Object) Object {
	r := &ArrayMap{}
	for k, v := range ensureNamespace(args, 0).aliases {
		r.Add(MakeSymbol(*k), v)
	}
	return r
}

var procNamespaceUnalias Proc = func(args []Object) Object {
	ns := ensureNamespace(args, 0)
	sym := ensureSymbol(args, 1)
	if sym.ns != nil {
		panic(RT.newError("Alias can't be namespace-qualified"))
	}
	delete(ns.aliases, sym.name)
	return NIL
}

var procSh Proc = func(args []Object) Object {
	strs := make([]string, len(args))
	for i, _ := range args {
		strs[i] = ensureString(args, i).s
	}
	cmd := exec.Command(strs[0], strs[1:len(strs)]...)
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(RT.newError(err.Error()))
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		panic(RT.newError(err.Error()))
	}
	if err = cmd.Start(); err != nil {
		panic(RT.newError(err.Error()))
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdoutReader)
	stdoutString := buf.String()
	buf = new(bytes.Buffer)
	buf.ReadFrom(stderrReader)
	stderrString := buf.String()
	if err = cmd.Wait(); err != nil {
		EmptyArrayMap().Assoc(MakeKeyword("success"), Bool{b: false})
	}
	res := EmptyArrayMap()
	res.Add(MakeKeyword("success"), Bool{b: true})
	res.Add(MakeKeyword("out"), String{s: stdoutString})
	res.Add(MakeKeyword("err"), String{s: stderrString})
	return res
}

var coreNamespace = GLOBAL_ENV.namespaces[MakeSymbol("gclojure.core").name]

func intern(name string, proc Proc) {
	coreNamespace.intern(MakeSymbol(name)).value = proc
}

func init() {
	intern("list**", procList)
	intern("cons*", procCons)
	intern("first*", procFirst)
	intern("next*", procNext)
	intern("rest*", procRest)
	intern("conj*", procConj)
	intern("seq*", procSeq)
	intern("instance?*", procIsInstance)
	intern("assoc*", procAssoc)
	intern("meta*", procMeta)
	intern("with-meta*", procWithMeta)
	intern("=*", procEquals)
	intern("count*", procCount)
	intern("subvec*", procSubvec)
	intern("cast*", procCast)
	intern("vec*", procVec)
	intern("hash-map*", procHashMap)
	intern("hash-set*", procHashSet)
	intern("str*", procStr)
	intern("symbol*", procSymbol)
	intern("gensym*", procGensym)
	intern("keyword*", procKeyword)
	intern("apply*", procApply)
	intern("lazy-seq*", procLazySeq)
	intern("delay*", procDelay)
	intern("force*", procForce)
	intern("identical*", procIdentical)
	intern("compare*", procCompare)
	intern("zero?*", procIsZero)
	intern("int*", procInt)
	intern("nth*", procNth)
	intern("<*", procLt)
	intern("<=*", procLte)
	intern(">*", procGt)
	intern(">=*", procGte)
	intern("==*", procEq)
	intern("inc'*", procIncEx)
	intern("inc*", procInc)
	intern("dec'*", procDecEx)
	intern("dec*", procDec)
	intern("add'*", procAddEx)
	intern("add*", procAdd)
	intern("multiply'*", procMultiplyEx)
	intern("multiply*", procMultiply)
	intern("divide*", procDivide)
	intern("subtract'*", procSubtractEx)
	intern("subtract*", procSubtract)
	intern("max*", procMax)
	intern("min*", procMin)
	intern("pos*", procIsPos)
	intern("neg*", procIsNeg)
	intern("quot*", procQuot)
	intern("rem*", procRem)
	intern("bit-not*", procBitNot)
	intern("bit-and*", procBitAnd)
	intern("bit-or*", procBitOr)
	intern("bit-xor*", procBitXor)
	intern("bit-and-not*", procBitAndNot)
	intern("bit-clear*", procBitClear)
	intern("bit-set*", procBitSet)
	intern("bit-flip*", procBitFlip)
	intern("bit-test*", procBitTest)
	intern("bit-shift-left*", procBitShiftLeft)
	intern("bit-shift-right*", procBitShiftRight)
	intern("unsigned-bit-shift-right*", procUnsignedBitShiftRight)
	intern("peek*", procPeek)
	intern("pop*", procPop)
	intern("contains?*", procContains)
	intern("get*", procGet)
	intern("dissoc*", procDissoc)
	intern("disj*", procDisj)
	intern("find*", procFind)
	intern("keys*", procKeys)
	intern("vals*", procVals)
	intern("rseq*", procRseq)
	intern("name*", procName)
	intern("namespace*", procNamespace)
	intern("find-var*", procFindVar)
	intern("sort*", procSort)
	intern("eval*", procEval)
	intern("type*", procType)
	intern("num*", procNumber)
	intern("double*", procDouble)
	intern("char*", procChar)
	intern("bool*", procBool)
	intern("numerator*", procNumerator)
	intern("denominator*", procDenominator)
	intern("bigint*", procBigInt)
	intern("bigfloat*", procBigFloat)
	intern("pr*", procPr)
	intern("newline*", procNewline)
	intern("print*", procPrint)
	intern("read*", procRead)
	intern("read-line*", procReadLine)
	intern("read-string*", procReadString)
	intern("nano-time*", procNanoTime)
	intern("macroexpand-1*", procMacroexpand1)
	intern("load-string*", procLoadString)
	intern("find-ns*", procFindNamespace)
	intern("create-ns*", procCreateNamespace)
	intern("remove-ns*", procRemoveNamespace)
	intern("all-ns*", procAllNamespaces)
	intern("ns-name*", procNamespaceName)
	intern("ns-map*", procNamespaceMap)
	intern("ns-unmap*", procNamespaceUnmap)
	intern("var-ns*", procVarNamespace)
	intern("refer*", procRefer)
	intern("alias*", procAlias)
	intern("ns-aliases*", procNamespaceAliases)
	intern("ns-unalias*", procNamespaceUnalias)

	intern("ex-info", procExInfo)
	intern("set-macro*", procSetMacro)
	intern("sh", procSh)

	currentNamespace := GLOBAL_ENV.currentNamespace
	GLOBAL_ENV.SetCurrentNamespace(coreNamespace)
	data, err := Asset("data/core.clj")
	if err != nil {
		panic(RT.newError("Could not load core.clj"))
	}
	reader := bytes.NewReader(data)
	processReader(NewReader(reader), EVAL)
	GLOBAL_ENV.SetCurrentNamespace(currentNamespace)
}
