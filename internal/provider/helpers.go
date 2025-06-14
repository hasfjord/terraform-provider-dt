// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// flattenStringList converts a list of strings to a list of string values.
func flattenStringListToAttr(ctx context.Context, input []string) (basetypes.ListValue, diag.Diagnostics) {
	if input == nil {
		return types.ListNull(types.StringType), nil
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, input)
	if diags.HasError() {
		return types.ListValueMust(types.StringType, []attr.Value{}), diags
	}
	return list, nil
}

func flattenStringSetToAttr(ctx context.Context, input []string) (basetypes.SetValue, diag.Diagnostics) {
	if input == nil {
		return types.SetNull(types.StringType), nil
	}
	set, diags := types.SetValueFrom(ctx, types.StringType, input)
	if diags.HasError() {
		return types.SetValueMust(types.StringType, []attr.Value{}), diags
	}
	return set, nil
}

// expandStringList converts a list of string values to a list of strings.
func expandStringList(ctx context.Context, listValue basetypes.ListValue) ([]string, diag.Diagnostics) {
	if listValue.IsNull() {
		return nil, nil
	}
	if listValue.IsUnknown() {
		return nil, nil
	}
	elements := make([]types.String, 0, len(listValue.Elements()))
	diags := listValue.ElementsAs(ctx, &elements, false)

	var result []string
	for _, element := range elements {
		result = append(result, element.ValueString())
	}

	return result, diags
}

func expandStringSet(ctx context.Context, setValue basetypes.SetValue) ([]string, diag.Diagnostics) {
	if setValue.IsNull() {
		return nil, nil
	}
	if setValue.IsUnknown() {
		return nil, nil
	}
	elements := make([]types.String, 0, len(setValue.Elements()))
	diags := setValue.ElementsAs(ctx, &elements, false)

	var result []string
	for _, element := range elements {
		result = append(result, element.ValueString())
	}

	return result, diags
}

var durationValidator = stringvalidator.RegexMatches(
	regexp.MustCompile(`^(\d+)([s])$`),
	"Duration must be in the format of <number><unit>, where unit is 's' (seconds).",
)
