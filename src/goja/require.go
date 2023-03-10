package goja

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"text/template"

	"merkaba/goja/parser"
)

type ModuleLoader func(*Runtime, *Object)

// SourceLoader represents a function that returns a file data at a given path.
// The function should return ModuleFileDoesNotExistError if the file either doesn't exist or is a directory.
// This error will be ignored by the resolver and the search will continue. Any other errors will be propagated.
type SourceLoader func(path string) ([]byte, error)

var (
	InvalidModuleError     = errors.New("Invalid module")
	IllegalModuleNameError = errors.New("Illegal module name")

	ModuleFileDoesNotExistError = errors.New("module file does not exist")
)

var native map[string]ModuleLoader

// Registry contains a cache of compiled modules which can be used by multiple Runtimes
type Registry struct {
	sync.Mutex
	native   map[string]ModuleLoader
	compiled map[string]*Program

	srcLoader     SourceLoader
	globalFolders []string
}

type RequireModule struct {
	r           *Registry
	runtime     *Runtime
	modules     map[string]*Object
	nodeModules map[string]*Object
}

func NewRegistry(opts ...Option) *Registry {
	r := &Registry{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func NewRegistryWithLoader(srcLoader SourceLoader) *Registry {
	return NewRegistry(WithLoader(srcLoader))
}

type Option func(*Registry)

// WithLoader sets a function which will be called by the require() function in order to get a source code for a
// module at the given path. The same function will be used to get external source maps.
// Note, this only affects the modules loaded by the require() function. If you need to use it as a source map
// loader for code parsed in a different way (such as runtime.RunString() or eval()), use (*Runtime).SetParserOptions()
func WithLoader(srcLoader SourceLoader) Option {
	return func(r *Registry) {
		r.srcLoader = srcLoader
	}
}

// WithGlobalFolders appends the given paths to the registry's list of
// global folders to search if the requested module is not found
// elsewhere.  By default, a registry's global folders list is empty.
// In the reference Node.js implementation, the default global folders
// list is $NODE_PATH, $HOME/.node_modules, $HOME/.node_libraries and
// $PREFIX/lib/node, see
// https://nodeorg/api/modules.html#modules_loading_from_the_global_folders.
func WithGlobalFolders(globalFolders ...string) Option {
	return func(r *Registry) {
		r.globalFolders = globalFolders
	}
}

// Enable adds the require() function to the specified runtime.
func (r *Registry) Enable(runtime *Runtime) *RequireModule {
	rrt := &RequireModule{
		r:           r,
		runtime:     runtime,
		modules:     make(map[string]*Object),
		nodeModules: make(map[string]*Object),
	}
	runtime.Set("require", rrt.require)
	return rrt
}

func (r *Registry) RegisterNativeModule(name string, loader ModuleLoader) {
	r.Lock()
	defer r.Unlock()

	if r.native == nil {
		r.native = make(map[string]ModuleLoader)
	}
	name = filepathClean(name)
	r.native[name] = loader
}

// DefaultSourceLoader is used if none was set (see WithLoader()). It simply loads files from the host's filesystem.
func DefaultSourceLoader(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		if os.IsNotExist(err) || errors.Is(err, syscall.EISDIR) {
			err = ModuleFileDoesNotExistError
		}
	}
	return data, err
}

func (r *Registry) getSource(p string) ([]byte, error) {
	srcLoader := r.srcLoader
	if srcLoader == nil {
		srcLoader = DefaultSourceLoader
	}
	return srcLoader(p)
}

func (r *Registry) getCompiledSource(p string) (*Program, error) {
	r.Lock()
	defer r.Unlock()

	prg := r.compiled[p]
	if prg == nil {
		buf, err := r.getSource(p)
		if err != nil {
			return nil, err
		}
		s := string(buf)

		if path.Ext(p) == ".json" {
			s = "module.exports = JSON.parse('" + template.JSEscapeString(s) + "')"
		}

		source := "(function(exports, require, module) {" + s + "\n})"
		parsed, err := Parse(p, source, parser.WithSourceMapLoader(r.srcLoader))
		if err != nil {
			return nil, err
		}
		prg, err = CompileAST(parsed, false)
		if err == nil {
			if r.compiled == nil {
				r.compiled = make(map[string]*Program)
			}
			r.compiled[p] = prg
		}
		return prg, err
	}
	return prg, nil
}

func (r *RequireModule) require(call FunctionCall) Value {
	ret, err := r.Require(call.Argument(0).String())
	if err != nil {
		if _, ok := err.(*Exception); !ok {
			panic(r.runtime.NewGoError(err))
		}
		panic(err)
	}
	return ret
}

func filepathClean(p string) string {
	return path.Clean(p)
}

// Require can be used to import modules from Go source (similar to JS require() function).
func (r *RequireModule) Require(p string) (ret Value, err error) {
	module, err := r.resolve(p)
	if err != nil {
		return
	}
	ret = module.Get("exports")
	return
}

func Require(runtime *Runtime, name string) Value {
	if r, ok := AssertFunction(runtime.Get("require")); ok {
		mod, err := r(Undefined(), runtime.ToValue(name))
		if err != nil {
			panic(err)
		}
		return mod
	}
	panic(runtime.NewTypeError("Please enable require for this runtime using new(require.Require).Enable(runtime)"))
}

func RegisterNativeModule(name string, loader ModuleLoader) {
	if native == nil {
		native = make(map[string]ModuleLoader)
	}
	name = filepathClean(name)
	native[name] = loader
}

func (r *RequireModule) clear() {
	r.modules = make(map[string]*Object)
}

func (r *Registry) clear() {
	r.compiled = make(map[string]*Program)
}
