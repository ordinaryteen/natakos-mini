from pydantic import BaseModel, Field

# Matches our signed InboundInteractionContract definition exactly
class InboundInteractionRequest(BaseModel):
    message_id: str = Field(..., description="Unique technical message identifier")
    sender_jid: str = Field(..., description="The phone string identifier of the sender")
    raw_text: str = Field(..., description="The raw message content text payload")

# Matches our signed ExecutionResponseContract definition exactly
class ExecutionResponse(BaseModel):
    status: str
    intent_detected: str
    reply_message: str