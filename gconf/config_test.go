package gconf

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type nestedFoo struct {
	Int            int     `yaml:"Int"`
	Float          float64 `yaml:"Float"`
	Int2           int     `yaml:"Int2"`
	EmbeddingText1 string  `file:"./text/1.txt"`
	EnvOverride    *string `yaml:"EnvOverride" env:"ENV_OVERRIDE"`
	PureEnvValue   *string `env:"PURE_ENV_VALUE"`
}

type nestedFoo3 struct {
	PureEnvValue *string `env:"PURE_ENV_VALUE2"`
}

type config struct {
	Foo        string      `yaml:"Foo"`
	NestedFoo  *nestedFoo  `yaml:"NestedFoo"`
	NestedFoo2 *nestedFoo  `yaml:"NestedFoo2"`
	NestedFoo3 *nestedFoo3 `yaml:"NestedFoo3"`
}

func TestUnmarshalConfFromDir(t *testing.T) {
	s := config{}
	os.Setenv("ENV_OVERRIDE", "greet")
	os.Setenv("PURE_ENV_VALUE", "pure")
	err := UnmarshalConfFromDir("uat", "env1", "../tests", &s)
	require.NoError(t, err)
	require.Equal(t, "env1", s.Foo)
	require.Equal(t, 11, s.NestedFoo.Int)
	require.Equal(t, 99.99, s.NestedFoo.Float)
	require.Equal(t, 100, s.NestedFoo.Int2)
	require.NotNil(t, s.NestedFoo.EnvOverride)
	require.Equal(t, "greet", *s.NestedFoo.EnvOverride)
	require.NotNil(t, s.NestedFoo.PureEnvValue)
	require.Equal(t, "pure", *s.NestedFoo.PureEnvValue)
	require.Nil(t, s.NestedFoo2)
	require.Nil(t, s.NestedFoo3)
	require.Equal(t, `hello
greeting!
让我们说中文？`, s.NestedFoo.EmbeddingText1)
}
