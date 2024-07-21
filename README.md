# ollama console syntax highlighting

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

Defaults to using language: fstab , because its very generic.

 