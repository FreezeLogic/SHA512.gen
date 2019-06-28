## How using

SHA512.gen.exe -path

path - may be filepath or folder. if path conatain folder this folder scanned recursively.

By default if file with 'filename.exe' dont contain 'filename.exe.sha512' pair in the same folder, tool generate new checksum file,  but if folder contain checksum file, tool verify him and print status in terminal.
This behavior may be overridden via cmd flags, flag **-c** force enable only checking mode.
Another flag **-b**, this flag produce bool logic for printing in terminal, tool respond onlu true or false, this may be need if you run this tool from another code use  *exec* or similar command.
