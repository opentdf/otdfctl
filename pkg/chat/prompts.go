package chat

// GetPrompts returns a map of questions and their corresponding prompts for the chat utility.
func GetPrompts() map[string]string {
	return map[string]string{
		"How do I initialize a new OpenTDF platform instance using the otdfctl CLI?":                             "How do I initialize a new OpenTDF platform instance using the otdfctl CLI?",
		"Can you explain how to define and apply an ABAC policy to a data object using otdfctl?":                 "Can you explain how to define and apply an ABAC policy to a data object using otdfctl?",
		"What are the key components of the ZTDF and how can I use otdfctl to utilize them?":                     "What are the key components of the ZTDF and how can I use otdfctl to utilize them?",
		"How can I encrypt a file using otdfctl and ensure my attributes travel with the data?":                  "How can I encrypt a file using otdfctl and ensure my attributes travel with the data?",
		"How do I enforce an attribute-based policy on a data object using otdfctl?":                             "How do I enforce an attribute-based policy on a data object using otdfctl?",
		"What commands are available in otdfctl for managing an OpenTDF platform instance?":                      "What commands are available in otdfctl for managing an OpenTDF platform instance?",
		"How can otdfctl and OpenTDF be used together to ensure compliance with data protection regulations?":    "How can otdfctl and OpenTDF be used together to ensure compliance with data protection regulations?",
		"How can I use otdfctl to decrypt a file and verify that the applied ABAC policy is correctly enforced?": "How can I use otdfctl to decrypt a file and verify that the applied ABAC policy is correctly enforced?",
	}
}
