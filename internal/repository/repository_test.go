package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepoSlug(t *testing.T) {
	for _, test := range [...]string{
		"git@gitlab.com:mintel/gitlab-mr-cli",
		"https://gitlab.com/mintel/gitlab-mr-cli",
		"git@gitlab.com:mintel/gitlab-mr-cli.git",
	} {
		t.Run(test, func(t *testing.T) {
			res, err := GetRepoSlug(test)
			assert.NoError(t, err)
			assert.Equal(t, "mintel/gitlab-mr-cli", res)

		})
	}
}
