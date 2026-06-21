# The structural contract telling the LLM how to parse and call our system hardware
ROOM_TOOL_DEFINITION = {
    "type": "function",
    "function": {
        "name": "check_room_status",
        "description": "Mengecek status ketersediaan live kamar kos berdasarkan nomor kamar secara real-time dari database.",
        "parameters": {
            "type": "object",
            "properties": {
                "room_number": {
                    "type": "string",
                    "description": "Nomor kamar kos yang ditanyakan oleh user. Contoh: '101', 'Kamar 1', 'Kamar 202'."
                }
            },
            "required": ["room_number"]
        }
    }
}