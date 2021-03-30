package render

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-kit/kit/log"
)

type Option func(*Renderer) error

// AddTemplates adds filenames for the next call to parseTempaltes
func AddTemplates(files ...string) Option {
	return func(r *Renderer) error {
		if len(files) == 0 {
			return fmt.Errorf("render: no templates passed")
		}
		r.templateFiles = files
		return nil
	}
}

func BaseTemplates(bases ...string) Option {
	return func(r *Renderer) error {
		r.baseTemplates = bases
		return nil
	}
}

// ErrorTemplate sets the filename of template that should be used for errors.
func ErrorTemplate(name string) Option {
	return func(r *Renderer) error {
		r.errorTemplate = name
		return nil
	}
}

func FuncMap(m template.FuncMap) Option {
	return func(r *Renderer) error {
		r.funcMap = m
		return nil
	}
}

// FuncInjector needs to return a function that is admissible for (html/template).FuncMap
type FuncInjector func(*http.Request) interface{}

// InjectTemplateFunc sets template functions, just like FuncMap.
// But in comparrison, these here are evaluated once with the request before being assigned to the template.
// This allows for accessing the requests cookies or headers, for instance.
func InjectTemplateFunc(name string, fn FuncInjector) Option {
	return func(r *Renderer) error {
		if _, has := r.tplFuncInjectors[name]; has {
			return fmt.Errorf("injection %s name already taken", name)
		}
		r.tplFuncInjectors[name] = fn
		return nil
	}
}

func SetLogger(l log.Logger) Option {
	return func(r *Renderer) error {
		if l == nil {
			return errors.New("render: nil logger passed")
		}
		r.log = l
		return nil
	}
}

// ErrorHandlerFunc is just like http.HandlerFunc but with an additional status code and error
type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, int, error)

// overwrites the default error handler and allows for inspecting errors before rendering them
func SetErrorHandler(fn ErrorHandlerFunc) Option {
	return func(r *Renderer) error {
		if fn == nil {
			return errors.New("render: nil ErrorHandlerFunc passed")
		}
		r.errHandler = fn
		return nil
	}
}
