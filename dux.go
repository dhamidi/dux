package dux

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// Environment provides read-only access to a key value store mapping strings to strings.
type Environment func(string) string

// SystemEnvironment provides access to the operating system's environment variables.
var (
	SystemEnvironment = os.Getenv
)

// Context provides information about the environment in which Dux runs.
type Context struct {
	BaseDir   string
	OutputDir string
	Logger    *log.Logger
	Env       Environment
	Data      interface{}
}

// NewContextFromEnvironment creates a new context with settings taking from the given environment.
func NewContextFromEnvironment(env Environment) *Context {
	return &Context{
		BaseDir:   env("PWD"),
		OutputDir: filepath.Join(env("PWD"), "tmp"),
		Logger:    log.New(os.Stdout, "", 0),
		Env:       env,
	}
}

// Log outputs a log message for the given module.
func (ctx *Context) Log(module, kind, message string) {
	ctx.Logger.Printf("%10s %10s %s", module, kind, message)
}

// MergeJSON deeply merges two raw deserialized JSON objects.
func MergeJSON(a, b interface{}) interface{} {
	result := map[string]interface{}{}

	if a == nil && b != nil {
		return b
	}

	aMap, aMapOk := a.(map[string]interface{})
	bMap, bMapOk := b.(map[string]interface{})
	if aMapOk && bMapOk {
		for key, value := range aMap {
			result[key] = value
		}
		for key, value := range bMap {
			result[key] = MergeJSON(result[key], value)
		}
		return result
	} else if !aMapOk && !bMapOk {
		return b
	}

	return a
}

// GatherData loads template data from the file system into the current context.
func (ctx *Context) GatherData(datafile string) error {
	filename := filepath.Join(ctx.BaseDir, datafile)
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	result := map[string]interface{}{}
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return err
	}

	ctx.Data = MergeJSON(ctx.Data, result)
	return nil
}

// Command describes an action that can be executed by a user of Dux
type Command interface {
	Execute(ctx *Context, args []string) error
	CommandName() string
	CommandDescription() string
}

// BlueprintArgument encodes the data about the arguments accepted by the blueprint.
type BlueprintArgument struct {
	Name string
	Type string
	Doc  *string
}

// BlueprintFileDescription describes which target file should be generated from which template.
type BlueprintFileDescription struct {
	// Template is the name of the template to use.
	Template string

	// Destination is the a template expression that is evaluated
	// to obtain the destination file name.
	Destination string
}

func (fd *BlueprintFileDescription) Render(ctx *Context, blueprint *Blueprint, data interface{}) error {
	destinationFileTemplate, err := template.New("main").Parse(fd.Destination)
	if err != nil {
		return fmt.Errorf("cannot parse destination file template %q: %s", fd.Destination, err)
	}
	destinationFileName := bytes.NewBufferString("")
	if err := destinationFileTemplate.Execute(destinationFileName, data); err != nil {
		return fmt.Errorf("failed to generated filename from %q: %s", fd.Destination, err)
	}

	t, err := blueprint.templates.Clone()
	if err != nil {
		return fmt.Errorf("failed to clone templates: %s", err)
	}
	outputFilePath := filepath.Join(ctx.OutputDir, destinationFileName.String())
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %s", outputFilePath, err)
	}
	defer outputFile.Close()
	if err := t.Execute(outputFile, data); err != nil {
		return fmt.Errorf("error rendering template %q: %s", outputFilePath, err)
	}
	ctx.Log(blueprint.Name, "create", destinationFileName.String())
	return nil
}

// Blueprint represents a set of files that need to be created on
// disk, followed by performing editing operations on existing files.
type Blueprint struct {
	// path is the location in which the blueprint is defined.
	path string

	// templates holds all templates used by this blueprint
	templates *template.Template

	// flags accepted by this blueprint
	flags *flag.FlagSet

	// Data holds the data that is provided in the blueprint
	// manifest or extracted from flags provided to the blueprint.
	Data map[string]interface{}

	// Name is the basename of the directory containing the blueprint
	Name string `json:"name"`

	// Args is a list of arguments that are accepted by this blueprint
	Args []*BlueprintArgument `json:"args"`

	// Files describes the list of files that should be generated
	// for this blueprint.
	Files []*BlueprintFileDescription `json:"files"`
}

// NewBlueprint creates a new blueprint instance.
func NewBlueprint(path string) *Blueprint {
	return &Blueprint{
		path:      path,
		templates: nil,
		Data:      map[string]interface{}{},
	}
}

// Render renders all templates of the blueprint
func (bp *Blueprint) Render(ctx *Context) error {
	data := MergeJSON(ctx.Data, bp.Data)
	for _, filedesc := range bp.Files {
		if err := filedesc.Render(ctx, bp, data); err != nil {
			ctx.Log(bp.Name, "error", err.Error())
		}
	}
	return nil
}

// LoadTemplates parses all of the blueprint's templates
func (bp *Blueprint) LoadTemplates() error {
	templates, err := template.ParseGlob(bp.TemplateDir() + "/*")
	if err != nil {
		return err
	}
	bp.templates = templates
	return nil
}

// LoadManifest parses the blueprints serialized description from JSON.
func (bp *Blueprint) LoadManifest() error {
	manifestFilename := filepath.Join(bp.path, "manifest.json")
	manifest, err := os.Open(manifestFilename)
	if err != nil {
		return err
	}
	return json.NewDecoder(manifest).Decode(bp)
}

// ValueTypes maps strings to value constructors
var ValueTypes = map[string]func() interface{}{
	"integer": func() interface{} {
		return int64(0)
	},
	"string": func() interface{} {
		return ""
	},
	"float": func() interface{} {
		return float64(0.0)
	},
	"bool": func() interface{} {
		return false
	},
	"duration": func() interface{} {
		return time.Duration(0)
	},
	"identifier": func() interface{} {
		return &Identifier{}
	},
}

// DefineFlags initializes the flag set for the blueprint based on the
// defined blueprint arguments.
func (bp *Blueprint) DefineFlags() error {
	bp.flags = flag.NewFlagSet(bp.Name, flag.ContinueOnError)
	for _, argument := range bp.Args {
		newValue, found := ValueTypes[argument.Type]
		if !found {
			return fmt.Errorf("Unknown argument type: %s", argument.Type)
		}
		rawValue := newValue()
		doc := ""
		if argument.Doc != nil {
			doc = *argument.Doc
		}
		switch v := rawValue.(type) {
		case int64:
			bp.flags.Int64(argument.Name, 0, doc)
		case float64:
			bp.flags.Float64(argument.Name, 0.0, doc)
		case string:
			bp.flags.String(argument.Name, "", doc)
		case time.Duration:
			bp.flags.Duration(argument.Name, 0*time.Millisecond, doc)
		case bool:
			bp.flags.Bool(argument.Name, false, doc)
		default:
			flagValue, ok := v.(flag.Value)
			if !ok {
				return fmt.Errorf("Cannot handle flag value of type %T", v)
			}
			bp.flags.Var(flagValue, argument.Name, doc)
		}

	}

	return nil
}

// ParseArgs parses the provided arguments according to the blueprint's flag set.
func (bp *Blueprint) ParseArgs(args []string) error {
	if err := bp.flags.Parse(args); err != nil {
		return err
	}
	bp.flags.VisitAll(func(f *flag.Flag) {
		bp.Data[f.Name] = f.Value.(flag.Getter).Get()
	})

	return nil
}

// TemplateDir returns the path to the directory that holds the templates for this blueprint.
func (bp *Blueprint) TemplateDir() string {
	return filepath.Join(bp.path, "templates")
}

// LoadBlueprint loads a blueprint from disk.
func (ctx *Context) LoadBlueprint(blueprintName string) (*Blueprint, error) {
	directoryName := filepath.Join(ctx.BaseDir, "blueprints", blueprintName)
	blueprint := NewBlueprint(directoryName)
	if err := blueprint.LoadTemplates(); err != nil {
		return nil, err
	}
	if err := blueprint.LoadManifest(); err != nil {
		return nil, err
	}
	if err := blueprint.DefineFlags(); err != nil {
		return nil, err
	}
	return blueprint, nil
}

// IdentifierStyle described a way of encoding an identifier as a string.
type IdentifierStyle interface {
	Parse(string) []string
	Original([]string) string
	Upper([]string) string
	Lower([]string) string
	Title([]string) string
}

// SeparatedIdentifier describes identifiers using a separator to
// distinguish between words.
type SeparatedIdentifier struct {
	Separator string
}

// Parse parses the identifier by splitting it according to the separator
func (s *SeparatedIdentifier) Parse(identifier string) []string {
	return strings.Split(identifier, s.Separator)
}

// Original renders the identifier as a string using the original casing
func (s *SeparatedIdentifier) Original(constituents []string) string {
	return strings.Join(constituents, s.Separator)
}

// Upper renders the identifier by converting all consituents to upper case before joining them.
func (s *SeparatedIdentifier) Upper(constituents []string) string {
	out := bytes.NewBufferString("")
	for i, c := range constituents {
		fmt.Fprintf(out, "%s", strings.ToUpper(c))
		if i > 0 && i < len(constituents)-1 {
			fmt.Fprintf(out, "%s", s.Separator)
		}
	}

	return out.String()
}

// Lower renders the identifier by converting all constituents to lower case before joining them.
func (s *SeparatedIdentifier) Lower(constituents []string) string {
	out := bytes.NewBufferString("")
	for i, c := range constituents {
		fmt.Fprintf(out, "%s", strings.ToLower(c))
		if i > 0 && i < len(constituents)-1 {
			fmt.Fprintf(out, "%s", s.Separator)
		}
	}

	return out.String()
}

// Title is an alias for Original.
func (s *SeparatedIdentifier) Title(constituents []string) string { return s.Original(constituents) }

// CasedIdentifier described identifier that distinguish constituents using letter casing.
type CasedIdentifier struct{}

// Parse parses the identifier by splitting on upper-case letters
func (s *CasedIdentifier) Parse(identifier string) []string {
	constituent := []rune{}
	result := []string{}
	for _, r := range identifier {
		if unicode.IsUpper(r) {
			result = append(result, string(constituent))
			constituent = []rune{r}
		} else {
			constituent = append(constituent, r)
		}
	}
	if len(constituent) > 0 {
		result = append(result, string(constituent))
	}

	return result
}

// Original renders the identifier preserving original casing.
func (s *CasedIdentifier) Original(constituents []string) string {
	return strings.Join(constituents, "")
}

// Upper renders the identifier by converting the first letter of each constituent to upper case.
func (s *CasedIdentifier) Upper(constituents []string) string {
	out := bytes.NewBufferString("")
	for _, c := range constituents {
		if len(c) == 0 {
			continue
		}
		part := []rune(c)
		part[0] = unicode.ToUpper(part[0])
		fmt.Fprintf(out, "%s", string(part))
	}
	return out.String()
}

// Lower renders the identifier by converting the first consituent to lower case and the rest to title case.
func (s *CasedIdentifier) Lower(constituents []string) string {
	out := bytes.NewBufferString("")
	for i, c := range constituents {
		if len(c) == 0 {
			continue
		}
		part := []rune(c)
		if i == 0 {
			part[0] = unicode.ToLower(part[0])
		} else {
			part[0] = unicode.ToUpper(part[0])
		}
		fmt.Fprintf(out, "%s", string(part))
	}
	return out.String()
}

// Title is an alias for upper.
func (s *CasedIdentifier) Title(constituents []string) string {
	return s.Upper(constituents)
}

// Identifier reprents a programming language identifier that can be
// expressed in various casing styles.
type Identifier struct {
	Constituents []string
	Style        IdentifierStyle
}

// String renders the identifier in the style that was detected during creation of the identifier.
func (i *Identifier) String() string {
	if i.Style == nil {
		return ""
	}
	return i.Style.Original(i.Constituents)
}

// Set implements flag.Value by parsing the identifier
func (i *Identifier) Set(s string) error {
	newIdentifier := ParseIdentifier(s)
	*i = *newIdentifier
	return nil
}

// Get implements flag.Value by returning the identifier itself.
func (i *Identifier) Get() interface{} {
	return i
}

var (
	// SnakeCaseStyle is an identifier that separates words using underscores.
	SnakeCaseStyle = &SeparatedIdentifier{Separator: "_"}

	// LispCaseStyle is an identifier that separates words using hyphens.
	LispCaseStyle = &SeparatedIdentifier{Separator: "-"}

	// CamelCasedStyle is an identifier that separates words using letter casing.
	CamelCasedStyle = new(CasedIdentifier)
)

// ParseIdentifier analyzes a string as an identifier.
func ParseIdentifier(identifier string) *Identifier {
	result := &Identifier{
		Constituents: []string{},
		Style:        CamelCasedStyle,
	}

runes:
	for _, r := range identifier {
		switch r {
		case '-':
			result.Style = LispCaseStyle
			break runes
		case '_':
			result.Style = SnakeCaseStyle
			break runes
		}
	}

	result.Constituents = result.Style.Parse(identifier)
	return result
}

// Upper returns the identifier in upper case
func (i *Identifier) Upper() string {
	return i.Style.Upper(i.Constituents)
}

// Lower returns the identifier in title case
func (i *Identifier) Lower() string {
	return i.Style.Lower(i.Constituents)
}

// Title returns the identifier in title case
func (i *Identifier) Title() string {
	return i.Style.Title(i.Constituents)
}

// ToSnake converts the identifier into snake case
func (i *Identifier) ToSnake() *Identifier {
	return &Identifier{
		Constituents: i.Constituents,
		Style:        SnakeCaseStyle,
	}
}
