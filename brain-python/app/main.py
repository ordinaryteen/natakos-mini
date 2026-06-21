from fastapi import FastAPI
from app.schemas import InboundInteractionRequest, ExecutionResponse
from app.storage.redis_cache import RedisCacheAdapter

app = FastAPI(title="Natakos Brain AI Engine - Dev Mode")

# Initialize our caching client utility adapter instance on startup
cache = RedisCacheAdapter()

@app.get("/")
def read_root():
    return {"status": "ONLINE", "service": "brain-python"}

@app.post("/process", response_model=ExecutionResponse)
def process_message_intent(payload: InboundInteractionRequest):
    print(f"\n📥 [BRAIN-PYTHON] Incoming network data from Go Gateway!")
    
    # 1. Fetch the short-term conversation logs for this specific user JID
    chat_history = cache.get_chat_history(payload.sender_jid)
    
    print(f"   • Target Sender JID: {payload.sender_jid}")
    print(f"   • Context History Count: {len(chat_history)} previous interactions saved.")
    
    # 2. Append the user's fresh prompt request message slice to the cache memory
    cache.append_message_to_history(payload.sender_jid, "user", payload.raw_text)

    # Simple mock simulation tracking if the user is in an active multi-turn session
    if len(chat_history) > 0:
        reply_msg = f"Sistem memori aktif! Saya ingat pesan terakhir kamu: '{chat_history[-1]['text']}'. Sekarang kamu mengirim: '{payload.raw_text}'"
    else:
        reply_msg = f"Sistem terintegrasi! Ini pesan pertama kamu dalam sesi ini: '{payload.raw_text}' (Saved to Redis cache)."

    # 3. Append the system's generated response slice to the cache thread log
    cache.append_message_to_history(payload.sender_jid, "assistant", reply_msg)

    return ExecutionResponse(
        status="SUCCESS",
        intent_detected="MULTI_TURN_CACHED_ROUTE",
        reply_message=reply_msg
    )