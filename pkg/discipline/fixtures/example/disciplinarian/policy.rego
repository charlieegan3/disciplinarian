package disciplinarian

deny[message] {
	not input.checks
	message := "disciplinarian config files must have checks set"
}
