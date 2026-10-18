package domain

import "context"

type ICDVerifier interface {
	VerifyAndLookupICD10(ctx context.Context, code string) (name string, err error)
	VerifyAndLookupICD9CM(ctx context.Context, code string) (name string, err error)
}
