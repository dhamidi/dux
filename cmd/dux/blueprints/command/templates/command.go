package main

// TODO: Command{{.name.Title }}{{$typeName := (printf "Command%s" .name.Title)}}
type {{$typeName}} struct {}

// CommandName implements dux.Command
func (c *{{$typeName}}) CommandName() string { return "{{.name.ToLisp.Lower }}" }

// CommandDescription implements dux.Command
func (c *{{$typeName}}) CommandDescription() string { return `{{.description}}` }

// Execute {{.description}}
func (c *{{$typeName}}) Execute(ctx *dux.Context, args []string) error {
	blueprintName := args[0]
	blueprint, err := ctx.LoadBlueprint(blueprintName)
	if err != nil {
		return err
	}
	if err := blueprint.ParseArgs(args[1:]); err != nil {
		return err
	}

	result := blueprint.Render(ctx)
	if result.HasError() {
		return result
	}

	return blueprint.CopyFilesToDestination(ctx, result)
}
