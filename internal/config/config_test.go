package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/ryantking/rudder/internal/kubes"
	"github.com/ryantking/rudder/internal/testutil"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

const (
	testDir   = "kube"
	testToken = "mykubestoken"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) SetupSuite() {
	os.Setenv(tokenVar, testToken)
}

func (suite *ConfigTestSuite) TearDownSuite() {
	os.Unsetenv(tokenVar)
}

func (suite *ConfigTestSuite) TearDownTest() {
	require := suite.Require()

	err := os.RemoveAll(configName)
	require.NoError(err)
	err = os.RemoveAll(testDir)
	require.NoError(err)
}

func (suite *ConfigTestSuite) TestLoad() {
	assert := suite.Assert()
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/main.yml")
	require.NoError(err)
	cfg, err := Load()
	require.NoError(err)
	assert.Equal(&testConfig, cfg)
}

func (suite *ConfigTestSuite) TestLoadDefaultTimeout() {
	assert := suite.Assert()
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/noTimeout.yml")
	require.NoError(err)
	cfg, err := Load()
	require.NoError(err)
	expected := testConfig
	expected.DockerTimeout = defaultTimeout
	assert.Equal(&expected, cfg)
}

func (suite *ConfigTestSuite) TestLoadDefaultNamespace() {
	assert := suite.Assert()
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/noNamespace.yml")
	require.NoError(err)
	cfg, err := Load()
	require.NoError(err)
	expected := testConfig
	expected.Deployments[0].KubeNamespace = defaultNamespace
	assert.Equal(&expected, cfg)
	expected.Deployments[0].KubeNamespace = "myproject"
}

func (suite *ConfigTestSuite) TestMissingConfig() {
	require := suite.Require()

	_, err := Load()
	require.True(os.IsNotExist(err))
}

func (suite *ConfigTestSuite) TestLoadBadYAML() {
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/badYAML.yml")
	require.NoError(err)
	_, err = Load()
	require.Error(err)
}

func (suite *ConfigTestSuite) TestLoadBadTimeout() {
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/badTimeout.yml")
	require.NoError(err)
	_, err = Load()
	require.Error(err)
}

func (suite *ConfigTestSuite) TestLoadMissingDockerImage() {
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/missingImage.yml")
	require.NoError(err)
	_, err = Load()
	require.EqualError(err, "required field missing: docker_image")
}

func (suite *ConfigTestSuite) TestLoadMissingName() {
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/missingName.yml")
	require.NoError(err)
	_, err = Load()
	require.EqualError(err, "required field missing: deployments[2].name")
}

func (suite *ConfigTestSuite) TestLoadMissingKubeServers() {
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/missingServers.yml")
	require.NoError(err)
	_, err = Load()
	require.EqualError(err, "required field missing: deployments[1].kube_servers")
}

func (suite *ConfigTestSuite) TestMakeConfig() {
	assert := suite.Assert()
	require := suite.Require()

	dply := Deployment{KubeServers: []string{"kubes.server.net"}, KubeNamespace: "myproj"}
	err := dply.MakeKubesConfig(testDir, 0)
	require.NoError(err)

	f, err := os.Open(fmt.Sprintf("%s/%s", testDir, kubesConfig))
	require.NoError(err)
	cfg := new(kubes.Config)
	err = yaml.NewDecoder(f).Decode(cfg)
	require.NoError(err)
	assert.Equal(testToken, cfg.Users[0].User.Token)
	assert.Equal(dply.KubeServers[0], cfg.Clusters[0].Cluster.Server)
	assert.Equal(dply.KubeNamespace, cfg.Contexts[0].Context.Namespace)
}

func TestConfigTestSuite(t *testing.T) {
	tests := new(ConfigTestSuite)
	suite.Run(t, tests)
}
