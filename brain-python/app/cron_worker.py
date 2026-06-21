import os
import sys
import time
import json
import pika

def run_automated_billing_cycle():
    print("\n⏰ [CRON WORKER] Internal alarm clock triggered! Initiating monthly billing cycle...")
    
    amqp_url = os.getenv("AMQP_SERVER_URL", "amqp://guest:guest@natakos-rabbitmq:5672/")
    params = pika.URLParameters(amqp_url)
    
    try:
        connection = pika.BlockingConnection(params)
        channel = connection.channel()
        channel.queue_declare(queue="interaction.queue", durable=True)
        
        # Simulating scanning our Postgres 'tenants' table for active billing targets
        # In a real run, this reads the unique JIDs directly from the DB!
        simulated_active_tenants = [
            {"name": "Budi", "phone_jid": "628123456789@s.whatsapp.net", "room": "Kamar 101", "role": "TENANT"}
        ]
        
        for tenant in simulated_active_tenants:
            print(f"   💸 Calculating monthly lease balance for {tenant['name']} ({tenant['room']})...")
            
            # Pack the automated system generation payload contract
            payload = {
                "message_id": f"CRON-INV-{int(time.time())}",
                "sender_jid": tenant["phone_jid"],
                "raw_text": f"SISTEM AUTOMATION: Halo {tenant['name']}, invoice sewa untuk {tenant['room']} telah diterbitkan. Sila cek tagihan kamu.",
                "tenant_name": tenant["name"],
                "room_number": tenant["room"]
            }
            
            # Drop the system event straight onto the event highway!
            correlation_id = f"req_cron_{int(time.time())}_{tenant['phone_jid'].split('@')[0]}"
            channel.basic_publish(
                exchange="",
                routing_key="interaction.queue",
                properties=pika.BasicProperties(
                    content_type="application/json",
                    delivery_mode=2, # Make message persistent on disk
                    headers={
                        "X-User-Role": tenant["role"], # Injected system authorization context
                        "X-Correlation-ID": correlation_id
                    }
                ),
                body=json.dumps(payload)
            )
            print(f"   📤 [CRON WORKER] Shot automated invoice event for '{tenant['name']}' into RabbitMQ queue. CorrelationID: {correlation_id}")
            
        connection.close()
        print("✅ [CRON WORKER] Billing cycle finalized successfully. Going back to sleep...")
        
    except Exception as e:
        print(f"❌ [CRON WORKER] Automation failure: {e}")

if __name__ == "__main__":
    print("🚀 Booting up Background Cron-Automation Daemon...")
    print("⏱️  Scheduler set to fire a simulated billing run every 10 seconds for local testing.")
    
    # Live execution loop tracking time intervals
    while True:
        run_automated_billing_cycle()
        time.sleep(10) # Wakes up every 10 seconds for quick local demonstration!