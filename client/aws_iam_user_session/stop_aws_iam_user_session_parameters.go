// Code generated by go-swagger; DO NOT EDIT.

package aws_iam_user_session

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
)

// NewStopAwsIamUserSessionParams creates a new StopAwsIamUserSessionParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewStopAwsIamUserSessionParams() *StopAwsIamUserSessionParams {
	return &StopAwsIamUserSessionParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewStopAwsIamUserSessionParamsWithTimeout creates a new StopAwsIamUserSessionParams object
// with the ability to set a timeout on a request.
func NewStopAwsIamUserSessionParamsWithTimeout(timeout time.Duration) *StopAwsIamUserSessionParams {
	return &StopAwsIamUserSessionParams{
		timeout: timeout,
	}
}

// NewStopAwsIamUserSessionParamsWithContext creates a new StopAwsIamUserSessionParams object
// with the ability to set a context for a request.
func NewStopAwsIamUserSessionParamsWithContext(ctx context.Context) *StopAwsIamUserSessionParams {
	return &StopAwsIamUserSessionParams{
		Context: ctx,
	}
}

// NewStopAwsIamUserSessionParamsWithHTTPClient creates a new StopAwsIamUserSessionParams object
// with the ability to set a custom HTTPClient for a request.
func NewStopAwsIamUserSessionParamsWithHTTPClient(client *http.Client) *StopAwsIamUserSessionParams {
	return &StopAwsIamUserSessionParams{
		HTTPClient: client,
	}
}

/* StopAwsIamUserSessionParams contains all the parameters to send to the API endpoint
   for the stop aws iam user session operation.

   Typically these are written to a http.Request.
*/
type StopAwsIamUserSessionParams struct {

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the stop aws iam user session params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StopAwsIamUserSessionParams) WithDefaults() *StopAwsIamUserSessionParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the stop aws iam user session params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StopAwsIamUserSessionParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) WithTimeout(timeout time.Duration) *StopAwsIamUserSessionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) WithContext(ctx context.Context) *StopAwsIamUserSessionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) WithHTTPClient(client *http.Client) *StopAwsIamUserSessionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) WithID(id string) *StopAwsIamUserSessionParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the stop aws iam user session params
func (o *StopAwsIamUserSessionParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *StopAwsIamUserSessionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
