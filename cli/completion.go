package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Tab completion for `anon send @<Tab>`. The shell scripts call the hidden
// `anon __members` command, which answers from a short-lived cache so Tab
// stays fast and still works offline.

const membersCacheTTL = 5 * time.Minute

func membersCachePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "anon", "members.cache"), nil
}

func writeMembersCache(members []string) {
	path, err := membersCachePath()
	if err != nil {
		return
	}
	os.WriteFile(path, []byte(strings.Join(members, "\n")), 0o600)
}

// completeMembers prints one "@username" per line. Errors are silent because
// the output goes straight into the shell's completion menu.
func completeMembers() {
	path, err := membersCachePath()
	if err != nil {
		return
	}
	data, readErr := os.ReadFile(path)
	stat, statErr := os.Stat(path)

	if readErr != nil || statErr != nil || time.Since(stat.ModTime()) > membersCacheTTL {
		client.Timeout = 3 * time.Second
		if members, err := fetchMembers(); err == nil {
			data = []byte(strings.Join(members, "\n"))
		}
	}

	for _, m := range strings.Fields(string(data)) {
		fmt.Println("@" + m)
	}
}

func cmdCompletion(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: anon completion <zsh|bash|powershell>")
	}
	switch args[0] {
	case "zsh":
		fmt.Print(zshCompletion)
	case "bash":
		fmt.Print(bashCompletion)
	case "powershell", "pwsh":
		fmt.Print(powershellCompletion)
	default:
		return fmt.Errorf("unsupported shell %q (use zsh, bash or powershell)", args[0])
	}
	return nil
}

const zshCompletion = `_anon() {
  if (( CURRENT == 2 )); then
    compadd -- send members config completion help
  elif (( CURRENT == 3 )) && [[ ${words[2]} == send ]]; then
    local -a names
    names=(${(f)"$(anon __members 2>/dev/null)"})
    if [[ $PREFIX == @* ]]; then
      compadd -- $names
    else
      compadd -- ${names#@}
    fi
  fi
}
(( $+functions[compdef] )) || { autoload -Uz compinit && compinit }
compdef _anon anon
`

// bash treats "@" as a word break by default, so the current word is read from
// COMP_LINE and the "@" is stripped from replies when bash splits on it.
const bashCompletion = `_anon() {
  local line="${COMP_LINE:0:COMP_POINT}"
  local -a parts
  read -ra parts <<< "$line"
  [[ "$line" == *" " ]] && parts+=("")
  local n=${#parts[@]} cur="${parts[${#parts[@]}-1]}"

  if (( n == 2 )); then
    COMPREPLY=($(compgen -W "send members config completion help" -- "$cur"))
  elif (( n == 3 )) && [[ "${parts[1]}" == send ]]; then
    local names
    names="$(anon __members 2>/dev/null)"
    [[ "$cur" != @* ]] && names="${names//@/}"
    COMPREPLY=($(compgen -W "$names" -- "$cur"))
    if [[ "$cur" == @* && "$COMP_WORDBREAKS" == *@* ]]; then
      COMPREPLY=("${COMPREPLY[@]#@}")
    fi
  fi
}
complete -F _anon anon
`

// In PowerShell a bare @name is splatting, so names complete without "@"
// (anon accepts both), or quoted as '@name' when the user typed "@".
const powershellCompletion = `Register-ArgumentCompleter -Native -CommandName anon, anon.exe -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $elements = @($commandAst.CommandElements | ForEach-Object { $_.ToString() })
    $index = $elements.Count
    if ($wordToComplete) { $index-- }

    if ($index -eq 1) {
        'send', 'members', 'config', 'completion', 'help' |
            Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_) }
    }
    elseif ($index -eq 2 -and $elements[1] -eq 'send') {
        $typed = $wordToComplete.Trim("'", '"').TrimStart('@')
        $quoted = $wordToComplete -match '^[''"]?@'
        anon __members 2>$null | ForEach-Object { $_.TrimStart('@') } |
            Where-Object { $_ -like "$typed*" } |
            ForEach-Object {
                $text = if ($quoted) { "'@$_'" } else { $_ }
                [System.Management.Automation.CompletionResult]::new($text, "@$_", 'ParameterValue', "@$_")
            }
    }
}
`
