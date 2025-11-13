from flask import Flask, request, Response
from transformers import VisionEncoderDecoderModel, ViTImageProcessor, AutoTokenizer
from PIL import Image
from deep_translator import GoogleTranslator
import torch
import json

app = Flask(__name__)

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
print("Используется устройство:", device)

model_name = "nlpconnect/vit-gpt2-image-captioning"
model = VisionEncoderDecoderModel.from_pretrained(model_name).to(device)
feature_extractor = ViTImageProcessor.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

@app.route("/image_description", methods=["POST"])
def caption():
    if 'image' not in request.files:
        return Response(
            json.dumps({"error": "No image provided"}, ensure_ascii=False),
            mimetype="application/json",
            status=400
        )

    image_file = request.files['image']
    image = Image.open(image_file.stream).convert("RGB")

    pixel_values = feature_extractor(images=image, return_tensors="pt").pixel_values.to(device)

    output_ids = model.generate(
        pixel_values,
        max_length=50,
        num_beams=4,
        do_sample=False
    )
    caption_en = tokenizer.decode(output_ids[0], skip_special_tokens=True)

    try:
        caption_ru = GoogleTranslator(source='en', target='ru').translate(caption_en)
    except Exception as e:
        print("Ошибка перевода:", e)
        caption_ru = caption_en

    return Response(
        json.dumps({"caption": caption_ru}, ensure_ascii=False),
        mimetype="application/json"
    )

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
