package main

import "github.com/dhamidi/dux"

// TODO: Command{{.name.Title }}{{$typeName := (printf "Command%s" .name.Title)}}
type {{$typeName}} struct {}

// CommandName implements dux.Command
func (c *{{$typeName}}) CommandName() string { return "{{.name.ToLisp.Lower }}" }

// CommandDescription implements dux.Command
func (c *{{$typeName}}) CommandDescription() string { return `{{.description}}` }

// Execute {{.description}}
func (c *{{$typeName}}) Execute(ctx *dux.Context, args []string) error {
        return fmt.Errorf("Command not implemented")
}
