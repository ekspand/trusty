package cis_test

import (
	"context"
	"os"
	"testing"

	"github.com/ekspand/trusty/api/v1/pb"
	"github.com/ekspand/trusty/backend/service/cis"
	"github.com/ekspand/trusty/client"
	"github.com/ekspand/trusty/client/embed"
	"github.com/ekspand/trusty/internal/appcontainer"
	"github.com/ekspand/trusty/internal/config"
	"github.com/ekspand/trusty/pkg/gserver"
	"github.com/ekspand/trusty/tests/mockpb"
	"github.com/ekspand/trusty/tests/testutils"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	projFolder = "../../../"
)

var (
	trustyServer *gserver.Server
	cisClient    client.CIClient
	raMock       = &mockpb.MockRAServer{}

	// serviceFactories provides map of trustyserver.ServiceFactory
	serviceFactories = map[string]gserver.ServiceFactory{
		cis.ServiceName: cis.Factory,
	}
)

func TestMain(m *testing.M) {
	var err error
	//	xlog.SetPackageLogLevel("github.com/go-phorce/dolly/xhttp", "retriable", xlog.DEBUG)

	cfg, err := testutils.LoadConfig(projFolder, "UNIT_TEST")
	if err != nil {
		panic(errors.Trace(err))
	}

	httpAddr := testutils.CreateURLs("http", "")

	httpcfg := &config.HTTPServer{
		ListenURLs: []string{httpAddr},
		Services:   []string{cis.ServiceName},
	}

	disco := appcontainer.NewDiscovery()
	disco.Register("MockRAServer", raMock)

	container, err := appcontainer.NewContainerFactory(nil).
		WithConfigurationProvider(func() (*config.Configuration, error) {
			return cfg, nil
		}).
		WithDiscoveryProvider(func() (appcontainer.Discovery, error) {
			return disco, nil
		}).
		CreateContainerWithDependencies()
	if err != nil {
		panic(errors.Trace(err))
	}

	trustyServer, err = gserver.Start("cis", httpcfg, container, serviceFactories)
	if err != nil || trustyServer == nil {
		panic(errors.Trace(err))
	}
	cisClient = embed.NewCIClient(trustyServer)

	err = trustyServer.Service("cis").(*cis.Service).OnStarted()
	if err != nil {
		panic(errors.Trace(err))
	}

	// Run the tests
	rc := m.Run()

	// cleanup
	trustyServer.Close()

	os.Exit(rc)
}

func TestRoots(t *testing.T) {
	raMock.SetResponse(&pb.RootsResponse{
		Roots: []*pb.RootCertificate{
			{
				Id: 1,
			},
		},
	})

	res, err := cisClient.GetRoots(context.Background(), &empty.Empty{})
	require.NoError(t, err)
	assert.Len(t, res.Roots, 1)

	raMock.SetError(errors.New("unable to read DB"))
	_, err = cisClient.GetRoots(context.Background(), &empty.Empty{})
	require.Error(t, err)
	assert.Equal(t, "failed to get Roots: unable to read DB", err.Error())
}
