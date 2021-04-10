// This file is generated by generate-std.joke script. Do not edit manually!

package os

import (
	. "github.com/candid82/joker/core"
	"io/ioutil"
	"os"
)

var __args__P ProcFn = __args_
var args_ Proc = Proc{Fn: __args__P, Name: "args_", Package: "std/os"}

func __args_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := commandArgs()
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __chdir__P ProcFn = __chdir_
var chdir_ Proc = Proc{Fn: __chdir__P, Name: "chdir_", Package: "std/os"}

func __chdir_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		dirname := ExtractString(_args, 0)
		_res := chdir(dirname)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __chmod__P ProcFn = __chmod_
var chmod_ Proc = Proc{Fn: __chmod__P, Name: "chmod_", Package: "std/os"}

func __chmod_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		name := ExtractString(_args, 0)
		mode := ExtractInt(_args, 1)
		err := os.Chmod(name, os.FileMode(mode))
		PanicOnErr(err)
		_res := NIL
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __close__P ProcFn = __close_
var close_ Proc = Proc{Fn: __close__P, Name: "close_", Package: "std/os"}

func __close_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		f := ExtractFile(_args, 0)
		err := f.Close()
		PanicOnErr(err)
		_res := NIL
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __create__P ProcFn = __create_
var create_ Proc = Proc{Fn: __create__P, Name: "create_", Package: "std/os"}

func __create_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		name := ExtractString(_args, 0)
		_res, err := os.Create(name)
		PanicOnErr(err)
		return MakeFile(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __create_temp__P ProcFn = __create_temp_
var create_temp_ Proc = Proc{Fn: __create_temp__P, Name: "create_temp_", Package: "std/os"}

func __create_temp_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		dir := ExtractString(_args, 0)
		pattern := ExtractString(_args, 1)
		_res, err := ioutil.TempFile(dir, pattern)
		PanicOnErr(err)
		return MakeFile(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __cwd__P ProcFn = __cwd_
var cwd_ Proc = Proc{Fn: __cwd__P, Name: "cwd_", Package: "std/os"}

func __cwd_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := getwd()
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __env__P ProcFn = __env_
var env_ Proc = Proc{Fn: __env__P, Name: "env_", Package: "std/os"}

func __env_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := env()
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __exec__P ProcFn = __exec_
var exec_ Proc = Proc{Fn: __exec__P, Name: "exec_", Package: "std/os"}

func __exec_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		name := ExtractString(_args, 0)
		opts := ExtractMap(_args, 1)
		_res := execute(name, opts)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isexists__P ProcFn = __isexists_
var isexists_ Proc = Proc{Fn: __isexists__P, Name: "isexists_", Package: "std/os"}

func __isexists_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		path := ExtractString(_args, 0)
		_res := exists(path)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __exit__P ProcFn = __exit_
var exit_ Proc = Proc{Fn: __exit__P, Name: "exit_", Package: "std/os"}

func __exit_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		code := ExtractInt(_args, 0)
		_res := NIL
		ExitJoker(code)
		return _res

	case _c == 0:
		_res := NIL
		ExitJoker(0)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __get_env__P ProcFn = __get_env_
var get_env_ Proc = Proc{Fn: __get_env__P, Name: "get_env_", Package: "std/os"}

func __get_env_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		key := ExtractString(_args, 0)
		_res := getEnv(key)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __hostname__P ProcFn = __hostname_
var hostname_ Proc = Proc{Fn: __hostname__P, Name: "hostname_", Package: "std/os"}

func __hostname_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res, err := os.Hostname()
		PanicOnErr(err)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __ls__P ProcFn = __ls_
var ls_ Proc = Proc{Fn: __ls__P, Name: "ls_", Package: "std/os"}

func __ls_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		dirname := ExtractString(_args, 0)
		_res := readDir(dirname)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __mkdir__P ProcFn = __mkdir_
var mkdir_ Proc = Proc{Fn: __mkdir__P, Name: "mkdir_", Package: "std/os"}

func __mkdir_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		name := ExtractString(_args, 0)
		perm := ExtractInt(_args, 1)
		_res := mkdir(name, perm)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __mkdir_temp__P ProcFn = __mkdir_temp_
var mkdir_temp_ Proc = Proc{Fn: __mkdir_temp__P, Name: "mkdir_temp_", Package: "std/os"}

func __mkdir_temp_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		dir := ExtractString(_args, 0)
		pattern := ExtractString(_args, 1)
		_res, err := ioutil.TempDir(dir, pattern)
		PanicOnErr(err)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __open__P ProcFn = __open_
var open_ Proc = Proc{Fn: __open__P, Name: "open_", Package: "std/os"}

func __open_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		name := ExtractString(_args, 0)
		_res, err := os.Open(name)
		PanicOnErr(err)
		return MakeFile(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __remove__P ProcFn = __remove_
var remove_ Proc = Proc{Fn: __remove__P, Name: "remove_", Package: "std/os"}

func __remove_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		name := ExtractString(_args, 0)
		err := os.Remove(name)
		PanicOnErr(err)
		_res := NIL
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __remove_all__P ProcFn = __remove_all_
var remove_all_ Proc = Proc{Fn: __remove_all__P, Name: "remove_all_", Package: "std/os"}

func __remove_all_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		path := ExtractString(_args, 0)
		err := os.RemoveAll(path)
		PanicOnErr(err)
		_res := NIL
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __set_env__P ProcFn = __set_env_
var set_env_ Proc = Proc{Fn: __set_env__P, Name: "set_env_", Package: "std/os"}

func __set_env_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		key := ExtractString(_args, 0)
		value := ExtractString(_args, 1)
		_res := setEnv(key, value)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sh__P ProcFn = __sh_
var sh_ Proc = Proc{Fn: __sh__P, Name: "sh_", Package: "std/os"}

func __sh_(_args []Object) Object {
	_c := len(_args)
	switch {
	case true:
		CheckArity(_args, 1, 999)
		name := ExtractString(_args, 0)
		arguments := ExtractStrings(_args, 1)
		_res := sh("", nil, nil, nil, name, arguments)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sh_from__P ProcFn = __sh_from_
var sh_from_ Proc = Proc{Fn: __sh_from__P, Name: "sh_from_", Package: "std/os"}

func __sh_from_(_args []Object) Object {
	_c := len(_args)
	switch {
	case true:
		CheckArity(_args, 2, 999)
		dir := ExtractString(_args, 0)
		name := ExtractString(_args, 1)
		arguments := ExtractStrings(_args, 2)
		_res := sh(dir, nil, nil, nil, name, arguments)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __stat__P ProcFn = __stat_
var stat_ Proc = Proc{Fn: __stat__P, Name: "stat_", Package: "std/os"}

func __stat_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		filename := ExtractString(_args, 0)
		_res := stat(filename)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __temp_dir__P ProcFn = __temp_dir_
var temp_dir_ Proc = Proc{Fn: __temp_dir__P, Name: "temp_dir_", Package: "std/os"}

func __temp_dir_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := os.TempDir()
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

func Init() {

	InternsOrThunks()
}

var osNamespace = GLOBAL_ENV.EnsureSymbolIsLib(MakeSymbol("joker.os"))

func init() {
	osNamespace.Lazy = Init
}
