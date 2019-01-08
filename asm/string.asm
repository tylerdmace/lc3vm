.ORIG x3000                         ; Address where program is loaded
LEA R0, TEST_STR                    ; Load the address of TEST_STR into R0
PUTs                                ; Ouput the string pointed to by R0 to STDOUT
HALT                                ; Halt execution
TEST_STR .STRINGZ "This is a test!" ; Our test string
.END                                ; EOF