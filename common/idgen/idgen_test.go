package idgen_test

import (
	"github.com/dcalsky/gogo/common/idgen"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	v := idgen.New[string]()
	require.NotEmpty(t, v)
}
