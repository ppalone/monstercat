package monstercat_test

import (
	"testing"

	"github.com/ppalone/monstercat"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	c := monstercat.NewClient(nil)
	assert.NotNil(t, c)
}
