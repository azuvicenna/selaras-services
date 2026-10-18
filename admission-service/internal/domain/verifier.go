package domain

import "context"

type DepartmentVerifier interface {
	VerifyActiveDepartment(ctx context.Context, departmentID string) error
}
