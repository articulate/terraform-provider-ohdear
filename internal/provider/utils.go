package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func diagErrorf(err error, format string, a ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf(format, a...),
			Detail:   err.Error(),
		},
	}
}
