package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/platform/protocol/go/common"
)

// WizardType defines the type of resource being created
type WizardType int

const (
	WizardTypeNamespace WizardType = iota
	WizardTypeAttribute
	WizardTypeAttributeValue
	WizardTypeKasRegistry
	WizardTypeProvider
)

// WizardStep defines the current step in the wizard
type WizardStep int

const (
	StepInit WizardStep = iota
	StepInput
	StepSelect
	StepConfirm
	StepExecuting
	StepComplete
	StepError
)

// SelectOption represents an option in a selection list
type SelectOption struct {
	Label string
	Value string
	ID    string
}

// WizardField represents a field to be filled in the wizard
type WizardField struct {
	Name        string
	Label       string
	Description string
	Value       string
	Required    bool
	FieldType   string // "input", "select", "multiselect"
	Options     []SelectOption
	Selected    int // For select fields
}

// Wizard is the Bubble Tea model for resource creation wizards
type Wizard struct {
	wizardType    WizardType
	step          WizardStep
	fields        []WizardField
	currentField  int
	textInput     textinput.Model
	selectIndex   int
	handler       handlers.Handler
	ctx           context.Context
	error         string
	result        string
	cancelled     bool

	// For multi-value input (like attribute values)
	multiValues   []string
	addingAnother bool

	// Context from shell (for context-aware creation)
	namespaceID   string
	namespaceName string
	attributeID   string
	attributeName string
}

// WizardResult is returned when the wizard completes
type WizardResult struct {
	Success   bool
	Message   string
	Cancelled bool
}

// Wizard styles
var (
	wizardTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Bold(true).
				MarginBottom(1)

	wizardLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true)

	wizardDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)

	wizardInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	wizardSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("78")).
				Bold(true)

	wizardUnselectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("250"))

	wizardErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("203")).
				Bold(true)

	wizardSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("78")).
				Bold(true)

	wizardHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))
)

// NewNamespaceWizard creates a wizard for creating a namespace
func NewNamespaceWizard(h handlers.Handler) *Wizard {
	ti := textinput.New()
	ti.Placeholder = "e.g., example.com"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return &Wizard{
		wizardType: WizardTypeNamespace,
		step:       StepInput,
		handler:    h,
		ctx:        context.Background(),
		textInput:  ti,
		fields: []WizardField{
			{
				Name:        "name",
				Label:       "Namespace Name",
				Description: "Enter the namespace name (e.g., example.com)",
				Required:    true,
				FieldType:   "input",
			},
		},
		currentField: 0,
	}
}

// NewAttributeWizard creates a wizard for creating an attribute
func NewAttributeWizard(h handlers.Handler, namespaceID, namespaceName string) *Wizard {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	w := &Wizard{
		wizardType:    WizardTypeAttribute,
		step:          StepInit,
		handler:       h,
		ctx:           context.Background(),
		textInput:     ti,
		namespaceID:   namespaceID,
		namespaceName: namespaceName,
		fields:        []WizardField{},
		currentField:  0,
		multiValues:   []string{},
	}

	// Build fields based on context
	w.buildAttributeFields()

	return w
}

// NewAttributeValueWizard creates a wizard for creating an attribute value
func NewAttributeValueWizard(h handlers.Handler, attributeID, attributeName, namespaceName string) *Wizard {
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return &Wizard{
		wizardType:    WizardTypeAttributeValue,
		step:          StepInput,
		handler:       h,
		ctx:           context.Background(),
		textInput:     ti,
		attributeID:   attributeID,
		attributeName: attributeName,
		namespaceName: namespaceName,
		fields: []WizardField{
			{
				Name:        "value",
				Label:       "Attribute Value",
				Description: fmt.Sprintf("Enter a value for attribute '%s'", attributeName),
				Required:    true,
				FieldType:   "input",
			},
		},
		currentField: 0,
	}
}

// NewKasRegistryWizard creates a wizard for creating a KAS Registry entry
func NewKasRegistryWizard(h handlers.Handler) *Wizard {
	ti := textinput.New()
	ti.Placeholder = "e.g., https://kas.example.com"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return &Wizard{
		wizardType: WizardTypeKasRegistry,
		step:       StepInput,
		handler:    h,
		ctx:        context.Background(),
		textInput:  ti,
		fields: []WizardField{
			{
				Name:        "uri",
				Label:       "KAS URI",
				Description: "Enter the URI for the Key Access Server (e.g., https://kas.example.com)",
				Required:    true,
				FieldType:   "input",
			},
			{
				Name:        "name",
				Label:       "Name",
				Description: "Enter a name for this KAS (optional, but recommended)",
				Required:    false,
				FieldType:   "input",
			},
		},
		currentField: 0,
	}
}

// NewProviderWizard creates a wizard for creating a Key Provider configuration
func NewProviderWizard(h handlers.Handler) *Wizard {
	ti := textinput.New()
	ti.Placeholder = "e.g., my-vault-provider"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return &Wizard{
		wizardType: WizardTypeProvider,
		step:       StepInput,
		handler:    h,
		ctx:        context.Background(),
		textInput:  ti,
		fields: []WizardField{
			{
				Name:        "name",
				Label:       "Provider Name",
				Description: "Enter a unique name for this provider configuration",
				Required:    true,
				FieldType:   "input",
			},
			{
				Name:        "manager",
				Label:       "Manager",
				Description: "Enter the key manager type (e.g., openbao, aws-kms)",
				Required:    true,
				FieldType:   "input",
			},
			{
				Name:        "config",
				Label:       "Configuration JSON",
				Description: "Enter the provider configuration as JSON (e.g., {\"url\":\"https://vault.example.com\"})",
				Required:    true,
				FieldType:   "input",
			},
		},
		currentField: 0,
	}
}

// buildAttributeFields builds the attribute wizard fields
func (w *Wizard) buildAttributeFields() {
	w.fields = []WizardField{}

	// If no namespace context, we need to select one
	if w.namespaceID == "" {
		w.fields = append(w.fields, WizardField{
			Name:        "namespace",
			Label:       "Namespace",
			Description: "Select the namespace for this attribute",
			Required:    true,
			FieldType:   "select",
			Options:     []SelectOption{}, // Will be populated dynamically
		})
	}

	// Attribute name
	w.fields = append(w.fields, WizardField{
		Name:        "name",
		Label:       "Attribute Name",
		Description: "Enter the attribute name",
		Required:    true,
		FieldType:   "input",
	})

	// Rule selection
	w.fields = append(w.fields, WizardField{
		Name:        "rule",
		Label:       "Attribute Rule",
		Description: "Select how attribute values are evaluated",
		Required:    true,
		FieldType:   "select",
		Options: []SelectOption{
			{Label: "All Of - Entity must have ALL values", Value: "ALL_OF"},
			{Label: "Any Of - Entity must have ANY value", Value: "ANY_OF"},
			{Label: "Hierarchy - Values form a hierarchy", Value: "HIERARCHY"},
		},
	})

	// Values (will be collected in a loop)
	w.fields = append(w.fields, WizardField{
		Name:        "values",
		Label:       "Attribute Values",
		Description: "Enter initial values for this attribute (optional)",
		Required:    false,
		FieldType:   "multiinput",
	})

	w.step = StepInput
}

// Init initializes the wizard
func (w *Wizard) Init() tea.Cmd {
	// If we need to load options for the first field, do it
	if len(w.fields) > 0 && w.fields[0].FieldType == "select" && len(w.fields[0].Options) == 0 {
		return w.loadOptionsCmd()
	}
	return textinput.Blink
}

// loadOptionsCmd loads options for the current select field
func (w *Wizard) loadOptionsCmd() tea.Cmd {
	return func() tea.Msg {
		field := &w.fields[w.currentField]

		switch field.Name {
		case "namespace":
			// Load namespaces
			resp, err := w.handler.ListNamespaces(w.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
			if err != nil {
				return wizardErrorMsg{err: fmt.Errorf("failed to load namespaces: %w", err)}
			}

			options := make([]SelectOption, 0, len(resp.GetNamespaces()))
			for _, ns := range resp.GetNamespaces() {
				options = append(options, SelectOption{
					Label: ns.GetName(),
					Value: ns.GetName(),
					ID:    ns.GetId(),
				})
			}

			if len(options) == 0 {
				return wizardErrorMsg{err: fmt.Errorf("no namespaces found - create a namespace first")}
			}

			return wizardOptionsLoadedMsg{fieldName: "namespace", options: options}
		}

		return nil
	}
}

// Message types
type wizardOptionsLoadedMsg struct {
	fieldName string
	options   []SelectOption
}

type wizardErrorMsg struct {
	err error
}

type wizardSuccessMsg struct {
	message string
}

// Update handles messages and updates the wizard
func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case wizardOptionsLoadedMsg:
		// Update the field with loaded options
		for i := range w.fields {
			if w.fields[i].Name == msg.fieldName {
				w.fields[i].Options = msg.options
				break
			}
		}
		w.step = StepSelect
		return w, nil

	case wizardErrorMsg:
		w.step = StepError
		w.error = msg.err.Error()
		return w, nil

	case wizardSuccessMsg:
		w.step = StepComplete
		w.result = msg.message
		return w, nil

	case tea.KeyMsg:
		return w.handleKeyMsg(msg)
	}

	// Update text input if in input mode
	if w.step == StepInput && w.currentFieldType() == "input" {
		var cmd tea.Cmd
		w.textInput, cmd = w.textInput.Update(msg)
		return w, cmd
	}

	return w, nil
}

// handleKeyMsg handles keyboard input
func (w *Wizard) handleKeyMsg(msg tea.KeyMsg) (*Wizard, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		w.cancelled = true
		w.step = StepComplete
		return w, nil

	case tea.KeyEnter:
		return w.handleEnter()

	case tea.KeyUp:
		if w.step == StepSelect || (w.step == StepInput && w.currentFieldType() == "select") {
			if w.selectIndex > 0 {
				w.selectIndex--
			}
		}
		return w, nil

	case tea.KeyDown:
		if w.step == StepSelect || (w.step == StepInput && w.currentFieldType() == "select") {
			field := w.currentFieldPtr()
			if field != nil && w.selectIndex < len(field.Options)-1 {
				w.selectIndex++
			}
		}
		return w, nil

	case tea.KeyTab:
		// Tab can be used to skip optional fields
		field := w.currentFieldPtr()
		if field != nil && !field.Required {
			return w.nextField()
		}
		return w, nil
	}

	// Pass to text input if in input mode
	if w.step == StepInput && w.currentFieldType() == "input" {
		var cmd tea.Cmd
		w.textInput, cmd = w.textInput.Update(msg)
		return w, cmd
	}

	// Handle character input for multiinput
	if w.step == StepInput && w.currentFieldType() == "multiinput" {
		var cmd tea.Cmd
		w.textInput, cmd = w.textInput.Update(msg)
		return w, cmd
	}

	return w, nil
}

// handleEnter processes the enter key based on current state
func (w *Wizard) handleEnter() (*Wizard, tea.Cmd) {
	field := w.currentFieldPtr()
	if field == nil {
		return w, nil
	}

	switch field.FieldType {
	case "input":
		value := strings.TrimSpace(w.textInput.Value())
		if value == "" && field.Required {
			w.error = "This field is required"
			return w, nil
		}
		w.error = ""
		field.Value = value
		return w.nextField()

	case "select":
		if len(field.Options) > 0 && w.selectIndex < len(field.Options) {
			field.Value = field.Options[w.selectIndex].Value
			field.Selected = w.selectIndex
			// Store ID if available (for namespace)
			if field.Name == "namespace" {
				w.namespaceID = field.Options[w.selectIndex].ID
				w.namespaceName = field.Options[w.selectIndex].Value
			}
		}
		return w.nextField()

	case "multiinput":
		value := strings.TrimSpace(w.textInput.Value())
		if value != "" {
			w.multiValues = append(w.multiValues, value)
			w.textInput.SetValue("")
			// Stay on this field for more values
			return w, nil
		}
		// Empty input means done adding values
		field.Value = strings.Join(w.multiValues, ", ")
		return w.nextField()
	}

	return w, nil
}

// nextField moves to the next field or executes if done
func (w *Wizard) nextField() (*Wizard, tea.Cmd) {
	w.currentField++
	w.selectIndex = 0

	if w.currentField >= len(w.fields) {
		// All fields complete, move to confirmation or execute
		return w.execute()
	}

	// Prepare for next field
	field := w.currentFieldPtr()
	if field == nil {
		return w.execute()
	}

	w.textInput.SetValue("")
	w.textInput.Placeholder = ""

	switch field.FieldType {
	case "input", "multiinput":
		w.step = StepInput
		w.textInput.Focus()
		return w, textinput.Blink

	case "select":
		if len(field.Options) == 0 {
			// Need to load options
			return w, w.loadOptionsCmd()
		}
		w.step = StepSelect
		return w, nil
	}

	return w, nil
}

// execute runs the create operation
func (w *Wizard) execute() (*Wizard, tea.Cmd) {
	w.step = StepExecuting

	return w, func() tea.Msg {
		var err error
		var successMsg string

		switch w.wizardType {
		case WizardTypeNamespace:
			successMsg, err = w.createNamespace()
		case WizardTypeAttribute:
			successMsg, err = w.createAttribute()
		case WizardTypeAttributeValue:
			successMsg, err = w.createAttributeValue()
		case WizardTypeKasRegistry:
			successMsg, err = w.createKasRegistry()
		case WizardTypeProvider:
			successMsg, err = w.createProvider()
		}

		if err != nil {
			return wizardErrorMsg{err: err}
		}
		return wizardSuccessMsg{message: successMsg}
	}
}

// createNamespace creates a new namespace
func (w *Wizard) createNamespace() (string, error) {
	name := w.getFieldValue("name")
	if name == "" {
		return "", fmt.Errorf("namespace name is required")
	}

	ns, err := w.handler.CreateNamespace(w.ctx, name, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create namespace: %w", err)
	}

	return fmt.Sprintf("Created namespace '%s' (ID: %s)", ns.GetName(), ns.GetId()), nil
}

// createAttribute creates a new attribute
func (w *Wizard) createAttribute() (string, error) {
	name := w.getFieldValue("name")
	rule := w.getFieldValue("rule")

	if name == "" {
		return "", fmt.Errorf("attribute name is required")
	}
	if rule == "" {
		return "", fmt.Errorf("attribute rule is required")
	}
	if w.namespaceID == "" {
		return "", fmt.Errorf("namespace is required")
	}

	attr, err := w.handler.CreateAttribute(w.ctx, name, rule, w.namespaceID, w.multiValues, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create attribute: %w", err)
	}

	valueCount := len(w.multiValues)
	if valueCount > 0 {
		return fmt.Sprintf("Created attribute '%s' with %d values (ID: %s)", attr.GetName(), valueCount, attr.GetId()), nil
	}
	return fmt.Sprintf("Created attribute '%s' (ID: %s)", attr.GetName(), attr.GetId()), nil
}

// createAttributeValue creates a new attribute value
func (w *Wizard) createAttributeValue() (string, error) {
	value := w.getFieldValue("value")
	if value == "" {
		return "", fmt.Errorf("value is required")
	}
	if w.attributeID == "" {
		return "", fmt.Errorf("attribute ID is required")
	}

	val, err := w.handler.CreateAttributeValue(w.ctx, w.attributeID, value, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create value: %w", err)
	}

	return fmt.Sprintf("Created value '%s' (ID: %s)", val.GetValue(), val.GetId()), nil
}

// createKasRegistry creates a new KAS Registry entry
func (w *Wizard) createKasRegistry() (string, error) {
	uri := w.getFieldValue("uri")
	if uri == "" {
		return "", fmt.Errorf("KAS URI is required")
	}

	name := w.getFieldValue("name")

	kas, err := w.handler.CreateKasRegistryEntry(w.ctx, uri, name, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create KAS Registry entry: %w", err)
	}

	if name != "" {
		return fmt.Sprintf("Created KAS '%s' at %s (ID: %s)", kas.GetName(), kas.GetUri(), kas.GetId()), nil
	}
	return fmt.Sprintf("Created KAS at %s (ID: %s)", kas.GetUri(), kas.GetId()), nil
}

// createProvider creates a new Key Provider configuration
func (w *Wizard) createProvider() (string, error) {
	name := w.getFieldValue("name")
	if name == "" {
		return "", fmt.Errorf("provider name is required")
	}

	manager := w.getFieldValue("manager")
	if manager == "" {
		return "", fmt.Errorf("manager is required")
	}

	config := w.getFieldValue("config")
	if config == "" {
		return "", fmt.Errorf("configuration is required")
	}

	provider, err := w.handler.CreateProviderConfig(w.ctx, name, manager, []byte(config), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create provider config: %w", err)
	}

	return fmt.Sprintf("Created provider '%s' (ID: %s)", provider.GetName(), provider.GetId()), nil
}

// getFieldValue returns the value for a field by name
func (w *Wizard) getFieldValue(name string) string {
	for _, f := range w.fields {
		if f.Name == name {
			return f.Value
		}
	}
	return ""
}

// currentFieldType returns the type of the current field
func (w *Wizard) currentFieldType() string {
	if w.currentField >= len(w.fields) {
		return ""
	}
	return w.fields[w.currentField].FieldType
}

// currentFieldPtr returns a pointer to the current field
func (w *Wizard) currentFieldPtr() *WizardField {
	if w.currentField >= len(w.fields) {
		return nil
	}
	return &w.fields[w.currentField]
}

// View renders the wizard
func (w *Wizard) View() string {
	var sb strings.Builder

	// Title based on wizard type
	switch w.wizardType {
	case WizardTypeNamespace:
		sb.WriteString(wizardTitleStyle.Render("Create Namespace") + "\n\n")
	case WizardTypeAttribute:
		sb.WriteString(wizardTitleStyle.Render("Create Attribute") + "\n")
		if w.namespaceName != "" {
			sb.WriteString(wizardDescStyle.Render(fmt.Sprintf("In namespace: %s", w.namespaceName)) + "\n")
		}
		sb.WriteString("\n")
	case WizardTypeAttributeValue:
		sb.WriteString(wizardTitleStyle.Render("Create Attribute Value") + "\n")
		if w.attributeName != "" {
			sb.WriteString(wizardDescStyle.Render(fmt.Sprintf("For attribute: %s", w.attributeName)) + "\n")
		}
		sb.WriteString("\n")
	case WizardTypeKasRegistry:
		sb.WriteString(wizardTitleStyle.Render("Create Key Access Server") + "\n\n")
	case WizardTypeProvider:
		sb.WriteString(wizardTitleStyle.Render("Create Provider Configuration") + "\n\n")
	}

	// Show completed fields
	for i := 0; i < w.currentField && i < len(w.fields); i++ {
		field := w.fields[i]
		sb.WriteString(wizardLabelStyle.Render(field.Label+": "))
		sb.WriteString(wizardSuccessStyle.Render(field.Value) + "\n")
	}

	// Current step content
	switch w.step {
	case StepInput:
		sb.WriteString(w.renderInputStep())

	case StepSelect:
		sb.WriteString(w.renderSelectStep())

	case StepExecuting:
		sb.WriteString(wizardDescStyle.Render("Creating resource...") + "\n")

	case StepComplete:
		if w.cancelled {
			sb.WriteString(wizardErrorStyle.Render("Cancelled") + "\n")
		} else {
			sb.WriteString(wizardSuccessStyle.Render("✓ "+w.result) + "\n")
		}

	case StepError:
		sb.WriteString(wizardErrorStyle.Render("Error: "+w.error) + "\n")
		sb.WriteString(wizardHintStyle.Render("Press Esc to cancel") + "\n")
	}

	// Help text
	if w.step != StepComplete && w.step != StepError && w.step != StepExecuting {
		sb.WriteString("\n")
		sb.WriteString(wizardHintStyle.Render("Press Enter to continue, Esc to cancel"))
	}

	return sb.String()
}

// renderInputStep renders an input field
func (w *Wizard) renderInputStep() string {
	field := w.currentFieldPtr()
	if field == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(wizardLabelStyle.Render(field.Label) + "\n")
	if field.Description != "" {
		sb.WriteString(wizardDescStyle.Render(field.Description) + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString(w.textInput.View() + "\n")

	if w.error != "" {
		sb.WriteString(wizardErrorStyle.Render(w.error) + "\n")
	}

	// For multiinput, show already entered values
	if field.FieldType == "multiinput" && len(w.multiValues) > 0 {
		sb.WriteString("\n" + wizardDescStyle.Render("Values entered:") + "\n")
		for _, v := range w.multiValues {
			sb.WriteString("  • " + wizardSuccessStyle.Render(v) + "\n")
		}
		sb.WriteString(wizardHintStyle.Render("Press Enter with empty input when done") + "\n")
	}

	return sb.String()
}

// renderSelectStep renders a selection list
func (w *Wizard) renderSelectStep() string {
	field := w.currentFieldPtr()
	if field == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(wizardLabelStyle.Render(field.Label) + "\n")
	if field.Description != "" {
		sb.WriteString(wizardDescStyle.Render(field.Description) + "\n")
	}
	sb.WriteString("\n")

	for i, opt := range field.Options {
		if i == w.selectIndex {
			sb.WriteString(wizardSelectedStyle.Render("> "+opt.Label) + "\n")
		} else {
			sb.WriteString(wizardUnselectedStyle.Render("  "+opt.Label) + "\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(wizardHintStyle.Render("Use ↑/↓ to navigate, Enter to select"))

	return sb.String()
}

// IsComplete returns true if the wizard has finished
func (w *Wizard) IsComplete() bool {
	return w.step == StepComplete || w.step == StepError
}

// WasCancelled returns true if the wizard was cancelled
func (w *Wizard) WasCancelled() bool {
	return w.cancelled
}

// GetResult returns the result message
func (w *Wizard) GetResult() string {
	if w.step == StepError {
		return w.error
	}
	return w.result
}

// GetError returns the error message if any
func (w *Wizard) GetError() string {
	return w.error
}

// Ensure Wizard can be used to fetch options dynamically
func (w *Wizard) LoadNamespaceOptions() ([]SelectOption, error) {
	resp, err := w.handler.ListNamespaces(w.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
	if err != nil {
		return nil, err
	}

	options := make([]SelectOption, 0, len(resp.GetNamespaces()))
	for _, ns := range resp.GetNamespaces() {
		options = append(options, SelectOption{
			Label: ns.GetName(),
			Value: ns.GetName(),
			ID:    ns.GetId(),
		})
	}
	return options, nil
}

// LoadAttributeOptions loads attribute options for a namespace
func (w *Wizard) LoadAttributeOptions(namespaceID string) ([]SelectOption, error) {
	resp, err := w.handler.ListAttributes(w.ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, 1000, 0)
	if err != nil {
		return nil, err
	}

	options := make([]SelectOption, 0)
	for _, attr := range resp.GetAttributes() {
		if attr.GetNamespace().GetId() == namespaceID || attr.GetNamespace().GetName() == namespaceID {
			options = append(options, SelectOption{
				Label: attr.GetName(),
				Value: attr.GetName(),
				ID:    attr.GetId(),
			})
		}
	}
	return options, nil
}

// Compile-time check that Wizard satisfies expected interface patterns
var _ tea.Model = (*Wizard)(nil)

// DeleteWizard handles resource deletion with safety confirmations
type DeleteWizard struct {
	resourceType string // "namespace", "attribute", "value"
	resourceName string
	resourceID   string
	resourceFQN  string
	step         DeleteStep
	handler      handlers.Handler
	ctx          context.Context
	selectIndex  int
	error        string
	result       string
	cancelled    bool

	// For showing what will be deleted
	childCount int // Number of child resources that will be affected
}

type DeleteStep int

const (
	DeleteStepConfirmType DeleteStep = iota // Choose deactivate vs delete
	DeleteStepConfirmAction                 // Final confirmation
	DeleteStepExecuting
	DeleteStepComplete
	DeleteStepError
)

// DeleteAction represents the type of deletion
type DeleteAction int

const (
	DeleteActionDeactivate DeleteAction = iota // Safe - just marks inactive
	DeleteActionDelete                         // Unsafe - permanent removal
)

// NewDeleteWizard creates a new delete wizard
func NewDeleteWizard(h handlers.Handler, resourceType, resourceName, resourceID, resourceFQN string, childCount int) *DeleteWizard {
	return &DeleteWizard{
		resourceType: resourceType,
		resourceName: resourceName,
		resourceID:   resourceID,
		resourceFQN:  resourceFQN,
		childCount:   childCount,
		step:         DeleteStepConfirmType,
		handler:      h,
		ctx:          context.Background(),
		selectIndex:  0, // Default to deactivate (safer option)
	}
}

// Init initializes the delete wizard
func (d *DeleteWizard) Init() tea.Cmd {
	return nil
}

// Update handles messages for the delete wizard
func (d *DeleteWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case deleteResultMsg:
		if msg.err != nil {
			d.step = DeleteStepError
			d.error = msg.err.Error()
		} else {
			d.step = DeleteStepComplete
			d.result = msg.message
		}
		return d, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			d.cancelled = true
			d.step = DeleteStepComplete
			return d, nil

		case tea.KeyUp:
			if d.step == DeleteStepConfirmType && d.selectIndex > 0 {
				d.selectIndex--
			}
			return d, nil

		case tea.KeyDown:
			if d.step == DeleteStepConfirmType && d.selectIndex < 1 {
				d.selectIndex++
			}
			return d, nil

		case tea.KeyEnter:
			return d.handleEnter()

		case tea.KeyRunes:
			// Handle 'y' or 'n' for confirmation
			if d.step == DeleteStepConfirmAction {
				switch string(msg.Runes) {
				case "y", "Y":
					return d.executeDelete()
				case "n", "N":
					d.cancelled = true
					d.step = DeleteStepComplete
					return d, nil
				}
			}
		}
	}

	return d, nil
}

func (d *DeleteWizard) handleEnter() (*DeleteWizard, tea.Cmd) {
	switch d.step {
	case DeleteStepConfirmType:
		d.step = DeleteStepConfirmAction
		return d, nil

	case DeleteStepConfirmAction:
		return d.executeDelete()
	}

	return d, nil
}

type deleteResultMsg struct {
	message string
	err     error
}

func (d *DeleteWizard) executeDelete() (*DeleteWizard, tea.Cmd) {
	d.step = DeleteStepExecuting

	return d, func() tea.Msg {
		var err error
		var message string
		isDeactivate := d.selectIndex == 0

		switch d.resourceType {
		case "namespace":
			if isDeactivate {
				_, err = d.handler.DeactivateNamespace(d.ctx, d.resourceID)
				if err == nil {
					message = fmt.Sprintf("Deactivated namespace '%s'", d.resourceName)
				}
			} else {
				err = d.handler.UnsafeDeleteNamespace(d.ctx, d.resourceID, d.resourceFQN)
				if err == nil {
					message = fmt.Sprintf("Permanently deleted namespace '%s'", d.resourceName)
				}
			}

		case "attribute":
			if isDeactivate {
				_, err = d.handler.DeactivateAttribute(d.ctx, d.resourceID)
				if err == nil {
					message = fmt.Sprintf("Deactivated attribute '%s'", d.resourceName)
				}
			} else {
				err = d.handler.UnsafeDeleteAttribute(d.ctx, d.resourceID, d.resourceFQN)
				if err == nil {
					message = fmt.Sprintf("Permanently deleted attribute '%s'", d.resourceName)
				}
			}

		case "value":
			if isDeactivate {
				_, err = d.handler.DeactivateAttributeValue(d.ctx, d.resourceID)
				if err == nil {
					message = fmt.Sprintf("Deactivated value '%s'", d.resourceName)
				}
			} else {
				err = d.handler.UnsafeDeleteAttributeValue(d.ctx, d.resourceID, d.resourceFQN)
				if err == nil {
					message = fmt.Sprintf("Permanently deleted value '%s'", d.resourceName)
				}
			}

		default:
			err = fmt.Errorf("unknown resource type: %s", d.resourceType)
		}

		return deleteResultMsg{message: message, err: err}
	}
}

// View renders the delete wizard
func (d *DeleteWizard) View() string {
	var sb strings.Builder

	// Warning header
	sb.WriteString(wizardErrorStyle.Render("⚠ DELETE RESOURCE") + "\n\n")

	// Show what will be deleted
	sb.WriteString(wizardLabelStyle.Render("Resource: "))
	sb.WriteString(wizardSelectedStyle.Render(d.resourceName) + "\n")
	sb.WriteString(wizardLabelStyle.Render("Type: "))
	sb.WriteString(outputStyle.Render(d.resourceType) + "\n")
	sb.WriteString(wizardLabelStyle.Render("ID: "))
	sb.WriteString(wizardHintStyle.Render(d.resourceID) + "\n")

	if d.childCount > 0 {
		sb.WriteString("\n")
		sb.WriteString(wizardErrorStyle.Render(fmt.Sprintf("Warning: This %s has %d child resource(s) that will be affected!", d.resourceType, d.childCount)))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	switch d.step {
	case DeleteStepConfirmType:
		sb.WriteString(wizardLabelStyle.Render("Choose action:") + "\n\n")

		// Deactivate option (safer)
		if d.selectIndex == 0 {
			sb.WriteString(wizardSelectedStyle.Render("> Deactivate (Recommended)") + "\n")
			sb.WriteString(wizardHintStyle.Render("    Marks as inactive, can be reactivated later") + "\n")
		} else {
			sb.WriteString(wizardUnselectedStyle.Render("  Deactivate (Recommended)") + "\n")
			sb.WriteString(wizardHintStyle.Render("    Marks as inactive, can be reactivated later") + "\n")
		}

		sb.WriteString("\n")

		// Delete option (dangerous)
		if d.selectIndex == 1 {
			sb.WriteString(wizardErrorStyle.Render("> Permanently Delete") + "\n")
			sb.WriteString(wizardErrorStyle.Render("    CANNOT BE UNDONE - removes from database") + "\n")
		} else {
			sb.WriteString(wizardUnselectedStyle.Render("  Permanently Delete") + "\n")
			sb.WriteString(wizardHintStyle.Render("    Cannot be undone - removes from database") + "\n")
		}

		sb.WriteString("\n")
		sb.WriteString(wizardHintStyle.Render("Use ↑/↓ to select, Enter to continue, Esc to cancel"))

	case DeleteStepConfirmAction:
		action := "deactivate"
		if d.selectIndex == 1 {
			action = "PERMANENTLY DELETE"
			sb.WriteString(wizardErrorStyle.Render("⚠ THIS ACTION CANNOT BE UNDONE ⚠") + "\n\n")
		}

		sb.WriteString(fmt.Sprintf("Are you sure you want to %s '%s'?\n\n", action, d.resourceName))
		sb.WriteString(wizardLabelStyle.Render("Type 'y' to confirm, 'n' to cancel: "))

	case DeleteStepExecuting:
		sb.WriteString(wizardDescStyle.Render("Processing...") + "\n")

	case DeleteStepComplete:
		if d.cancelled {
			sb.WriteString(wizardHintStyle.Render("Cancelled") + "\n")
		} else {
			sb.WriteString(wizardSuccessStyle.Render("✓ "+d.result) + "\n")
		}

	case DeleteStepError:
		sb.WriteString(wizardErrorStyle.Render("Error: "+d.error) + "\n")
	}

	return sb.String()
}

// IsComplete returns true if the wizard has finished
func (d *DeleteWizard) IsComplete() bool {
	return d.step == DeleteStepComplete || d.step == DeleteStepError
}

// WasCancelled returns true if the wizard was cancelled
func (d *DeleteWizard) WasCancelled() bool {
	return d.cancelled
}

// GetResult returns the result message
func (d *DeleteWizard) GetResult() string {
	if d.step == DeleteStepError {
		return d.error
	}
	return d.result
}

// GetError returns the error message if any
func (d *DeleteWizard) GetError() string {
	return d.error
}

// ============================================================================
// Key Assignment Wizard
// ============================================================================

// KeyAssignWizard is the wizard for assigning keys to resources
type KeyAssignWizard struct {
	resourceType string // "namespace", "attribute", "value"
	resourceName string
	resourceID   string
	step         KeyAssignStep
	handler      handlers.Handler
	ctx          context.Context
	selectIndex  int
	keys         []KeyOption
	error        string
	result       string
	cancelled    bool
}

type KeyAssignStep int

const (
	KeyAssignStepLoading KeyAssignStep = iota
	KeyAssignStepSelect
	KeyAssignStepConfirm
	KeyAssignStepExecuting
	KeyAssignStepComplete
	KeyAssignStepError
)

// KeyOption represents a key available for assignment
type KeyOption struct {
	ID        string // System ID of the key
	KeyID     string // User-facing key identifier
	KasURI    string
	Algorithm string
	Status    string
}

// NewKeyAssignWizard creates a new key assignment wizard
func NewKeyAssignWizard(h handlers.Handler, resourceType, resourceName, resourceID string) *KeyAssignWizard {
	return &KeyAssignWizard{
		resourceType: resourceType,
		resourceName: resourceName,
		resourceID:   resourceID,
		step:         KeyAssignStepLoading,
		handler:      h,
		ctx:          context.Background(),
		selectIndex:  0,
		keys:         []KeyOption{},
	}
}

// keyAssignLoadedMsg is sent when keys are loaded
type keyAssignLoadedMsg struct {
	keys []KeyOption
}

// keyAssignErrorMsg is sent on error
type keyAssignErrorMsg struct {
	err error
}

// keyAssignResultMsg is sent when assignment completes
type keyAssignResultMsg struct {
	message string
	err     error
}

// Init initializes the key assignment wizard
func (k *KeyAssignWizard) Init() tea.Cmd {
	return k.loadKeysCmd()
}

// loadKeysCmd loads available keys from the API
func (k *KeyAssignWizard) loadKeysCmd() tea.Cmd {
	return func() tea.Msg {
		resp, err := k.handler.ListKasKeys(k.ctx, 100, 0, 0, handlers.KasIdentifier{}, nil)
		if err != nil {
			return keyAssignErrorMsg{err: fmt.Errorf("failed to load keys: %w", err)}
		}

		kasKeys := resp.GetKasKeys()
		if len(kasKeys) == 0 {
			return keyAssignErrorMsg{err: fmt.Errorf("no keys found in the system - create a key first")}
		}

		keys := make([]KeyOption, 0, len(kasKeys))
		for _, kasKey := range kasKeys {
			key := kasKey.GetKey()
			// Only show active keys
			if key.GetKeyStatus().String() == "KEY_STATUS_ACTIVE" {
				keys = append(keys, KeyOption{
					ID:        key.GetId(),
					KeyID:     key.GetKeyId(),
					KasURI:    kasKey.GetKasUri(),
					Algorithm: key.GetKeyAlgorithm().String(),
					Status:    "active",
				})
			}
		}

		if len(keys) == 0 {
			return keyAssignErrorMsg{err: fmt.Errorf("no active keys found in the system")}
		}

		return keyAssignLoadedMsg{keys: keys}
	}
}

// Update handles messages for the key assignment wizard
func (k *KeyAssignWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case keyAssignLoadedMsg:
		k.keys = msg.keys
		k.step = KeyAssignStepSelect
		return k, nil

	case keyAssignErrorMsg:
		k.step = KeyAssignStepError
		k.error = msg.err.Error()
		return k, nil

	case keyAssignResultMsg:
		if msg.err != nil {
			k.step = KeyAssignStepError
			k.error = msg.err.Error()
		} else {
			k.step = KeyAssignStepComplete
			k.result = msg.message
		}
		return k, nil

	case tea.KeyMsg:
		return k.handleKeyMsg(msg)
	}

	return k, nil
}

// handleKeyMsg handles key messages
func (k *KeyAssignWizard) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		k.cancelled = true
		k.step = KeyAssignStepComplete
		return k, nil

	case tea.KeyUp:
		if k.step == KeyAssignStepSelect && k.selectIndex > 0 {
			k.selectIndex--
		}
		return k, nil

	case tea.KeyDown:
		if k.step == KeyAssignStepSelect && k.selectIndex < len(k.keys)-1 {
			k.selectIndex++
		}
		return k, nil

	case tea.KeyEnter:
		return k.handleEnter()
	}

	return k, nil
}

// handleEnter handles enter key press
func (k *KeyAssignWizard) handleEnter() (tea.Model, tea.Cmd) {
	switch k.step {
	case KeyAssignStepSelect:
		k.step = KeyAssignStepConfirm
		return k, nil

	case KeyAssignStepConfirm:
		k.step = KeyAssignStepExecuting
		return k, k.executeAssignment()

	case KeyAssignStepComplete, KeyAssignStepError:
		return k, nil
	}

	return k, nil
}

// executeAssignment performs the key assignment
func (k *KeyAssignWizard) executeAssignment() tea.Cmd {
	return func() tea.Msg {
		if k.selectIndex >= len(k.keys) {
			return keyAssignResultMsg{err: fmt.Errorf("invalid key selection")}
		}

		selectedKey := k.keys[k.selectIndex]

		var err error
		switch k.resourceType {
		case "namespace":
			_, err = k.handler.AssignKeyToAttributeNamespace(k.ctx, k.resourceID, selectedKey.ID)
		case "attribute":
			_, err = k.handler.AssignKeyToAttribute(k.ctx, k.resourceID, selectedKey.ID)
		case "value":
			_, err = k.handler.AssignKeyToAttributeValue(k.ctx, k.resourceID, selectedKey.ID)
		default:
			return keyAssignResultMsg{err: fmt.Errorf("unknown resource type: %s", k.resourceType)}
		}

		if err != nil {
			return keyAssignResultMsg{err: fmt.Errorf("failed to assign key: %w", err)}
		}

		return keyAssignResultMsg{
			message: fmt.Sprintf("Assigned key '%s' to %s '%s'", selectedKey.KeyID, k.resourceType, k.resourceName),
		}
	}
}

// View renders the key assignment wizard
func (k *KeyAssignWizard) View() string {
	var sb strings.Builder

	sb.WriteString(wizardTitleStyle.Render("Assign Key") + "\n")
	sb.WriteString(wizardDescStyle.Render(fmt.Sprintf("To %s: %s", k.resourceType, k.resourceName)) + "\n\n")

	switch k.step {
	case KeyAssignStepLoading:
		sb.WriteString(wizardDescStyle.Render("Loading available keys...") + "\n")

	case KeyAssignStepSelect:
		sb.WriteString(wizardLabelStyle.Render("Select a key to assign:") + "\n\n")
		for i, key := range k.keys {
			cursor := "  "
			style := wizardUnselectedStyle
			if i == k.selectIndex {
				cursor = "▸ "
				style = wizardSelectedStyle
			}
			sb.WriteString(cursor)
			sb.WriteString(style.Render(fmt.Sprintf("%s (%s) - %s", key.KeyID, key.Algorithm, key.KasURI)))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
		sb.WriteString(wizardHintStyle.Render("↑/↓ to navigate, Enter to select, Esc to cancel"))

	case KeyAssignStepConfirm:
		selectedKey := k.keys[k.selectIndex]
		sb.WriteString(wizardLabelStyle.Render("Confirm Assignment") + "\n\n")
		sb.WriteString(fmt.Sprintf("Key: %s\n", wizardSelectedStyle.Render(selectedKey.KeyID)))
		sb.WriteString(fmt.Sprintf("KAS: %s\n", wizardDescStyle.Render(selectedKey.KasURI)))
		sb.WriteString(fmt.Sprintf("Algorithm: %s\n", wizardDescStyle.Render(selectedKey.Algorithm)))
		sb.WriteString("\n")
		sb.WriteString(wizardDescStyle.Render(fmt.Sprintf("Assign to %s '%s'?", k.resourceType, k.resourceName)) + "\n\n")
		sb.WriteString(wizardHintStyle.Render("Press Enter to confirm, Esc to cancel"))

	case KeyAssignStepExecuting:
		sb.WriteString(wizardDescStyle.Render("Assigning key...") + "\n")

	case KeyAssignStepComplete:
		if k.cancelled {
			sb.WriteString(wizardErrorStyle.Render("Cancelled") + "\n")
		} else {
			sb.WriteString(wizardSuccessStyle.Render("✓ "+k.result) + "\n")
		}

	case KeyAssignStepError:
		sb.WriteString(wizardErrorStyle.Render("Error: "+k.error) + "\n")
		sb.WriteString(wizardHintStyle.Render("Press Esc to go back") + "\n")
	}

	return sb.String()
}

// IsComplete returns true if the wizard has finished
func (k *KeyAssignWizard) IsComplete() bool {
	return k.step == KeyAssignStepComplete || k.step == KeyAssignStepError
}

// WasCancelled returns true if the wizard was cancelled
func (k *KeyAssignWizard) WasCancelled() bool {
	return k.cancelled
}

// GetResult returns the result message
func (k *KeyAssignWizard) GetResult() string {
	if k.step == KeyAssignStepError {
		return k.error
	}
	return k.result
}

// GetError returns the error message if any
func (k *KeyAssignWizard) GetError() string {
	return k.error
}

// ============================================================================
// Simple Delete Wizard (for resources without deactivate option)
// ============================================================================

// SimpleDeleteWizard handles resource deletion with a simple confirmation
// Used for resources like KAS, providers that don't have deactivate flow
type SimpleDeleteWizard struct {
	resourceType string // "kas", "provider", etc.
	resourceName string
	resourceID   string
	step         SimpleDeleteStep
	handler      handlers.Handler
	ctx          context.Context
	error        string
	result       string
	cancelled    bool
}

type SimpleDeleteStep int

const (
	SimpleDeleteStepConfirm SimpleDeleteStep = iota
	SimpleDeleteStepExecuting
	SimpleDeleteStepComplete
	SimpleDeleteStepError
)

// NewSimpleDeleteWizard creates a new simple delete wizard
func NewSimpleDeleteWizard(h handlers.Handler, resourceType, resourceName, resourceID string) *SimpleDeleteWizard {
	return &SimpleDeleteWizard{
		resourceType: resourceType,
		resourceName: resourceName,
		resourceID:   resourceID,
		step:         SimpleDeleteStepConfirm,
		handler:      h,
		ctx:          context.Background(),
	}
}

// simpleDeleteResultMsg is sent when deletion completes
type simpleDeleteResultMsg struct {
	message string
	err     error
}

// Init initializes the simple delete wizard
func (d *SimpleDeleteWizard) Init() tea.Cmd {
	return nil
}

// Update handles messages for the simple delete wizard
func (d *SimpleDeleteWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case simpleDeleteResultMsg:
		if msg.err != nil {
			d.step = SimpleDeleteStepError
			d.error = msg.err.Error()
		} else {
			d.step = SimpleDeleteStepComplete
			d.result = msg.message
		}
		return d, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			d.cancelled = true
			d.step = SimpleDeleteStepComplete
			return d, nil

		case tea.KeyRunes:
			if d.step == SimpleDeleteStepConfirm {
				switch string(msg.Runes) {
				case "y", "Y":
					return d.executeDelete()
				case "n", "N":
					d.cancelled = true
					d.step = SimpleDeleteStepComplete
					return d, nil
				}
			}

		case tea.KeyEnter:
			// Enter alone doesn't confirm - need explicit y/n
			return d, nil
		}
	}

	return d, nil
}

func (d *SimpleDeleteWizard) executeDelete() (*SimpleDeleteWizard, tea.Cmd) {
	d.step = SimpleDeleteStepExecuting

	return d, func() tea.Msg {
		var err error
		var message string

		switch d.resourceType {
		case "kas":
			_, err = d.handler.DeleteKasRegistryEntry(d.ctx, d.resourceID)
			if err == nil {
				message = fmt.Sprintf("Deleted KAS '%s'", d.resourceName)
			}

		case "provider":
			err = d.handler.DeleteProviderConfig(d.ctx, d.resourceID)
			if err == nil {
				message = fmt.Sprintf("Deleted provider '%s'", d.resourceName)
			}

		default:
			err = fmt.Errorf("unknown resource type: %s", d.resourceType)
		}

		return simpleDeleteResultMsg{message: message, err: err}
	}
}

// View renders the simple delete wizard
func (d *SimpleDeleteWizard) View() string {
	var sb strings.Builder

	// Warning header
	sb.WriteString(wizardErrorStyle.Render("⚠ DELETE RESOURCE") + "\n\n")

	// Show what will be deleted
	sb.WriteString(wizardLabelStyle.Render("Resource: "))
	sb.WriteString(wizardSelectedStyle.Render(d.resourceName) + "\n")
	sb.WriteString(wizardLabelStyle.Render("Type: "))
	sb.WriteString(outputStyle.Render(d.resourceType) + "\n")
	sb.WriteString(wizardLabelStyle.Render("ID: "))
	sb.WriteString(wizardHintStyle.Render(d.resourceID) + "\n\n")

	switch d.step {
	case SimpleDeleteStepConfirm:
		sb.WriteString(wizardErrorStyle.Render("⚠ THIS ACTION CANNOT BE UNDONE ⚠") + "\n\n")
		sb.WriteString(fmt.Sprintf("Are you sure you want to delete '%s'?\n\n", d.resourceName))
		sb.WriteString(wizardLabelStyle.Render("Type 'y' to confirm, 'n' to cancel: "))

	case SimpleDeleteStepExecuting:
		sb.WriteString(wizardDescStyle.Render("Deleting...") + "\n")

	case SimpleDeleteStepComplete:
		if d.cancelled {
			sb.WriteString(wizardHintStyle.Render("Cancelled") + "\n")
		} else {
			sb.WriteString(wizardSuccessStyle.Render("✓ "+d.result) + "\n")
		}

	case SimpleDeleteStepError:
		sb.WriteString(wizardErrorStyle.Render("Error: "+d.error) + "\n")
	}

	return sb.String()
}

// IsComplete returns true if the wizard has finished
func (d *SimpleDeleteWizard) IsComplete() bool {
	return d.step == SimpleDeleteStepComplete || d.step == SimpleDeleteStepError
}

// WasCancelled returns true if the wizard was cancelled
func (d *SimpleDeleteWizard) WasCancelled() bool {
	return d.cancelled
}

// GetResult returns the result message
func (d *SimpleDeleteWizard) GetResult() string {
	if d.step == SimpleDeleteStepError {
		return d.error
	}
	return d.result
}

// GetError returns the error message if any
func (d *SimpleDeleteWizard) GetError() string {
	return d.error
}

// ============================================================================
// KAS Key Wizard (Multi-step wizard for creating keys)
// ============================================================================

// KasKeyMode represents the key creation mode
type KasKeyMode int

const (
	KasKeyModeLocal KasKeyMode = iota
	KasKeyModeProvider
	KasKeyModeRemote
	KasKeyModePublicKeyOnly
)

func (m KasKeyMode) String() string {
	switch m {
	case KasKeyModeLocal:
		return "local"
	case KasKeyModeProvider:
		return "provider"
	case KasKeyModeRemote:
		return "remote"
	case KasKeyModePublicKeyOnly:
		return "public_key"
	default:
		return "unknown"
	}
}

// KasKeyWizardStep represents the current step in the wizard
type KasKeyWizardStep int

const (
	KasKeyStepMode KasKeyWizardStep = iota
	KasKeyStepConfig
	KasKeyStepModeSpecific
	KasKeyStepConfirm
	KasKeyStepExecuting
	KasKeyStepComplete
	KasKeyStepError
)

// KasKeyWizard handles multi-step key creation
type KasKeyWizard struct {
	kasID       string
	kasName     string
	step        KasKeyWizardStep
	handler     handlers.Handler
	ctx         context.Context
	textInput   textinput.Model
	selectIndex int
	error       string
	result      string
	cancelled   bool

	// Key configuration
	mode      KasKeyMode
	keyID     string
	algorithm string // e.g., "rsa:2048", "ec:secp256r1"

	// Mode-specific fields
	wrappingKeyID  string
	wrappingKeyHex string
	publicKeyPEM   string
	privateKeyPEM  string
	providerID     string

	// Dynamic data
	algorithms []SelectOption
	providers  []SelectOption

	// Current field being edited
	currentField int
	fields       []string // field names for current step
}

// NewKasKeyWizard creates a new KAS key creation wizard
func NewKasKeyWizard(h handlers.Handler, kasID, kasName string) *KasKeyWizard {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 1024
	ti.Width = 60

	return &KasKeyWizard{
		kasID:       kasID,
		kasName:     kasName,
		step:        KasKeyStepMode,
		handler:     h,
		ctx:         context.Background(),
		textInput:   ti,
		selectIndex: 0,
		algorithms: []SelectOption{
			{Label: "RSA 2048", Value: "rsa:2048"},
			{Label: "RSA 4096", Value: "rsa:4096"},
			{Label: "EC secp256r1", Value: "ec:secp256r1"},
			{Label: "EC secp384r1", Value: "ec:secp384r1"},
			{Label: "EC secp521r1", Value: "ec:secp521r1"},
		},
	}
}

// kasKeyResultMsg is sent when key creation completes
type kasKeyResultMsg struct {
	message string
	err     error
}

// kasKeyProvidersLoadedMsg is sent when providers are loaded
type kasKeyProvidersLoadedMsg struct {
	providers []SelectOption
}

// Init initializes the KAS key wizard
func (k *KasKeyWizard) Init() tea.Cmd {
	return nil
}

// Update handles messages for the KAS key wizard
func (k *KasKeyWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case kasKeyResultMsg:
		if msg.err != nil {
			k.step = KasKeyStepError
			k.error = msg.err.Error()
		} else {
			k.step = KasKeyStepComplete
			k.result = msg.message
		}
		return k, nil

	case kasKeyProvidersLoadedMsg:
		k.providers = msg.providers
		return k, nil

	case tea.KeyMsg:
		return k.handleKeyMsg(msg)
	}

	// Update text input if in text input mode
	if k.isTextInputStep() {
		var cmd tea.Cmd
		k.textInput, cmd = k.textInput.Update(msg)
		return k, cmd
	}

	return k, nil
}

func (k *KasKeyWizard) isTextInputStep() bool {
	switch k.step {
	case KasKeyStepConfig:
		// Key ID is text input
		return k.currentField == 0
	case KasKeyStepModeSpecific:
		// Most mode-specific fields are text input
		return true
	}
	return false
}

func (k *KasKeyWizard) handleKeyMsg(msg tea.KeyMsg) (*KasKeyWizard, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		k.cancelled = true
		k.step = KasKeyStepComplete
		return k, nil

	case tea.KeyUp:
		if k.selectIndex > 0 {
			k.selectIndex--
		}
		return k, nil

	case tea.KeyDown:
		maxIndex := k.getMaxSelectIndex()
		if k.selectIndex < maxIndex {
			k.selectIndex++
		}
		return k, nil

	case tea.KeyEnter:
		return k.handleEnter()
	}

	// Pass to text input if appropriate
	if k.isTextInputStep() {
		var cmd tea.Cmd
		k.textInput, cmd = k.textInput.Update(msg)
		return k, cmd
	}

	return k, nil
}

func (k *KasKeyWizard) getMaxSelectIndex() int {
	switch k.step {
	case KasKeyStepMode:
		return 3 // 4 modes: 0-3
	case KasKeyStepConfig:
		if k.currentField == 1 { // algorithm selection
			return len(k.algorithms) - 1
		}
	}
	return 0
}

func (k *KasKeyWizard) handleEnter() (*KasKeyWizard, tea.Cmd) {
	switch k.step {
	case KasKeyStepMode:
		k.mode = KasKeyMode(k.selectIndex)
		k.step = KasKeyStepConfig
		k.currentField = 0
		k.selectIndex = 0
		k.textInput.SetValue("")
		k.textInput.Placeholder = "e.g., my-encryption-key"
		k.textInput.Focus()
		return k, textinput.Blink

	case KasKeyStepConfig:
		if k.currentField == 0 {
			// Key ID entered, move to algorithm
			k.keyID = strings.TrimSpace(k.textInput.Value())
			if k.keyID == "" {
				k.error = "Key ID is required"
				return k, nil
			}
			k.error = ""
			k.currentField = 1
			k.selectIndex = 0
			return k, nil
		} else {
			// Algorithm selected, move to mode-specific step
			k.algorithm = k.algorithms[k.selectIndex].Value
			k.step = KasKeyStepModeSpecific
			k.currentField = 0
			k.textInput.SetValue("")
			k.setupModeSpecificFields()
			return k, textinput.Blink
		}

	case KasKeyStepModeSpecific:
		return k.handleModeSpecificEnter()

	case KasKeyStepConfirm:
		return k.executeCreate()

	case KasKeyStepComplete, KasKeyStepError:
		return k, nil
	}

	return k, nil
}

func (k *KasKeyWizard) setupModeSpecificFields() {
	switch k.mode {
	case KasKeyModeLocal:
		k.fields = []string{"wrappingKeyID", "wrappingKeyHex"}
		k.textInput.Placeholder = "e.g., local-wrapping-key"
	case KasKeyModeProvider:
		k.fields = []string{"providerID", "publicKeyPEM", "privateKeyPEM", "wrappingKeyID"}
		k.textInput.Placeholder = "Provider config name or ID"
	case KasKeyModeRemote:
		k.fields = []string{"providerID", "publicKeyPEM", "wrappingKeyID"}
		k.textInput.Placeholder = "Provider config name or ID"
	case KasKeyModePublicKeyOnly:
		k.fields = []string{"publicKeyPEM"}
		k.textInput.Placeholder = "Base64-encoded PEM public key"
	}
}

func (k *KasKeyWizard) handleModeSpecificEnter() (*KasKeyWizard, tea.Cmd) {
	value := strings.TrimSpace(k.textInput.Value())

	// Store the value based on current field
	if k.currentField < len(k.fields) {
		fieldName := k.fields[k.currentField]
		switch fieldName {
		case "wrappingKeyID":
			k.wrappingKeyID = value
		case "wrappingKeyHex":
			k.wrappingKeyHex = value
		case "publicKeyPEM":
			k.publicKeyPEM = value
		case "privateKeyPEM":
			k.privateKeyPEM = value
		case "providerID":
			k.providerID = value
		}
	}

	k.currentField++
	k.textInput.SetValue("")

	// Check if we've completed all fields
	if k.currentField >= len(k.fields) {
		k.step = KasKeyStepConfirm
		return k, nil
	}

	// Set up next field
	nextField := k.fields[k.currentField]
	switch nextField {
	case "wrappingKeyID":
		k.textInput.Placeholder = "e.g., wrapping-key-id"
	case "wrappingKeyHex":
		k.textInput.Placeholder = "Hex-encoded wrapping key"
	case "publicKeyPEM":
		k.textInput.Placeholder = "Base64-encoded PEM public key"
	case "privateKeyPEM":
		k.textInput.Placeholder = "Base64-encoded PEM private key"
	case "providerID":
		k.textInput.Placeholder = "Provider config name or ID"
	}

	return k, textinput.Blink
}

func (k *KasKeyWizard) executeCreate() (*KasKeyWizard, tea.Cmd) {
	k.step = KasKeyStepExecuting

	return k, func() tea.Msg {
		// Note: This is a simplified implementation
		// A full implementation would construct the proper PublicKeyCtx and PrivateKeyCtx
		// based on the mode and provided values

		// For now, return an informative error about what would be created
		return kasKeyResultMsg{
			err: fmt.Errorf("key creation via wizard not yet fully implemented. Use CLI: otdfctl policy kas-registry key create --kas-id %s --key-id %s --algorithm %s --mode %s",
				k.kasID, k.keyID, k.algorithm, k.mode.String()),
		}
	}
}

// View renders the KAS key wizard
func (k *KasKeyWizard) View() string {
	var sb strings.Builder

	sb.WriteString(wizardTitleStyle.Render("Create KAS Key") + "\n")
	sb.WriteString(wizardDescStyle.Render(fmt.Sprintf("For KAS: %s", k.kasName)) + "\n\n")

	switch k.step {
	case KasKeyStepMode:
		sb.WriteString(k.renderModeStep())

	case KasKeyStepConfig:
		sb.WriteString(k.renderConfigStep())

	case KasKeyStepModeSpecific:
		sb.WriteString(k.renderModeSpecificStep())

	case KasKeyStepConfirm:
		sb.WriteString(k.renderConfirmStep())

	case KasKeyStepExecuting:
		sb.WriteString(wizardDescStyle.Render("Creating key...") + "\n")

	case KasKeyStepComplete:
		if k.cancelled {
			sb.WriteString(wizardHintStyle.Render("Cancelled") + "\n")
		} else {
			sb.WriteString(wizardSuccessStyle.Render("✓ "+k.result) + "\n")
		}

	case KasKeyStepError:
		sb.WriteString(wizardErrorStyle.Render("Error: "+k.error) + "\n")
		sb.WriteString("\n" + wizardHintStyle.Render("Press Esc to go back") + "\n")
	}

	return sb.String()
}

func (k *KasKeyWizard) renderModeStep() string {
	var sb strings.Builder

	sb.WriteString(wizardLabelStyle.Render("Select Key Mode:") + "\n\n")

	modes := []struct {
		label string
		desc  string
	}{
		{"Local", "KAS generates key, you provide wrapping key"},
		{"Provider", "External provider stores wrapping key"},
		{"Remote", "Key stored entirely at external provider"},
		{"Public Key Only", "Import public key (no private key)"},
	}

	for i, mode := range modes {
		if i == k.selectIndex {
			sb.WriteString(wizardSelectedStyle.Render("> "+mode.label) + "\n")
			sb.WriteString(wizardHintStyle.Render("    "+mode.desc) + "\n")
		} else {
			sb.WriteString(wizardUnselectedStyle.Render("  "+mode.label) + "\n")
			sb.WriteString(wizardHintStyle.Render("    "+mode.desc) + "\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(wizardHintStyle.Render("↑/↓ to navigate, Enter to select, Esc to cancel"))

	return sb.String()
}

func (k *KasKeyWizard) renderConfigStep() string {
	var sb strings.Builder

	sb.WriteString(wizardLabelStyle.Render("Key Configuration") + "\n")
	sb.WriteString(wizardDescStyle.Render(fmt.Sprintf("Mode: %s", k.mode.String())) + "\n\n")

	if k.currentField == 0 {
		// Key ID input
		sb.WriteString(wizardLabelStyle.Render("Key ID:") + "\n")
		sb.WriteString(wizardDescStyle.Render("A unique identifier for this key") + "\n\n")
		sb.WriteString(k.textInput.View() + "\n")
		if k.error != "" {
			sb.WriteString(wizardErrorStyle.Render(k.error) + "\n")
		}
	} else {
		// Show entered Key ID
		sb.WriteString(wizardLabelStyle.Render("Key ID: "))
		sb.WriteString(wizardSuccessStyle.Render(k.keyID) + "\n\n")

		// Algorithm selection
		sb.WriteString(wizardLabelStyle.Render("Algorithm:") + "\n\n")
		for i, alg := range k.algorithms {
			if i == k.selectIndex {
				sb.WriteString(wizardSelectedStyle.Render("> "+alg.Label) + "\n")
			} else {
				sb.WriteString(wizardUnselectedStyle.Render("  "+alg.Label) + "\n")
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString(wizardHintStyle.Render("Enter to continue, Esc to cancel"))

	return sb.String()
}

func (k *KasKeyWizard) renderModeSpecificStep() string {
	var sb strings.Builder

	sb.WriteString(wizardLabelStyle.Render(fmt.Sprintf("Mode: %s - Additional Configuration", k.mode.String())) + "\n\n")

	// Show already entered values
	for i := 0; i < k.currentField && i < len(k.fields); i++ {
		fieldName := k.fields[i]
		var value string
		switch fieldName {
		case "wrappingKeyID":
			value = k.wrappingKeyID
		case "wrappingKeyHex":
			value = "***" // Don't show key material
		case "publicKeyPEM":
			value = "(provided)"
		case "privateKeyPEM":
			value = "(provided)"
		case "providerID":
			value = k.providerID
		}
		sb.WriteString(wizardLabelStyle.Render(fieldName+": "))
		sb.WriteString(wizardSuccessStyle.Render(value) + "\n")
	}

	// Current field
	if k.currentField < len(k.fields) {
		sb.WriteString("\n")
		fieldName := k.fields[k.currentField]
		sb.WriteString(wizardLabelStyle.Render(fieldName+":") + "\n")
		sb.WriteString(k.textInput.View() + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(wizardHintStyle.Render("Enter to continue, Esc to cancel"))

	return sb.String()
}

func (k *KasKeyWizard) renderConfirmStep() string {
	var sb strings.Builder

	sb.WriteString(wizardLabelStyle.Render("Confirm Key Creation") + "\n\n")

	sb.WriteString(fmt.Sprintf("%s: %s\n", wizardLabelStyle.Render("KAS"), wizardDescStyle.Render(k.kasName)))
	sb.WriteString(fmt.Sprintf("%s: %s\n", wizardLabelStyle.Render("Mode"), wizardDescStyle.Render(k.mode.String())))
	sb.WriteString(fmt.Sprintf("%s: %s\n", wizardLabelStyle.Render("Key ID"), wizardSelectedStyle.Render(k.keyID)))
	sb.WriteString(fmt.Sprintf("%s: %s\n", wizardLabelStyle.Render("Algorithm"), wizardDescStyle.Render(k.algorithm)))

	if k.wrappingKeyID != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", wizardLabelStyle.Render("Wrapping Key ID"), wizardDescStyle.Render(k.wrappingKeyID)))
	}
	if k.providerID != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", wizardLabelStyle.Render("Provider"), wizardDescStyle.Render(k.providerID)))
	}

	sb.WriteString("\n")
	sb.WriteString(wizardHintStyle.Render("Press Enter to create, Esc to cancel"))

	return sb.String()
}

// IsComplete returns true if the wizard has finished
func (k *KasKeyWizard) IsComplete() bool {
	return k.step == KasKeyStepComplete || k.step == KasKeyStepError
}

// WasCancelled returns true if the wizard was cancelled
func (k *KasKeyWizard) WasCancelled() bool {
	return k.cancelled
}

// GetResult returns the result message
func (k *KasKeyWizard) GetResult() string {
	if k.step == KasKeyStepError {
		return k.error
	}
	return k.result
}

// GetError returns the error message if any
func (k *KasKeyWizard) GetError() string {
	return k.error
}
