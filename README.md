# supreme-fishstick

## Requirements
- Golang ^1.8 installed

## Setup
- clone the repo

```bash
 git clone https://github.com/Likeaprayer/supreme-fishstick.git

```
-open it
```bash
    cd supreme-fishstick/c
```
## Concurrent Log Processor

Input:
- Log File: A file containing lines of logs (sample content shown below).
- Keywords: A list of keywords you need to track and count.

Instructions 
- change directory into the `conlogpro` folder
```bash
    cd conlogpro
```
- run the program with the path to the file and the keywords
```bash
    go run . <absolute-path-to-file> <keyword1> <keyword2> ...
```

example input 
```bash
    go run . s.txt info error debug  
```

example output
```bash
Keyword Frequency:
---------------
info : 250
error : 167
debug : 83
---------------
```

## Prime Palindrome Concurrency

Input
- A single integer, `N` (1 ≤ N ≤ 50), representing the number of prime palindromic numbers to find.


Output:
- A single integer, representing the sum of the first `N` prime palindromic numbers.

Instuctions 
- change directory into the `pripacon` folder
```bash
    cd pripacon
```
- run the program with the integer N
```bash
    go run . <Integer>
```
