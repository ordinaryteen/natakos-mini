import os
import sys
import time
import json
import pika
from openai import OpenAI
from app.storage.redis_cache import RedisCacheAdapter
from app.storage.vector_store import VectorStoreAdapter
from app.tools.manifest import ROOM_TOOL_DEFINITION

import datetime

def log_structured(level: str, correlation_id: str, message: str):
    log_payload = {
        "timestamp": datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z"),
        "level": level,
        "service": "brain-python",
        "correlation_id": correlation_id,
        "message": message
    }
    print(json.dumps(log_payload), flush=True)

def main():
    print("🚀 Booting up Resilient Brain-Python Background Worker Engine...")
    
    cache = RedisCacheAdapter()
    vector_vault = VectorStoreAdapter()

    # Initialize the live DeepSeek client pointing to their official base URL endpoint
    # We use your live token injected through the orchestrator environmental variables
    ai_client = OpenAI(
        api_key=os.getenv("AI_API_KEY"),
        base_url="https://api.deepseek.com/v1"
    )

    amqp_url = os.getenv("AMQP_SERVER_URL", "amqp://guest:guest@natakos-rabbitmq:5672/")
    params = pika.URLParameters(amqp_url)
    
    connection = None
    for i in range(5):
        try:
            connection = pika.BlockingConnection(params)
            break
        except pika.exceptions.AMQPConnectionError:
            print("⏳ Broker not ready yet, retrying in 3 seconds...")
            time.sleep(3)

    if not connection:
        print("❌ Critical: Could not connect to RabbitMQ broker highway.")
        sys.exit(1)

    channel = connection.channel()
    channel.queue_declare(queue="interaction.queue", durable=True)

    print("⚡ Listening live for incoming events inside 'interaction.queue'...")

    def callback(ch, method, properties, body):
        # Unpack transit payload data
        data = json.loads(body.decode("utf-8"))
        sender_jid = data.get("sender_jid")
        raw_text = data.get("raw_text")
        tenant_name = data.get("tenant_name")
        room_number = data.get("room_number")

        headers = properties.headers or {}
        user_role = headers.get("X-User-Role", "UNKNOWN")
        correlation_id = headers.get("X-Correlation-ID", "UNKNOWN_REQ")
        
        log_structured("INFO", correlation_id, f"Sucked message from queue. Role verified: {user_role}")
        
        final_reply = ""

        try:
            if user_role != "TENANT":
                final_reply = "Akses ditolak. Tingkatan role tidak dikenal."
                log_structured("WARN", correlation_id, f"Access denied. Inbound role was {user_role}")
            else:
                # 1. Fetch memory state context
                chat_history = cache.get_chat_history(sender_jid)
                cache.append_message_to_history(sender_jid, "user", raw_text)

                log_structured("INFO", correlation_id, "Executing Pass 1 intent classification and tool matching...")
                
                system_instruction = "You are the AI Property Manager for Kos Natakos. Use tools if the user asks about live availability."
                messages_payload = [{"role": "system", "content": system_instruction}]
                for msg in chat_history[-3:]:
                    messages_payload.append({"role": msg["role"], "content": msg["text"]})
                messages_payload.append({"role": "user", "content": raw_text})

                # FIRST API PASS: Hand the model the data AND the tools list blueprints
                response = ai_client.chat.completions.create(
                    model="deepseek-chat",
                    messages=messages_payload,
                    tools=[ROOM_TOOL_DEFINITION], # Injected system tools
                    tool_choice="auto",
                    temperature=0.1, # Keep it deterministic for structural parsing
                    timeout=10.0
                )
                
                response_message = response.choices[0].message
                tool_calls = response_message.tool_calls

                # Check if the model decided it needs hardware tool lookups!
                if tool_calls:
                    log_structured("INFO", correlation_id, "Tool intercept: DeepSeek requested database execution access")
                    messages_payload.append(response_message)
                    
                    for tool_call in tool_calls:
                        if tool_call.function.name == "check_room_status":
                            # Extract parameters computed by the LLM
                            args = json.loads(tool_call.function.arguments)
                            target_room = args.get("room_number")
                            log_structured("INFO", correlation_id, f"Querying Postgres table for room: '{target_room}'")
                            
                            # Fire actual SQL command line
                            db_status = vector_vault.check_room_status_db(target_room)
                            log_structured("INFO", correlation_id, f"Relational DB Return: Status is '{db_status}'")
                            
                            messages_payload.append({
                                "role": "tool",
                                "tool_call_id": tool_call.id,
                                "name": "check_room_status",
                                "content": json.dumps({"status": db_status, "room": target_room})
                            })
                    
                    log_structured("INFO", correlation_id, "Pass 2 synthesis processing initiated against DeepSeek model core")
                    final_response = ai_client.chat.completions.create(
                        model="deepseek-chat",
                        messages=messages_payload,
                        timeout=10.0
                    )
                    final_reply = final_response.choices[0].message.content
                else:
                    # No tools needed; Fall back to default prompt response generation
                    final_reply = response_message.content
                
                log_structured("INFO", correlation_id, "Transaction lifecycle finalized cleanly.")

        except Exception as e:
            log_structured("ERROR", correlation_id, f"Critical System Intercept caught crash: {e}")
            final_reply = f"Halo Kak {tenant_name}, layanan kami sedang optimalisasi. Silakan coba sesaat lagi."

        # =======================================================
        # FINALIZE AND ACKNOWLEDGE: THIS KILLS THE LOOP!
        # =======================================================
        cache.append_message_to_history(sender_jid, "assistant", final_reply)
        log_structured("INFO", correlation_id, "Outbound notification job dispatched back to client.")
        
        # Even if things crashed, we generated a response and handled the state. 
        # Tell RabbitMQ it is 100% safe to drop this message from disk!
        ch.basic_ack(delivery_tag=method.delivery_tag)

    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(queue="interaction.queue", on_message_callback=callback)
    channel.start_consuming()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("Stopping worker loop...")
        sys.exit(0)