---
Title: Dotfiles
StyleSheets:
  - /static/styles/utils/code.css
Scripts:
  - /static/scripts/utils/copy-code.js
---

# Dotfiles

## ~/.zshrc

```bash
export ZSH="$HOME/.oh-my-zsh"

ZSH_THEME="agnoster"

HYPHEN_INSENSITIVE="true"

HIST_STAMPS="dd/mm/yyyy"

plugins=(git)

source $ZSH/oh-my-zsh.sh

SCRIPTS_DIR="$HOME/.scripts" # Path to the scripts directory
# Check if the directory exists
if [ -d "$SCRIPTS_DIR" ]; then
    for script in "$SCRIPTS_DIR"/*; do
        # Check if it's a regular file and executable
        if [ -f "$script" ] && [ -x "$script" ]; then
            script_name=$(basename "$script")
            # Create an alias for each script
            alias "$script_name"="$script"
        fi
    done
fi

export NVM_DIR="$HOME/.nvm"
[ -s "$HOMEBREW_PREFIX/opt/nvm/nvm.sh" ] && \. "$HOMEBREW_PREFIX/opt/nvm/nvm.sh" # This loads nvm
[ -s "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm" ] && \. "$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm" # This loads nvm bash_completion

####
# ZSH function to auto-switch to correct Node version
# - https://gist.github.com/callumlocke/30990e247e52ab6ac1aa98e5f0e5bbf5
####
auto-switch-node-version() {
  NVMRC_PATH=$(nvm_find_nvmrc)
  CURRENT_NODE_VERSION=$(nvm version)

  if [[ ! -z "$NVMRC_PATH" ]]; then
    REQUESTED_NODE_VERSION=$(cat $NVMRC_PATH)
    MATCHED_NODE_VERSION=$(nvm_match_version $REQUESTED_NODE_VERSION)
    if [[ ! -z "$MATCHED_NODE_VERSION" && $MATCHED_NODE_VERSION != "N/A" ]]; then
      unset AUTOSWITCH_NODE_SUPPRESS_WARNING
      if [[ $CURRENT_NODE_VERSION != $MATCHED_NODE_VERSION ]]; then
        nvm use $REQUESTED_NODE_VERSION
      fi
    else
      if [[ $AUTOSWITCH_NODE_SUPPRESS_WARNING == $NVMRC_PATH ]]; then
        return
      fi
      echo "\nWARNING\n  [Found: $REQUESTED_NODE_VERSION] !=  [Current: $CURRENT_NODE_VERSION]\n\n  You might want to run \"nvm install\""
      export AUTOSWITCH_NODE_SUPPRESS_WARNING=$NVMRC_PATH
    fi
  else
    unset AUTOSWITCH_NODE_SUPPRESS_WARNING
    if [[ $CURRENT_NODE_VERSION != $(nvm version default)  ]]; then
      nvm use default
    fi
  fi
}
add-zsh-hook chpwd auto-switch-node-version
auto-switch-node-version

export PATH="$HOME/.cargo/bin:$PATH"
export PATH="$HOME/go/bin:$PATH"
export PATH="$HOME/bin:$PATH"
```

## ~/.vimrc

```bash
set encoding=utf-8
set updatetime=250
set t_Co=256
set number relativenumber

let g:dracula_colorterm = 0
let g:dracula_italic = 0

call plug#begin()

Plug 'vim-airline/vim-airline'
Plug 'tpope/vim-sensible'
Plug 'tpope/vim-fugitive'
Plug 'tpope/vim-commentary'
Plug 'tpope/vim-surround'
Plug 'ycm-core/YouCompleteMe'
Plug 'rust-lang/rust.vim'
Plug 'airblade/vim-gitgutter'
Plug 'ervandew/supertab'
Plug 'dracula/vim', { 'as': 'dracula' }

call plug#end()

colorscheme dracula

syntax enable
filetype plugin indent on

let g:rustfmt_autosave = 1
let g:ycm_rust_src_path = '/usr/local/rust/rustc-1.5.0/src'

let g:gitgutter_sign_added = '+'
let g:gitgutter_sign_modified = '>'
let g:gitgutter_sign_removed = '-'
let g:gitgutter_sign_removed_first_line = '^'
let g:gitgutter_sign_modified_removed = '<'
```

## ~/.gitconfig

```ini
[user]
	name = Nathan Doolan
	email = <my-email>

[filter "lfs"]
	clean = git-lfs clean -- %f
	smudge = git-lfs smudge -- %f
	process = git-lfs filter-process
	required = true
```

## ~/.config/zed/settings.json

```json
{
  "formatter": "language_server",
  "languages": {
    "TypeScript": {
      "formatter": "prettier",
      "code_actions_on_format": {
        "source.fixAll.eslint": true
      }
    },
    "TSX": {
      "formatter": "prettier",
      "code_actions_on_format": {
        "source.fixAll.eslint": true
      }
    }
  },
  "autosave": "on_focus_change",
  "theme": "Dracula",
  "tab_size": 2,
  "telemetry": {
    "diagnostics": false,
    "metrics": false
  },
  "vim_mode": false,
  "terminal": {
    "line_height": "comfortable",
    "font_family": "Noto Mono for Powerline"
  }
}
```
