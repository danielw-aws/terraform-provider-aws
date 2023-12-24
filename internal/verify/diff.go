package verify

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-aws/names"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// Find JSON diff functions in the json.go file.

// SetTagsDiff sets the new plan difference with the result of
// merging resource tags on to those defined at the provider-level;
// returns an error if unsuccessful or if the resource tags are identical
// to those configured at the provider-level to avoid non-empty plans
// after resource READ operations as resource and provider-level tags
// will be indistinguishable when returned from an AWS API.
func SetTagsDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig
	fmt.Println("[Triage] diff.go SetTagsDiff ,diff=", diff)
	fmt.Println("[Triage] diff.go SetTagsDiff ,defaultTagsConfig=", defaultTagsConfig)
	fmt.Println("[Triage] diff.go SetTagsDiff ,ignoreTagsConfig=", ignoreTagsConfig)

	resourceTags := tftags.New(ctx, diff.Get("tags").(map[string]interface{}))

	fmt.Println("[Triage] diff.go SetTagsDiff ,resourceTags=", resourceTags)

	allTags := defaultTagsConfig.MergeTags(resourceTags).IgnoreConfig(ignoreTagsConfig)
	fmt.Println("[Triage] diff.go SetTagsDiff ,allTags=", allTags)
	// To ensure "tags_all" is correctly computed, we explicitly set the attribute diff
	// when the merger of resource-level tags onto provider-level tags results in n > 0 tags,
	// otherwise we mark the attribute as "Computed" only when there is a known diff (excluding an empty map)
	// or a change for "tags_all".
	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/18366
	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/19005
	//fmt.Println("[Triage] diff.go SetTagsDiff ,!diff.GetRawPlan().GetAttr(\"tags\").IsWhollyKnown()=", !diff.GetRawPlan().GetAttr("tags").IsWhollyKnown())
	if !diff.GetRawPlan().GetAttr("tags").IsWhollyKnown() {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
		return nil
	}

	if diff.HasChange("tags") {
		fmt.Println("[Triage] diff.go line 51 SetTagsDiff ,diff.HasChange(\"tags\")=", diff.HasChange("tags"))
		_, n := diff.GetChange("tags")
		fmt.Println("[Triage] diff.go diff.GetChange ,n=", n)
		newTags := tftags.New(ctx, n.(map[string]interface{}))
		fmt.Println("[Triage] line 55 diff.go newTags)=", newTags)

		fmt.Println("[Triage] diff.go before Set New Computed SetTagsDiff ,diff.Get(names.AttrTagsAll)=", diff.Get(names.AttrTagsAll))
		if newTags.HasZeroValue() {
			fmt.Println("[Triage] diff.go newTags.HasZeroValue()=", newTags.HasZeroValue())
			if err := diff.SetNewComputed("tags_all"); err != nil {
				return fmt.Errorf("error setting tags_all to computed: %w", err)
			}
		}
		fmt.Println("[Triage] diff.go after Set New Computed SetTagsDiff ,diff.Get(names.AttrTagsAll)=", diff.Get(names.AttrTagsAll))
		fmt.Println("[Triage] diff.go after Set New Computed SetTagsDiff ,diff.Get(names.AttrTags)=", diff.Get(names.AttrTags))
		if len(allTags) > 0 && (!newTags.HasZeroValue() || !allTags.HasZeroValue()) {
			fmt.Println("[Triage] diff.go SetTagsDiff ,len(allTags) > 0 && (!newTags.HasZeroValue() || !allTags.HasZeroValue())=", len(allTags) > 0 && (!newTags.HasZeroValue() || !allTags.HasZeroValue()))
			if err := diff.SetNew("tags_all", allTags.Map()); err != nil {
				return fmt.Errorf("error setting new tags_all diff: %w", err)
			}
		}

		if len(allTags) == 0 {
			fmt.Println("[Triage] diff.go SetTagsDiff ,len(allTags) == 0", len(allTags) == 0)
			if err := diff.SetNewComputed("tags_all"); err != nil {
				return fmt.Errorf("error setting tags_all to computed: %w", err)
			}
		}
	} else if !diff.HasChange("tags") {
		fmt.Println("[Triage] diff.go SetTagsDiff ,!diff.HasChange(\"tags\") ", !diff.HasChange("tags"))
		if len(allTags) > 0 && !allTags.HasZeroValue() {
			fmt.Println("[Triage] diff.go SetTagsDiff ,len(allTags) > 0 && !allTags.HasZeroValue()", len(allTags) > 0 && !allTags.HasZeroValue())
			if err := diff.SetNew("tags_all", allTags.Map()); err != nil {
				return fmt.Errorf("error setting new tags_all diff: %w", err)
			}
			return nil
		}

		var ta tftags.KeyValueTags
		if tagsAll, ok := diff.Get("tags_all").(map[string]interface{}); ok {
			fmt.Println("[Triage] diff.go SetTagsDiff line 90 ,tagsAll=", tagsAll)
			ta = tftags.New(ctx, tagsAll)
			fmt.Println("[Triage] diff.go SetTagsDiff ,ta=", ta)
		}
		if len(allTags) > 0 && !ta.DeepEqual(allTags) && allTags.HasZeroValue() {
			fmt.Println("[Triage] diff.go SetTagsDiff ,len(allTags) > 0 && !ta.DeepEqual(allTags) && allTags.HasZeroValue()=", len(allTags) > 0 && !ta.DeepEqual(allTags) && allTags.HasZeroValue())
			if err := diff.SetNewComputed("tags_all"); err != nil {
				return fmt.Errorf("error setting tags_all to computed: %w", err)
			}
			return nil
		}
	} else if tagsAll, ok := diff.Get("tags_all").(map[string]interface{}); ok {
		fmt.Println("[Triage] diff.go SetTagsDiff line 102,tagsAll=", tagsAll)
		ta := tftags.New(ctx, tagsAll)
		if !ta.DeepEqual(allTags) {
			fmt.Println("[Triage] diff.go SetTagsDiff !ta.DeepEqual(allTags)=", !ta.DeepEqual(allTags))
			if allTags.HasZeroValue() {
				fmt.Println("[Triage] diff.go SetTagsDiff allTags.HasZeroValue()=", allTags.HasZeroValue())
				if err := diff.SetNewComputed("tags_all"); err != nil {
					return fmt.Errorf("error setting tags_all to computed: %w", err)
				}
			}
		}
	} else if len(diff.Get("tags_all").(map[string]interface{})) > 0 {
		fmt.Println("[Triage] diff.go SetTagsDiff len(diff.Get(\"tags_all\").(map[string]interface{})) > 0=", len(diff.Get("tags_all").(map[string]interface{})) > 0)
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	} else if diff.HasChange("tags_all") {
		fmt.Println("[Triage] diff.go SetTagsDiff line 119 diff.HasChange(\"tags_all\")=", diff.HasChange("tags_all"))
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	}
	fmt.Println("[Triage] diff.go SetTagsDiff before nil")
	return nil
}

// SuppressEquivalentStringCaseInsensitive provides custom difference suppression
// for strings that are equal under case-insensitivity.
func SuppressEquivalentStringCaseInsensitive(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

// SuppressEquivalentRoundedTime returns a difference suppression function that compares
// two time value with the specified layout rounded to the specified duration.
func SuppressEquivalentRoundedTime(layout string, d time.Duration) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, _ *schema.ResourceData) bool {
		if old, err := time.Parse(layout, old); err == nil {
			if new, err := time.Parse(layout, new); err == nil {
				return old.Round(d).Equal(new.Round(d))
			}
		}

		return false
	}
}

// SuppressMissingOptionalConfigurationBlock handles configuration block attributes in the following scenario:
//   - The resource schema includes an optional configuration block with defaults
//   - The API response includes those defaults to refresh into the Terraform state
//   - The operator's configuration omits the optional configuration block
func SuppressMissingOptionalConfigurationBlock(k, old, new string, d *schema.ResourceData) bool {
	return old == "1" && new == "0"
}

// DiffStringMaps returns the set of keys and values that must be created, the set of keys
// and values that must be destroyed, and the set of keys and values that are unchanged.
func DiffStringMaps(oldMap, newMap map[string]interface{}) (map[string]*string, map[string]*string, map[string]*string) {
	// First, we're creating everything we have
	add := map[string]*string{}
	for k, v := range newMap {
		add[k] = aws.String(v.(string))
	}

	// Build the maps of what to remove and what is unchanged
	remove := map[string]*string{}
	unchanged := map[string]*string{}
	for k, v := range oldMap {
		old, ok := add[k]
		if !ok || aws.StringValue(old) != v.(string) {
			// Delete it!
			remove[k] = aws.String(v.(string))
		} else if ok {
			unchanged[k] = aws.String(v.(string))
			// already present so remove from new
			delete(add, k)
		}
	}

	return add, remove, unchanged
}
