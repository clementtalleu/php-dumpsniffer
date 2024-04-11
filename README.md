# DumpSniffer

DumpSniffer is a command-line tool written in Go for detecting occurrences of common debugging functions like `var_dump()`, `dump()`, `die`, etc., in PHP files. 
It can be useful for identifying and cleaning up debugging code in your PHP projects.

## Features

- Detects occurrences of `var_dump()`, `dump()`, `die`, etc., in PHP files.
- Supports both single files and directories for analysis.
- Recursive scanning for files in nested directories.
- Outputs file paths and line numbers where occurrences are found.


## Installation
Before using DumpSniffer, ensure you have Go installed on your machine.
If you haven't installed Go yet, you can download and install it from the
[official Go website](https://go.dev/dl/)

Once you have Go installed, you can build the DumpSniffer executable using the following steps:
 

1. Clone this repository to your local machine:
```
git clone https://github.com/clementtalleu/php-dumpsniffer.git
```

2. Navigate to the project directory:
```
cd php-dumpsniffer
```

3. Build the executable using the Go compiler:
```
go build dumpsniffer.go
```
This command will generate an executable file named dumpsniffer in the current directory.


4. After building the executable, you can copy it to make it globally accessible from the command line. 
Run the following command:

For linux and macos users:
```
sudo cp dumpsniffer /usr/local/bin/
```

For windows users:
```
copy dumpsniffer.exe C:\Windows\System32
```

Now, DumpSniffer is installed on your system and ready to use. You can invoke it from any directory using the dumpsniffer command in your terminal.





## Usage

To analyze a directory or a single PHP file, simply provide the path as a command-line argument:

```bash
dumpsniffer /path/to/your/file.php
dumpsniffer /path/to/your/directory
```