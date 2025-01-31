package middleware_test

import (
	"fmt"
	"github.com/dcalsky/gogo"
	"github.com/dcalsky/gogo/base"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatch(t *testing.T) {
	t.Run("unknown exception", func(t *testing.T) {
		exception := base.Catch(func() {
			panic(fmt.Errorf("test error"))
		})

		require.NotNil(t, exception)
		require.Equal(t, gogo.InternalError.Message, exception.Message)
	})

	t.Run("defined exception", func(t *testing.T) {
		exception := base.Catch(func() {
			panic(gogo.InvalidParamErr)
		})

		require.NotNil(t, exception)
		require.Equal(t, gogo.InvalidParamErr.Message, exception.Message)
	})

	t.Run("no exception", func(t *testing.T) {
		exception := base.Catch(func() {
			// ok
		})

		require.Nil(t, exception)
	})

}

func TestExceptionToError(t *testing.T) {
	// non-nil exception
	except := &base.Exception{}
	var err error
	err = except
	assert.NotNil(t, err)

	// nil exception
	except = nil
	err = except
	assert.Nil(t, err)

	// func return a non-nil exception
	except = func() *base.Exception {
		return &base.Exception{}
	}()
	err = except
	assert.NotNil(t, err)

	// func return a nil exception
	except = func() *base.Exception {
		return nil
	}()
	err = except
	assert.Nil(t, err)

	// func return a non-nil exception as error
	err = func() error {
		return &base.Exception{}
	}()
	assert.NotNil(t, err)

	// func return a nil exception as error
	err = func() error {
		var exception *base.Exception
		return exception
	}()
	assert.Nil(t, err)
}
