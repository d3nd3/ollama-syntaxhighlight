package main

import (
    "fmt"
    "io"
    "os"
    "os/exec"
)


func StripANSI(input []byte) []byte {
    output := make([]byte, 0, len(input))
    inSequence := false
    for _, b := range input {
        if b == 0x1b { // ESC character
            inSequence = true
            continue
        }
        if inSequence {
            if (b >= 0x40 && b <= 0x5a) || (b >= 0x61 && b <= 0x7a) {
                inSequence = false
            }
            continue
        }
        output = append(output, b)
    }
    return output
}
func main() {
    buffer := make([]byte, 1024) // Adjust buffer size as needed
    userInputMode := false
    llmMsg := []byte{} //empty byte slice
    llmMsg_size := 0

    for {
        
        n, err := os.Stdin.Read(buffer)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
            os.Exit(1)
        }

        /*
            Our job is to append a newline to everything, so it gets instantly written out by batcat.
            But new issue: if its instantly written out, it doens't have syntax highlighting context.
            Since syntax highlighting is applied in chunks. Somehow. Or it depends on newlines.
        */
        if n > 0 {

            // fmt.Printf("Size: %d %02X",n,buffer[0])

            // if (n>20) {
            //     fmt.Println()
            //     fmt.Println()
            // }
            // fmt.Printf("%02X\n",buffer[0]);
            // Write the chunk
            msg := buffer[:n]
            
            //Inbetween or LLM
            if ( /*(buffer[0] != 0x1B || n==4)*/ !userInputMode ) {
                if ( n>20 ) {
                    if (llmMsg_size >0) {
                        //syntaxhighlight dump here

                        stripped := string(StripANSI(llmMsg[:llmMsg_size]))
                        language := os.Getenv("OL_LANG")
                        if (language == "") {
                            language = "fstab"
                        }
                        batCmd := exec.Command("/usr/bin/pipetty","/usr/bin/batcat", "-ppP", "-l"+language)
                        // batCmd := exec.Command("/usr/bin/xxd")
                        stdin, err := batCmd.StdinPipe()
                        if err != nil {
                            fmt.Println("Error obtaining stdin pipe:", err)
                            return
                        }

                        batStdout, err := batCmd.StdoutPipe()
                        if err != nil {
                            fmt.Println("Error obtaining stdbatStdoutout pipe:", err)
                            return
                        }

                        if err := batCmd.Start(); err != nil {
                            fmt.Println("Error starting command:", err)
                            return
                        }

                        // Send input to the process
                        go func() {
                            defer stdin.Close()
                            io.WriteString(stdin, stripped)
                        }()

                        // Read the output from the process
                        batOutput, err := io.ReadAll(batStdout)
                        if err != nil {
                            fmt.Println("Error reading batStdout:", err)
                            return
                        }

                        if err := batCmd.Wait(); err != nil {
                            fmt.Println("Error waiting for command to finish:", err)
                            return
                        }
                        fmt.Println()
                        fmt.Println("==========================================")
                        fmt.Println()
                        fmt.Println("==========================================")
                        os.Stdout.Write(batOutput)
                    }
                    //--This is a inbetween message.--
                    //clear slice
                    llmMsg_size = 0
                    // llmMsg = nil
                    userInputMode = true
                    
                    // we can't flush on each word, cos of syntax highlighting.
                    // we could flush every comma and full stop mb.

                    // fmt.Println("Writing Inbetween")

                    os.Stdout.Write(msg)
                    // fmt.Println()
                    // os.Stdout.Sync()
                } else { 
                    //----LLM Output----
                    os.Stdout.Write(msg)
                    //First output will be considered LLM, but isnt'. (>>>)
                    if (buffer[0] != 0x1B) {
                        
                        // What if we save this and parse to batcat?
                        llmMsg = append(llmMsg[:llmMsg_size],msg...)

                        llmMsg_size += n
                        //fmt.Println()
                        //os.Stdout.Sync()
                    } else {
                        // fmt.Println("Writing >>>")
                        // for i, b := range msg {
                        //     fmt.Printf("byte[%d] = %02X\n", i, b)
                        // }
                    }
                }
            } else {
                //----User is typing----
                if (n==1 ||n==4) {
                    os.Stdout.Write(msg)
                    // fmt.Println()
                    // os.Stdout.Sync()
                    if (buffer[0] == '\n') {
                        //Restore LLM mode
                        // fmt.Println("Restoring LLM detected newline")
                        userInputMode = false
                    }
                } else {
                    os.Stdout.Write(msg)
                    /*
                    fmt.Println("Restoring LLM for unrecognised sequence")
                    //Restore LLM mode
                    os.Stdout.Write(msg)
                    userInputMode = false
                    */
                }
            }
            

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