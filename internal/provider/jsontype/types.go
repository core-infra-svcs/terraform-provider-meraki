package jsontype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var BoolType = boolType{
	BoolType: types.BoolType,
}

type boolType struct {
	basetypes.BoolType
}

func (mst boolType) Equal(o attr.Type) bool {
	_, ok := o.(boolType)
	return ok
}

func (mst boolType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mst.BoolType.ValueFromTerraform(ctx, in)

	return Bool{BoolValue: val.(types.Bool)}, err
}

var StringType = stringType{
	types.StringType,
}

// stringType -
type stringType struct {
	basetypes.StringType
}

func (mst stringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mst.StringType.ValueFromTerraform(ctx, in)

	return String{StringValue: val.(types.String)}, err
}

func (mst stringType) Equal(o attr.Type) bool {
	_, ok := o.(stringType)
	return ok
}
