package stuff

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestMyThing(t *testing.T) {
	assert.Equal(t, 5 ,MyThing(2,3))
}
