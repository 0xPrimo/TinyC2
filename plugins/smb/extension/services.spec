x64:
	load "bin/services.x64.o"
		merge
	dfr "resolve" "ror13" "KERNEL32, NTDLL"
	dfr "resolve_ext" "strings"
