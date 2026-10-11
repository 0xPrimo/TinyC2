#ifndef _STDLIB_H_
#define _STDLIB_H_

#include <windows.h>

CHAR*  Base64Encode( const BYTE* input, DWORD size );
BYTE*  Base64Decode( const char* input, DWORD* outLen );
SIZE_T CharStringToWCharString( _Inout_ PWCHAR Destination, _In_ PCHAR Source, SIZE_T _In_ MaximumAllowed );
BOOL   FsFileRead( CONST CHAR* path, LPVOID* buffer, DWORD* buffsize );
BOOL   ReadPipeAll( HANDLE hPipe, BYTE** out_buf, DWORD* out_len );
BOOL   WritePipeAll( HANDLE hPipe, const void* data, size_t len );

#endif // _STDLIB_H_
