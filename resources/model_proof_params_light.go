/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ProofParamsLight struct {
	Key
	Attributes ProofParamsLightAttributes `json:"attributes"`
}
type ProofParamsLightResponse struct {
	Data     ProofParamsLight `json:"data"`
	Included Included         `json:"included"`
}

type ProofParamsLightListResponse struct {
	Data     []ProofParamsLight `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustProofParamsLight - returns ProofParamsLight from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustProofParamsLight(key Key) *ProofParamsLight {
	var proofParamsLight ProofParamsLight
	if c.tryFindEntry(key, &proofParamsLight) {
		return &proofParamsLight
	}
	return nil
}
