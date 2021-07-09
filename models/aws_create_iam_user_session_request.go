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

// AwsCreateIamUserSessionRequest aws create iam user session request
//
// swagger:model AwsCreateIamUserSessionRequest
type AwsCreateIamUserSessionRequest struct {

	// the account number of the aws account related to the role
	// Required: true
	AccountNumber *string `json:"accountNumber"`

	// aws access key Id
	AwsAccessKeyID string `json:"awsAccessKeyId,omitempty"`

	// aws secret key
	AwsSecretKey string `json:"awsSecretKey,omitempty"`

	// mfa device
	MfaDevice string `json:"mfaDevice,omitempty"`

	// profile name
	ProfileName string `json:"profileName,omitempty"`

	// the region on which the session will be initialized
	// Required: true
	Region *string `json:"region"`

	// the name which will be displayed
	// Required: true
	SessionName *string `json:"sessionName"`

	// user name
	UserName string `json:"userName,omitempty"`
}

// Validate validates this aws create iam user session request
func (m *AwsCreateIamUserSessionRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccountNumber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRegion(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSessionName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AwsCreateIamUserSessionRequest) validateAccountNumber(formats strfmt.Registry) error {

	if err := validate.Required("accountNumber", "body", m.AccountNumber); err != nil {
		return err
	}

	return nil
}

func (m *AwsCreateIamUserSessionRequest) validateRegion(formats strfmt.Registry) error {

	if err := validate.Required("region", "body", m.Region); err != nil {
		return err
	}

	return nil
}

func (m *AwsCreateIamUserSessionRequest) validateSessionName(formats strfmt.Registry) error {

	if err := validate.Required("sessionName", "body", m.SessionName); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this aws create iam user session request based on context it is used
func (m *AwsCreateIamUserSessionRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AwsCreateIamUserSessionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AwsCreateIamUserSessionRequest) UnmarshalBinary(b []byte) error {
	var res AwsCreateIamUserSessionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
