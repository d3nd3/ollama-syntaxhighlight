# ollama console syntax highlighting :rocket:

## Usage
`git clone https://github.com/d3nd3/ollama-syntaxhighlight`  
`cd ollama-syntaxhighlight`  
`go build chunk-reader.go`  

`nano ~/.bash_aliases`
```bash
llama() {
    OL_LANG="${1:-fstab}"
    ollama run llama3 | OL_LANG="$OL_LANG" ~/ollama-syntaxhighlight/chunk-reader
}
```

```bash
sudo apt install batcat
sudo apt install bat

batcat -L
llama python
llama C++
llama javascript
```

Defaults to using language: **fstab** if no language specified , because its very generic.


## Creation Journey
So at first, I expected it to just work by using stdbuf -i0 -o0 and simple pipes. However it wasn't so simple because:
* There were ansi terminal sequences that had to be stripped
* The syntax highlighter program does not write out instantly what it reads, it waits for newlines etc. and ollama first line doesn't use a newline.
* The syntax highlighter needs the entire context/buffer in memory at once, so it can't do incremental write out.

So I then tried appending a newline to every output of ollama, so that it forces batcat to instantly dump everything it receives. This means that during LLM incremental printing, every word has its own line.  So a further post-processing step had to be performed to remove the lines again.  

As final solution because of above problems, I called batcat directly inside of the golang chunk-reader program.  Now after the incremetal output of the LLM, it also does a full output with syntax highlighting too!

I am happy with the result, but **what** a journey.

 