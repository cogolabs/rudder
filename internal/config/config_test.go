package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryantking/rudder/internal/kubes"
	"github.com/ryantking/rudder/internal/testutil"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

const (
	testConfigPath = "kube/config"
	testToken      = "mykubestoken"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TearDownTest() {
	require := suite.Require()

	matches, err := filepath.Glob(fmt.Sprintf("./%s.*", configBase))
	require.NoError(err)
	for _, match := range matches {
		err := os.RemoveAll(match)
		require.NoError(err)
	}
	err = os.RemoveAll(filepath.Dir(testConfigPath))
	require.NoError(err)
}

func (suite *ConfigTestSuite) TestLoad() {
	assert := suite.Assert()
	require := suite.Require()
	tests := []struct {
		path string
		ext  string
	}{{"main.yml", "yml"}, {"main.json", "json"}}

	for _, tt := range tests {
		path := fmt.Sprintf("../../test/configs/%s", tt.path)
		toPath := fmt.Sprintf("%s.%s", configBase, tt.ext)
		err := testutil.WriteConfigTo(path, toPath)
		require.NoError(err)
		cfg, err := Load()
		require.NoError(err)
		assert.Equal(&testConfig, cfg)
		err = os.RemoveAll(toPath)
		require.NoError(err)
	}
}

func (suite *ConfigTestSuite) TestLoadBadEXT() {
	require := suite.Require()

	path := "../../test/configs/main.yml"
	toPath := fmt.Sprintf("%s.gson", configBase)
	err := testutil.WriteConfigTo(path, toPath)
	require.NoError(err)
	_, err = Load()
	require.EqualError(err, "unsupported config format: .gson")
}

func (suite *ConfigTestSuite) TestLoadDefaultTimeout() {
	assert := suite.Assert()
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/noTimeout.yml")
	require.NoError(err)
	cfg, err := Load()
	require.NoError(err)
	expected := testConfig
	expected.Containers[0].Timeout = defaultTimeout
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
	require.Equal(ErrConfigNotFound, err)
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
	require.EqualError(err, "required field missing: containers[0].image")
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

func (suite *ConfigTestSuite) TestLoadMissingKubeServerEndpoint() {
	require := suite.Require()

	err := testutil.WriteConfig("../../test/configs/missingServerEndpoint.yml")
	require.NoError(err)
	_, err = Load()
	require.EqualError(err, "required field missing: deployments[0].kube_servers[1].server")
}

func (suite *ConfigTestSuite) TestMakeConfig() {
	assert := suite.Assert()
	require := suite.Require()

	dply := Deployment{KubeServers: []KubeServer{{Server: "kubes.server.net"}}, KubeNamespace: "myproj"}
	err := dply.MakeKubesConfig(&User{Token: testToken}, testConfigPath, 0)
	require.NoError(err)

	f, err := os.Open(testConfigPath)
	require.NoError(err)
	cfg := new(kubes.Config)
	err = yaml.NewDecoder(f).Decode(cfg)
	require.NoError(err)
	assert.Equal(testToken, cfg.Users[0].User.Token)
	assert.Equal(dply.KubeServers[0].Server, cfg.Clusters[0].Cluster.Server)
	assert.Equal(dply.KubeNamespace, cfg.Contexts[0].Context.Namespace)
}

func (suite *ConfigTestSuite) TestShouldDeploy() {
	assert := suite.Assert()
	tests := []struct {
		dply         *Deployment
		branch       string
		tag          string
		shouldDeploy bool
	}{
		{&Deployment{Branch: "master"}, "master", "", true},
		{&Deployment{Branch: "master"}, "dev", "", false},
		{&Deployment{Branch: "master", OnlyTags: true}, "master", "", false},
		{&Deployment{Branch: "master", OnlyTags: true, tagsRegex: "^(multi-v.*)$"}, "master", "multi-v0.1", true},
		{&Deployment{Branch: "master", OnlyTags: true, tagsRegex: "^(multi-v.*)$"}, "master", "canary-v0.1", false},
	}

	for _, tt := range tests {
		shouldDeploy := tt.dply.ShouldDeploy(tt.branch, tt.tag)
		assert.Equal(tt.shouldDeploy, shouldDeploy)
	}
}

func TestConfigTestSuite(t *testing.T) {
	tests := new(ConfigTestSuite)
	suite.Run(t, tests)
}
