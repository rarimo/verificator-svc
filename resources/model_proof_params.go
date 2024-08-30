/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ProofParams struct {
	Key
	Attributes ProofParamsAttributes `json:"attributes"`
}
type ProofParamsResponse struct {
	Data     ProofParams `json:"data"`
	Included Included    `json:"included"`
}

type ProofParamsListResponse struct {
	Data     []ProofParams `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustProofParams - returns ProofParams from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustProofParams(key Key) *ProofParams {
	var proofParams ProofParams
	if c.tryFindEntry(key, &proofParams) {
		return &proofParams
	}
	return nil
}
