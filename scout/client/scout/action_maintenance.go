// generated code; DO NOT EDIT; 2020-07-17 11:13:40.57595 +0200 CEST m=+0.029825291"
//
// Client for Choria RPC Agent 'scout'' Version 0.0.1 generated using Choria version 0.14.0

package scoutclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/choria-io/go-choria/protocol"
	rpcclient "github.com/choria-io/go-choria/providers/agent/mcorpc/client"
)

// MaintenanceRequester performs a RPC request to scout#maintenance
type MaintenanceRequester struct {
	r    *requester
	outc chan *MaintenanceOutput
}

// MaintenanceOutput is the output from the maintenance action
type MaintenanceOutput struct {
	details *ResultDetails
	reply   map[string]interface{}
}

// MaintenanceResult is the result from a maintenance action
type MaintenanceResult struct {
	stats   *rpcclient.Stats
	outputs []*MaintenanceOutput
}

// Stats is the rpc request stats
func (d *MaintenanceResult) Stats() Stats {
	return d.stats
}

// ResultDetails is the details about the request
func (d *MaintenanceOutput) ResultDetails() *ResultDetails {
	return d.details
}

// HashMap is the raw output data
func (d *MaintenanceOutput) HashMap() map[string]interface{} {
	return d.reply
}

// JSON is the JSON representation of the output data
func (d *MaintenanceOutput) JSON() ([]byte, error) {
	return json.Marshal(d.reply)
}

// ParseOutput parses the result value from the Maintenance action into target
func (d *MaintenanceOutput) ParseMaintenanceOutput(target interface{}) error {
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
func (d *MaintenanceRequester) Do(ctx context.Context) (*MaintenanceResult, error) {
	dres := &MaintenanceResult{}

	handler := func(pr protocol.Reply, r *rpcclient.RPCReply) {
		output := &MaintenanceOutput{
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

		if d.outc != nil {
			d.outc <- output
			return
		}

		dres.outputs = append(dres.outputs, output)
	}

	res, err := d.r.do(ctx, handler)
	if err != nil {
		return nil, err
	}

	dres.stats = res

	return dres, nil
}

// EachOutput iterates over all results received
func (d *MaintenanceResult) EachOutput(h func(r *MaintenanceOutput)) {
	for _, resp := range d.outputs {
		h(resp)
	}
}

// Checks is an optional input to the maintenance action
//
// Description: Check to pause, empty means all
func (d *MaintenanceRequester) Checks(v []interface{}) *MaintenanceRequester {
	d.r.args["checks"] = v

	return d
}

// Failed is the value of the failed output
//
// Description: List of checks that could not be paused
func (d *MaintenanceOutput) Failed() []interface{} {
	val := d.reply["failed"]
	return val.([]interface{})
}

// Skipped is the value of the skipped output
//
// Description: List of checks that was skipped
func (d *MaintenanceOutput) Skipped() []interface{} {
	val := d.reply["skipped"]
	return val.([]interface{})
}

// Transitioned is the value of the transitioned output
//
// Description: List of checks that were paused
func (d *MaintenanceOutput) Transitioned() []interface{} {
	val := d.reply["transitioned"]
	return val.([]interface{})
}
