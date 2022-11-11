# Base16 Terminal Sexy

This script reads in a JSON exported theme from [terminal.sexy](https://terminal.sexy)
and converts it automatically to respective sh and vim files to be used in conjunction with
[https://github.com/chriskempson/base16-vim](https://github.com/chriskempson/base16-vim).

With some [configuration](https://browntreelabs.com/base-16-shell-and-why-its-so-awsome/), you can sync your terminal and vim themes.
This script makes it easy to do it with any of my custom made themes! Note, the color alignment is not always perfect, so be prepared
fiddle with some of the locations for colors.

## How to run

1. Install Go
2. Run main.go
	`go run main.go --file <json exported theme>`
	  - optional: `--neovim-out <path to output for neovim file>`
	    DEFAULT: ~/.local/share/nvim/site/pack/packer/start/base16-vim/colors
	  - optional: `--terminal-out <path to output for terminal file>`
	    DEFAULT: ~/.config/base16-shell/scripts 
3. Done :) Restart your terminal for changes to pick up
