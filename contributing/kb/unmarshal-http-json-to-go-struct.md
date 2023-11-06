# Unmarshal HTTP JSON Payload into a Go Struct

## Overview
In this kb article you will learn a technique to extract a json payload in an http client library 
into a golang struct.



## Problem

Manually extracting deeply nested or multiple structs can be cumbersome to develop and troubleshoot depending on the size.

```go


// First you define the internal data model 
type Car struct {
	Wheels Wheels
	Paint Paint
	Interior Interior
}

type Wheels struct {
	Size string
	
}

type Paint struct {
	Color string
}

type Interior struct {
	...
}

// Then you unmarshal each (hopefully) strongly typed http response manually into new objects


var wheels Wheels
var paint Paint
var interior Interior 

wheels.Size = jsontypes.StringValue(inlineResp.Car.Wheels.Size)
paint.Color = jsontypes.StringValue(inlineResp.Car.Paint.Color)
interior.X.Y.Z = jsontypes.StringValue(inlineResp.Car.Interior.X.Y.Z)

var car Car{
    Wheels wheels
    Paint paint
    Interior interior
}

// Finally pull "Car" into internal go data model
data.Car = car

```
The larger the data model, the more cumbersome this task becomes requiring more time to develop 
and possibly duplicate code as HTTP calls are often reused for resources that use the same HTTP 
POST/PUT for Create/Update/Delete calls.

To reduce duplicate code developers may attempt to create a func to perform this payload extraction 
and then return a fully rendered struct but this can lead to internal state consistency issues within Terraform
and you will likely see an error message such as:

```text
 Error: Provider returned invalid result object after apply
        
        After the apply operation, the provider still indicated an unknown value for
        provider_name.test.car.paint.color. All values must be known
        after apply, so this is always a bug in the provider and should be reported
        in the provider's own repository. Terraform will still save the other known
        object values in the state.
```

## Solution

Leveraging tags in the struct itself can provide enough information to decode raw json data into a go struct 
by matching the tag names to the JSON keys. As long as the values are typed correctly this will unmarshal 
with very little effort:

```go

// First you define the internal data model WITH TAGS!
type Car struct {
    Wheels Wheels `tfsdk:"wheels" json:"Wheels"`
    Paint Paint `tfsdk:"paint" json:"Paint"`
    Interior Interior `tfsdk:"interior" json:"Interior"`
}

type Wheels struct {
Size string `tfsdk:"size" json:"Size"`

}

type Paint struct {
Color string `tfsdk:"color" json:"Color"`
}

type Interior struct {
... `tfsdk:"." json:"."`
}


// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}
	
```

This solution is easy to implement as we only need to define the data model once and allows us to extract the data 
regardless of if the HTTP client returns a strongly typed response or just casts payload data to an empty interface{}
which happens when there is incomplete api schema object data in the openapi spec.