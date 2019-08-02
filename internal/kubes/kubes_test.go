package kubes

import (
	"fmt"
	"os"
	"testing"

	"github.com/ryantking/rudder/internal/config"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

const (
	testDir   = "kube"
	testToken = "mykubestoken"
)

type KubesTestSuite struct {
	suite.Suite
}

func (suite *KubesTestSuite) SetupSuite() {
	os.Setenv(tokenVar, testToken)
}

func (suite *KubesTestSuite) TearDownSuite() {
	os.Unsetenv(tokenVar)
}

func (suite *KubesTestSuite) TearDownTest() {
	require := suite.Require()

	err := os.RemoveAll(testDir)
	require.NoError(err)
}

func (suite *KubesTestSuite) TestMakeConfig() {
	assert := suite.Assert()
	require := suite.Require()

	dply := config.Deployment{KubeServers: []string{"kubes.server.net"}, KubeNamespace: "myproj"}
	err := MakeConfig(testDir, dply, 0)
	require.NoError(err)

	f, err := os.Open(fmt.Sprintf("%s/%s", testDir, configName))
	require.NoError(err)
	cfg := new(Config)
	err = yaml.NewDecoder(f).Decode(cfg)
	require.NoError(err)
	assert.Equal(testToken, cfg.Users[0].User.Token)
	assert.Equal(dply.KubeServers[0], cfg.Clusters[0].Cluster.Server)
	assert.Equal(dply.KubeNamespace, cfg.Contexts[0].Context.Namespace)
}

func TestKubesTestSuite(t *testing.T) {
	tests := new(KubesTestSuite)
	suite.Run(t, tests)
}
