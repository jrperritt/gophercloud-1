#! /bin/bash

_cli_bash_autocomplete() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
	if [[ ${cur} == -* || ${prev} != -* || ${prev} == "--debug" || ${prev} == "--no-cache" ]]; then
		opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
		COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
	fi
	return 0
}

complete -o default -F _cli_bash_autocomplete stack
