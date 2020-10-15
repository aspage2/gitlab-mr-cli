package stuff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyThing(t *testing.T) {
	assert.Equal(t, int32(5), MyThing(2, 3))
}
