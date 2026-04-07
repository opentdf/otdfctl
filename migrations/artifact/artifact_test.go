package artifact

import (
	"bytes"
	"testing"

	"github.com/Masterminds/semver/v3"
	artifactmetadata "github.com/opentdf/otdfctl/migrations/artifact/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRejectsUnsupportedSchemaVersion(t *testing.T) {
	t.Parallel()

	_, err := New(ArtifactOpts{
		Version: semver.MustParse("v2.0.0"),
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnsupportedSchemaVersion)
}

func TestNewDefaultsCurrentVersion(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	doc, err := New(ArtifactOpts{Writer: &buf})
	require.NoError(t, err)

	require.NoError(t, doc.Write())
	assert.Contains(t, buf.String(), `"schema": "v1.0.0"`)
	assert.Contains(t, buf.String(), `"name": "`+artifactmetadata.ArtifactName+`"`)
}

func TestArtifactSummaryReturnsEncodedJSON(t *testing.T) {
	t.Parallel()

	doc, err := New(ArtifactOpts{})
	require.NoError(t, err)

	summary, err := doc.Summary()
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"counts": {
			"namespaces": 0,
			"actions": 0,
			"subject_condition_sets": 0,
			"subject_mappings": 0,
			"registered_resources": 0,
			"obligation_triggers": 0,
			"skipped": 0
		}
	}`, string(summary))
}

func TestArtifactBuildAndCommitAreNotImplemented(t *testing.T) {
	t.Parallel()

	doc, err := New(ArtifactOpts{})
	require.NoError(t, err)

	buildErr := doc.Build()
	require.Error(t, buildErr)
	assert.ErrorContains(t, buildErr, "not implemented")

	commitErr := doc.Commit()
	require.Error(t, commitErr)
	assert.ErrorContains(t, commitErr, "not implemented")
}
