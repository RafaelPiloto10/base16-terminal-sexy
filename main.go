package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Base16JSON struct {
	Name       string
	Author     string
	Color      []string
	Foreground string
	Background string
}

var fileName = flag.String("file", "", "read base16 JSON file exported from https://terminal.sexy")
var base16NeoVimDir = flag.String("neovim-out", ".local/share/nvim/site/pack/packer/start/base16-vim/colors", "neovim base16 output folder")
var base16TerminalDir = flag.String("terminal-out", ".config/base16-shell/scripts", "terminal base16 output folder")

func main() {

	flag.Parse()

	if *fileName == "" {
		panic("invalid base16 file")
	} else if !strings.HasSuffix(*fileName, ".json") {
		panic(fmt.Sprintf("invalid base16 file %s; expecting json file", *fileName))
	}

	buf, err := os.ReadFile(*fileName)
	if err != nil {
		panic(err.Error())
	}

	colorscheme := Base16JSON{}
	err = json.Unmarshal(buf, &colorscheme)
	if err != nil {
		panic(err.Error())
	}

	tokens := strings.Split(*fileName, "/")
	name := ""

	if len(tokens) < 1 {
		name = strings.Split((*fileName), ".json")[0]
	} else {
		*fileName = tokens[len(tokens)-1]
		name = strings.Split((*fileName), ".json")[0]
	}

	neovimScheme := generateNeovimScheme(colorscheme)
	terminalScheme := generateTerminalScheme(colorscheme)

	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	
	vimLoc := fmt.Sprintf("%s/%s/base16-%s.vim", home, *base16NeoVimDir, name)
	err = os.WriteFile(vimLoc, neovimScheme, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write; %v", err))
	}

	fmt.Printf("wrote vim scheme to %s\n", vimLoc)

	terminalLoc := fmt.Sprintf("%s/%s/base16-%s.sh", home, *base16TerminalDir, name)	
	err = os.WriteFile(terminalLoc, terminalScheme, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write; %v", err))
	}
	fmt.Printf("wrote terminal scheme to %s\n", terminalLoc)
}

func generateTerminalScheme(bj Base16JSON) []byte {
	template := "" +
		"#!/bin/sh\n" +
		fmt.Sprintln(fmt.Sprintf(`color00="%s" # Base 00 - Black`, formatColorSlash(bj.Background))) +
		fmt.Sprintln(fmt.Sprintf(`color01="%s" # Base 08 - Red`, formatColorSlash(bj.Color[1]))) +
		fmt.Sprintln(fmt.Sprintf(`color02="%s" # Base 0B - Green`, formatColorSlash(bj.Color[2]))) +
		fmt.Sprintln(fmt.Sprintf(`color03="%s" # Base 0A - Yellow`, formatColorSlash(bj.Color[3]))) +
		fmt.Sprintln(fmt.Sprintf(`color04="%s" # Base 0D - Blue`, formatColorSlash(bj.Color[4]))) +
		fmt.Sprintln(fmt.Sprintf(`color05="%s" # Base 0E - Magenta`, formatColorSlash(bj.Color[5]))) +
		fmt.Sprintln(fmt.Sprintf(`color06="%s" # Base 0C - Cyan`, formatColorSlash(bj.Color[6]))) +
		fmt.Sprintln(fmt.Sprintf(`color07="%s" # Base 05 - White`, formatColorSlash(bj.Color[7]))) +
		fmt.Sprintln(fmt.Sprintf(`color08="%s" # Base 03 - Bright Black`, formatColorSlash(bj.Color[8]))) +
		"color09=$color01 # Base 08 - Bright Red\n" +
		"color10=$color02 # Base 0B - Bright Green\n" +
		"color11=$color03 # Base 0A - Bright Yellow\n" +
		"color12=$color04 # Base 0D - Bright Blue\n" +
		"color13=$color05 # Base 0E - Bright Magenta\n" +
		"color14=$color06 # Base 0C - Bright Cyan\n" +
		fmt.Sprintln(fmt.Sprintf(`color15="%s" # Base 07 - Bright White`, formatColorSlash(bj.Foreground))) +
		fmt.Sprintln(fmt.Sprintf(`color16="%s" # Base 09`, formatColorSlash(bj.Color[4]))) +
		fmt.Sprintln(fmt.Sprintf(`color17="%s" # Base 0F`, formatColorSlash(bj.Color[2]))) +
		fmt.Sprintln(fmt.Sprintf(`color18="%s" # Base 01`, formatColorSlash(bj.Color[5]))) +
		fmt.Sprintln(fmt.Sprintf(`color19="%s" # Base 02`, formatColorSlash(bj.Background))) +
		fmt.Sprintln(fmt.Sprintf(`color20="%s" # Base 04`, formatColorSlash(bj.Color[7]))) +
		fmt.Sprintln(fmt.Sprintf(`color21="%s" # Base 06`, formatColorSlash(bj.Foreground))) +
		fmt.Sprintln(fmt.Sprintf(`color_foreground="%s" # Base 05`, formatColorSlash(bj.Foreground))) +
		fmt.Sprintln(fmt.Sprintf(`color_background="%s" # Base 00`, formatColorSlash(bj.Background))) +
		fmt.Sprintln(`if [ -n "$TMUX" ]; then`) +
		"\t# Tell tmux to pass the escape sequences through\n" +
		"\t# (Source: http://permalink.gmane.org/gmane.comp.terminal-emulators.tmux.user/1324)\n" +
		"\tput_template() { printf '\\033Ptmux;\\033\\033]4;%d;rgb:%s\\033\\033\\\\\\033\\\\' $@; }\n" +
		"\tput_template_var() { printf '\\033Ptmux;\\033\\033]%d;rgb:%s\\033\\033\\\\\\033\\\\' $@; }\n" +
		"\tput_template_custom() { printf '\\033Ptmux;\\033\\033]%s%s\\033\\033\\\\\\033\\\\' $@; }\n" +
		fmt.Sprintln(`elif [ "${TERM%%[-.]*}" = 'screen' ]; then`) +
		"\t# GNU screen (screen, screen-256color, screen-256color-bce)\n" +
		"\tput_template() { printf '\\033P\\033]4;%d;rgb:%s\\007\\033\\\\' $@; }\n" +
		"\tput_template_var() { printf '\\033P\\033]%d;rgb:%s\\007\\033\\\\' $@; }\n" +
		"\tput_template_custom() { printf '\\033P\\033]%s%s\\007\\033\\\\' $@; }\n" +
		fmt.Sprintln(`elif [ "${TERM%%-*}" = 'linux' ]; then`) +
		"\tput_template() { [ $1 -lt 16 ] && printf '\\e]P%x%s' $1 $(echo $2 | sed 's/\\///g'); }\n" +
		"\tput_template_var() { true; }\n" +
		"\tput_template_custom() { true; }\n" +
		"else\n" +
		"\tput_template() { printf '\\033]4;%d;rgb:%s\\033\\\\' $@; }\n" +
		"\tput_template_var() { printf '\\033]%d;rgb:%s\\033\\\\' $@; }\n" +
		"\tput_template_custom() { printf '\\033]%s%s\\033\\\\' $@; }\n" +
		"fi\n" +
		"# 16 color space\n" +
		"put_template 0  $color00\n" +
		"put_template 1  $color01\n" +
		"put_template 2  $color02\n" +
		"put_template 3  $color03\n" +
		"put_template 4  $color04\n" +
		"put_template 5  $color05\n" +
		"put_template 6  $color06\n" +
		"put_template 7  $color07\n" +
		"put_template 8  $color08\n" +
		"put_template 9  $color09\n" +
		"put_template 10 $color10\n" +
		"put_template 11 $color11\n" +
		"put_template 12 $color12\n" +
		"put_template 13 $color13\n" +
		"put_template 14 $color14\n" +
		"put_template 15 $color15\n" +
		"# 256 color space\n" +
		"put_template 16 $color16\n" +
		"put_template 17 $color17\n" +
		"put_template 18 $color18\n" +
		"put_template 19 $color19\n" +
		"put_template 20 $color20\n" +
		"put_template 21 $color21\n" +
		"# foreground / background / cursor color\n" +
		fmt.Sprintln(`if [ -n "$ITERM_SESSION_ID" ]; then`) +
		"\t# iTerm2 proprietary escape codes\n" +
		"\tput_template_custom Pg f8f8f2 # foreground\n" +
		"\tput_template_custom Ph 272822 # background\n" +
		"\tput_template_custom Pi f8f8f2 # bold color\n" +
		"\tput_template_custom Pj 49483e # selection color\n" +
		"\tput_template_custom Pk f8f8f2 # selected text color\n" +
		"\tput_template_custom Pl f8f8f2 # cursor\n" +
		"\tput_template_custom Pm 272822 # cursor text\n" +
		"else\n" +
		"\tput_template_var 10 $color_foreground\n" +
		fmt.Sprintln(`	if [ "$BASE16_SHELL_SET_BACKGROUND" != false ]; then`) +
		"\t\tput_template_var 11 $color_background\n" +
		fmt.Sprintln(`		if [ "${TERM%%-*}" = "rxvt" ]; then`) +
		"\t\t\tput_template_var 708 $color_background # internal border (rxvt)\n" +
		"\t\tfi\n" +
		"\tfi\n" +
		fmt.Sprintln(`	put_template_custom 12 ";7" # cursor (reverse video)`) +
		"fi\n" +
		"# clean up\n" +
		"unset -f put_template\n" +
		"unset -f put_template_var\n" +
		"unset -f put_template_custom\n" +
		"unset color00\n" +
		"unset color01\n" +
		"unset color02\n" +
		"unset color03\n" +
		"unset color04\n" +
		"unset color05\n" +
		"unset color06\n" +
		"unset color07\n" +
		"unset color08\n" +
		"unset color09\n" +
		"unset color10\n" +
		"unset color11\n" +
		"unset color12\n" +
		"unset color13\n" +
		"unset color14\n" +
		"unset color15\n" +
		"unset color16\n" +
		"unset color17\n" +
		"unset color18\n" +
		"unset color19\n" +
		"unset color20\n" +
		"unset color21\n" +
		"unset color_foreground\n" +
		"unset color_background\n"

	return []byte(template)
}

func formatColorSlash(hex string) string {
	color := []string{}
	i := 0
	for _, c := range strings.Split(hex, "") {
		if c != "#" {
			if i%2 == 0 {
				color = append(color, "/")
			}
			color = append(color, c)
			i += 1
		}
	}

	return strings.Join(color, "")[1:]
}

func formatClean(hex string) string {
	return hex[1:]
}

func generateNeovimScheme(bj Base16JSON) []byte {
	template := `
" vi:syntax=vim
if !has("gui_running")
  if exists("g:base16_shell_path")
    execute "silent !/bin/sh ".g:base16_shell_path."/base16-polybar.sh"
  endif
endif

" GUI color definitions
let s:gui00        = "` + fmt.Sprintf("%s", formatClean(bj.Color[0])) + `"
let g:base16_gui00 = "` + fmt.Sprintf("%s", formatClean(bj.Color[0])) + `"
let s:gui01        = "` + fmt.Sprintf("%s", formatClean(bj.Color[12])) + `"
let g:base16_gui01 = "` + fmt.Sprintf("%s", formatClean(bj.Color[12])) + `"
let s:gui02        = "` + fmt.Sprintf("%s", formatClean(bj.Color[13])) + `"
let g:base16_gui02 = "` + fmt.Sprintf("%s", formatClean(bj.Color[13])) + `"
let s:gui03        = "` + fmt.Sprintf("%s", formatClean(bj.Color[8])) + `"
let g:base16_gui03 = "` + fmt.Sprintf("%s", formatClean(bj.Color[8])) + `"
let s:gui04        = "` + fmt.Sprintf("%s", formatClean(bj.Color[14])) + `"
let g:base16_gui04 = "` + fmt.Sprintf("%s", formatClean(bj.Color[14])) + `"
let s:gui05        = "` + fmt.Sprintf("%s", formatClean(bj.Color[7])) + `"
let g:base16_gui05 = "` + fmt.Sprintf("%s", formatClean(bj.Color[7])) + `"
let s:gui06        = "` + fmt.Sprintf("%s", formatClean(bj.Color[15])) + `"
let g:base16_gui06 = "` + fmt.Sprintf("%s", formatClean(bj.Color[15])) + `"
let s:gui07        = "` + fmt.Sprintf("%s", formatClean(bj.Color[9])) + `"
let g:base16_gui07 = "` + fmt.Sprintf("%s", formatClean(bj.Color[9])) + `"
let s:gui08        = "` + fmt.Sprintf("%s", formatClean(bj.Color[1])) + `"
let g:base16_gui08 = "` + fmt.Sprintf("%s", formatClean(bj.Color[1])) + `"
let s:gui09        = "` + fmt.Sprintf("%s", formatClean(bj.Color[10])) + `"
let g:base16_gui09 = "` + fmt.Sprintf("%s", formatClean(bj.Color[10])) + `"
let s:gui0A        = "` + fmt.Sprintf("%s", formatClean(bj.Color[3])) + `"
let g:base16_gui0A = "` + fmt.Sprintf("%s", formatClean(bj.Color[3])) + `"
let s:gui0B        = "` + fmt.Sprintf("%s", formatClean(bj.Color[2])) + `"
let g:base16_gui0B = "` + fmt.Sprintf("%s", formatClean(bj.Color[2])) + `"
let s:gui0C        = "` + fmt.Sprintf("%s", formatClean(bj.Color[6])) + `"
let g:base16_gui0C = "` + fmt.Sprintf("%s", formatClean(bj.Color[6])) + `"
let s:gui0D        = "` + fmt.Sprintf("%s", formatClean(bj.Color[4])) + `"
let g:base16_gui0D = "` + fmt.Sprintf("%s", formatClean(bj.Color[4])) + `"
let s:gui0E        = "` + fmt.Sprintf("%s", formatClean(bj.Color[5])) + `"
let g:base16_gui0E = "` + fmt.Sprintf("%s", formatClean(bj.Color[5])) + `"
let s:gui0F        = "` + fmt.Sprintf("%s", formatClean(bj.Color[11])) + `"
let g:base16_gui0F = "` + fmt.Sprintf("%s", formatClean(bj.Color[11])) + `"

" Terminal color definitions
let s:cterm00        = "00"
let g:base16_cterm00 = "00"
let s:cterm03        = "08"
let g:base16_cterm03 = "08"
let s:cterm05        = "07"
let g:base16_cterm05 = "07"
let s:cterm07        = "15"
let g:base16_cterm07 = "15"
let s:cterm08        = "01"
let g:base16_cterm08 = "01"
let s:cterm0A        = "03"
let g:base16_cterm0A = "03"
let s:cterm0B        = "02"
let g:base16_cterm0B = "02"
let s:cterm0C        = "06"
let g:base16_cterm0C = "06"
let s:cterm0D        = "04"
let g:base16_cterm0D = "04"
let s:cterm0E        = "05"
let g:base16_cterm0E = "05"
if exists("base16colorspace") && base16colorspace == "256"
  let s:cterm01        = "18"
  let g:base16_cterm01 = "18"
  let s:cterm02        = "19"
  let g:base16_cterm02 = "19"
  let s:cterm04        = "20"
  let g:base16_cterm04 = "20"
  let s:cterm06        = "21"
  let g:base16_cterm06 = "21"
  let s:cterm09        = "16"
  let g:base16_cterm09 = "16"
  let s:cterm0F        = "17"
  let g:base16_cterm0F = "17"
else
  let s:cterm01        = "10"
  let g:base16_cterm01 = "10"
  let s:cterm02        = "11"
  let g:base16_cterm02 = "11"
  let s:cterm04        = "12"
  let g:base16_cterm04 = "12"
  let s:cterm06        = "13"
  let g:base16_cterm06 = "13"
  let s:cterm09        = "09"
  let g:base16_cterm09 = "09"
  let s:cterm0F        = "14"
  let g:base16_cterm0F = "14"
endif

" Neovim terminal colours
if has("nvim")
  let g:terminal_color_0 =  "` + fmt.Sprintf("%s", bj.Color[0]) + `"
  let g:terminal_color_1 =  "` + fmt.Sprintf("%s", bj.Color[1]) + `"
  let g:terminal_color_2 =  "` + fmt.Sprintf("%s", bj.Color[2]) + `"
  let g:terminal_color_3 =  "` + fmt.Sprintf("%s", bj.Color[3]) + `"
  let g:terminal_color_4 =  "` + fmt.Sprintf("%s", bj.Color[4]) + `"
  let g:terminal_color_5 =  "` + fmt.Sprintf("%s", bj.Color[5]) + `"
  let g:terminal_color_6 =  "` + fmt.Sprintf("%s", bj.Color[6]) + `"
  let g:terminal_color_7 =  "` + fmt.Sprintf("%s", bj.Color[7]) + `"
  let g:terminal_color_8 =  "` + fmt.Sprintf("%s", bj.Color[8]) + `"
  let g:terminal_color_9 =  "` + fmt.Sprintf("%s", bj.Color[1]) + `"
  let g:terminal_color_10 = "` + fmt.Sprintf("%s", bj.Color[2]) + `"
  let g:terminal_color_11 = "` + fmt.Sprintf("%s", bj.Color[3]) + `"
  let g:terminal_color_12 = "` + fmt.Sprintf("%s", bj.Color[4]) + `"
  let g:terminal_color_13 = "` + fmt.Sprintf("%s", bj.Color[5]) + `"
  let g:terminal_color_14 = "` + fmt.Sprintf("%s", bj.Color[6]) + `"
  let g:terminal_color_15 = "` + fmt.Sprintf("%s", bj.Color[15]) + `"
  let g:terminal_color_background = g:terminal_color_0
  let g:terminal_color_foreground = g:terminal_color_5
  if &background == "light"
    let g:terminal_color_background = g:terminal_color_7
    let g:terminal_color_foreground = g:terminal_color_2
  endif
elseif has("terminal")
  let g:terminal_ansi_colors = [
        \ "` + fmt.Sprintf("%s", bj.Color[0]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[1]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[2]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[3]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[4]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[5]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[6]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[7]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[8]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[1]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[2]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[3]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[4]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[5]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[6]) + `",
        \ "` + fmt.Sprintf("%s", bj.Color[15]) + `",
        \ ]
endif

" Theme setup
hi clear
syntax reset
let g:colors_name = "base16-polybar"

" Highlighting function
" Optional variables are attributes and guisp
function! g:Base16hi(group, guifg, guibg, ctermfg, ctermbg, ...)
  let l:attr = get(a:, 1, "")
  let l:guisp = get(a:, 2, "")

  if a:guifg != ""
    exec "hi " . a:group . " guifg=#" . a:guifg
  endif
  if a:guibg != ""
    exec "hi " . a:group . " guibg=#" . a:guibg
  endif
  if a:ctermfg != ""
    exec "hi " . a:group . " ctermfg=" . a:ctermfg
  endif
  if a:ctermbg != ""
    exec "hi " . a:group . " ctermbg=" . a:ctermbg
  endif
  if l:attr != ""
    exec "hi " . a:group . " gui=" . l:attr . " cterm=" . l:attr
  endif
  if l:guisp != ""
    exec "hi " . a:group . " guisp=#" . l:guisp
  endif
endfunction


fun <sid>hi(group, guifg, guibg, ctermfg, ctermbg, attr, guisp)
  call g:Base16hi(a:group, a:guifg, a:guibg, a:ctermfg, a:ctermbg, a:attr, a:guisp)
endfun

" Vim editor colors
call <sid>hi("Normal",        s:gui05, s:gui00, s:cterm05, s:cterm00, "", "")
call <sid>hi("Bold",          "", "", "", "", "bold", "")
call <sid>hi("Debug",         s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("Directory",     s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("Error",         s:gui00, s:gui08, s:cterm00, s:cterm08, "", "")
call <sid>hi("ErrorMsg",      s:gui08, s:gui00, s:cterm08, s:cterm00, "", "")
call <sid>hi("Exception",     s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("FoldColumn",    s:gui0C, s:gui01, s:cterm0C, s:cterm01, "", "")
call <sid>hi("Folded",        s:gui03, s:gui01, s:cterm03, s:cterm01, "", "")
call <sid>hi("IncSearch",     s:gui01, s:gui09, s:cterm01, s:cterm09, "none", "")
call <sid>hi("Italic",        "", "", "", "", "none", "")
call <sid>hi("Macro",         s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("MatchParen",    "", s:gui03, "", s:cterm03,  "", "")
call <sid>hi("ModeMsg",       s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("MoreMsg",       s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("Question",      s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("Search",        s:gui01, s:gui0A, s:cterm01, s:cterm0A,  "", "")
call <sid>hi("Substitute",    s:gui01, s:gui0A, s:cterm01, s:cterm0A, "none", "")
call <sid>hi("SpecialKey",    s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("TooLong",       s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("Underlined",    s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("Visual",        "", s:gui02, "", s:cterm02, "", "")
call <sid>hi("VisualNOS",     s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("WarningMsg",    s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("WildMenu",      s:gui08, s:gui0A, s:cterm08, "", "", "")
call <sid>hi("Title",         s:gui0D, "", s:cterm0D, "", "none", "")
call <sid>hi("Conceal",       s:gui0D, s:gui00, s:cterm0D, s:cterm00, "", "")
call <sid>hi("Cursor",        s:gui00, s:gui05, s:cterm00, s:cterm05, "", "")
call <sid>hi("NonText",       s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("LineNr",        s:gui03, s:gui01, s:cterm03, s:cterm01, "", "")
call <sid>hi("SignColumn",    s:gui03, s:gui01, s:cterm03, s:cterm01, "", "")
call <sid>hi("StatusLine",    s:gui04, s:gui02, s:cterm04, s:cterm02, "none", "")
call <sid>hi("StatusLineNC",  s:gui03, s:gui01, s:cterm03, s:cterm01, "none", "")
call <sid>hi("VertSplit",     s:gui02, s:gui02, s:cterm02, s:cterm02, "none", "")
call <sid>hi("ColorColumn",   "", s:gui01, "", s:cterm01, "none", "")
call <sid>hi("CursorColumn",  "", s:gui01, "", s:cterm01, "none", "")
call <sid>hi("CursorLine",    "", s:gui01, "", s:cterm01, "none", "")
call <sid>hi("CursorLineNr",  s:gui04, s:gui01, s:cterm04, s:cterm01, "", "")
call <sid>hi("QuickFixLine",  "", s:gui01, "", s:cterm01, "none", "")
call <sid>hi("PMenu",         s:gui05, s:gui01, s:cterm05, s:cterm01, "none", "")
call <sid>hi("PMenuSel",      s:gui01, s:gui05, s:cterm01, s:cterm05, "", "")
call <sid>hi("TabLine",       s:gui03, s:gui01, s:cterm03, s:cterm01, "none", "")
call <sid>hi("TabLineFill",   s:gui03, s:gui01, s:cterm03, s:cterm01, "none", "")
call <sid>hi("TabLineSel",    s:gui0B, s:gui01, s:cterm0B, s:cterm01, "none", "")

" Standard syntax highlighting
call <sid>hi("Boolean",      s:gui09, "", s:cterm09, "", "", "")
call <sid>hi("Character",    s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("Comment",      s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("Conditional",  s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("Constant",     s:gui09, "", s:cterm09, "", "", "")
call <sid>hi("Define",       s:gui0E, "", s:cterm0E, "", "none", "")
call <sid>hi("Delimiter",    s:gui0F, "", s:cterm0F, "", "", "")
call <sid>hi("Float",        s:gui09, "", s:cterm09, "", "", "")
call <sid>hi("Function",     s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("Identifier",   s:gui08, "", s:cterm08, "", "none", "")
call <sid>hi("Include",      s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("Keyword",      s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("Label",        s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("Number",       s:gui09, "", s:cterm09, "", "", "")
call <sid>hi("Operator",     s:gui05, "", s:cterm05, "", "none", "")
call <sid>hi("PreProc",      s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("Repeat",       s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("Special",      s:gui0C, "", s:cterm0C, "", "", "")
call <sid>hi("SpecialChar",  s:gui0F, "", s:cterm0F, "", "", "")
call <sid>hi("Statement",    s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("StorageClass", s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("String",       s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("Structure",    s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("Tag",          s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("Todo",         s:gui0A, s:gui01, s:cterm0A, s:cterm01, "", "")
call <sid>hi("Type",         s:gui0A, "", s:cterm0A, "", "none", "")
call <sid>hi("Typedef",      s:gui0A, "", s:cterm0A, "", "", "")

" C highlighting
call <sid>hi("cOperator",   s:gui0C, "", s:cterm0C, "", "", "")
call <sid>hi("cPreCondit",  s:gui0E, "", s:cterm0E, "", "", "")

" C# highlighting
call <sid>hi("csClass",                 s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("csAttribute",             s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("csModifier",              s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("csType",                  s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("csUnspecifiedStatement",  s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("csContextualStatement",   s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("csNewDecleration",        s:gui08, "", s:cterm08, "", "", "")

" CSS highlighting
call <sid>hi("cssBraces",      s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("cssClassName",   s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("cssColor",       s:gui0C, "", s:cterm0C, "", "", "")

" Diff highlighting
call <sid>hi("DiffAdd",      s:gui0B, s:gui01,  s:cterm0B, s:cterm01, "", "")
call <sid>hi("DiffChange",   s:gui03, s:gui01,  s:cterm03, s:cterm01, "", "")
call <sid>hi("DiffDelete",   s:gui08, s:gui01,  s:cterm08, s:cterm01, "", "")
call <sid>hi("DiffText",     s:gui0D, s:gui01,  s:cterm0D, s:cterm01, "", "")
call <sid>hi("DiffAdded",    s:gui0B, s:gui00,  s:cterm0B, s:cterm00, "", "")
call <sid>hi("DiffFile",     s:gui08, s:gui00,  s:cterm08, s:cterm00, "", "")
call <sid>hi("DiffNewFile",  s:gui0B, s:gui00,  s:cterm0B, s:cterm00, "", "")
call <sid>hi("DiffLine",     s:gui0D, s:gui00,  s:cterm0D, s:cterm00, "", "")
call <sid>hi("DiffRemoved",  s:gui08, s:gui00,  s:cterm08, s:cterm00, "", "")

" Git highlighting
call <sid>hi("gitcommitOverflow",       s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("gitcommitSummary",        s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("gitcommitComment",        s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("gitcommitUntracked",      s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("gitcommitDiscarded",      s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("gitcommitSelected",       s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("gitcommitHeader",         s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("gitcommitSelectedType",   s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("gitcommitUnmergedType",   s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("gitcommitDiscardedType",  s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("gitcommitBranch",         s:gui09, "", s:cterm09, "", "bold", "")
call <sid>hi("gitcommitUntrackedFile",  s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("gitcommitUnmergedFile",   s:gui08, "", s:cterm08, "", "bold", "")
call <sid>hi("gitcommitDiscardedFile",  s:gui08, "", s:cterm08, "", "bold", "")
call <sid>hi("gitcommitSelectedFile",   s:gui0B, "", s:cterm0B, "", "bold", "")

" GitGutter highlighting
call <sid>hi("GitGutterAdd",     s:gui0B, s:gui01, s:cterm0B, s:cterm01, "", "")
call <sid>hi("GitGutterChange",  s:gui0D, s:gui01, s:cterm0D, s:cterm01, "", "")
call <sid>hi("GitGutterDelete",  s:gui08, s:gui01, s:cterm08, s:cterm01, "", "")
call <sid>hi("GitGutterChangeDelete",  s:gui0E, s:gui01, s:cterm0E, s:cterm01, "", "")

" HTML highlighting
call <sid>hi("htmlBold",    s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("htmlItalic",  s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("htmlEndTag",  s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("htmlTag",     s:gui05, "", s:cterm05, "", "", "")

" JavaScript highlighting
call <sid>hi("javaScript",        s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("javaScriptBraces",  s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("javaScriptNumber",  s:gui09, "", s:cterm09, "", "", "")
" pangloss/vim-javascript highlighting
call <sid>hi("jsOperator",          s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("jsStatement",         s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("jsReturn",            s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("jsThis",              s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("jsClassDefinition",   s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("jsFunction",          s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("jsFuncName",          s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("jsFuncCall",          s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("jsClassFuncName",     s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("jsClassMethodType",   s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("jsRegexpString",      s:gui0C, "", s:cterm0C, "", "", "")
call <sid>hi("jsGlobalObjects",     s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("jsGlobalNodeObjects", s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("jsExceptions",        s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("jsBuiltins",          s:gui0A, "", s:cterm0A, "", "", "")

" Mail highlighting
call <sid>hi("mailQuoted1",  s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("mailQuoted2",  s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("mailQuoted3",  s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("mailQuoted4",  s:gui0C, "", s:cterm0C, "", "", "")
call <sid>hi("mailQuoted5",  s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("mailQuoted6",  s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("mailURL",      s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("mailEmail",    s:gui0D, "", s:cterm0D, "", "", "")

" Markdown highlighting
call <sid>hi("markdownCode",              s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("markdownError",             s:gui05, s:gui00, s:cterm05, s:cterm00, "", "")
call <sid>hi("markdownCodeBlock",         s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("markdownHeadingDelimiter",  s:gui0D, "", s:cterm0D, "", "", "")

" NERDTree highlighting
call <sid>hi("NERDTreeDirSlash",  s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("NERDTreeExecFile",  s:gui05, "", s:cterm05, "", "", "")

" PHP highlighting
call <sid>hi("phpMemberSelector",  s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("phpComparison",      s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("phpParent",          s:gui05, "", s:cterm05, "", "", "")
call <sid>hi("phpMethodsVar",      s:gui0C, "", s:cterm0C, "", "", "")

" Python highlighting
call <sid>hi("pythonOperator",  s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("pythonRepeat",    s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("pythonInclude",   s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("pythonStatement", s:gui0E, "", s:cterm0E, "", "", "")

" Ruby highlighting
call <sid>hi("rubyAttribute",               s:gui0D, "", s:cterm0D, "", "", "")
call <sid>hi("rubyConstant",                s:gui0A, "", s:cterm0A, "", "", "")
call <sid>hi("rubyInterpolationDelimiter",  s:gui0F, "", s:cterm0F, "", "", "")
call <sid>hi("rubyRegexp",                  s:gui0C, "", s:cterm0C, "", "", "")
call <sid>hi("rubySymbol",                  s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("rubyStringDelimiter",         s:gui0B, "", s:cterm0B, "", "", "")

" SASS highlighting
call <sid>hi("sassidChar",     s:gui08, "", s:cterm08, "", "", "")
call <sid>hi("sassClassChar",  s:gui09, "", s:cterm09, "", "", "")
call <sid>hi("sassInclude",    s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("sassMixing",     s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("sassMixinName",  s:gui0D, "", s:cterm0D, "", "", "")

" Signify highlighting
call <sid>hi("SignifySignAdd",     s:gui0B, s:gui01, s:cterm0B, s:cterm01, "", "")
call <sid>hi("SignifySignChange",  s:gui0D, s:gui01, s:cterm0D, s:cterm01, "", "")
call <sid>hi("SignifySignDelete",  s:gui08, s:gui01, s:cterm08, s:cterm01, "", "")

" Spelling highlighting
call <sid>hi("SpellBad",     "", "", "", "", "undercurl", s:gui08)
call <sid>hi("SpellLocal",   "", "", "", "", "undercurl", s:gui0C)
call <sid>hi("SpellCap",     "", "", "", "", "undercurl", s:gui0D)
call <sid>hi("SpellRare",    "", "", "", "", "undercurl", s:gui0E)

" Startify highlighting
call <sid>hi("StartifyBracket",  s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("StartifyFile",     s:gui07, "", s:cterm07, "", "", "")
call <sid>hi("StartifyFooter",   s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("StartifyHeader",   s:gui0B, "", s:cterm0B, "", "", "")
call <sid>hi("StartifyNumber",   s:gui09, "", s:cterm09, "", "", "")
call <sid>hi("StartifyPath",     s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("StartifySection",  s:gui0E, "", s:cterm0E, "", "", "")
call <sid>hi("StartifySelect",   s:gui0C, "", s:cterm0C, "", "", "")
call <sid>hi("StartifySlash",    s:gui03, "", s:cterm03, "", "", "")
call <sid>hi("StartifySpecial",  s:gui03, "", s:cterm03, "", "", "")

" Java highlighting
call <sid>hi("javaOperator",     s:gui0D, "", s:cterm0D, "", "", "")

" Remove functions
delf <sid>hi

" Remove color variables
unlet s:gui00 s:gui01 s:gui02 s:gui03  s:gui04  s:gui05  s:gui06  s:gui07  s:gui08  s:gui09 s:gui0A  s:gui0B  s:gui0C  s:gui0D  s:gui0E  s:gui0F
unlet s:cterm00 s:cterm01 s:cterm02 s:cterm03 s:cterm04 s:cterm05 s:cterm06 s:cterm07 s:cterm08 s:cterm09 s:cterm0A s:cterm0B s:cterm0C s:cterm0D s:cterm0E s:cterm0F
	`
	return []byte(template)
}
