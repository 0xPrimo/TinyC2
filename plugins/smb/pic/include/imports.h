#ifndef _IMPORTS_H_
#define _IMPORTS_H_

#include <winsock2.h>
#include <ws2tcpip.h>
#include <windows.h>

WINBASEAPI DWORD WINAPI KERNEL32$GetLastError( VOID );
#define GetLastError KERNEL32$GetLastError

WINBASEAPI HANDLE KERNEL32$GetProcessHeap();
#define GetProcessHeap KERNEL32$GetProcessHeap

WINBASEAPI PVOID NTAPI NTDLL$RtlAllocateHeap( HANDLE HeapHandle, ULONG Flags, SIZE_T Size );
#define RtlAllocateHeap NTDLL$RtlAllocateHeap

WINBASEAPI WINBOOL WINAPI KERNEL32$CloseHandle( HANDLE hObject );
#define CloseHandle KERNEL32$CloseHandle

WINBASEAPI PVOID NTAPI
NTDLL$RtlReAllocateHeap( HANDLE HeapHandle, ULONG Flags, PVOID BaseAddress, SIZE_T Size );
#define RtlReAllocateHeap NTDLL$RtlReAllocateHeap

WINBASEAPI
BOOL NTAPI NTDLL$RtlFreeHeap( HANDLE HeapHandle, _In_opt_ ULONG Flags, PVOID BaseAddress );
#define RtlFreeHeap NTDLL$RtlFreeHeap

WINBASEAPI int KERNEL32$lstrlenA( LPCSTR lpString );
#define lstrlenA KERNEL32$lstrlenA

WINBASEAPI void __cdecl MSVCRT$memset( void* dest, int c, size_t count );
#define memset MSVCRT$memset

WINBASEAPI void* __cdecl
MSVCRT$memcpy( void* __restrict__ _Dst, const void* __restrict__ _Src, size_t _MaxCount );
#define memcpy MSVCRT$memcpy

WINBASEAPI char* __cdecl MSVCRT$_strdup( const char* _Src ) __MINGW_ATTRIB_DEPRECATED_MSVC2005;
#define strdup MSVCRT$_strdup

WINBASEAPI WINBOOL WINAPI KERNEL32$WriteFile(
    HANDLE hFile, LPCVOID lpBuffer, DWORD nNumberOfBytesToWrite, LPDWORD lpNumberOfBytesWritten,
    LPOVERLAPPED lpOverlapped );
#define WriteFile KERNEL32$WriteFile

WINBASEAPI WINBOOL WINAPI KERNEL32$PeekNamedPipe(
    HANDLE hNamedPipe, LPVOID lpBuffer, DWORD nBufferSize, LPDWORD lpBytesRead,
    LPDWORD lpTotalBytesAvail, LPDWORD lpBytesLeftThisMessage );
#define PeekNamedPipe KERNEL32$PeekNamedPipe

WINBASEAPI WINBOOL WINAPI KERNEL32$ReadFile(
    HANDLE hFile, LPVOID lpBuffer, DWORD nNumberOfBytesToRead, LPDWORD lpNumberOfBytesRead,
    LPOVERLAPPED lpOverlapped );
#define ReadFile KERNEL32$ReadFile

WINBASEAPI HANDLE WINAPI KERNEL32$CreateNamedPipeA(
    LPCSTR lpName, DWORD dwOpenMode, DWORD dwPipeMode, DWORD nMaxInstances, DWORD nOutBufferSize,
    DWORD nInBufferSize, DWORD nDefaultTimeOut, LPSECURITY_ATTRIBUTES lpSecurityAttributes );
#define CreateNamedPipeA KERNEL32$CreateNamedPipeA

WINBASEAPI WINBOOL WINAPI KERNEL32$ConnectNamedPipe( HANDLE hNamedPipe, LPOVERLAPPED lpOverlapped );
#define ConnectNamedPipe KERNEL32$ConnectNamedPipe

#endif // _IMPORTS_H_
