package cmd

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
    "os/user"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
)

func getShell() (string, error) {
    shellEnv := os.Getenv("SHELL")
    if shellEnv == "" {
        return "", errors.New("SHELL environment variable not set")
    }
    parts := strings.Split(shellEnv, "/")
    return parts[len(parts)-1], nil
}

func homeDir() string {
    if u, err := user.Current(); err == nil {
        return u.HomeDir
    }
    // fallback
    return os.Getenv("HOME")
}

func appendIfNotExists(file, content string) error {
    data, err := os.ReadFile(file)
    if err != nil && !errors.Is(err, fs.ErrNotExist) {
        return err
    }
    if strings.Contains(string(data), content) {
        return nil
    }
    f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()
    if _, err := f.WriteString("\n" + content + "\n"); err != nil {
        return err
    }
    return nil
}

var autoCompletionCmd = &cobra.Command{
    Use:   "auto",
    Short: "Setup shell completion automatically",
    Long: `Detect your shell and install completion script + config:

- Writes completion file to a user directory
- Adds sourcing/config lines to your shell rc file
- You may need to restart your shell after this`,
    RunE: func(cmd *cobra.Command, args []string) error {
        shell, err := getShell()
        if err != nil {
            return err
        }
        fmt.Printf("Detected shell: %s\n", shell)
        home := homeDir()

        switch shell {
        case "bash":
            dir := filepath.Join(home, ".bash_completion.d")
            if err := os.MkdirAll(dir, 0755); err != nil {
                return err
            }
            compFile := filepath.Join(dir, "ghp")
            f, err := os.Create(compFile)
            if err != nil {
                return err
            }
            defer f.Close()
            if err := cmd.Root().GenBashCompletion(f); err != nil {
                return err
            }
            bashrc := filepath.Join(home, ".bashrc")
            sourceLine := fmt.Sprintf("source %s", compFile)
            if err := appendIfNotExists(bashrc, sourceLine); err != nil {
                return err
            }
            fmt.Printf("Bash completion installed to %s and sourced in %s\n", compFile, bashrc)
            fmt.Println("Restart your shell or run `source ~/.bashrc`")

        case "zsh":
            compDir := filepath.Join(home, ".zsh", "completions")
            if err := os.MkdirAll(compDir, 0755); err != nil {
                return err
            }
            compFile := filepath.Join(compDir, "_ghp")
            f, err := os.Create(compFile)
            if err != nil {
                return err
            }
            defer f.Close()
            if err := cmd.Root().GenZshCompletion(f); err != nil {
                return err
            }
            zshrc := filepath.Join(home, ".zshrc")
            fpathLine := fmt.Sprintf("fpath+=(%s)", compDir)
            compinitLines := `autoload -Uz compinit
compinit`
            if err := appendIfNotExists(zshrc, fpathLine); err != nil {
                return err
            }
            if err := appendIfNotExists(zshrc, compinitLines); err != nil {
                return err
            }
            fmt.Printf("Zsh completion installed to %s and configured in %s\n", compFile, zshrc)
            fmt.Println("Restart your shell or run `source ~/.zshrc`")

        case "fish":
            compDir := filepath.Join(home, ".config", "fish", "completions")
            if err := os.MkdirAll(compDir, 0755); err != nil {
                return err
            }
            compFile := filepath.Join(compDir, "ghp.fish")
            f, err := os.Create(compFile)
            if err != nil {
                return err
            }
            defer f.Close()
            if err := cmd.Root().GenFishCompletion(f, true); err != nil {
                return err
            }
            fmt.Printf("Fish completion installed to %s\n", compFile)
            fmt.Println("Restart your shell or open a new shell session")

        case "powershell":
            compDir := filepath.Join(home, "Documents", "WindowsPowerShell", "Modules", "ghp")
            if err := os.MkdirAll(compDir, 0755); err != nil {
                return err
            }
            compFile := filepath.Join(compDir, "ghp.ps1")
            f, err := os.Create(compFile)
            if err != nil {
                return err
            }
            defer f.Close()
            if err := cmd.Root().GenPowerShellCompletionWithDesc(f); err != nil {
                return err
            }
            fmt.Printf("PowerShell completion installed to %s\n", compFile)
            fmt.Println("Add `Import-Module ghp` to your PowerShell profile to enable completion")

        default:
            return fmt.Errorf("unsupported shell: %s", shell)
        }

        return nil
    },
}

func init() {
	rootCmd.AddCommand(autoCompletionCmd)
}