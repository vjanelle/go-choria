// generated code; DO NOT EDIT

package rpcutilclient

import (
	"time"

	coreclient "github.com/choria-io/go-choria/client/client"
	rpcclient "github.com/choria-io/go-choria/providers/agent/mcorpc/client"
)

// OptionReset resets the client options to use across requests to an empty list
func (p *RpcutilClient) OptionReset() *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = []rpcclient.RequestOption{}
	p.ns = p.clientOpts.ns
	p.targets = []string{}
	p.filters = []FilterFunc{
		FilterFunc(coreclient.AgentFilter("rpcutil")),
	}

	return p
}

// OptionIdentityFilter adds an identity filter
func (p *RpcutilClient) OptionIdentityFilter(f ...string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	for _, i := range f {
		p.filters = append(p.filters, FilterFunc(coreclient.IdentityFilter(i)))
	}

	p.ns.Reset()

	return p
}

// OptionClassFilter adds a class filter
func (p *RpcutilClient) OptionClassFilter(f ...string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	for _, i := range f {
		p.filters = append(p.filters, FilterFunc(coreclient.ClassFilter(i)))
	}

	p.ns.Reset()

	return p
}

// OptionFactFilter adds a fact filter
func (p *RpcutilClient) OptionFactFilter(f ...string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	for _, i := range f {
		p.filters = append(p.filters, FilterFunc(coreclient.FactFilter(i)))
	}

	p.ns.Reset()

	return p
}

// OptionAgentFilter adds an agent filter
func (p *RpcutilClient) OptionAgentFilter(a ...string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	for _, f := range a {
		p.filters = append(p.filters, FilterFunc(coreclient.AgentFilter(f)))
	}

	p.ns.Reset()

	return p
}

// OptionCombinedFilter adds a combined filter
func (p *RpcutilClient) OptionCombinedFilter(f ...string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	for _, i := range f {
		p.filters = append(p.filters, FilterFunc(coreclient.CombinedFilter(i)))
	}

	p.ns.Reset()

	return p
}

// OptionCompoundFilter adds a compound filter
func (p *RpcutilClient) OptionCompoundFilter(f ...string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	for _, i := range f {
		p.filters = append(p.filters, FilterFunc(coreclient.CompoundFilter(i)))
	}

	p.ns.Reset()

	return p
}

// OptionCollective sets the collective to target
func (p *RpcutilClient) OptionCollective(c string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = append(p.clientRPCOpts, rpcclient.Collective(c))
	return p
}

// OptionInBatches performs requests in batches
func (p *RpcutilClient) OptionInBatches(size int, sleep int) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = append(p.clientRPCOpts, rpcclient.InBatches(size, sleep))
	return p
}

// OptionDiscoveryTimeout configures the request discovery timeout, defaults to configured discovery timeout
func (p *RpcutilClient) OptionDiscoveryTimeout(t time.Duration) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = append(p.clientRPCOpts, rpcclient.DiscoveryTimeout(t))
	return p
}

// OptionLimitMethod configures the method to use when limiting targets - "random" or "first"
func (p *RpcutilClient) OptionLimitMethod(m string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = append(p.clientRPCOpts, rpcclient.LimitMethod(m))
	return p
}

// OptionLimitSize sets limits on the targets, either a number of a percentage like "10%"
func (p *RpcutilClient) OptionLimitSize(s string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = append(p.clientRPCOpts, rpcclient.LimitSize(s))
	return p
}

// OptionLimitSeed sets the random seed used to select targets when limiting and limit method is "random"
func (p *RpcutilClient) OptionLimitSeed(s int64) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.clientRPCOpts = append(p.clientRPCOpts, rpcclient.LimitSeed(s))
	return p
}

// OptionTargets sets specific node targets which would avoid discovery for all action calls until reset
func (p *RpcutilClient) OptionTargets(t []string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.targets = t
	return p
}

// OptionWorkers sets how many worker connections should be started to the broker
func (p *RpcutilClient) OptionWorkers(w int) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.workers = w
	return p
}

// OptionExprFilter sets a filter expression that will remove results from the result set
func (p *RpcutilClient) OptionExprFilter(f string) *RpcutilClient {
	p.Lock()
	defer p.Unlock()

	p.exprFilter = f
	return p
}
