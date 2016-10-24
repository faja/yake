_yake() {
  YAKEFILE="Yakefile"
  local opts
  case "${COMP_WORDS[COMP_CWORD-1]}" in
    *)
      opts=""
      if test -f $YAKEFILE
      then
        opts=$(awk -F: '/^\S/ {print $1}' $YAKEFILE)
      fi
      ;;
  esac
  COMPREPLY=($(compgen -W "${opts}" -- ${COMP_WORDS[COMP_CWORD]}))
  return 0
}

complete -F _yake yake
