package email

import (
	"github.com/supme/service/pkg/dns"
	"github.com/supme/service/pkg/goroutine"
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/idna"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	// ToDo use github.com/opennota/re2dfa for speed up
	splitEmailRe = regexp.MustCompile(`^([A-Z0-9a-z._%+\-]+)@(.+\..{2,9})$`)
	dnsCache     *dns.MXCache
)

type Validator struct {
	broker *goroutine.Broker
}

func NewValidator(maxWorkers, dnsCacheExpirationSecond int) *Validator {
	broker := goroutine.NewBroker(int64(maxWorkers))
	dnsCache = dns.NewMXCache(dnsCacheExpirationSecond)
	return &Validator{
		broker: broker,
	}
}

func (e *Validator) StreamValidate(in proto.Email_StreamValidateServer) error {
	wg := &sync.WaitGroup{}
	for {
		v, err := in.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		wg.Add(1)
		e.broker.Next()
		go func(v *proto.EmailValidateRequest) {
			canonical, protoErr := e.validate(v.Email)
			r := proto.EmailValidateReply{Id: v.Id, Canonical: canonical, Valid: protoErr == proto.EmailValidateError_NO_ERROR, Error: protoErr}
			err = in.Send(&r)
			if err != nil {
				log.Println(err)
			}
			e.broker.Ready()
			wg.Done()
		}(v)
	}
	wg.Wait()
	return nil
}

func (e *Validator) Validate(ctx context.Context, in *proto.EmailValidateRequest) (*proto.EmailValidateReply, error) {
	e.broker.Next()
	defer e.broker.Ready()
	canonical, err := e.validate(in.Email)
	return &proto.EmailValidateReply{Id: in.Id, Canonical: canonical, Valid: err == proto.EmailValidateError_NO_ERROR, Error: err}, nil
}

func (e *Validator) validate(email string) (string, proto.EmailValidateError) {
	var eml, domain string
	s := strings.TrimSpace(email)
	if m := splitEmailRe.FindStringSubmatch(s); m != nil && len(m) == 3 {
		eml = strings.ToLower(strings.TrimSpace(m[1]))
		domain = strings.TrimRight(strings.ToLower(strings.TrimSpace(m[2])), ".")
	} else {
		return "", proto.EmailValidateError_BAD_EMAIL_FORMAT
	}

	punycode, err := idna.ToASCII(domain)
	if err != nil {
		return "", proto.EmailValidateError_BAD_EMAIL_FORMAT
	}

	canonicalizeEmail := strings.ToLower(eml + "@" + domain)

	_, err = dnsCache.Get(punycode)
	protoErr := dnsErrorToProtoError(err)
	//_, protoErr := e.domainCache.checkMX(punycode)

	return canonicalizeEmail, protoErr
}

func dnsErrorToProtoError(err error) proto.EmailValidateError {
	switch err {
	case nil:
		return proto.EmailValidateError_NO_ERROR
	case dns.ErrorHostDoesNotHaveMX:
		return proto.EmailValidateError_HOST_DOES_NOT_HAVE_MX
	case dns.ErrorMXHostInReservedIPRange:
		return proto.EmailValidateError_MX_HOST_IN_RESERVED_IP_RANGE
	case dns.ErrorOther:
		return proto.EmailValidateError_OTHER_ERROR
	}
	return proto.EmailValidateError_OTHER_ERROR
}
