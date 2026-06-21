# 🤫 Natakos Mini: Chaos Brainfart

Waddup twin, welcome to the simulation of internal engine room of **0xNTK$**. This shi is a highly experimental, slightly unhinged, asynchronous event-driven AI Property Manager prototype (vibecoded, ofcourse, what do you expect?). It was built under heavy time-wasting pressure in Surabaya, Indonesia.

It is held together by pure architectural intuition, a fallback ring stuff like that, and a small can of caramel machiato. Take it away, honey.

## What is this
Basically I am trying to make a chatbot that can know which room is available in that moment, and maybe if someone ask about a boarding house rule, that bot can answer it based on the document (yea classic RAG actually). Basically, i wanna make a chatbot that pretends as a business operator for landlords, a classic payment reminder would be a good example to think of. Someday it'll be able to do accounting or maybe reporting things. My main focus on this project is simply to understand how backend works, not the businses case itself (yet).

## How tf is this even running

* **Gateway-Go:** The border control agent. It catches WhatsApp webhooks from tenant, kicks back an instant response (while thinking "damn this UI is ugly, might as well change it later"), and drops the event suitcase into the transport highway before anyone can notice.
* **Natakos-RabbitMQ-Broker:** The emotional buffer wheel. It absorbs rapid message flooding so the core engine never explodes (while hoping that this shi dont crash and burn).
* **Natakos-Redis-Cache:** The short-term memory tank. Uses strict, native Redis Lists (`RPUSH`/`LRANGE`) because string-serialization was giving us empty arrays and making the CEO cry.
* **Natakos-Postgres + pgvector:** The relational semantic vault. It holds static tables and performs high-dimensional vector math (`<=>` cosine similarity matrices) with aggressive, explicit parameters (`%s::vector`) because Postgres doesn't speak standard Python lists naturally. Splendid!
* **Brain-Python-Worker:** The center of intelligence. Intercepts incoming requests, executes a high-stakes 2-Pass Handshake with DeepSeek, calls native database tools dynamically when someone asks if Kamar 1 is empty, and handles structural errors without dropping the ball, bla bla bla bla...

## Why We Don't Crash(out)

If the database throws a violent syntax fit (like when we invented the illegal `IS NOT EXISTS NULL` phrase) or DeepSeek throws a 404 embedding tantrum, **the system does not die.** The Python worker instantly triggers a **Tier 2 Hardcoded Fallback String Defense Ring**, tricks the client into thinking everything is fine, clears the queue message, and keeps the cluster completely loop-free.

Yea, i think that's all I can say.

- N
