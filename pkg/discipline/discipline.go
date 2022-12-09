package discipline

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"

	"github.com/charlieegan3/disciplinarian/pkg/config"
)

type Result struct {
	File     string
	Messages []string
}

func Run(ctx context.Context, cfg *config.Config) ([]Result, error) {
	var results []Result

	for _, c := range cfg.Checks {
		// modules is a simple mapping containing unbundled policy files
		modules := make(map[string]string)

		// we're only interested in the results of these rules in the loaded policy files
		regoArgs := []func(*rego.Rego){rego.Query("data.disciplinarian.deny")}

		// load any directories as policy bundles
		for _, p := range c.Policies {
			if p.IsDir {
				regoArgs = append(regoArgs, rego.LoadBundle(p.Path))
			} else {
				content, err := os.ReadFile(p.Path)
				if err != nil {
					return nil, fmt.Errorf("failed to read policy file %s: %w", p.Path, err)
				}
				modules[p.Path] = string(content)
			}
		}

		compiler, err := ast.CompileModules(modules)
		if err != nil {
			return nil, fmt.Errorf("failed to compile polciy files: %w", err)
		}
		regoArgs = append(regoArgs, rego.Compiler(compiler))

		// partial of the check to run against of the check's sources
		check, err := rego.New(regoArgs...).PartialResult(ctx)
		if err != nil {
			log.Fatalf("failed to compute partial result: %s", err)
		}

		for _, s := range c.Sources {
			if !s.IsValid {
				fmt.Fprintf(os.Stderr, "skipping invalid source file: %q\n", s.Path)
				continue
			}

			var paths []string
			// handle a directory of files by non-recursively listing the files
			if s.IsDir {
				files, err := os.ReadDir(s.Path)
				if err != nil {
					return nil, fmt.Errorf("failed to list the files in the source directory %s: %w", s.Path, err)
				}

				for _, f := range files {
					paths = append(paths, fmt.Sprintf("%s/%s", s.Path, f.Name()))
				}
			} else {
				paths = append(paths, s.Path)
			}

			for _, p := range paths {
				b, err := os.ReadFile(p)
				if err != nil {
					return nil, fmt.Errorf("failed to read source file %s: %w", s.Path, err)
				}

				// TODO: support other formats of structured data
				var input map[string]interface{}
				if err := yaml.Unmarshal(b, &input); err != nil {
					return nil, fmt.Errorf("failed to unmarshal source file %s: %w", s.Path, err)
				}

				resultSet, err := check.Rego(rego.Input(input)).Eval(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate policy: %w", err)
				}

				// aggregate the results of the policy evaluation
				var messages []string
				for _, r := range resultSet {
					for _, e := range r.Expressions {
						for _, v := range e.Value.([]interface{}) {
							str, ok := v.(string)
							if ok {
								messages = append(messages, str)
							} else {
								b, err := json.Marshal(e)
								if err != nil {
									return nil, fmt.Errorf("failed to marshal expression for error message: %w", err)
								}
								return results, fmt.Errorf("failed parse policy output: %s", string(b))
							}
						}
					}
				}

				if len(messages) > 0 {
					results = append(results, Result{
						File:     p,
						Messages: messages,
					})
				}
			}
		}
	}

	return results, nil
}
