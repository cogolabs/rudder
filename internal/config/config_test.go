package config

import (
	"os"
	"testing"

	"github.com/ryantking/rudder/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TearDownTest() {
	require := suite.Require()

	err := os.RemoveAll(configName)
	require.NoError(err)
}

func (suite *ConfigTestSuite) TestLoadConfig() {
	assert := suite.Assert()
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/main.yaml")
	require.NoError(err)
	cfg, err := Load()
	require.NoError(err)
	assert.Equal(&testConfig, cfg)
}

func TestConfigTestSuite(t *testing.T) {
	tests := new(ConfigTestSuite)
	suite.Run(t, tests)
}
