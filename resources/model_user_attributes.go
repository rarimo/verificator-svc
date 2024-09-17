/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserAttributes struct {
	// Lower user age limit
	AgeLowerBound *int32 `json:"age_lower_bound,omitempty"`
	// Event ID of event
	EventId *string `json:"event_id,omitempty"`
	// User nationality
	Nationality *string `json:"nationality,omitempty"`
	// Enable verification of sex param
	Sex *bool `json:"sex,omitempty"`
	// Parameters for checking user uniqueness
	Uniqueness *bool `json:"uniqueness,omitempty"`
}
