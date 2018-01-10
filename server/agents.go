package server

import (
	"context"
	"fmt"

	"github.com/choria-io/go-choria/agents/choriautil"
	"github.com/choria-io/go-choria/agents/discovery"
	"github.com/choria-io/go-choria/agents/provision"
	"github.com/choria-io/go-choria/build"
)

func (srv *Instance) setupCoreAgents(ctx context.Context) error {
	da, err := discovery.New(srv.agents)
	if err != nil {
		return fmt.Errorf("Could not setup initial agents: %s", err.Error())
	}

	srv.agents.RegisterAgent(ctx, "discovery", da, srv.connector)

	cu, err := choriautil.New(srv.agents)
	if err != nil {
		return fmt.Errorf("Could not setup choria_util agent: %s", err.Error())
	}

	srv.agents.RegisterAgent(ctx, "choria_util", cu, srv.connector)

	if build.ProvisionBrokerURLs != "" {
		pa, err := provision.New(srv.agents)
		if err != nil {
			return fmt.Errorf("Could not setup choria_provision agent: %s", err.Error())
		}

		srv.agents.RegisterAgent(ctx, "choria_provision", pa, srv.connector)
	}

	return nil
}
