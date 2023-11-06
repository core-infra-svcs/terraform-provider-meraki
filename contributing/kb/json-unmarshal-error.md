# JSON Unmarshal Error

## Overview
In this kb article you will learn a technique to deal with situations were the HTTP server returns a json payload in a different format 
than expected from the HTTP client, the resources internal data model, or in rare cases the openAPI specification. 



## Problem

In rare cases you may encounter a server-side error in which the json response is malformed. 
In the below example we encountered a json response in which the server fails to strip an extra 
double-quote from the end of the array:

`{"tags": ["a", "b", "c""]}`

We could just pass this malformed data as a string however this is inconsistent with the openAPI spec 
and the HTTP client which expects an array of strings.

```go
// A simple data model which defines an array of strings.
type NetworkResourceModel struct {
Id                      types.String                    `tfsdk:"id" json:"-"`
Tags                    []jsonTypes.String               `tfsdk:"tags" json:"tags"`
}
```

## Solution

The long-term solution can only be addressed server-side, however an interim workaround is to create a new type 
and extend it with a function to handle the error and return the correct tag format:

```go
type Tag string

func (t *Tag) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*t = Tag(strings.Trim(s, `"`))
	return nil
}

// Use the new tag type in our struct
type NetworkResourceModel struct {
Id                      types.String        `tfsdk:"id" json:"-"`
Tags                    []Tag               `tfsdk:"tags" json:"tags"`
}

// Decode as usual
if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
resp.Diagnostics.AddError(
"JSON decoding error",
fmt.Sprintf("%v\n", err.Error()),
)
return
}

This approuch is easy to manage and can be removed once the server-side issue is resolved.

```