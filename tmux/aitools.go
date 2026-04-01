package tmux

// AITool describes a known AI CLI tool with its display metadata.
type AITool struct {
	Name  string
	Icon  string
	Color string // hex color code
}

// aiToolMap is the single source of truth for known AI CLI tools.
var aiToolMap = map[string]AITool{
	"claude": {Name: "claude", Icon: "✦", Color: "#F59E0B"},
	"codex":  {Name: "codex", Icon: "◈", Color: "#60A5FA"},
	"aider":  {Name: "aider", Icon: "⬡", Color: "#34D399"},
	"gemini": {Name: "gemini", Icon: "✧", Color: "#A78BFA"},
}

// IsAICommand reports whether cmd is a known AI CLI process.
func IsAICommand(cmd string) bool {
	_, ok := aiToolMap[cmd]
	return ok
}

// LookupAITool returns the AITool for the given command name, if known.
func LookupAITool(cmd string) (AITool, bool) {
	t, ok := aiToolMap[cmd]
	return t, ok
}
