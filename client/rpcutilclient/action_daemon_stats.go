// generated code; DO NOT EDIT; 2021-01-12 09:27:23.500027 +0100 CET m=+0.015344331"
//
// Client for Choria RPC Agent 'rpcutil'' Version 0.19.0 generated using Choria version 0.18.0

package rpcutilclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/choria-io/go-choria/protocol"
	rpcclient "github.com/choria-io/go-choria/providers/agent/mcorpc/client"
	"github.com/choria-io/go-choria/providers/agent/mcorpc/ddl/agent"
	"github.com/choria-io/go-choria/providers/agent/mcorpc/replyfmt"
)

// DaemonStatsRequester performs a RPC request to rpcutil#daemon_stats
type DaemonStatsRequester struct {
	r    *requester
	outc chan *DaemonStatsOutput
}

// DaemonStatsOutput is the output from the daemon_stats action
type DaemonStatsOutput struct {
	details *ResultDetails
	reply   map[string]interface{}
}

// DaemonStatsResult is the result from a daemon_stats action
type DaemonStatsResult struct {
	ddl        *agent.DDL
	stats      *rpcclient.Stats
	outputs    []*DaemonStatsOutput
	rpcreplies []*replyfmt.RPCReply
	mu         sync.Mutex
}

func (d *DaemonStatsResult) RenderResults(w io.Writer, format RenderFormat, displayMode DisplayMode, verbose bool, silent bool, colorize bool, log Log) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.stats == nil {
		return fmt.Errorf("result stats is not set, result was not completed")
	}

	results := &replyfmt.RPCResults{
		Agent:   d.stats.Agent(),
		Action:  d.stats.Action(),
		Replies: d.rpcreplies,
		Stats:   d.stats,
	}

	addl, err := d.ddl.ActionInterface(d.stats.Action())
	if err != nil {
		return err
	}

	switch format {
	case JSONFormat:
		return results.RenderJSON(w, addl)
	case TableFormat:
		return results.RenderTable(w, addl)
	case TXTFooter:
		results.RenderTXTFooter(w, verbose)
		return nil
	default:
		return results.RenderTXT(w, addl, verbose, silent, replyfmt.DisplayMode(displayMode), colorize, log)
	}
}

// Stats is the rpc request stats
func (d *DaemonStatsResult) Stats() Stats {
	return d.stats
}

// ResultDetails is the details about the request
func (d *DaemonStatsOutput) ResultDetails() *ResultDetails {
	return d.details
}

// HashMap is the raw output data
func (d *DaemonStatsOutput) HashMap() map[string]interface{} {
	return d.reply
}

// JSON is the JSON representation of the output data
func (d *DaemonStatsOutput) JSON() ([]byte, error) {
	return json.Marshal(d.reply)
}

// ParseOutput parses the result value from the DaemonStats action into target
func (d *DaemonStatsOutput) ParseDaemonStatsOutput(target interface{}) error {
	j, err := d.JSON()
	if err != nil {
		return fmt.Errorf("could not access payload: %s", err)
	}

	err = json.Unmarshal(j, target)
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON payload: %s", err)
	}

	return nil
}

// Do performs the request
func (d *DaemonStatsRequester) Do(ctx context.Context) (*DaemonStatsResult, error) {
	dres := &DaemonStatsResult{ddl: d.r.client.ddl}

	handler := func(pr protocol.Reply, r *rpcclient.RPCReply) {
		// filtered by expr filter
		if r == nil {
			return
		}

		output := &DaemonStatsOutput{
			reply: make(map[string]interface{}),
			details: &ResultDetails{
				sender:  pr.SenderID(),
				code:    int(r.Statuscode),
				message: r.Statusmsg,
				ts:      pr.Time(),
			},
		}

		err := json.Unmarshal(r.Data, &output.reply)
		if err != nil {
			d.r.client.errorf("Could not decode reply from %s: %s", pr.SenderID(), err)
		}

		// caller wants a channel so we dont return a resulset too (lots of memory etc)
		// this is unused now, no support for setting a channel
		if d.outc != nil {
			d.outc <- output
			return
		}

		// else prepare our result set
		dres.mu.Lock()
		dres.outputs = append(dres.outputs, output)
		dres.rpcreplies = append(dres.rpcreplies, &replyfmt.RPCReply{
			Sender:   pr.SenderID(),
			RPCReply: r,
		})
		dres.mu.Unlock()
	}

	res, err := d.r.do(ctx, handler)
	if err != nil {
		return nil, err
	}

	dres.stats = res

	return dres, nil
}

// EachOutput iterates over all results received
func (d *DaemonStatsResult) EachOutput(h func(r *DaemonStatsOutput)) {
	for _, resp := range d.outputs {
		h(resp)
	}
}

// Agents is the value of the agents output
//
// Description: List of agents loaded
func (d *DaemonStatsOutput) Agents() interface{} {
	val := d.reply["agents"]
	return val.(interface{})
}

// Configfile is the value of the configfile output
//
// Description: Config file used to start the daemon
func (d *DaemonStatsOutput) Configfile() interface{} {
	val := d.reply["configfile"]
	return val.(interface{})
}

// Filtered is the value of the filtered output
//
// Description: Didn't pass filter checks
func (d *DaemonStatsOutput) Filtered() interface{} {
	val := d.reply["filtered"]
	return val.(interface{})
}

// Passed is the value of the passed output
//
// Description: Passed filter checks
func (d *DaemonStatsOutput) Passed() interface{} {
	val := d.reply["passed"]
	return val.(interface{})
}

// Pid is the value of the pid output
//
// Description: Process ID of the daemon
func (d *DaemonStatsOutput) Pid() interface{} {
	val := d.reply["pid"]
	return val.(interface{})
}

// Replies is the value of the replies output
//
// Description: Replies sent back to clients
func (d *DaemonStatsOutput) Replies() interface{} {
	val := d.reply["replies"]
	return val.(interface{})
}

// Starttime is the value of the starttime output
//
// Description: Time the server started
func (d *DaemonStatsOutput) Starttime() interface{} {
	val := d.reply["starttime"]
	return val.(interface{})
}

// Threads is the value of the threads output
//
// Description: List of threads active in the daemon
func (d *DaemonStatsOutput) Threads() interface{} {
	val := d.reply["threads"]
	return val.(interface{})
}

// Times is the value of the times output
//
// Description: Processor time consumed by the daemon
func (d *DaemonStatsOutput) Times() interface{} {
	val := d.reply["times"]
	return val.(interface{})
}

// Total is the value of the total output
//
// Description: Total messages received
func (d *DaemonStatsOutput) Total() interface{} {
	val := d.reply["total"]
	return val.(interface{})
}

// Ttlexpired is the value of the ttlexpired output
//
// Description: Messages that did pass TTL checks
func (d *DaemonStatsOutput) Ttlexpired() interface{} {
	val := d.reply["ttlexpired"]
	return val.(interface{})
}

// Unvalidated is the value of the unvalidated output
//
// Description: Messages that failed security validation
func (d *DaemonStatsOutput) Unvalidated() interface{} {
	val := d.reply["unvalidated"]
	return val.(interface{})
}

// Validated is the value of the validated output
//
// Description: Messages that passed security validation
func (d *DaemonStatsOutput) Validated() interface{} {
	val := d.reply["validated"]
	return val.(interface{})
}

// Version is the value of the version output
//
// Description: MCollective Version
func (d *DaemonStatsOutput) Version() interface{} {
	val := d.reply["version"]
	return val.(interface{})
}
