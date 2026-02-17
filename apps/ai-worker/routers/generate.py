from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
import httpx

router = APIRouter()

class GenerateRequest(BaseModel):
    query: str
    context: list[str]
    max_tokens: int = 500
    
    
class GenerateResponse(BaseModel):
    answer: str
    model: str
    

@router.post("/generate", response_model=GenerateResponse)
async def generate_answer(request: GenerateRequest):
    try:
        context_text = "\n\n".join([
            f"[文書 {i+1}]\n{text}"
            for i, text in enumerate(request.context)
        ])
        
        prompt = f"""以下の文書を参考に、質問に答えてください。
        
        [参考文書]
        {context_text}
        
        [質問]
        {request.query}
        
        [回答]
        """
        
        async with httpx.AsyncClient(timeout=120.0) as client:
            response = await client.post(
                "http://localhost:11434/api/generate",
                json={
                    "model": "qwen2.5:7b",
                    "prompt": prompt,
                    "stream": False,
                    "options": {
                        "temperature": 0.7,
                        "num_predict": 200 # request.max_tokens=500
                    }
                }
            )
            
            if response.status_code != 200:
                raise HTTPException(
                    status_code=500, 
                    detail=f"Ollama returned {response.status_code}"
                )
            
            result = response.json()
            
            return GenerateResponse(
                answer=result["response"],
                model="qwen2.5:7b"
            )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
