package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Deansie/SJPI24-degree-project/internal/logging"
	"github.com/Deansie/SJPI24-degree-project/internal/validation"
	"github.com/spf13/cobra"
	"github.com/yannh/kubeconform/pkg/validator"
)

type Input struct {
	Name string
	Data []byte
}

var validateCmd = &cobra.Command{
	Use:   "validate [files or directories]",
	Short: "Validate Kubernetes manifests against schemas",
	Long: appLogo + `Validate one or more YAML files or directories containing Kubernetes manifests.
If no arguments are provided, manifests are read from stdin.`,
	Args: cobra.ArbitraryArgs,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().String("kubernetes-version", "1.29.0", "Kubernetes version to validate against")
	validateCmd.Flags().Bool("strict", true, "Disallow additional properties not in the schema")
	validateCmd.Flags().Bool("ignore-missing-schemas", false, "Skip validation for resource kinds without a schema")
}

// Command entrypoint

func runValidate(cmd *cobra.Command, args []string) error {
	logging.L().Info("Starting validation", "args", args)

	opts, err := readValidatorOptions(cmd)
	if err != nil {
		logging.L().Error("Failed to read options", "err", err)
		return err
	}

	inputs, err := collectInputs(args)
	if err != nil {
		logging.L().Error("Failed to collect inputs", "err", err)
		return err
	}

	logging.L().Debug("Inputs collected", "count", len(inputs))

	results := validateInputsPerFile(opts, inputs)

	if len(results) == 0 {
		fmt.Println("No YAML manifests found to validate")
		logging.L().Warn("No YAML manifests found")
		return nil
	}

	if errors := reportResults(results); errors > 0 {
		logging.L().Error("Validation failed", "errors", errors)
		fmt.Println("\nRule validation skipped due to schema validation errors")
		return fmt.Errorf("%d validation error(s) found", errors)
	}

	// Parse resources for rule validation

	var allResources []validation.Resource

	for _, in := range inputs {
		resources, err := validation.ParseResources(bytes.NewReader(in.Data))
		if err != nil {
			logging.L().Error("Failed to parse resources", "err", err)
			return fmt.Errorf("failed to parse resources: %w", err)
		}
		allResources = append(allResources, resources...)
	}

	logging.L().Debug("Resources parsed", "count", len(allResources))

	// Apply rules

	violations := validation.ValidateRules(allResources)
	if len(violations) > 0 {
		validation.ReportRuleViolations(violations)
		logging.L().Error("Rule violations found", "count", len(violations))
		return fmt.Errorf("%d rule violation(s) found", len(violations))
	}

	fmt.Println("✓ Rule validation passed (all resources comply with enforced rules)")
	logging.L().Info("Validation completed successfully")
	return nil
}

// Validator-setup

type validatorOptions struct {
	k8sVersion           string
	strict               bool
	ignoreMissingSchemas bool
}

func readValidatorOptions(cmd *cobra.Command) (validatorOptions, error) {
	k8sVersion, err := cmd.Flags().GetString("kubernetes-version")
	if err != nil {
		return validatorOptions{}, err
	}
	strict, err := cmd.Flags().GetBool("strict")
	if err != nil {
		return validatorOptions{}, err
	}
	ignoreMissingSchemas, err := cmd.Flags().GetBool("ignore-missing-schemas")
	if err != nil {
		return validatorOptions{}, err
	}

	return validatorOptions{
		k8sVersion:           k8sVersion,
		strict:               strict,
		ignoreMissingSchemas: ignoreMissingSchemas,
	}, nil
}

func newValidator(opts validatorOptions) (validator.Validator, error) {
	return validator.New(
		[]string{"default"},
		validator.Opts{
			KubernetesVersion:    opts.k8sVersion,
			Strict:               opts.strict,
			IgnoreMissingSchemas: opts.ignoreMissingSchemas,
			SkipTLS:              true,
		},
	)
}

// Input-collection

func collectInputs(args []string) ([]Input, error) {
	if len(args) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		return []Input{
			{
				Name: "stdin",
				Data: data,
			},
		}, nil
	}

	var inputs []Input

	for _, path := range args {
		files, err := expandPath(path)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}

			inputs = append(inputs, Input{
				Name: file,
				Data: content,
			})
		}
	}

	return inputs, nil
}

func expandPath(path string) ([]string, error) {
	if hasGlob(path) {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		return filterYAML(matches), nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if isYAML(path) {
			return []string{path}, nil
		}
		return nil, nil
	}

	var files []string
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isYAML(p) {
			return nil
		}
		files = append(files, p)
		return nil
	})

	return files, err
}

func hasGlob(path string) bool {
	for _, c := range path {
		switch c {
		case '*', '?', '[':
			return true
		}
	}
	return false
}

// Validation

func validateInputsPerFile(
	opts validatorOptions,
	inputs []Input,
) []validator.Result {

	var allResults []validator.Result

	for _, in := range inputs {
		v, err := newValidator(opts)
		if err != nil {
			allResults = append(allResults, validator.Result{
				Status: validator.Error,
				Err:    err,
			})
			continue
		}

		results := v.Validate(in.Name, io.NopCloser(bytes.NewReader(in.Data)))
		allResults = append(allResults, results...)
	}

	return allResults
}

// Reporting

func reportResults(results []validator.Result) int {
	errors := 0
	for _, res := range results {
		errors += printResult(res)
	}
	return errors
}

func printResult(res validator.Result) int {
	kind, name := "Unknown", "Unknown"
	if sig, err := res.Resource.Signature(); err == nil && sig != nil {
		kind = sig.Kind
		name = sig.Name
	}

	switch res.Status {
	case validator.Valid:
		fmt.Printf("✓ %s %s\n", kind, name)

	case validator.Invalid:
		fmt.Printf("✗ %s %s\n", kind, name)
		for _, ve := range res.ValidationErrors {
			fmt.Printf("  - %s: %s\n", ve.Path, ve.Msg)
		}
		return len(res.ValidationErrors)

	case validator.Error:
		fmt.Printf("! error: %v\n", res.Err)
		return 1
	}

	return 0
}

// Helpers

func isYAML(path string) bool {
	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return true
	default:
		return false
	}
}

func filterYAML(paths []string) []string {
	var out []string
	for _, p := range paths {
		if isYAML(p) {
			out = append(out, p)
		}
	}
	return out
}
