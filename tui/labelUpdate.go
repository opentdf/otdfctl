package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
)

type LabelUpdate struct {
	label  LabelItem
	update Update
	attr   *policy.Attribute
	sdk    handlers.Handler
}

func InitLabelUpdate(label LabelItem, attr *policy.Attribute, sdk handlers.Handler) LabelUpdate {
	// label := attr.Metadata.Labels[labelIdx]
	return LabelUpdate{
		label:  label,
		update: InitUpdate([]string{"Key", "Value"}, []string{label.title, label.description}),
		attr:   attr,
		sdk:    sdk,
	}
}

func (m LabelUpdate) Init() tea.Cmd {
	return nil
}

func (m LabelUpdate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// return InitLabelList(m.attr, m.sdk)
			if m.update.focusIndex == len(m.update.inputs) {
				// update the label
				metadata := &common.MetadataMutable{Labels: m.attr.Metadata.Labels}
				oldKey := m.label.title
				// oldVal := m.label.description
				newKey := m.update.inputs[0].Value()
				newVal := m.update.inputs[1].Value()
				if oldKey != newKey {
					delete(metadata.Labels, oldKey)
				}
				metadata.Labels[newKey] = newVal
				// metadata := common.MetadataMutable{Labels: map[string]string{"abc": "def"}}
				behavior := common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_REPLACE
				// behavior := common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_EXTEND
				attr, err := m.sdk.UpdateAttribute(m.attr.Id, metadata, behavior)
				if err != nil {
					// return error view
				}
				return InitLabelList(attr, m.sdk)
			}
		}
	}
	update, cmd := m.update.Update(msg)
	m.update = update.(Update)
	return m, cmd
}

func (m LabelUpdate) View() string {
	return m.update.View()
}
