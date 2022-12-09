package discipline

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/charlieegan3/disciplinarian/pkg/config"
)

func TestRun(t *testing.T) {
	cfg := &config.Config{
		Checks: []config.Check{
			{
				Sources: []config.Source{
					{
						Path: "fixtures/config.yaml", // self config file which is invalid
					},
				},
				Policies: []config.Source{
					{
						IsDir: true,
						Path:  "fixtures/example/disciplinarian",
					},
				},
			},
		},
	}

	results, err := Run(context.Background(), cfg)
	require.NoError(t, err)

	expectedResults := []Result{
		{
			File: "fixtures/config.yaml",
			Messages: []string{
				"disciplinarian config files must have checks set",
			},
		},
	}

	assert.Equal(t, expectedResults, results)
}
