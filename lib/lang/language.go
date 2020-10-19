package lang

// Language ...
type Language int

const (
	// JavaScript ...
	JavaScript = iota
	// Go ...
	Go
)

// LanguageProperty ...
type LanguageProperty struct {
	Extension string
	Commands  []string
}

// LanguageConfig ...
var LanguageConfig = map[Language]LanguageProperty{
	JavaScript: {
		Extension: "js",
		Commands:  []string{"node"},
	},
	Go: {
		Extension: "go",
		Commands:  []string{"go", "run"},
	},
}
