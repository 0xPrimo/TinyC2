#include "Command.h"
#include <iostream>

// -> Gurkirat Singh
//      - https://tbhaxor.com/dumping-token-information-in-windows/
//
string GetProcessTokenUser(ULONG pid) {
	HANDLE          hProcess    = NULL;
    HANDLE          hToken      = NULL;
    ULONG           size        = 0;
    PTOKEN_USER     tu          = NULL;
    DWORD           szAccount    = MAX_PATH;
    DWORD           szDomain     = MAX_PATH;
    string          sTokenUser  = "";
    CHAR            account[MAX_PATH];
    CHAR            domain[MAX_PATH];
    SID_NAME_USE    snu;

    // open handle to process object
	hProcess = OpenProcess(PROCESS_QUERY_INFORMATION, FALSE, pid);
	if (hProcess == NULL || hProcess == INVALID_HANDLE_VALUE) {
		hProcess = OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, FALSE, pid);
		if (hProcess == NULL || hProcess == INVALID_HANDLE_VALUE) {
			printf("OpenProcess failed: ( %d )\n", GetLastError());
            goto Cleanup;
        }
	}


    // open handle to process token object
	if (!OpenProcessToken(hProcess, TOKEN_QUERY, &hToken)) {
        printf("OpenProcessToken failed: ( %d )\n", GetLastError());
        goto Cleanup;
    }

    // get token struct size
    GetTokenInformation(hToken, TokenUser, NULL, 0x0, &size);

    tu = (PTOKEN_USER)RtlAllocateHeap(GetProcessHeap(), HEAP_ZERO_MEMORY, size);
    if (tu == NULL) {
		printf("RtlAllocateHeap failed: ( %d )\n", GetLastError());
		goto Cleanup;
    }
    
    if (!GetTokenInformation(hToken, TokenUser, tu, size, &size)) {
        printf("GetTokenInformation failed: ( %d )\n", GetLastError());
        goto Cleanup;
    }

    // get account name and domain name
	if (!LookupAccountSidA(NULL, tu->User.Sid, account, &szAccount, domain, &szDomain, &snu)) {
		printf("LookupAccountSidA failed: ( %d )\n", GetLastError());
        goto Cleanup;
	}

    sTokenUser += domain + string("/") + account;
Cleanup:
    if (hProcess)
        CloseHandle(hProcess);
    if (tu)
        RtlFreeHeap(GetProcessHeap(), 0, tu);


    return sTokenUser;
}

// -> @hasherezade
//      - https://gist.github.com/hasherezade/c3f82fb3099fb5d1afd84c9e8831af1e
//
BOOL CommandPs(json& args, string artifact, json& result) {
    PSYSTEM_PROCESS_INFORMATION pi      = NULL;
    ULONG                       size    = 0;
    json                        pslist  = json::array();

    NtQuerySystemInformation(SystemProcessInformation, 0, 0, &size );

    pi = (PSYSTEM_PROCESS_INFORMATION)RtlAllocateHeap(GetProcessHeap(), HEAP_ZERO_MEMORY, size);
    if (pi == NULL) {
        printf("RtlAllocateHeap failed\n");
        return FALSE;
    }

    if (NtQuerySystemInformation(SystemProcessInformation, pi, size, &size) == STATUS_SUCCESS) {
        while (true) {
            json process = json::object();
            
            process["ppid"]     = (ULONG)pi->InheritedFromUniqueProcessId;
            process["pid"]      = (ULONG)pi->UniqueProcessId;
            if (pi->ImageName.Buffer == NULL) {
                process["name"]     = "";
            } else {
                process["name"]     = std::wstring(pi->ImageName.Buffer).c_str();
            }

            process["account"]  = GetProcessTokenUser(pi->InheritedFromUniqueProcessId);
            pslist.push_back(process);


            if (!pi->NextEntryOffset) {
                break;
            }

            pi = (SYSTEM_PROCESS_INFORMATION*)((ULONG_PTR)pi + pi->NextEntryOffset);
        }

        result["name"]      = "ps"; 
        result["artifact"]  = pslist.dump();

        return TRUE;
    }

    result["name"]      = "ps"; 
    result["artifact"]  = "";

    printf("NtQuerySystemInformation failed\n");
    return FALSE;
}
