package verification

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateEmail(t *testing.T) {
	t.Parallel()
	errTest := errors.New("test")
	ctxBackground := context.Background()

	tests := map[string]struct {
		email          string
		makeMXLookuper func(ctrl *gomock.Controller) MXLookuper
		errWrapped     error
		errMessage     string
	}{
		"Invalid email format": {
			email:      "aa",
			errWrapped: ErrEmailFormatNotValid,
			errMessage: "email format is not valid: aa",
		},
		"Valid email format but not existing": {
			email: "aa@aa.aa",
			makeMXLookuper: func(ctrl *gomock.Controller) MXLookuper {
				lookuper := NewMockMXLookuper(ctrl)
				lookuper.EXPECT().LookupMX(ctxBackground, "aa.aa").Return(nil, errTest)
				return lookuper
			},
			errWrapped: ErrEmailHostUnreachable,
			errMessage: "email host is not reachable: for host aa.aa: test",
		},
		"no_mx_record_found": {
			email: "aa@aa.aa",
			makeMXLookuper: func(ctrl *gomock.Controller) MXLookuper {
				lookuper := NewMockMXLookuper(ctrl)
				lookuper.EXPECT().LookupMX(ctxBackground, "aa.aa").Return([]*net.MX{}, nil)
				return lookuper
			},
			errWrapped: ErrEmailHostUnreachable,
			errMessage: "email host is not reachable: for host aa.aa: no MX record found",
		},
		"Valid email format and existing": {
			email: "aa@aa.aa",
			makeMXLookuper: func(ctrl *gomock.Controller) MXLookuper {
				lookuper := NewMockMXLookuper(ctrl)
				lookuper.EXPECT().LookupMX(ctxBackground, "aa.aa").Return([]*net.MX{{}}, nil)
				return lookuper
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			var mxLookuper MXLookuper
			if tc.makeMXLookuper != nil {
				mxLookuper = tc.makeMXLookuper(ctrl)
			}

			err := ValidateEmail(ctxBackground, tc.email, mxLookuper)

			assert.ErrorIs(t, err, tc.errWrapped)
			if tc.errWrapped != nil {
				assert.EqualError(t, err, tc.errMessage)
			}
		})
	}
}
