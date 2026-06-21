import os
import sys
from app.storage.vector_store import VectorStoreAdapter

def run_data_extraction_and_sync():
    print("💎 [PIPELINE] Scanning PostgreSQL database for unmapped text rules...")
    store = VectorStoreAdapter()
    
    try:
        with store.conn.cursor() as cur:
            cur.execute("SELECT id, rule_text FROM property_rules WHERE embedding IS NULL;")
            unmapped_rows = cur.fetchall()
            
            if not unmapped_rows:
                print("✨ [PIPELINE] All property rules have already been fully vectorized and mapped!")
                return
                
            print(f"📦 [PIPELINE] Found {len(unmapped_rows)} text rules needing mathematical vectorization.")
            
            for row in unmapped_rows:
                rule_id, rule_text = row[0], row[1]
                print(f"   • Processing Row ID {rule_id}: '{rule_text[:40]}...'")
                
                # Compute vector matrix arrays
                vec = store.get_embedding(rule_text)
                
                # Stamp the array directly back into Postgres!
                store.update_rule_embedding(rule_id, vec)
                print(f"   ✅ Row ID {rule_id} fully vectorized into 1536-dimension space!")
                
        print("🎉 [PIPELINE] Data extraction pipeline finalized successfully.")
    except Exception as e:
        print(f"❌ [PIPELINE] Sync failed: {e}")
    finally:
        store.close()

if __name__ == "__main__":
    run_data_extraction_and_sync()