// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// GcpEditIamUserAccountOauthSessionURIRequest gcp edit iam user account oauth session Uri request
//
// swagger:model GcpEditIamUserAccountOauthSessionUriRequest
type GcpEditIamUserAccountOauthSessionURIRequest struct {

	// in: path
	// Required: true
	ID *string `json:"id"`
}

// Validate validates this gcp edit iam user account oauth session Uri request
func (m *GcpEditIamUserAccountOauthSessionURIRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GcpEditIamUserAccountOauthSessionURIRequest) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this gcp edit iam user account oauth session Uri request based on context it is used
func (m *GcpEditIamUserAccountOauthSessionURIRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GcpEditIamUserAccountOauthSessionURIRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GcpEditIamUserAccountOauthSessionURIRequest) UnmarshalBinary(b []byte) error {
	var res GcpEditIamUserAccountOauthSessionURIRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
