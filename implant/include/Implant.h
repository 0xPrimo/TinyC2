#pragma once

#include "Ntdll.h"
#include <vector>
#include <string>
#include <nlohmann/json.hpp>

#include "Channel.h"

#pragma comment(lib, "ntdll.lib")

using json = nlohmann::json;
using string = std::string;

#define JOB_TYPE_THREAD 	0x00000000
#define JOB_TYPE_PROCESS 	0x00000001

extern "C" {
	_Ret_range_(<= , MAXLONG)
		NTSYSAPI
		ULONG
		NTAPI
		RtlRandomEx(
			_Inout_ PULONG Seed
		);
}

inline ULONG RandomUint32() {
	ULONG seed = GetTickCount();
	return RtlRandomEx(&seed);
}

typedef struct {
	LIST_ENTRY	ListEntry;
	DWORD 		ID;
	HANDLE 		hProcess;
	HANDLE		hThread;
	HANDLE 		hAnonPipe;
	DWORD 		Type;
	DWORD 		Status;
} JOB, *PJOB;

// IImplant interface
//
typedef struct {
	// implant id
	DWORD 		SessionID;
	LIST_ENTRY	JobList;
} IImplant;


BOOL ImplantInitialize();
void ImplantRegister();
void ImplantLoop();

extern IImplant g_Implant;
