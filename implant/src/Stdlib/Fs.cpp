#include "Stdlib.h"

#include <stdio.h>

BOOL FsFileRead(CONST CHAR *path, LPVOID *buffer, DWORD *buffsize) {
    HANDLE hFile = NULL;
	LPVOID buf = NULL;
	SIZE_T size = 0;

	hFile = CreateFileA( path, GENERIC_READ, 0, 0, OPEN_EXISTING, 0, 0 );
	if ( hFile == INVALID_HANDLE_VALUE ) {
		printf( "[-] CreateFileA: ( %d )\n", GetLastError( ) );
		goto Cleanup;
	}

	size = GetFileSize( hFile, 0 );
	if ( size == INVALID_FILE_SIZE ) {
		printf( "[-] GetFileSize: ( %d )\n", GetLastError( ) );
		goto Cleanup;
	}

	buf = HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, size );
	if ( !buf ) {
		printf( "[-] LocalAlloc: ( %d )\n", GetLastError( ) );
		goto Cleanup;
	}

	if ( !ReadFile( hFile, buf, size, NULL, NULL ) ) {
		printf( "[-] ReadFile: ( %d )\n", GetLastError( ) );
		goto Cleanup;
	}

	if ( hFile )
		CloseHandle( hFile );

	*buffer     = buf;
    *buffsize   = size;
    return TRUE;
Cleanup:
	if ( hFile )
		CloseHandle( hFile );
	if ( buf )
		LocalFree( buf );

	return FALSE;
}