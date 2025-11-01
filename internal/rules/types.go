package rules

type ActionSpec struct {
	Type    string `yaml:"type"`    // "move" | "copy" | "delete" | "rename"
	Target  string `yaml:"target"`  // for move/copy
	Pattern string `yaml:"pattern"` // for rename
}
