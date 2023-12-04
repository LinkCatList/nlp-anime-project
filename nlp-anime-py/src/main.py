import src.get as gt
from fastapi import FastAPI
from fastapi.encoders import jsonable_encoder
from fastapi import Request
import json
import requests
from googletrans import Translator

def get_result_model(prompt):
    prompt_tokenize = gt.tokenize_text(prompt)

    prompt_extra = []
    for i in prompt_tokenize:
      if i in gt.mini_vocab:
        prompt_extra.append(i)

    prompt_encode = gt.EncodeTokens(prompt_extra, gt.mini_vocab)
    prediction = gt.my_model.predict(gt.tf.constant([prompt_encode]))

    return gt.labels[gt.np.argmax(prediction)]




app = FastAPI()

@app.post('/')
async def input_request(request: Request):
    bebra = (await request.body())
    data = json.loads(bebra)
    prompt = data["query"]

    #translator = Translator()
    #if translator.detect(prompt).lang != 'en':
      #prompt = translator.translate(prompt, src='ru', dest='en').text
    ans = get_result_model(prompt)
    return {"key":ans}



#uvicorn main:app --port 8080