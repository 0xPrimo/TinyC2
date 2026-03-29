#ifndef _IMPORTS_H_
#define _IMPORTS_H_

#include <winsock2.h>
#include <ws2tcpip.h>
#include <windows.h>

WINBASEAPI DWORD WINAPI KERNEL32$GetLastError(VOID);
#define GetLastError KERNEL32$GetLastError

WINBASEAPI HANDLE KERNEL32$GetProcessHeap();
#define GetProcessHeap KERNEL32$GetProcessHeap

WINBASEAPI PVOID NTAPI NTDLL$RtlAllocateHeap(HANDLE HeapHandle, ULONG Flags, SIZE_T Size);
#define RtlAllocateHeap NTDLL$RtlAllocateHeap

WINBASEAPI PVOID NTAPI NTDLL$RtlReAllocateHeap(HANDLE HeapHandle, ULONG Flags, PVOID BaseAddress, SIZE_T Size);
#define RtlReAllocateHeap NTDLL$RtlReAllocateHeap

WINBASEAPI BOOL NTAPI NTDLL$RtlFreeHeap(HANDLE HeapHandle,_In_opt_ ULONG Flags, PVOID BaseAddress);
#define RtlFreeHeap NTDLL$RtlFreeHeap

WINSOCK_API_LINKAGE int WSAAPI WS2_32$send(SOCKET s,const char *buf,int len,int flags);
#define send WS2_32$send

WINSOCK_API_LINKAGE int WSAAPI WS2_32$closesocket(SOCKET s);
#define closesocket WS2_32$closesocket

WINSOCK_API_LINKAGE int WSAAPI WS2_32$connect(SOCKET s,const struct sockaddr *name,int namelen);
#define connect WS2_32$connect

WINSOCK_API_LINKAGE int WSAAPI WS2_32$ioctlsocket(SOCKET s,__LONG32 cmd,u_long *argp);
#define ioctlsocket WS2_32$ioctlsocket

WINSOCK_API_LINKAGE int WSAAPI WS2_32$recv(SOCKET s,char *buf,int len,int flags);
#define recv WS2_32$recv

WINSOCK_API_LINKAGE INT WSAAPI WS2_32$inet_pton(INT Family, LPCSTR pStringBuf, PVOID pAddr);
#define inet_pton WS2_32$inet_pton

WINSOCK_API_LINKAGE u_long WSAAPI WS2_32$ntohl(u_long netlong);
#define ntohl WS2_32$ntohl

WINSOCK_API_LINKAGE int WSAAPI WS2_32$WSAStartup(WORD wVersionRequested,LPWSADATA lpWSAData);
#define WSAStartup WS2_32$WSAStartup

WINSOCK_API_LINKAGE u_long WSAAPI WS2_32$htonl(u_long hostlong);
#define htonl WS2_32$htonl

WINSOCK_API_LINKAGE u_short WSAAPI WS2_32$htons(u_short hostshort);
#define htons WS2_32$htons

WINSOCK_API_LINKAGE int WSAAPI WS2_32$WSACleanup(void);
#define WSACleanup WS2_32$WSACleanup

WINSOCK_API_LINKAGE int WSAAPI WS2_32$WSAGetLastError(void);
#define WSAGetLastError WS2_32$WSAGetLastError

WINSOCK_API_LINKAGE SOCKET WSAAPI WS2_32$socket(int af,int type,int protocol);
#define socket WS2_32$socket

WINBASEAPI int KERNEL32$lstrlenA(LPCSTR lpString);
#define lstrlenA KERNEL32$lstrlenA

WINBASEAPI void __cdecl MSVCRT$memset(void* dest, int c, size_t count);
#define memset MSVCRT$memset

WINBASEAPI void* __cdecl MSVCRT$memcpy(void* __restrict__ _Dst, const void* __restrict__ _Src, size_t _MaxCount);
#define memcpy MSVCRT$memcpy

#endif // _IMPORTS_H_