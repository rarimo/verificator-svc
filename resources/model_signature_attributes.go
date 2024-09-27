/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type SignatureAttributes struct {
	// Generated proof's public signals, see which element corresponds to a certain one pub-signal: https://github.com/rarimo/passport-zk-circuits#query-circuit-public-signals
	PubSignals []string `json:"pub_signals"`
	// Signature, must be len(signature) == 64
	Signature string `json:"signature"`
}
