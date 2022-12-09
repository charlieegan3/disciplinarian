package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	expectedConfig := &Config{
		Checks: []Check{
			{
				Sources: []Source{
					{
						Path:    "config.yaml",
						IsValid: true,
					},
					{
						Path:    "",
						IsValid: false,
					},
				},
				Policies: []Source{
					{
						IsDir: true,
						Path:  "example/disciplinarian",
					},
				},
			},
		},
	}

	config, err := Load("./fixtures/config.yaml")
	require.NoError(t, err)

	assert.Equal(t, expectedConfig, config)
}
