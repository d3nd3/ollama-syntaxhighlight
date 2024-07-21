package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    buffer := make([]byte, 1024) // Adjust buffer size as needed

    for {
        n, err := os.Stdin.Read(buffer)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
            os.Exit(1)
        }

        if n > 0 {
            //The inserted color by batcat creates much larger 'words' when colored.
            // fmt.Printf("Size: %d",n);

            bufferAsString := buffer[:n]
            // words := strings.Split(bufferAsString)

            // for _, word := range words {
            //     if (len(word) <= 20) {
            //         bufferAsString = strings.Replace(bufferAsString,word+"\n",word)
            //     }
            // }


            //All received input is new-line ended.
            //strip newline from small
            lenBuf := len(bufferAsString)
            if ( (lenBuf>1 || bufferAsString[0] != '\n') && bufferAsString[lenBuf-1] == '\n' ) {
                bufferAsString = bufferAsString[:lenBuf-1]
            }
            // Write the chunk
            os.Stdout.Write(bufferAsString)
            os.Stdout.Sync()
            
            // if ( n > 20 ) {
            //     fmt.Println()
            // } else {
            //      os.Stdout.WriteString("\x03")
            //      // Ensure it's flushed to stdout
            //      os.Stdout.Sync()
            // }
            
        }
    }
}   