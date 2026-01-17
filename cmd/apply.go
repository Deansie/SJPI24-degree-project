package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	applyFile       string
	applyDir        string
	applyNamespace  string
	applyKubeconfig string
	applyDryRun     bool
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Kubernetes manifests to the cluster",
	Long: appLogo + `
Applies validated YAML manifests to the Kubernetes cluster using client-go.
Supports kubeconfig for authentication and namespace/context overrides.

Examples:
  k8s-deploy apply --file manifest.yaml --namespace myapp
  k8s-deploy apply --dir manifests/ --dry-run`,
	RunE: runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVar(&applyFile, "file", "", "Path to YAML manifest file")
	applyCmd.Flags().StringVar(&applyDir, "dir", "", "Directory containing YAML manifests")
	applyCmd.Flags().StringVar(&applyNamespace, "namespace", "", "Namespace to apply to (overrides manifest)")
	applyCmd.Flags().StringVar(&applyKubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "Path to kubeconfig")
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "Simulate apply without changes")
	_ = applyCmd.MarkFlagRequired("file")
}

func runApply(cmd *cobra.Command, args []string) error {
	config, err := clientcmd.BuildConfigFromFlags("", applyKubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}
	var files []string
	if applyFile != "" {
		files = append(files, applyFile)
	} else if applyDir != "" {
		err = filepath.Walk(applyDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Ext(path) != ".yaml" {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to walk dir: %w", err)
		}
	} else {
		return fmt.Errorf("must specify --file or --dir")
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}
		dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		obj := &unstructured.Unstructured{}
		_, gvk, err := dec.Decode(data, nil, obj)
		if err != nil {
			return fmt.Errorf("failed to decode %s: %w", file, err)
		}
		if applyNamespace != "" {
			obj.SetNamespace(applyNamespace)
		}
		resourceName := strings.ToLower(gvk.Kind) + "s"
		gvr := schema.GroupVersionResource{
			Group:    gvk.Group,
			Version:  gvk.Version,
			Resource: resourceName,
		}
		resource := dynClient.Resource(gvr).Namespace(obj.GetNamespace())
		if applyDryRun {
			fmt.Printf("Dry-run: Would apply %s %s in %s\n", obj.GetKind(), obj.GetName(), obj.GetNamespace())
			continue
		}
		_, err = resource.Apply(context.Background(), obj.GetName(), obj, metav1.ApplyOptions{FieldManager: "k8s-deploy"})
		if err != nil {
			return fmt.Errorf("failed to apply %s: %w", file, err)
		}
		fmt.Printf("Applied %s %s in %s\n", obj.GetKind(), obj.GetName(), obj.GetNamespace())
	}
	return nil
}
