/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AdvancedVerificationAttributes struct {
	// Lower user age limit
	AgeLowerBound *int32 `json:"age_lower_bound,omitempty"`
	// Birth date lower bound in hex format
	BirthDateLowerBound *string `json:"birth_date_lower_bound,omitempty"`
	// Birth date upper bound in hex format
	BirthDateUpperBound *string `json:"birth_date_upper_bound,omitempty"`
	// User citezenship mask
	CitizenshipMask *string `json:"citizenship_mask,omitempty"`
	// Event data in hex format
	EventData *string `json:"event_data,omitempty"`
	// Event ID
	EventId string `json:"event_id"`
	// Expiration date lower bound in hex format
	ExpirationDateLowerBound *string `json:"expiration_date_lower_bound,omitempty"`
	// Expiration date upper bound in hex format
	ExpirationDateUpperBound *string `json:"expiration_date_upper_bound,omitempty"`
	// Enable verification of expiration lower bound param
	ExpirationLowerBound *bool `json:"expiration_lower_bound,omitempty"`
	// Identity counter
	IdentityCounter *int64 `json:"identity_counter,omitempty"`
	// Identity counter lower bound
	IdentityCounterLowerBound *int64 `json:"identity_counter_lower_bound,omitempty"`
	// Identity counter upper bound
	IdentityCounterUpperBound *int64 `json:"identity_counter_upper_bound,omitempty"`
	// Selector value
	Selector string `json:"selector"`
	// User sex
	Sex *string `json:"sex,omitempty"`
	// Timestamp lower bound
	TimestampLowerBound *int64 `json:"timestamp_lower_bound,omitempty"`
	// Timestamp upper bound
	TimestampUpperBound *int64 `json:"timestamp_upper_bound,omitempty"`
}
