#include "Stdlib.h"
#include <stdio.h>

#define INITIAL_CAPACITY 4096
#define CHUNK_SIZE 4096

BOOL ReadPipeAll( HANDLE hPipe, BYTE** out_buf, DWORD* out_len ) {
    HANDLE hHeap    = GetProcessHeap();
    SIZE_T capacity = INITIAL_CAPACITY;
    DWORD  total    = 0;

    BYTE* buf = (BYTE*)HeapAlloc( hHeap, 0, capacity );
    if (!buf)
        return FALSE;

    for (;;) {
        if (capacity - total < CHUNK_SIZE) {
            SIZE_T new_cap = capacity * 2;
            if (new_cap < capacity) { // overflow guard
                HeapFree( hHeap, 0, buf );
                SetLastError( ERROR_ARITHMETIC_OVERFLOW );
                return FALSE;
            }
            BYTE* nbuf = (BYTE*)HeapReAlloc( hHeap, 0, buf, new_cap );
            if (!nbuf) {
                DWORD err = GetLastError();
                HeapFree( hHeap, 0, buf ); // original still valid on failure
                SetLastError( err );
                return FALSE;
            }
            buf      = nbuf;
            capacity = new_cap;
        }

        DWORD got = 0;
        BOOL  ok  = ReadFile( hPipe, buf + total, (DWORD)( capacity - total ), &got, NULL );
        total += got;

        if (ok) {
            break;
        }

        DWORD err = GetLastError();
        if (err == ERROR_MORE_DATA) {
            continue;
        }

        HeapFree( hHeap, 0, buf );
        SetLastError( err );
        return FALSE;
    }

    *out_buf = buf;
    *out_len = total;
    return TRUE;
}

BOOL WritePipeAll( HANDLE hPipe, const void* data, size_t len ) {
    const BYTE* p = static_cast<const BYTE*>( data );
    while (len > 0) {
        DWORD chunk   = ( len > ( 1u << 20 ) ) ? ( 1u << 20 ) : static_cast<DWORD>( len );
        DWORD written = 0;
        if (!WriteFile( hPipe, p, chunk, &written, nullptr )) {
            printf( "WriteFile failed: %ld\n", GetLastError() );
            return FALSE;
        }
        if (written == 0) {
            SetLastError( ERROR_WRITE_FAULT );
            return FALSE;
        }
        p += written;
        len -= written;
    }
    return TRUE;
}
