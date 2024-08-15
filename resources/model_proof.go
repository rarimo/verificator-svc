/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Proof struct {
	Key
	Attributes ProofAttributes `json:"attributes"`
}
type ProofRequest struct {
	Data     Proof    `json:"data"`
	Included Included `json:"included"`
}

type ProofListRequest struct {
	Data     []Proof  `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustProof - returns Proof from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustProof(key Key) *Proof {
	var proof Proof
	if c.tryFindEntry(key, &proof) {
		return &proof
	}
	return nil
}
