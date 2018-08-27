package phone

import (
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (e *Validator) StreamValidate(in proto.Phone_StreamValidateServer) error {
	return nil
}

func (e *Validator) Validate(ctx context.Context, in *proto.PhoneValidateRequest) (*proto.PhoneValidateReply, error) {
	return &proto.PhoneValidateReply{}, nil
}
