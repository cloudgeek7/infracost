package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	hcl2 "github.com/infracost/infracost/internal/hcl/block"
)

func LoadBlocksFromFile(file *hcl.File) (hcl.Blocks, error) {
	contents, diagnostics := file.Body.Content(hcl2.TerraformSchema012)
	if diagnostics != nil && diagnostics.HasErrors() {
		return nil, diagnostics
	}

	if contents == nil {
		return nil, fmt.Errorf("file contents is empty")
	}

	return contents.Blocks, nil
}