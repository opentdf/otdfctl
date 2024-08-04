package chat

// GetPrompts returns a map of questions and their corresponding prompts for the chat utility.
func GetPrompts() map[string]string {
	return map[string]string{
		"How do I initialize a new OpenTDF platform instance using the otdfctl CLI?":                             "To initialize a new OpenTDF platform instance using the otdfctl CLI, use the command `otdfctl init` and follow the prompts to configure your instance.",
		"Can you explain how to define and apply an ABAC policy to a data object using otdfctl?":                 "To define and apply an ABAC policy to a data object using otdfctl, use the `otdfctl policy create` command to define the policy and `otdfctl policy apply` to apply it to the data object.",
		"What are the key components of the ZTDF and how can I use otdfctl to utilize them?":                     "The key components of the ZTDF include the policy engine, attribute store, and enforcement point. Use otdfctl commands like `otdfctl policy`, `otdfctl attribute`, and `otdfctl enforce` to interact with these components.",
		"How can I encrypt a file using otdfctl and ensure my attributes travel with the data?":                  "To encrypt a file using otdfctl and ensure attributes travel with the data, use the `otdfctl encrypt` command and specify the attributes using the `--attributes` flag.",
		"How do I enforce an attribute-based policy on a data object using otdfctl?":                             "To enforce an attribute-based policy on a data object using otdfctl, first define the policy with `otdfctl policy create`, then apply it using `otdfctl policy apply`.",
		"What commands are available in otdfctl for managing an OpenTDF platform instance?":                      "Commands available in otdfctl for managing an OpenTDF platform instance include `otdfctl init`, `otdfctl policy`, `otdfctl attribute`, `otdfctl encrypt`, `otdfctl decrypt`, and `otdfctl enforce`.",
		"How can otdfctl and OpenTDF be used together to ensure compliance with data protection regulations?":    "otdfctl and OpenTDF can be used together to ensure compliance with data protection regulations by defining and enforcing policies that control access based on attributes, ensuring data is encrypted and decrypted according to these policies.",
		"How can I use otdfctl to decrypt a file and verify that the applied ABAC policy is correctly enforced?": "To decrypt a file and verify that the applied ABAC policy is correctly enforced, use the `otdfctl decrypt` command and check the policy enforcement logs to ensure compliance.",
	}
}
