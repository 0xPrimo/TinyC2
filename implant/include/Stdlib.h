#ifndef _STDLIB_H_
#define _STDLIB_H_

#include <windows.h>

CHAR*  Base64Encode( const BYTE* input, DWORD size );
BYTE*  Base64Decode( const char* input, DWORD* outLen );
SIZE_T CharStringToWCharString( _Inout_ PWCHAR Destination, _In_ PCHAR Source, SIZE_T _In_ MaximumAllowed );
BOOL   FsFileRead( CONST CHAR* path, LPVOID* buffer, DWORD* buffsize );

#endif // _STDLIB_H_
