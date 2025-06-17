/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserV2Attributes struct {
	// Birth date lower bound in hex format
	BirthDateLowerBound string `json:"birth_date_lower_bound"`
	// Birth date upper bound in hex format
	BirthDateUpperBound string `json:"birth_date_upper_bound"`
	// Citizenship mask in hex format
	CitizenshipMask string `json:"citizenship_mask"`
	// Event data in hex format
	EventData string `json:"event_data"`
	// Event ID
	EventId string `json:"event_id"`
	// Expiration date lower bound in hex format
	ExpirationDateLowerBound string `json:"expiration_date_lower_bound"`
	// Expiration date upper bound in hex format
	ExpirationDateUpperBound string `json:"expiration_date_upper_bound"`
	// Identity counter
	IdentityCounter int32 `json:"identity_counter"`
	// Identity counter lower bound
	IdentityCounterLowerBound int32 `json:"identity_counter_lower_bound"`
	// Identity counter upper bound
	IdentityCounterUpperBound int32 `json:"identity_counter_upper_bound"`
	// Selector value
	Selector string `json:"selector"`
	// Timestamp lower bound
	TimestampLowerBound string `json:"timestamp_lower_bound"`
	// Timestamp upper bound
	TimestampUpperBound string `json:"timestamp_upper_bound"`
}
