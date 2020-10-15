package stuff

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestMyThing(t *testing.T) {
	assert.Equal(t, 5, MyThing(2, 3))
}
