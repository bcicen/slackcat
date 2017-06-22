" vim-slackcat.vim
" Ridiculously simple plugin to send a visual selection to an Slack channel
"
" Copyright © 2016 Paco Esteban <paco@onna.be>

" Permission is hereby granted, free of charge, to any person obtaining
" a copy of this software and associated documentation files (the 'Software'),
" to deal in the Software without restriction, including without limitation
" the rights to use, copy, modify, merge, publish, distribute, sublicense,
" and/or sell copies of the Software, and to permit persons to whom the
" Software is furnished to do so, subject to the following conditions:

" The above copyright notice and this permission notice shall be included
" in all copies or substantial portions of the Software.

" THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
" EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
" OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
" IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
" DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
" TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
" OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

" it accepts g:slackcat_default_channel
" slackcat (http://slackcat.chat/) must be configured beforehand.

if !exists("g:slackcat_default_channel")
    let g:slackcat_default_channel = ""
endif

" send selection to slack
vnoremap <Leader>s :<C-u>call SendToSlack()<CR>

function! SendToSlack()
    call inputsave()
    let s_channel = input("Slack Channel? ", g:slackcat_default_channel)
    let s_lang = input("lang? ", &filetype)
    call inputrestore()
    echo "\rSending to Slack ..."
    let s_selection = s:escapeTildes(s:getVisualSelection())
    if empty(s_lang)
        let s_lang = 'txt'
    endif
    let return = system("echo '". s_selection ."' |slackcat -c " . s_channel . " --filetype " . s_lang)
    echo "\rSent !"
endfunction

function! s:escapeTildes(text)
    return substitute(a:text, "'", "'\"'\"'", 'g')
endfunction

function! s:getVisualSelection()
    " Why is this not a built-in Vim script function?!
    let [lnum1, col1] = getpos("'<")[1:2]
    let [lnum2, col2] = getpos("'>")[1:2]
    let lines = getline(lnum1, lnum2)
    let lines[-1] = lines[-1][: col2 - (&selection == 'inclusive' ? 1 : 2)]
    let lines[0] = lines[0][col1 - 1:]
    return join(lines, "\n")
endfunction
