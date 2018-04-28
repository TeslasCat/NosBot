package modules

import (
    "../types"
)

var modules map[string]func(*types.Message) types.Response

func init () {
    modules = make(map[string]func(*types.Message) types.Response)
}

func RegisterModule (name string, handler func(*types.Message) types.Response) {
    modules[name] = handler
}

func Get (name string) func(*types.Message) types.Response {
    return modules[name]
}