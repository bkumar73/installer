// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewV2GetNextStepsParams creates a new V2GetNextStepsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV2GetNextStepsParams() *V2GetNextStepsParams {
	return &V2GetNextStepsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV2GetNextStepsParamsWithTimeout creates a new V2GetNextStepsParams object
// with the ability to set a timeout on a request.
func NewV2GetNextStepsParamsWithTimeout(timeout time.Duration) *V2GetNextStepsParams {
	return &V2GetNextStepsParams{
		timeout: timeout,
	}
}

// NewV2GetNextStepsParamsWithContext creates a new V2GetNextStepsParams object
// with the ability to set a context for a request.
func NewV2GetNextStepsParamsWithContext(ctx context.Context) *V2GetNextStepsParams {
	return &V2GetNextStepsParams{
		Context: ctx,
	}
}

// NewV2GetNextStepsParamsWithHTTPClient creates a new V2GetNextStepsParams object
// with the ability to set a custom HTTPClient for a request.
func NewV2GetNextStepsParamsWithHTTPClient(client *http.Client) *V2GetNextStepsParams {
	return &V2GetNextStepsParams{
		HTTPClient: client,
	}
}

/* V2GetNextStepsParams contains all the parameters to send to the API endpoint
   for the v2 get next steps operation.

   Typically these are written to a http.Request.
*/
type V2GetNextStepsParams struct {

	/* DiscoveryAgentVersion.

	   The software version of the discovery agent that is retrieving instructions.
	*/
	DiscoveryAgentVersion *string

	/* HostID.

	   The host that is retrieving instructions.

	   Format: uuid
	*/
	HostID strfmt.UUID

	/* InfraEnvID.

	   The infra-env of the host that is retrieving instructions.

	   Format: uuid
	*/
	InfraEnvID strfmt.UUID

	/* Timestamp.

	   The time on the host as seconds since the Unix epoch.
	*/
	Timestamp *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v2 get next steps params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V2GetNextStepsParams) WithDefaults() *V2GetNextStepsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v2 get next steps params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V2GetNextStepsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v2 get next steps params
func (o *V2GetNextStepsParams) WithTimeout(timeout time.Duration) *V2GetNextStepsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v2 get next steps params
func (o *V2GetNextStepsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v2 get next steps params
func (o *V2GetNextStepsParams) WithContext(ctx context.Context) *V2GetNextStepsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v2 get next steps params
func (o *V2GetNextStepsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v2 get next steps params
func (o *V2GetNextStepsParams) WithHTTPClient(client *http.Client) *V2GetNextStepsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v2 get next steps params
func (o *V2GetNextStepsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDiscoveryAgentVersion adds the discoveryAgentVersion to the v2 get next steps params
func (o *V2GetNextStepsParams) WithDiscoveryAgentVersion(discoveryAgentVersion *string) *V2GetNextStepsParams {
	o.SetDiscoveryAgentVersion(discoveryAgentVersion)
	return o
}

// SetDiscoveryAgentVersion adds the discoveryAgentVersion to the v2 get next steps params
func (o *V2GetNextStepsParams) SetDiscoveryAgentVersion(discoveryAgentVersion *string) {
	o.DiscoveryAgentVersion = discoveryAgentVersion
}

// WithHostID adds the hostID to the v2 get next steps params
func (o *V2GetNextStepsParams) WithHostID(hostID strfmt.UUID) *V2GetNextStepsParams {
	o.SetHostID(hostID)
	return o
}

// SetHostID adds the hostId to the v2 get next steps params
func (o *V2GetNextStepsParams) SetHostID(hostID strfmt.UUID) {
	o.HostID = hostID
}

// WithInfraEnvID adds the infraEnvID to the v2 get next steps params
func (o *V2GetNextStepsParams) WithInfraEnvID(infraEnvID strfmt.UUID) *V2GetNextStepsParams {
	o.SetInfraEnvID(infraEnvID)
	return o
}

// SetInfraEnvID adds the infraEnvId to the v2 get next steps params
func (o *V2GetNextStepsParams) SetInfraEnvID(infraEnvID strfmt.UUID) {
	o.InfraEnvID = infraEnvID
}

// WithTimestamp adds the timestamp to the v2 get next steps params
func (o *V2GetNextStepsParams) WithTimestamp(timestamp *int64) *V2GetNextStepsParams {
	o.SetTimestamp(timestamp)
	return o
}

// SetTimestamp adds the timestamp to the v2 get next steps params
func (o *V2GetNextStepsParams) SetTimestamp(timestamp *int64) {
	o.Timestamp = timestamp
}

// WriteToRequest writes these params to a swagger request
func (o *V2GetNextStepsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.DiscoveryAgentVersion != nil {

		// header param discovery_agent_version
		if err := r.SetHeaderParam("discovery_agent_version", *o.DiscoveryAgentVersion); err != nil {
			return err
		}
	}

	// path param host_id
	if err := r.SetPathParam("host_id", o.HostID.String()); err != nil {
		return err
	}

	// path param infra_env_id
	if err := r.SetPathParam("infra_env_id", o.InfraEnvID.String()); err != nil {
		return err
	}

	if o.Timestamp != nil {

		// query param timestamp
		var qrTimestamp int64

		if o.Timestamp != nil {
			qrTimestamp = *o.Timestamp
		}
		qTimestamp := swag.FormatInt64(qrTimestamp)
		if qTimestamp != "" {

			if err := r.SetQueryParam("timestamp", qTimestamp); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
