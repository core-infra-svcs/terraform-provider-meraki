package tools

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/http/httputil"
	"regexp"
)

func obfuscateMerakiApiKey(data string) string {
	re := regexp.MustCompile("(?m)^.*X-Cisco-Meraki-Api-Key.*$[\r\n]+")
	obfuscated := re.ReplaceAllString(data, "X-Cisco-Meraki-Api-Key: ***OBFUSCATED***\n")
	return obfuscated
}

// CollectHttpDiagnostics - responsible for gathering and logging HTTP driven events
func CollectHttpDiagnostics(ctx context.Context, diags *diag.Diagnostics, httpResp *http.Response) *diag.Diagnostics {

	// Collect HTTP request diagnostics
	reqDump, err := httputil.DumpRequestOut(httpResp.Request, false) // already consumed
	if err != nil {
		diags.AddWarning(
			"Failed to gather HTTP request diagnostics", fmt.Sprintf("\n%s", err),
		)
	}

	// Collect HTTP response diagnostics
	respDump, err := httputil.DumpResponse(httpResp, true)
	if err != nil {
		diags.AddWarning(
			"Failed to gather HTTP inlineResp diagnostics", fmt.Sprintf("\n%s", err),
		)
	}

	// Check for errors after diagnostics collected
	if diags.HasError() {

		diags.AddError(
			"Request Diagnostics:",
			fmt.Sprintf("\n%s", obfuscateMerakiApiKey(string(reqDump))),
		)

		diags.AddError(
			"Response Diagnostics:",
			fmt.Sprintf("\n%s", respDump),
		)
	} else {
		// Write logs
		tflog.Trace(ctx, fmt.Sprintf("Request Diagnostics\n%s", obfuscateMerakiApiKey(string(reqDump))))
		tflog.Trace(ctx, fmt.Sprintf("Response Diagnostics\n%s", respDump))
	}

	return diags

}
