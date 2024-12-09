/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserAttributes struct {
	// Lower user age limit
	AgeLowerBound *int32 `json:"age_lower_bound,omitempty"`
	// Event ID of event
	EventId *string `json:"event_id,omitempty"`
	// Enable verification of expiration lower bound param
	ExpirationLowerBound *bool `json:"expiration_lower_bound,omitempty"`
	// User nationality
	Nationality *string `json:"nationality,omitempty"`
	// You can use this instead of 'nationality' params, it will check nationality bit in selector
	NationalityCheck *bool `json:"nationality_check,omitempty"`
	// Enable verification of sex param
	Sex *bool `json:"sex,omitempty"`
	// Parameters for checking user uniqueness
	Uniqueness *bool `json:"uniqueness,omitempty"`
}
