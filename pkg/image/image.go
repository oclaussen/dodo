package image

import (
	"io"
	"net"

	"github.com/docker/docker/api/types"
	dodotypes "github.com/oclaussen/dodo/pkg/types"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var (
	errMissingImageID = errors.New(
		"build complete, but the server did not send an image id")
)

// Image represents the data necessary to build a docker image
type Image struct {
	config      *ImageConfig
	client      Client
	authConfigs map[string]types.AuthConfig
	session     session
}

// ImageConfig represents the build configuration for a docker image
type ImageConfig struct {
	Name         string
	Context      string
	Dockerfile   string
	Steps        []string
	Args         dodotypes.KeyValueList
	NoCache      bool
	ForceRebuild bool
	ForcePull    bool
	PrintOutput  bool
}

// Client represents a docker client that can do everything this package
// needs
type Client interface {
	DialSession(context.Context, string, map[string][]string) (net.Conn, error)
	ImageBuild(context.Context, io.Reader, types.ImageBuildOptions) (types.ImageBuildResponse, error)
}

// NewImage initializes and validates a new Image object
func NewImage(client Client, authConfigs map[string]types.AuthConfig, config *ImageConfig) (*Image, error) {
	if client == nil {
		return nil, errors.New("client may not be nil")
	}

	// TODO: do this somewhere else
	if config.Context == "" {
		config.Context = "."
	}

	session, err := prepareSession(config.Context)
	if err != nil {
		return nil, err
	}

	return &Image{
		client:      client,
		authConfigs: authConfigs,
		config:      config,
		session:     session,
	}, nil
}
