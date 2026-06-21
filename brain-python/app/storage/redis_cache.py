import os
import json
import redis

class RedisCacheAdapter:
    def __init__(self):
        # Read internal environment parameters set by Docker Compose
        host = os.getenv("REDIS_HOST", "natakos-redis")
        port = int(os.getenv("REDIS_PORT", 6379))
        
        # Open a resilient network connection pool client
        self.client = redis.Redis(host=host, port=port, decode_responses=True)
        print(f"🔌 [BRAIN-PYTHON] Connected successfully to in-memory cache vault at {host}:{port}")

    def get_chat_history(self, sender_jid: str) -> list:
        """Extracts the rolling history list from the cache vault"""
        cache_key = f"chat:history:{sender_jid}"
        try:
            raw_list = self.client.lrange(cache_key, 0, -1)
            if not raw_list:
                return []
            return [json.loads(item) for item in raw_list]
        except Exception as e:
            print(f"⚠️ [REDIS ERROR] Failed to fetch chat history: {e}")
            return []

    def append_message_to_history(self, sender_jid: str, role: str, text: str, max_history: int = 6):
        """Pushes a new message slice to memory and sets a rolling 30-minute expiration"""
        cache_key = f"chat:history:{sender_jid}"
        try:
            # Append the new conversation record item as JSON string
            item = json.dumps({"role": role, "text": text})
            self.client.rpush(cache_key, item)
            
            # Trim to prevent unbounded growth
            self.client.ltrim(cache_key, -max_history, -1)
            
            # Slide the rolling expiration
            self.client.expire(cache_key, 1800) # 1800 seconds = 30 mins
        except Exception as e:
            print(f"⚠️ [REDIS ERROR] Failed to append message to history: {e}")