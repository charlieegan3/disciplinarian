package disciplinarian

deny[message] {
	not input.checks
	message := "disciplinarian config files must have checks set"
}

deny[message] {
	s := input.checks[_].sources[_]
	object.get(s, "path", "") == ""

	message := "sources cannot have an empty path"
}

deny[message] {
	s := input.checks[_].sources[_]
	keys := object.keys(s)

	keys != {"path"}

	message := sprintf("unexpected keys in source: %v", [keys])
}
