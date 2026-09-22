package hooks

// Config configures the container-hooks mode. When Enabled is true, the runner
// stops calling the Docker API for job container operations and instead 
// invokes Cmd (optionally through Entrypoint) for every hook event.
type Config struct {
    Enabled    bool   `yaml:"enabled"`    // Enabled turns container-hooks mode on. When false, the runner uses the built-in Docker backend and ignores Entrypoint and Cmd.
    Entrypoint string `yaml:"entrypoint"` // Entrypoint is an optional interpreter used to launch Cmd (for example, "/bin/sh" on Linux or "cmd.exe" on Windows). Empty means Cmd is executed directly.
    Cmd        string `yaml:"cmd"`        // Cmd is the path to the single hook entrypoint the runner invokes for every hook event.
}

// ContainerStepSpecs describes a step that runs in its own container
// (uses: docker://...). Collected via stepDescriber.describe() before 
// execution starts and passed to the prepare hook inside prepare-input.json.
type ContainerStepSpecs struct {
	StepNumber  int
	StepID      string
	StepName    string
	Image       string
	Entrypoint  []string
	Args        []string
	WorkingDir  string
	Env         map[string]string
	Credentials *struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}