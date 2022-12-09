package disciplinarian

import future.keywords.in

deny[message] {
	thing := input.list[i]

	keys := object.keys(thing)

	unsupported_keys := keys - {"name", "url", "type"}

	count(unsupported_keys) > 0

	message := sprintf("unsupported keys for thing %d: %v", [i, unsupported_keys])
}

deny[message] {
	thing := input.list[i]

	not thing.type in {"fruit"}

	message := sprintf("thing %d is not a fruit", [i])
}
