#include "Stdlib.h"

#include <wincrypt.h>

#pragma comment(lib, "crypt32.lib")

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

BYTE* Base64Decode(const char* input, DWORD* outLen) {
	DWORD dwDecodedLen = 0;

	if (!CryptStringToBinaryA(input, 0, CRYPT_STRING_BASE64, NULL, &dwDecodedLen, NULL, NULL)) {
		return NULL;
	}

	BYTE* decodedData = (BYTE*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, dwDecodedLen);
	if (!decodedData) return NULL;

	if (!CryptStringToBinaryA(input, 0, CRYPT_STRING_BASE64, decodedData, &dwDecodedLen, NULL, NULL)) {
		HeapFree(GetProcessHeap(), 0, decodedData);
		return NULL;
	}

	*outLen = dwDecodedLen;
	return decodedData;
}