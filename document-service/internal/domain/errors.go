package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDocumentVoid      = errors.New("document is void")
	ErrTemplateNotFound  = errors.New("document template not found")
	ErrCompilationFailed = errors.New("failed to compile document to PDF")
)
