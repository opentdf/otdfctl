package tui

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	osprofiles "github.com/jrschumacher/go-osprofiles"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
)

const (
	helpText = "Type 'help' for commands, 'exit' or Ctrl+C to quit."
)

// Styles
var (
	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Bright blue
			Bold(true)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")) // Light gray

	outputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")) // Medium gray

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")). // Bright red
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("78")) // Bright green

	shellHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")). // Dark gray
			Italic(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")). // Bright pink
			Bold(true)
)

// HistoryEntry represents a command and its output
type HistoryEntry struct {
	Prompt  string
	Command string
	Output  string
	IsError bool
}

// PathSegment represents a single level in the navigation hierarchy
type PathSegment struct {
	Type  string // "root", "namespaces", "namespace", "attribute-definitions", "attribute", "attribute-values", "value", "registered-resources", "resource"
	Name  string // The name/id of the resource at this level
	ID    string // The actual ID if different from name
	Value string // Display value (for ls output)
}

// ShellContext tracks the current location in the resource hierarchy
type ShellContext struct {
	Path []PathSegment // Current path from root
}

// Shell is the Bubble Tea model for the interactive shell
type Shell struct {
	textInput       textinput.Model
	history         []HistoryEntry      // Display history (command + output pairs)
	commandHistory  []string            // Command history for up/down navigation
	historyPosition int                 // Current position in command history (-1 = not navigating)
	tempInput       string              // Temporary storage for current input when navigating history
	context         ShellContext
	handler         handlers.Handler
	ctx             context.Context
	profileName     string
	width           int
	height          int

	// Wizard mode
	wizard       *Wizard
	wizardActive bool

	// Delete wizard mode
	deleteWizard       *DeleteWizard
	deleteWizardActive bool
}

// NewShell creates a new Shell model
func NewShell(h handlers.Handler, profileName string) Shell {
	ti := textinput.New()
	ti.Placeholder = "type a command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80

	// Add welcome message to history
	history := []HistoryEntry{
		{
			Prompt:  "",
			Command: "",
			Output:  titleStyle.Render("Welcome to otdfctl interactive shell!") + "\n" + shellHelpStyle.Render(helpText),
			IsError: false,
		},
	}

	return Shell{
		textInput:       ti,
		history:         history,
		commandHistory:  []string{},
		historyPosition: -1,
		tempInput:       "",
		context: ShellContext{
			Path: []PathSegment{{Type: "root", Name: "/"}},
		},
		handler:     h,
		profileName: profileName,
		ctx:         context.Background(),
	}
}

// Init initializes the Bubble Tea model
func (s Shell) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (s Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// If wizard is active, delegate to it
	if s.wizardActive && s.wizard != nil {
		wizardModel, wizardCmd := s.wizard.Update(msg)
		s.wizard = wizardModel.(*Wizard)

		// Check if wizard completed
		if s.wizard.IsComplete() {
			s.wizardActive = false

			// Add result to history
			result := s.wizard.GetResult()
			isError := s.wizard.GetError() != ""
			if s.wizard.WasCancelled() {
				result = "Cancelled"
			}

			if result != "" {
				s.history = append(s.history, HistoryEntry{
					Prompt:  "",
					Command: "",
					Output:  result,
					IsError: isError,
				})
			}

			s.wizard = nil
			s.textInput.Focus()
			return s, textinput.Blink
		}

		return s, wizardCmd
	}

	// If delete wizard is active, delegate to it
	if s.deleteWizardActive && s.deleteWizard != nil {
		deleteModel, deleteCmd := s.deleteWizard.Update(msg)
		s.deleteWizard = deleteModel.(*DeleteWizard)

		// Check if delete wizard completed
		if s.deleteWizard.IsComplete() {
			s.deleteWizardActive = false

			// Add result to history
			result := s.deleteWizard.GetResult()
			isError := s.deleteWizard.GetError() != ""
			if s.deleteWizard.WasCancelled() {
				result = "Cancelled"
			}

			if result != "" {
				s.history = append(s.history, HistoryEntry{
					Prompt:  "",
					Command: "",
					Output:  result,
					IsError: isError,
				})
			}

			// If deletion was successful and we deleted current resource, navigate up
			if !s.deleteWizard.WasCancelled() && s.deleteWizard.GetError() == "" {
				// Go up one level after successful deletion
				if len(s.context.Path) > 1 {
					s.context.Path = s.context.Path[:len(s.context.Path)-1]
				}
			}

			s.deleteWizard = nil
			s.textInput.Focus()
			return s, textinput.Blink
		}

		return s, deleteCmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit

		case tea.KeyTab:
			// Handle tab completion
			s.handleTabCompletion()
			return s, nil

		case tea.KeyEnter:
			input := s.textInput.Value()
			s.textInput.SetValue("")

			if input == "" {
				return s, nil
			}

			// Add command to history (avoid duplicate consecutive commands)
			if len(s.commandHistory) == 0 || s.commandHistory[len(s.commandHistory)-1] != input {
				s.commandHistory = append(s.commandHistory, input)
			}

			// Reset history navigation
			s.historyPosition = -1
			s.tempInput = ""

			// Save current prompt
			currentPrompt := s.getPrompt()

			// Parse command
			parts := strings.Fields(input)
			command := parts[0]
			args := parts[1:]

			// Check if this is a create command that should launch a wizard
			if command == "create" {
				wizardCmd := s.startCreateWizard(args)
				if wizardCmd != nil {
					// Add command to history
					s.history = append(s.history, HistoryEntry{
						Prompt:  currentPrompt,
						Command: input,
						Output:  "",
						IsError: false,
					})
					return s, wizardCmd
				}
			}

			// Check if this is a delete/rm command that should launch delete wizard
			if command == "rm" || command == "delete" {
				deleteCmd := s.startDeleteWizard()
				if deleteCmd != nil {
					// Add command to history
					s.history = append(s.history, HistoryEntry{
						Prompt:  currentPrompt,
						Command: input,
						Output:  "",
						IsError: false,
					})
					return s, deleteCmd
				}
			}

			// Execute command
			output := s.executeCommand(command, args)

			// Check if it's an error message
			isError := strings.HasPrefix(output, "Error:") ||
			           strings.HasPrefix(output, "Unknown command:") ||
			           strings.HasPrefix(output, "No matching") ||
			           strings.HasPrefix(output, "No matches") ||
			           strings.Contains(output, "not found")

			// Add to history
			s.history = append(s.history, HistoryEntry{
				Prompt:  currentPrompt,
				Command: input,
				Output:  output,
				IsError: isError,
			})

			return s, nil

		case tea.KeyUp:
			// Navigate to previous command in history
			if len(s.commandHistory) == 0 {
				return s, nil
			}

			// If not currently navigating history, save current input
			if s.historyPosition == -1 {
				s.tempInput = s.textInput.Value()
				s.historyPosition = len(s.commandHistory)
			}

			// Move to previous command
			if s.historyPosition > 0 {
				s.historyPosition--
				s.textInput.SetValue(s.commandHistory[s.historyPosition])
				s.textInput.SetCursor(len(s.commandHistory[s.historyPosition]))
			}

			return s, nil

		case tea.KeyDown:
			// Navigate to next command in history
			if s.historyPosition == -1 {
				// Not navigating history, nothing to do
				return s, nil
			}

			s.historyPosition++

			// If we've gone past the end, restore the temp input
			if s.historyPosition >= len(s.commandHistory) {
				s.historyPosition = -1
				s.textInput.SetValue(s.tempInput)
				s.textInput.SetCursor(len(s.tempInput))
			} else {
				// Load next command
				s.textInput.SetValue(s.commandHistory[s.historyPosition])
				s.textInput.SetCursor(len(s.commandHistory[s.historyPosition]))
			}

			return s, nil
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

// startCreateWizard starts the appropriate wizard based on context and args
func (s *Shell) startCreateWizard(args []string) tea.Cmd {
	currentType := s.getCurrentType()

	// Determine what to create based on args or context
	var resourceType string
	if len(args) > 0 {
		resourceType = args[0]
	} else {
		// Infer from context
		switch currentType {
		case "root", "namespaces":
			resourceType = "namespace"
		case "namespace", "attribute-definitions":
			resourceType = "attribute"
		case "attribute", "attribute-values":
			resourceType = "value"
		default:
			return nil
		}
	}

	switch resourceType {
	case "namespace", "ns":
		s.wizard = NewNamespaceWizard(s.handler)
		s.wizardActive = true
		return s.wizard.Init()

	case "attribute", "attr":
		// Use namespace from context if available
		namespaceID := s.getNamespaceID()
		namespaceName := ""
		for _, seg := range s.context.Path {
			if seg.Type == "namespace" {
				namespaceName = seg.Name
				break
			}
		}
		s.wizard = NewAttributeWizard(s.handler, namespaceID, namespaceName)
		s.wizardActive = true
		return s.wizard.Init()

	case "value", "val":
		// Must be in an attribute context
		attributeID := s.getAttributeID()
		if attributeID == "" {
			s.history = append(s.history, HistoryEntry{
				Prompt:  "",
				Command: "",
				Output:  errorStyle.Render("Error: Navigate to an attribute first, or use 'create value' from within an attribute"),
				IsError: true,
			})
			return nil
		}

		attributeName := ""
		namespaceName := ""
		for _, seg := range s.context.Path {
			if seg.Type == "attribute" {
				attributeName = seg.Name
			}
			if seg.Type == "namespace" {
				namespaceName = seg.Name
			}
		}

		s.wizard = NewAttributeValueWizard(s.handler, attributeID, attributeName, namespaceName)
		s.wizardActive = true
		return s.wizard.Init()

	default:
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  errorStyle.Render(fmt.Sprintf("Unknown resource type: %s\nAvailable: namespace, attribute, value", resourceType)),
			IsError: true,
		})
		return nil
	}
}

// startDeleteWizard starts the delete wizard for the current resource
func (s *Shell) startDeleteWizard() tea.Cmd {
	currentType := s.getCurrentType()

	var resourceType, resourceName, resourceID, resourceFQN string
	var childCount int

	switch currentType {
	case "namespace":
		resourceType = "namespace"
		resourceID = s.getNamespaceID()
		for _, seg := range s.context.Path {
			if seg.Type == "namespace" {
				resourceName = seg.Name
				break
			}
		}
		// Get namespace to find FQN and count attributes
		if ns, err := s.handler.GetNamespace(s.ctx, resourceID); err == nil {
			resourceFQN = ns.GetFqn()
			// Count child attributes
			if attrs, err := s.handler.ListAttributes(s.ctx, 0, 1000, 0); err == nil {
				for _, attr := range attrs.GetAttributes() {
					if attr.GetNamespace().GetId() == resourceID {
						childCount++
					}
				}
			}
		}

	case "attribute":
		resourceType = "attribute"
		resourceID = s.getAttributeID()
		for _, seg := range s.context.Path {
			if seg.Type == "attribute" {
				resourceName = seg.Name
				break
			}
		}
		// Get attribute to find FQN and count values
		if attr, err := s.handler.GetAttribute(s.ctx, resourceID); err == nil {
			resourceFQN = attr.GetFqn()
			childCount = len(attr.GetValues())
		}

	case "value":
		resourceType = "value"
		resourceID = s.getValueID()
		for _, seg := range s.context.Path {
			if seg.Type == "value" {
				resourceName = seg.Name
				break
			}
		}
		// Get value to find FQN
		if val, err := s.handler.GetAttributeValue(s.ctx, resourceID); err == nil {
			resourceFQN = val.GetFqn()
		}

	default:
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  errorStyle.Render("Error: Navigate to a specific resource (namespace, attribute, or value) to delete it"),
			IsError: true,
		})
		return nil
	}

	if resourceID == "" {
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  errorStyle.Render("Error: Could not determine resource ID"),
			IsError: true,
		})
		return nil
	}

	s.deleteWizard = NewDeleteWizard(s.handler, resourceType, resourceName, resourceID, resourceFQN, childCount)
	s.deleteWizardActive = true
	return s.deleteWizard.Init()
}

// View renders the shell interface
func (s Shell) View() string {
	var sb strings.Builder

	// Render all history entries
	for _, entry := range s.history {
		if entry.Prompt != "" && entry.Command != "" {
			// Show prompt and command
			sb.WriteString(promptStyle.Render(entry.Prompt))
			sb.WriteString(commandStyle.Render(entry.Command))
			sb.WriteString("\n")
		}

		// Show output if present
		if entry.Output != "" {
			if entry.IsError {
				sb.WriteString(errorStyle.Render(entry.Output))
			} else if strings.Contains(entry.Output, "✓") || strings.Contains(entry.Output, "Success") {
				sb.WriteString(successStyle.Render(entry.Output))
			} else {
				sb.WriteString(outputStyle.Render(entry.Output))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("\n") // Extra spacing between commands
	}

	// If wizard is active, render it instead of the normal prompt
	if s.wizardActive && s.wizard != nil {
		sb.WriteString(s.wizard.View())
		return sb.String()
	}

	// If delete wizard is active, render it instead of the normal prompt
	if s.deleteWizardActive && s.deleteWizard != nil {
		sb.WriteString(s.deleteWizard.View())
		return sb.String()
	}

	// Current prompt and input
	sb.WriteString(promptStyle.Render(s.getPrompt()))
	sb.WriteString(s.textInput.View())
	sb.WriteString("\n\n")

	// Help text at bottom
	sb.WriteString(shellHelpStyle.Render(helpText))

	return sb.String()
}

// getPrompt returns the current prompt string
func (s Shell) getPrompt() string {
	path := s.getPathString()
	return fmt.Sprintf("%s:%s> ", s.profileName, path)
}

// getPathString converts the current path to a string representation
func (s Shell) getPathString() string {
	if len(s.context.Path) == 0 || (len(s.context.Path) == 1 && s.context.Path[0].Type == "root") {
		return "/"
	}

	var parts []string
	for i, seg := range s.context.Path {
		if i == 0 && seg.Type == "root" {
			continue
		}
		if seg.Name != "" {
			parts = append(parts, seg.Name)
		}
	}

	return "/" + strings.Join(parts, "/")
}

// executeCommand processes and executes the given command
func (s *Shell) executeCommand(command string, args []string) string {
	switch command {
	case "help", "?":
		return s.cmdHelp()
	case "ls", "list":
		return s.cmdLs()
	case "cd":
		return s.cmdCd(args)
	case "pwd":
		return s.cmdPwd()
	case "get", "show":
		return s.cmdGet()
	case "profile":
		return s.cmdProfile(args)
	case "keys":
		return s.cmdKeys()
	case "key":
		return s.cmdKey(args)
	case "clear", "cls":
		// Clear history
		s.history = []HistoryEntry{}
		return ""
	case "exit", "quit":
		// This should trigger quit, but for now just return a message
		return "Use Ctrl+C to exit"
	default:
		return fmt.Sprintf("Unknown command: %s\nType 'help' for available commands.", command)
	}
}

// cmdHelp displays help information based on current context
func (s *Shell) cmdHelp() string {
	var sb strings.Builder

	sb.WriteString("Available commands:\n\n")
	sb.WriteString("Navigation:\n")
	sb.WriteString("  ls, list        List items in current location\n")
	sb.WriteString("  cd <path>       Change directory\n")
	sb.WriteString("  cd ..           Go up one level\n")
	sb.WriteString("  cd /            Go to root\n")
	sb.WriteString("  pwd             Print working directory\n\n")

	sb.WriteString("Information:\n")
	sb.WriteString("  help, ?         Show this help\n")
	sb.WriteString("  clear, cls      Clear output\n\n")

	sb.WriteString("Profile Management:\n")
	sb.WriteString("  profile         Show current profile\n")
	sb.WriteString("  profile list    List all profiles\n")
	sb.WriteString("  profile use <name>  Switch to a different profile\n\n")

	sb.WriteString("Resource Details:\n")
	sb.WriteString("  get, show       Show detailed information about current resource\n\n")

	sb.WriteString("Resource Creation:\n")
	sb.WriteString("  create          Start creation wizard (context-aware)\n")
	sb.WriteString("  create namespace    Create a new namespace\n")
	sb.WriteString("  create attribute    Create a new attribute\n")
	sb.WriteString("  create value        Create a new attribute value\n\n")

	sb.WriteString("Resource Deletion:\n")
	sb.WriteString("  rm, delete      Delete current resource (with confirmation)\n")
	sb.WriteString("                  Offers deactivate (safe) or permanent delete\n\n")

	sb.WriteString("Shell:\n")
	sb.WriteString("  exit, quit      Exit shell (or use Ctrl+C)\n")

	// Context-specific hints
	currentType := s.getCurrentType()
	switch currentType {
	case "root", "namespaces":
		sb.WriteString("\nHint: Use 'create' to start a namespace creation wizard\n")
	case "namespace":
		sb.WriteString("\nHint: Use 'create' for attributes, 'rm' to delete this namespace\n")
	case "attribute-definitions":
		sb.WriteString("\nHint: Use 'create' to start an attribute creation wizard\n")
	case "attribute":
		sb.WriteString("\nHint: Use 'create' for values, 'rm' to delete this attribute\n")
	case "attribute-values":
		sb.WriteString("\nHint: Use 'create' to start a value creation wizard\n")
	case "value":
		sb.WriteString("\nHint: Use 'rm' to delete this value\n")
	}

	return sb.String()
}

// cmdPwd displays the current path
func (s *Shell) cmdPwd() string {
	return s.getPathString()
}

// cmdProfile handles profile management commands
func (s *Shell) cmdProfile(args []string) string {
	if len(args) == 0 {
		// Show current profile
		return fmt.Sprintf("Current profile: %s", successStyle.Render(s.profileName))
	}

	subcommand := args[0]
	switch subcommand {
	case "list", "ls":
		return s.profileList()
	case "use", "switch":
		if len(args) < 2 {
			return errorStyle.Render("Error: profile name required\nUsage: profile use <name>")
		}
		return s.profileSwitch(args[1])
	case "current", "show":
		return fmt.Sprintf("Current profile: %s", successStyle.Render(s.profileName))
	default:
		return errorStyle.Render(fmt.Sprintf("Unknown profile subcommand: %s\nAvailable: list, use <name>, current", subcommand))
	}
}

// profileList lists all available profiles
func (s *Shell) profileList() string {
	profiler, err := profiles.CreateProfiler(profiles.ProfileDriverFileSystem)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error loading profiles: %v", err))
	}

	globalCfg := osprofiles.GetGlobalConfig(profiler)
	defaultProfile := globalCfg.GetDefaultProfile()
	allProfiles := osprofiles.ListProfiles(profiler)

	if len(allProfiles) == 0 {
		return "No profiles found"
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Available profiles:") + "\n\n")

	for _, p := range allProfiles {
		if p == s.profileName {
			sb.WriteString(successStyle.Render("> " + p + " (current)") + "\n")
		} else if p == defaultProfile {
			sb.WriteString(outputStyle.Render("* " + p + " (default)") + "\n")
		} else {
			sb.WriteString(outputStyle.Render("  " + p) + "\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(shellHelpStyle.Render("> current  * default"))

	return sb.String()
}

// profileSwitch switches to a different profile
func (s *Shell) profileSwitch(profileName string) string {
	// Load the new profile
	store, err := profiles.LoadOtdfctlProfileStore(profiles.ProfileDriverFileSystem, profileName)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error loading profile '%s': %v", profileName, err))
	}

	// Create a new handler with the new profile
	newHandler, err := handlers.New(handlers.WithProfile(store))
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error initializing handler for profile '%s': %v", profileName, err))
	}

	// Update the shell with the new profile and handler
	s.profileName = profileName
	s.handler = newHandler

	// Reset to root when switching profiles
	s.context.Path = []PathSegment{{Type: "root", Name: "/"}}

	return successStyle.Render(fmt.Sprintf("✓ Switched to profile: %s", profileName))
}

// cmdGet displays detailed information about the current resource
func (s *Shell) cmdGet() string {
	currentType := s.getCurrentType()

	switch currentType {
	case "root":
		return errorStyle.Render("Error: Nothing to get at root level\nNavigate to a specific resource first")
	case "namespaces":
		return errorStyle.Render("Error: Navigate into a specific namespace first")
	case "namespace":
		return s.getNamespace()
	case "attribute-definitions":
		return errorStyle.Render("Error: Navigate into a specific attribute first")
	case "attribute":
		return s.getAttribute()
	case "attribute-values":
		return errorStyle.Render("Error: Navigate into a specific value first")
	case "value":
		return s.getValue()
	case "registered-resources":
		return errorStyle.Render("Error: Navigate into a specific resource first")
	case "resource":
		return s.getResource()
	default:
		return errorStyle.Render(fmt.Sprintf("Error: Cannot get details for type: %s", currentType))
	}
}

// getNamespace displays detailed information about a namespace
func (s *Shell) getNamespace() string {
	namespaceID := s.getNamespaceID()
	if namespaceID == "" {
		return errorStyle.Render("Error: Could not determine namespace ID")
	}

	ns, err := s.handler.GetNamespace(s.ctx, namespaceID)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error fetching namespace: %v", err))
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Namespace Details") + "\n\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Name"), successStyle.Render(ns.GetName())))
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("ID"), shellHelpStyle.Render(ns.GetId())))

	if ns.GetFqn() != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("FQN"), outputStyle.Render(ns.GetFqn())))
	}

	if ns.GetActive() != nil {
		sb.WriteString(fmt.Sprintf("%s: %v\n", outputStyle.Render("Active"), ns.GetActive().GetValue()))
	}

	return sb.String()
}

// getAttribute displays detailed information about an attribute definition
func (s *Shell) getAttribute() string {
	attrID := s.getAttributeID()
	if attrID == "" {
		return errorStyle.Render("Error: Could not determine attribute ID")
	}

	attr, err := s.handler.GetAttribute(s.ctx, attrID)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error fetching attribute: %v", err))
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Attribute Definition Details") + "\n\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Name"), successStyle.Render(attr.GetName())))
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("ID"), shellHelpStyle.Render(attr.GetId())))
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Namespace"), outputStyle.Render(attr.GetNamespace().GetName())))
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Rule"), outputStyle.Render(attr.GetRule().String())))

	if attr.GetActive() != nil {
		sb.WriteString(fmt.Sprintf("%s: %v\n", outputStyle.Render("Active"), attr.GetActive().GetValue()))
	}

	// Show values if any
	values := attr.GetValues()
	if len(values) > 0 {
		sb.WriteString(fmt.Sprintf("\n%s: %d\n", outputStyle.Render("Values"), len(values)))
		for _, v := range values {
			sb.WriteString(fmt.Sprintf("  • %s\n", outputStyle.Render(v.GetValue())))
		}
	}

	return sb.String()
}

// getValue displays detailed information about an attribute value
func (s *Shell) getValue() string {
	valueID := s.getValueID()
	if valueID == "" {
		return errorStyle.Render("Error: Could not determine value ID")
	}

	val, err := s.handler.GetAttributeValue(s.ctx, valueID)
	if err != nil {
		return errorStyle.Render(fmt.Sprintf("Error fetching value: %v", err))
	}

	// Get attribute name from path
	attributeName := ""
	for _, seg := range s.context.Path {
		if seg.Type == "attribute" {
			attributeName = seg.Name
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Attribute Value Details") + "\n\n")
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Value"), successStyle.Render(val.GetValue())))
	sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("ID"), shellHelpStyle.Render(val.GetId())))

	if attributeName != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Attribute"), outputStyle.Render(attributeName)))
	} else if val.GetAttribute() != nil {
		// Fallback to API response if available
		sb.WriteString(fmt.Sprintf("%s: %s\n", outputStyle.Render("Attribute"), outputStyle.Render(val.GetAttribute().GetName())))
	}

	if val.GetActive() != nil {
		sb.WriteString(fmt.Sprintf("%s: %v\n", outputStyle.Render("Active"), val.GetActive().GetValue()))
	}

	return sb.String()
}

// getResource displays detailed information about a registered resource (placeholder)
func (s *Shell) getResource() string {
	return shellHelpStyle.Render("Resource details not yet implemented")
}

// cmdLs lists items at the current location
func (s *Shell) cmdLs() string {
	currentType := s.getCurrentType()

	switch currentType {
	case "root":
		return s.lsRoot()
	case "namespaces":
		return s.lsNamespaces()
	case "namespace":
		return s.lsNamespace()
	case "attribute-definitions":
		return s.lsAttributeDefinitions()
	case "attribute":
		return s.lsAttribute()
	case "attribute-values":
		return s.lsAttributeValues()
	case "registered-resources":
		return s.lsRegisteredResources()
	default:
		return fmt.Sprintf("Cannot list items at this location (type: %s)", currentType)
	}
}

// lsRoot lists top-level categories
func (s *Shell) lsRoot() string {
	return "namespaces/\nregistered-resources/"
}

// lsNamespaces lists all namespaces
func (s *Shell) lsNamespaces() string {
	resp, err := s.handler.ListNamespaces(s.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
	if err != nil {
		return fmt.Sprintf("Error listing namespaces: %v", err)
	}

	if len(resp.GetNamespaces()) == 0 {
		return "No namespaces found"
	}

	var sb strings.Builder
	for _, ns := range resp.GetNamespaces() {
		name := ns.GetName()
		if name == "" {
			name = ns.GetId()
		}
		sb.WriteString(fmt.Sprintf("%s/\n", name))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// lsNamespace lists contents of a namespace
func (s *Shell) lsNamespace() string {
	return "attribute-definitions/"
}

// lsAttributeDefinitions lists attributes in a namespace
func (s *Shell) lsAttributeDefinitions() string {
	// Get namespace ID from path
	namespaceID := s.getNamespaceID()
	if namespaceID == "" {
		return "Error: could not determine namespace"
	}

	resp, err := s.handler.ListAttributes(s.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
	if err != nil {
		return fmt.Sprintf("Error listing attributes: %v", err)
	}

	// Filter by namespace
	var sb strings.Builder
	count := 0
	for _, attr := range resp.GetAttributes() {
		if attr.GetNamespace().GetId() == namespaceID || attr.GetNamespace().GetName() == namespaceID {
			name := attr.GetName()
			if name == "" {
				name = attr.GetId()
			}
			sb.WriteString(fmt.Sprintf("%s/\n", name))
			count++
		}
	}

	if count == 0 {
		return "No attributes found in this namespace"
	}

	return strings.TrimRight(sb.String(), "\n")
}

// lsAttribute lists contents of an attribute
func (s *Shell) lsAttribute() string {
	return "attribute-values/"
}

// lsAttributeValues lists values of an attribute
func (s *Shell) lsAttributeValues() string {
	// Get attribute ID from path
	attributeID := s.getAttributeID()
	if attributeID == "" {
		return "Error: could not determine attribute"
	}

	attr, err := s.handler.GetAttribute(s.ctx, attributeID)
	if err != nil {
		return fmt.Sprintf("Error getting attribute: %v", err)
	}

	values := attr.GetValues()
	if len(values) == 0 {
		return "No values found for this attribute"
	}

	var sb strings.Builder
	for _, val := range values {
		name := val.GetValue()
		if name == "" {
			name = val.GetId()
		}
		sb.WriteString(fmt.Sprintf("%s\n", name))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// lsRegisteredResources lists registered resources
func (s *Shell) lsRegisteredResources() string {
	// TODO: Implement when SDK supports registered resources listing
	return "Registered resources listing not yet implemented"
}

// cmdCd changes the current directory
func (s *Shell) cmdCd(args []string) string {
	if len(args) == 0 {
		// cd with no args goes to root
		s.context.Path = []PathSegment{{Type: "root", Name: "/"}}
		return ""
	}

	path := args[0]

	// Handle special cases
	switch path {
	case "/":
		s.context.Path = []PathSegment{{Type: "root", Name: "/"}}
		return ""
	case "..":
		return s.cdUp()
	case ".":
		return ""
	default:
		return s.cdTo(path)
	}
}

// cdUp goes up one level in the hierarchy
func (s *Shell) cdUp() string {
	if len(s.context.Path) <= 1 {
		return "Already at root"
	}
	s.context.Path = s.context.Path[:len(s.context.Path)-1]
	return ""
}

// cdTo navigates to a specific path (relative or absolute)
func (s *Shell) cdTo(path string) string {
	// Handle absolute paths
	if strings.HasPrefix(path, "/") {
		return s.cdAbsolute(path)
	}

	// Handle relative paths
	return s.cdRelative(path)
}

// cdAbsolute navigates to an absolute path
func (s *Shell) cdAbsolute(path string) string {
	// Start from root
	s.context.Path = []PathSegment{{Type: "root", Name: "/"}}

	// Remove leading slash and split
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	return s.cdParts(parts)
}

// cdRelative navigates relative to current location
func (s *Shell) cdRelative(path string) string {
	parts := strings.Split(path, "/")
	return s.cdParts(parts)
}

// cdParts navigates through a series of path parts
func (s *Shell) cdParts(parts []string) string {
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			if err := s.cdUp(); err != "" {
				return err
			}
			continue
		}

		// Try to cd into this part
		if err := s.cdInto(part); err != "" {
			return err
		}
	}
	return ""
}

// cdInto navigates into a specific child from current location
func (s *Shell) cdInto(name string) string {
	currentType := s.getCurrentType()

	switch currentType {
	case "root":
		return s.cdIntoRoot(name)
	case "namespaces":
		return s.cdIntoNamespaces(name)
	case "namespace":
		return s.cdIntoNamespace(name)
	case "attribute-definitions":
		return s.cdIntoAttributeDefinitions(name)
	case "attribute":
		return s.cdIntoAttribute(name)
	case "attribute-values":
		return s.cdIntoAttributeValue(name)
	case "registered-resources":
		return s.cdIntoRegisteredResources(name)
	default:
		return fmt.Sprintf("Cannot navigate from current location (type: %s)", currentType)
	}
}

// cdIntoRoot navigates from root into a top-level category
func (s *Shell) cdIntoRoot(name string) string {
	name = strings.TrimSuffix(name, "/")

	switch name {
	case "namespaces":
		s.context.Path = append(s.context.Path, PathSegment{Type: "namespaces", Name: "namespaces"})
		return ""
	case "registered-resources":
		s.context.Path = append(s.context.Path, PathSegment{Type: "registered-resources", Name: "registered-resources"})
		return ""
	default:
		return fmt.Sprintf("Unknown directory: %s (available: namespaces, registered-resources)", name)
	}
}

// cdIntoNamespaces navigates into a specific namespace
func (s *Shell) cdIntoNamespaces(name string) string {
	name = strings.TrimSuffix(name, "/")

	// Try to get the namespace to verify it exists
	ns, err := s.handler.GetNamespace(s.ctx, name)
	if err != nil {
		// If direct lookup failed, try to find by listing all namespaces
		resp, listErr := s.handler.ListNamespaces(s.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
		if listErr != nil {
			return fmt.Sprintf("Namespace not found: %s", name)
		}

		// Find namespace with matching name
		var found *policy.Namespace
		for _, namespace := range resp.GetNamespaces() {
			if namespace.GetName() == name || namespace.GetId() == name {
				found = namespace
				break
			}
		}

		if found == nil {
			return fmt.Sprintf("Namespace not found: %s", name)
		}

		ns = found
	}

	s.context.Path = append(s.context.Path, PathSegment{
		Type: "namespace",
		Name: ns.GetName(),
		ID:   ns.GetId(),
	})
	return ""
}

// cdIntoNamespace navigates from namespace into a subcategory
func (s *Shell) cdIntoNamespace(name string) string {
	name = strings.TrimSuffix(name, "/")

	switch name {
	case "attribute-definitions":
		s.context.Path = append(s.context.Path, PathSegment{Type: "attribute-definitions", Name: "attribute-definitions"})
		return ""
	default:
		return fmt.Sprintf("Unknown directory: %s (available: attribute-definitions)", name)
	}
}

// cdIntoAttributeDefinitions navigates into a specific attribute
func (s *Shell) cdIntoAttributeDefinitions(name string) string {
	name = strings.TrimSuffix(name, "/")

	namespaceID := s.getNamespaceID()
	if namespaceID == "" {
		return "Error: could not determine namespace"
	}

	// Try to get the attribute to verify it exists
	// Build FQN: namespace/attribute
	fqn := fmt.Sprintf("https://%s/attr/%s", namespaceID, name)

	attr, err := s.handler.GetAttribute(s.ctx, fqn)
	if err != nil {
		// Try with just the name
		attr, err = s.handler.GetAttribute(s.ctx, name)
		if err != nil {
			// If direct lookups failed, list all attributes and find by name
			resp, listErr := s.handler.ListAttributes(s.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
			if listErr != nil {
				return fmt.Sprintf("Attribute not found: %s", name)
			}

			// Find attribute with matching name in the current namespace
			var found *policy.Attribute
			for _, attribute := range resp.GetAttributes() {
				nsID := attribute.GetNamespace().GetId()
				nsName := attribute.GetNamespace().GetName()
				attrName := attribute.GetName()
				attrID := attribute.GetId()

				// Check if this attribute is in the current namespace and matches the name
				if (nsID == namespaceID || nsName == namespaceID) &&
					(attrName == name || attrID == name) {
					found = attribute
					break
				}
			}

			if found == nil {
				return fmt.Sprintf("Attribute not found: %s", name)
			}

			attr = found
		}
	}

	s.context.Path = append(s.context.Path, PathSegment{
		Type: "attribute",
		Name: attr.GetName(),
		ID:   attr.GetId(),
	})
	return ""
}

// cdIntoAttribute navigates from attribute into subcategory
func (s *Shell) cdIntoAttribute(name string) string {
	name = strings.TrimSuffix(name, "/")

	switch name {
	case "attribute-values":
		s.context.Path = append(s.context.Path, PathSegment{Type: "attribute-values", Name: "attribute-values"})
		return ""
	default:
		return fmt.Sprintf("Unknown directory: %s (available: attribute-values)", name)
	}
}

// cdIntoAttributeValue navigates into a specific attribute value
func (s *Shell) cdIntoAttributeValue(name string) string {
	name = strings.TrimSuffix(name, "/")

	// Get the attribute ID from path
	attributeID := s.getAttributeID()
	if attributeID == "" {
		return "Error: could not determine attribute"
	}

	// Get attribute to access its values
	attr, err := s.handler.GetAttribute(s.ctx, attributeID)
	if err != nil {
		return fmt.Sprintf("Error getting attribute: %v", err)
	}

	// Find the value by name or ID
	var foundValue *policy.Value
	for _, val := range attr.GetValues() {
		if val.GetValue() == name || val.GetId() == name {
			foundValue = val
			break
		}
	}

	if foundValue == nil {
		return fmt.Sprintf("Value not found: %s", name)
	}

	// Add value to path
	s.context.Path = append(s.context.Path, PathSegment{
		Type: "value",
		Name: foundValue.GetValue(),
		ID:   foundValue.GetId(),
	})
	return ""
}

// cdIntoRegisteredResources navigates into a specific resource
func (s *Shell) cdIntoRegisteredResources(name string) string {
	// TODO: Implement when SDK supports registered resources
	return "Registered resources navigation not yet implemented"
}

// getCurrentType returns the type of the current location
func (s *Shell) getCurrentType() string {
	if len(s.context.Path) == 0 {
		return "root"
	}
	return s.context.Path[len(s.context.Path)-1].Type
}

// getNamespaceID gets the namespace ID from the current path
func (s *Shell) getNamespaceID() string {
	for _, seg := range s.context.Path {
		if seg.Type == "namespace" {
			if seg.ID != "" {
				return seg.ID
			}
			return seg.Name
		}
	}
	return ""
}

// getAttributeID gets the attribute ID from the current path
func (s *Shell) getAttributeID() string {
	for _, seg := range s.context.Path {
		if seg.Type == "attribute" {
			if seg.ID != "" {
				return seg.ID
			}
			return seg.Name
		}
	}
	return ""
}

// getValueID gets the value ID from the current path
func (s *Shell) getValueID() string {
	for _, seg := range s.context.Path {
		if seg.Type == "value" {
			if seg.ID != "" {
				return seg.ID
			}
			return seg.Name
		}
	}
	return ""
}

// handleTabCompletion handles tab key press for autocompletion
func (s *Shell) handleTabCompletion() {
	input := s.textInput.Value()
	if input == "" {
		// Show available commands
		commands := s.getAvailableCommands()
		msg := shellHelpStyle.Render("Available commands: ") + outputStyle.Render(strings.Join(commands, ", "))
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  msg,
			IsError: false,
		})
		return
	}

	// Parse the current input
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	// If cursor is at the end and input ends with space, we're completing a new argument
	cursorAtEnd := s.textInput.Position() == len(input)
	hasTrailingSpace := strings.HasSuffix(input, " ")

	if len(parts) == 1 && !hasTrailingSpace {
		// Completing command name
		s.completeCommand(parts[0])
	} else {
		// Completing arguments (paths for cd, ls, etc.)
		command := parts[0]
		var prefix string
		if hasTrailingSpace || !cursorAtEnd {
			prefix = ""
		} else {
			prefix = parts[len(parts)-1]
		}
		s.completeArgument(command, prefix, parts)
	}
}

// completeCommand completes a partially typed command
func (s *Shell) completeCommand(prefix string) {
	commands := s.getAvailableCommands()
	matches := findCompletions(prefix, commands)

	if len(matches) == 0 {
		msg := errorStyle.Render("No matching commands")
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  msg,
			IsError: true,
		})
	} else if len(matches) == 1 {
		// Single match - complete it
		s.textInput.SetValue(matches[0] + " ")
		s.textInput.SetCursor(len(matches[0]) + 1)
	} else {
		// Multiple matches - show them
		msg := shellHelpStyle.Render("Possible commands:\n") + outputStyle.Render(strings.Join(matches, "\n"))
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  msg,
			IsError: false,
		})
		// Set to common prefix
		commonPrefix := findCommonPrefix(matches)
		if len(commonPrefix) > len(prefix) {
			s.textInput.SetValue(commonPrefix)
			s.textInput.SetCursor(len(commonPrefix))
		}
	}
}

// completeArgument completes arguments for commands (typically paths)
func (s *Shell) completeArgument(command string, prefix string, parts []string) {
	var matches []string

	// Handle different commands
	switch command {
	case "cd", "ls":
		// Get available items at current location
		items := s.getAvailableItems()
		matches = findCompletions(prefix, items)
	case "profile":
		// If this is the first argument after "profile", complete subcommands
		if len(parts) == 2 {
			subcommands := []string{"list", "use", "current"}
			matches = findCompletions(prefix, subcommands)
		} else if len(parts) == 3 && parts[1] == "use" {
			// If this is "profile use <prefix>", complete with profile names
			profiler, err := profiles.CreateProfiler(profiles.ProfileDriverFileSystem)
			if err == nil {
				profileNames := osprofiles.ListProfiles(profiler)
				matches = findCompletions(prefix, profileNames)
			}
		}
	case "create":
		// Complete resource types
		if len(parts) == 2 {
			resourceTypes := []string{"namespace", "attribute", "value"}
			matches = findCompletions(prefix, resourceTypes)
		}
	default:
		return
	}

	if len(matches) == 0 {
		msg := errorStyle.Render("No matches")
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  msg,
			IsError: true,
		})
	} else if len(matches) == 1 {
		// Single match - complete it
		// Rebuild command with completion
		newInput := command + " " + matches[0]
		s.textInput.SetValue(newInput)
		s.textInput.SetCursor(len(newInput))
	} else {
		// Multiple matches - show them
		msg := shellHelpStyle.Render("Possible completions:\n") + outputStyle.Render(strings.Join(matches, "\n"))
		s.history = append(s.history, HistoryEntry{
			Prompt:  "",
			Command: "",
			Output:  msg,
			IsError: false,
		})
		// Set to common prefix
		commonPrefix := findCommonPrefix(matches)
		if len(commonPrefix) > len(prefix) {
			newInput := command + " " + commonPrefix
			s.textInput.SetValue(newInput)
			s.textInput.SetCursor(len(newInput))
		}
	}
}

// getAvailableCommands returns list of available commands based on context
func (s *Shell) getAvailableCommands() []string {
	// Base commands available everywhere
	commands := []string{"help", "ls", "cd", "pwd", "clear", "profile", "create", "exit", "quit"}

	// Add get/show and rm/delete commands when on a specific resource
	currentType := s.getCurrentType()
	if currentType != "root" && currentType != "namespaces" && currentType != "attribute-definitions" &&
		currentType != "attribute-values" && currentType != "registered-resources" {
		commands = append(commands, "get", "show", "rm", "delete")
	}

	return commands
}

// getAvailableItems returns list of items at current location (for path completion)
func (s *Shell) getAvailableItems() []string {
	currentType := s.getCurrentType()

	switch currentType {
	case "root":
		return []string{"namespaces", "registered-resources"}
	case "namespaces":
		return s.getNamespaceNames()
	case "namespace":
		return []string{"attribute-definitions"}
	case "attribute-definitions":
		return s.getAttributeNames()
	case "attribute":
		return []string{"attribute-values"}
	case "attribute-values":
		return s.getAttributeValueNames()
	case "registered-resources":
		// TODO: Implement when registered resources are supported
		return []string{}
	default:
		return []string{}
	}
}

// getNamespaceNames fetches and returns namespace names
func (s *Shell) getNamespaceNames() []string {
	resp, err := s.handler.ListNamespaces(s.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, ns := range resp.GetNamespaces() {
		name := ns.GetName()
		if name == "" {
			name = ns.GetId()
		}
		names = append(names, name)
	}
	return names
}

// getAttributeNames fetches and returns attribute names for current namespace
func (s *Shell) getAttributeNames() []string {
	namespaceID := s.getNamespaceID()
	if namespaceID == "" {
		return []string{}
	}

	resp, err := s.handler.ListAttributes(s.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, attr := range resp.GetAttributes() {
		if attr.GetNamespace().GetId() == namespaceID || attr.GetNamespace().GetName() == namespaceID {
			name := attr.GetName()
			if name == "" {
				name = attr.GetId()
			}
			names = append(names, name)
		}
	}
	return names
}

// getAttributeValueNames fetches and returns attribute value names
func (s *Shell) getAttributeValueNames() []string {
	attributeID := s.getAttributeID()
	if attributeID == "" {
		return []string{}
	}

	attr, err := s.handler.GetAttribute(s.ctx, attributeID)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, val := range attr.GetValues() {
		name := val.GetValue()
		if name == "" {
			name = val.GetId()
		}
		names = append(names, name)
	}
	return names
}

// findCompletions finds all items that start with the given prefix
func findCompletions(prefix string, items []string) []string {
	if prefix == "" {
		return items
	}

	var matches []string
	lowerPrefix := strings.ToLower(prefix)
	for _, item := range items {
		if strings.HasPrefix(strings.ToLower(item), lowerPrefix) {
			matches = append(matches, item)
		}
	}
	return matches
}

// findCommonPrefix finds the longest common prefix among strings
func findCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

// StartShell starts the interactive shell
func StartShell(h handlers.Handler, profileName string) {
	// Clear the terminal screen before starting the shell
	fmt.Print("\033[2J\033[H")

	p := tea.NewProgram(NewShell(h, profileName))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
