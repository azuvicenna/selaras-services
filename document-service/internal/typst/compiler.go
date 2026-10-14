package typst

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Compiler interface {
	Compile(ctx context.Context, source string) ([]byte, error)
	IsAvailable() bool
}

type CLICompiler struct {
	binPath string
}

func NewCLICompiler(binPath string) *CLICompiler {
	if binPath == "" {
		binPath = "typst"
	}
	return &CLICompiler{binPath: binPath}
}

func (c *CLICompiler) IsAvailable() bool {
	_, err := exec.LookPath(c.binPath)
	return err == nil
}

func (c *CLICompiler) Compile(ctx context.Context, source string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "typst-doc-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	inFile := filepath.Join(tmpDir, "document.typ")
	outFile := filepath.Join(tmpDir, "document.pdf")

	if err := os.WriteFile(inFile, []byte(source), 0600); err != nil {
		return nil, fmt.Errorf("failed to write typst source file: %w", err)
	}

	cmd := exec.CommandContext(ctx, c.binPath, "compile", inFile, outFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("typst compilation error: %s (%w)", stderr.String(), err)
	}

	pdfBytes, err := os.ReadFile(outFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read output pdf: %w", err)
	}

	return pdfBytes, nil
}

// MockCompiler provides a dummy PDF generator for testing environments
type MockCompiler struct{}

func NewMockCompiler() *MockCompiler {
	return &MockCompiler{}
}

func (m *MockCompiler) IsAvailable() bool {
	return true
}

func (m *MockCompiler) Compile(ctx context.Context, source string) ([]byte, error) {
	// Minimal valid PDF structure
	dummyPDF := []byte("%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n3 0 obj<</Type/Page/MediaBox[0 0 595 842]/Parent 2 0 R/Resources<<>>>>endobj\nxref\n0 4\n0000000000 65535 f \n0000000009 00000 n \n0000000052 00000 n \n0000000101 00000 n \ntrailer<</Size 4/Root 1 0 R>>\nstartxref\n180\n%%EOF\n")
	return dummyPDF, nil
}
