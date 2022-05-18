package cmd

type state struct {
	Debug    bool
	Quiet    bool
	Kind     string
	name     string
	isInside bool
}

type commandValues struct {
	use     string
	short   string
	long    string
	options []string
}

var chemotionValues = commandValues{
	use:   "chemotion",
	short: "CLI for Chemotion ELN",
	long: `Chemotion ELN is an Electronic Lab Notebook solution.
	Developed for, and by, researchers, the software aims
	to work for you. See, https://www.chemotion.net.`,
	options: []string{"instance", "system"},
}
