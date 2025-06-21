package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

// completionCmd generates shell completion scripts for your CLI tool.
var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion scripts",
    Long: `Generate shell completion scripts for your CLI tool.

To load completions:

Bash:

  $ source <(yourcli completion bash)

Zsh:

  $ source <(yourcli completion zsh)

Fish:

  $ yourcli completion fish | source

PowerShell:

  PS> yourcli completion powershell | Out-String | Invoke-Expression
`,
    Args: cobra.ExactValidArgs(1),
    ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
    RunE: func(cmd *cobra.Command, args []string) error {
        switch args[0] {
        case "bash":
            return cmd.Root().GenBashCompletion(os.Stdout)
        case "zsh":
            return cmd.Root().GenZshCompletion(os.Stdout)
        case "fish":
            return cmd.Root().GenFishCompletion(os.Stdout, true)
        case "powershell":
            return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
        default:
            return fmt.Errorf("unsupported shell type: %s", args[0])
        }
    },
}

func init() {
    rootCmd.AddCommand(completionCmd)
}