from flask import Flask, request, jsonify
from transformers import VisionEncoderDecoderModel, ViTImageProcessor, AutoTokenizer
from PIL import Image
import io

app = Flask(__name__)

# Загружаем модель и токенайзер
model_name = "tuman/vit-rugpt2-image-captioning"
model = VisionEncoderDecoderModel.from_pretrained(model_name)
feature_extractor = ViTImageProcessor.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

@app.route("/caption", methods=["POST"])
def caption():
    if 'image' not in request.files:
        return jsonify({"error": "No image provided"}), 400

    # Получаем изображение
    image_file = request.files['image']
    image = Image.open(image_file.stream).convert("RGB")

    # Преобразуем в формат, подходящий для модели
    pixel_values = feature_extractor(images=image, return_tensors="pt").pixel_values

    # Генерируем подпись
    output_ids = model.generate(pixel_values, max_length=50, num_beams=4)
    caption = tokenizer.decode(output_ids[0], skip_special_tokens=True)

    return jsonify({"caption": caption})

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=5000)
