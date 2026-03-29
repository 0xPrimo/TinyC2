#ifndef _IMPORTS_H_
#define _IMPORTS_H_

#include <windows.h>
#include <wininet.h>

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

WINBASEAPI int KERNEL32$lstrlenA(LPCSTR lpString);
#define lstrlenA KERNEL32$lstrlenA

INTERNETAPI_(HINTERNET) WININET$InternetOpenA(LPCSTR lpszAgent,DWORD dwAccessType,LPCSTR lpszProxy,LPCSTR lpszProxyBypass,DWORD dwFlags);
#define InternetOpenA WININET$InternetOpenA

INTERNETAPI_(HINTERNET) WININET$InternetConnectA(HINTERNET hInternet,LPCSTR lpszServerName,INTERNET_PORT nServerPort,LPCSTR lpszUserName,LPCSTR lpszPassword,DWORD dwService,DWORD dwFlags,DWORD_PTR dwContext);
#define InternetConnectA WININET$InternetConnectA

BOOLAPI WININET$InternetCloseHandle(HINTERNET hInternet);
#define InternetCloseHandle WININET$InternetCloseHandle

INTERNETAPI_(HINTERNET) WININET$HttpOpenRequestA(HINTERNET hConnect,LPCSTR lpszVerb,LPCSTR lpszObjectName,LPCSTR lpszVersion,LPCSTR lpszReferrer,LPCSTR *lplpszAcceptTypes,DWORD dwFlags,DWORD_PTR dwContext);
#define HttpOpenRequestA WININET$HttpOpenRequestA

BOOLAPI WININET$HttpSendRequestA(HINTERNET hRequest,LPCSTR lpszHeaders,DWORD dwHeadersLength,LPVOID lpOptional,DWORD dwOptionalLength);
#define HttpSendRequestA WININET$HttpSendRequestA

BOOLAPI WININET$InternetReadFile(HINTERNET hFile,LPVOID lpBuffer,DWORD dwNumberOfBytesToRead,LPDWORD lpdwNumberOfBytesRead);
#define InternetReadFile WININET$InternetReadFile

WINBASEAPI void __cdecl MSVCRT$memset(void* dest, int c, size_t count);
#define memset MSVCRT$memset

WINBASEAPI void* __cdecl MSVCRT$memcpy(void* __restrict__ _Dst, const void* __restrict__ _Src, size_t _MaxCount);
#define memcpy MSVCRT$memcpy

#endif // _IMPORTS_H_