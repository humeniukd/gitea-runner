package schema

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGiteaSchemaFactory(t *testing.T) {
	schema := GetGiteaWorkflowSchema()
	_ = schema

	data, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile("gitea_workflow_schema.json", append(data, "\n"...), 0o600)
	require.NoError(t, err)
}
