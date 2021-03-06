package appcontainer

import (
	"os"
	"path"
	"testing"

	"github.com/ekspand/trusty/internal/config"
	"github.com/go-phorce/dolly/algorithms/guid"
	"github.com/go-phorce/dolly/audit"
	"github.com/go-phorce/dolly/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContainerFactory(t *testing.T) {
	output := path.Join(os.TempDir(), "tests", "trusty", guid.MustCreate())
	os.MkdirAll(output, 0777)
	defer os.Remove(output)

	tcases := []struct {
		name string
		err  string
		cfg  *config.Configuration
	}{
		{
			name: "no_logs",
			cfg: &config.Configuration{
				Logs:  config.Logger{Directory: "/dev/null"},
				Audit: config.Logger{Directory: "/dev/null"},
			},
		},
		{
			name: "with_logs",
			cfg: &config.Configuration{
				/*
					Metrics: {
						Disabled: &falseVal,
					},
				*/
				Logs:  config.Logger{Directory: output},
				Audit: config.Logger{Directory: output},
			},
		},
	}

	for _, tc := range tcases {

		t.Run(tc.name, func(t *testing.T) {

			container, err := NewContainerFactory(nil).
				WithConfigurationProvider(func() (*config.Configuration, error) {
					return tc.cfg, nil
				}).
				CreateContainerWithDependencies()
			require.NoError(t, err)

			err = container.Invoke(func(cfg *config.Configuration,
				auditor audit.Auditor,
				scheduler tasks.Scheduler,
			) {
			})
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			}
		})
	}
}
