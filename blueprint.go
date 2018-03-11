package dux

// Blueprint collects information about files to generate.
type Blueprint struct {
	Name  string            // The ID of the blueprint
	Files map[string]string // Files maps destination file names to template file names.
}

// DefineFile adds an entry for destinationFileName into the list of
// the blueprint's files.
func (bp *Blueprint) DefineFile(destinationFileName, templateFileName string) *Blueprint {
	if bp.Files == nil {
		bp.Files = map[string]string{}
	}
	bp.Files[destinationFileName] = templateFileName
	return bp
}
