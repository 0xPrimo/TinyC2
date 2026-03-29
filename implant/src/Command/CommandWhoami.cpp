#include "Command.h"

#pragma comment(lib, "Secur32.lib")

/*
 * This is a beacon object file copied from trustedsec CS-Situational-Awareness-BOF repo.
 * https://github.com/trustedsec/CS-Situational-Awareness-BOF/blob/master/src/SA/whoami/entry.c
 */

void internal_printf(std::string &output, const char* format, ...) {
    char buffer[4096];

    va_list args;
    va_start(args, format);

    vsnprintf(buffer, sizeof(buffer), format, args);

    va_end(args);

    output += buffer;
}


#include <windows.h>
#define SECURITY_WIN32
#include <security.h>
#include <sddl.h>

typedef struct
{
    UINT Rows;
    UINT Cols;
    LPWSTR Content[1];
} WhoamiTable;

char* WhoamiGetUser(EXTENDED_NAME_FORMAT NameFormat)
{
    char* UsrBuf = (char*)malloc(MAX_PATH);
    ULONG UsrSiz = MAX_PATH;

    if (UsrBuf == NULL)
        return NULL;

    if (GetUserNameExA(NameFormat, UsrBuf, &UsrSiz))
    {
        return UsrBuf;
    }

    free(UsrBuf);
    return NULL;
}

VOID* WhoamiGetTokenInfo(std::string &output, TOKEN_INFORMATION_CLASS TokenType)
{
    HANDLE hToken = 0;
    DWORD dwLength = 0;
    VOID* pTokenInfo = 0;
    VOID* pResult = 0;

    if (OpenProcessToken(GetCurrentProcess(), TOKEN_READ, &hToken))
    {
        GetTokenInformation(hToken,
            (_TOKEN_INFORMATION_CLASS)TokenType,
            NULL,
            dwLength,
            &dwLength);

        if (GetLastError() == ERROR_INSUFFICIENT_BUFFER)
        {
            pTokenInfo = malloc(dwLength);
            if (pTokenInfo == NULL)
            {
                internal_printf(output, "ERROR: not enough memory to allocate the token structure.\r\n");
                goto cleanup;
            }
        }
        else {
            goto cleanup;
        }

        if (!GetTokenInformation(hToken, (_TOKEN_INFORMATION_CLASS)TokenType,
            (LPVOID)pTokenInfo,
            dwLength,
            &dwLength))
        {
            internal_printf(output, "ERROR 0x%x: could not get token information.\r\n", GetLastError());
            goto cleanup;
        }
        pResult = pTokenInfo;
        pTokenInfo = NULL;
    }

cleanup:
    if (hToken) {
        CloseHandle(hToken);
    }
    if (pTokenInfo) {
        free(pTokenInfo);
    }

    return pResult;
}

int WhoamiUser(std::string& output)
{
    PTOKEN_USER pUserInfo = (PTOKEN_USER)WhoamiGetTokenInfo(output, TokenUser);
    char* pUserStr = NULL;
    char* pSidStr = NULL;
    WhoamiTable* UserTable = NULL;
    int retval = 0;

    if (pUserInfo == NULL)
    {
        retval = 1;
        goto end;
    }

    pUserStr = WhoamiGetUser(NameSamCompatible);
    if (pUserStr == NULL)
    {
        retval = 1;
        goto end;
    }

    internal_printf(output, "\nUserName\t\tSID\n");
    internal_printf(output, "====================== ====================================\n");

    if (ConvertSidToStringSidA(pUserInfo->User.Sid, &pSidStr)) {
        internal_printf(output, "%s\t%s\n\n", pUserStr, pSidStr);
        LocalFree(pSidStr);
        pSidStr = NULL;
    };



    /* cleanup our allocations */
end:
    if (pSidStr) { LocalFree(pSidStr); }
    if (pUserInfo) { free(pUserInfo); }
    if (pUserStr) { free(pUserStr); };

    return retval;
}

int WhoamiGroups(std::string& output)
{
    DWORD dwIndex = 0;
    char* pSidStr = NULL;

    char szGroupName[255] = { 0 };
    char szDomainName[255] = { 0 };

    DWORD cchGroupName = _countof(szGroupName);
    DWORD cchDomainName = _countof(szDomainName);

    SID_NAME_USE Use = (SID_NAME_USE)0;

    PTOKEN_GROUPS pGroupInfo = (PTOKEN_GROUPS)WhoamiGetTokenInfo(output, TokenGroups);
    WhoamiTable* GroupTable = NULL;

    if (pGroupInfo == NULL)
    {
        return 1;
    }

    /* the header is the first (0) row, so we start in the second one (1) */


    internal_printf(output, "\n%-50s%-25s%-45s%-25s\n", "GROUP INFORMATION", "Type", "SID", "Attributes");
    internal_printf(output, "================================================= ===================== ============================================= ==================================================\n");

    for (dwIndex = 0; dwIndex < pGroupInfo->GroupCount; dwIndex++)
    {
        if (LookupAccountSidA(NULL,
            pGroupInfo->Groups[dwIndex].Sid,
            (LPSTR)&szGroupName,
            &cchGroupName,
            (LPSTR)&szDomainName,
            &cchDomainName,
            &Use) == 0)
        {
            //If we fail lets try to get the next entry
            continue;
        }

        /* the original tool seems to limit the list to these kind of SID items */
        if ((Use == SidTypeWellKnownGroup || Use == SidTypeAlias ||
            Use == SidTypeLabel || Use == SidTypeGroup) && !(pGroupInfo->Groups[dwIndex].Attributes & SE_GROUP_LOGON_ID))
        {
            char tmpBuffer[1024] = { 0 };

            /* looks like windows treats 0x60 as 0x7 for some reason, let's just nod and call it a day:
               0x60 is SE_GROUP_INTEGRITY | SE_GROUP_INTEGRITY_ENABLED
               0x07 is SE_GROUP_MANDATORY | SE_GROUP_ENABLED_BY_DEFAULT | SE_GROUP_ENABLED */

            if (pGroupInfo->Groups[dwIndex].Attributes == 0x60)
                pGroupInfo->Groups[dwIndex].Attributes = 0x07;

            /* 1- format it as DOMAIN\GROUP if the domain exists, or just GROUP if not */
            sprintf((char*)&tmpBuffer, "%s%s%s", szDomainName, cchDomainName ? "\\" : "", szGroupName);
            internal_printf(output, "%-50s", tmpBuffer);

            /* 2- let's find out the group type by using a simple lookup table for lack of a better method */
            if (Use == SidTypeWellKnownGroup) {
                internal_printf(output, "%-25s", "Well-known group ");
            }
            else if (Use == SidTypeAlias) {
                internal_printf(output, "%-25s", "Alias ");
            }
            else if (Use == SidTypeLabel) {
                internal_printf(output, "%-25s", "Label ");
            }
            else if (Use == SidTypeGroup) {
                internal_printf(output, "%-25s", "Group ");
            }
            /* 3- turn that SID into text-form */
            if (ConvertSidToStringSidA(pGroupInfo->Groups[dwIndex].Sid, &pSidStr)) {

                //WhoamiSetTable(GroupTable, pSidStr, PrintingRow, 2);
                internal_printf(output, "%-45s ", pSidStr);

                LocalFree(pSidStr);
                pSidStr = NULL;

            }

            /* 4- reuse that buffer for appending the attributes in text-form at the very end */
            ZeroMemory(tmpBuffer, sizeof(tmpBuffer));

            if (pGroupInfo->Groups[dwIndex].Attributes & SE_GROUP_MANDATORY)
                internal_printf(output, "Mandatory group, ");
            if (pGroupInfo->Groups[dwIndex].Attributes & SE_GROUP_ENABLED_BY_DEFAULT)
                internal_printf(output, "Enabled by default, ");
            if (pGroupInfo->Groups[dwIndex].Attributes & SE_GROUP_ENABLED)
                internal_printf(output, "Enabled group, ");
            if (pGroupInfo->Groups[dwIndex].Attributes & SE_GROUP_OWNER)
                internal_printf(output, "Group owner, ");
            internal_printf(output, "\n");
        }
        /* reset the buffers so that we can reuse them */
        ZeroMemory(szGroupName, sizeof(szGroupName));
        ZeroMemory(szDomainName, sizeof(szDomainName));

        cchGroupName = 255;
        cchDomainName = 255;
    }


    /* cleanup our allocations */
    free(pGroupInfo);

    return 0;
}

int WhoamiPriv(std::string& output)
{
    PTOKEN_PRIVILEGES pPrivInfo = (PTOKEN_PRIVILEGES)WhoamiGetTokenInfo(output, TokenPrivileges);
    DWORD dwResult = 0, dwIndex = 0;
    WhoamiTable* PrivTable = NULL;

    if (pPrivInfo == NULL)
    {
        return 1;
    }

    internal_printf(output, "\n\n%-30s%-50s%-30s\n", "Privilege Name", "Description", "State");
    internal_printf(output, "============================= ================================================= ===========================\n");

    for (dwIndex = 0; dwIndex < pPrivInfo->PrivilegeCount; dwIndex++)
    {
        char* PrivName = NULL;
        char* DispName = NULL;
        DWORD PrivNameSize = 0, DispNameSize = 0;
        BOOL ret = FALSE;

        LookupPrivilegeNameA(NULL,
            &pPrivInfo->Privileges[dwIndex].Luid,
            NULL,
            &PrivNameSize); // getting size

        PrivName = (char*)malloc(++PrivNameSize);

        if (LookupPrivilegeNameA(NULL,
            &pPrivInfo->Privileges[dwIndex].Luid,
            PrivName,
            &PrivNameSize) == 0)
        {
            if (PrivName) { free(PrivName); PrivName = NULL; }
            continue; // try to get next
        }

        //WhoamiSetTableDyn(PrivTable, PrivName, dwIndex + 1, 0);
        internal_printf(output, "%-30s", PrivName);


        /* try to grab the size of the string, also, beware, as this call is
           unimplemented in ReactOS/Wine at the moment */

        LookupPrivilegeDisplayNameA(NULL, PrivName, NULL, &DispNameSize, &dwResult);

        DispName = (char*)malloc(++DispNameSize);

        ret = LookupPrivilegeDisplayNameA(NULL, PrivName, DispName, &DispNameSize, &dwResult);
        if (PrivName != NULL)
            free(PrivName);
        if (ret && DispName)
        {
            internal_printf(output, "%-50s", DispName);
        }
        else
        {
            internal_printf(output, "%-50s", "???");
        }
        if (DispName != NULL)
            free(DispName);

        if (pPrivInfo->Privileges[dwIndex].Attributes & SE_PRIVILEGE_ENABLED)
            internal_printf(output, "%-30s\n", "Enabled");
        else
            internal_printf(output, "%-30s\n", "Disabled");
    }


    /* cleanup our allocations */
    if (pPrivInfo) { free(pPrivInfo); }

    return 0;
}

// CommandWhoami get implant context
//
BOOL CommandWhoami(json& args, string artifact, json& result) {
    std::string output;

    (void)WhoamiUser(output);
    (void)WhoamiGroups(output);
    (void)WhoamiPriv(output);
	
	result["name"] = "whoami";
    result["output"] = output;

    return TRUE;
}