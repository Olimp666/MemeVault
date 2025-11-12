from flask import Flask, request, jsonify
import os
import cv2
from transformers import VisionEncoderDecoderModel, ViTImageProcessor, AutoTokenizer
from PIL import Image
from deep_translator import GoogleTranslator

app = Flask(__name__)

# Настройка модели
model_name = "nlpconnect/vit-gpt2-image-captioning"
model = VisionEncoderDecoderModel.from_pretrained(model_name)
feature_extractor = ViTImageProcessor.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

def extract_keyframes(video_path, num_frames=5):
    cap = cv2.VideoCapture(video_path)
    total = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    frames = []
    for i in range(num_frames):
        idx = int((i / num_frames) * total)
        cap.set(cv2.CAP_PROP_POS_FRAMES, idx)
        ret, frame = cap.read()
        if not ret:
            break
        img = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        pil = Image.fromarray(img)
        frames.append(pil)
    cap.release()
    return frames

@app.route("/caption_video", methods=["POST"])
def caption_video():
    if 'video' not in request.files:
        return jsonify({"error": "No video provided"}), 400

    video_file = request.files['video']
    tmp_path = "temp_video.mp4"
    video_file.save(tmp_path)

    frames = extract_keyframes(tmp_path, num_frames=5)

    captions = []
    for img in frames:
        pixel_values = feature_extractor(images=img, return_tensors="pt").pixel_values
        output_ids = model.generate(pixel_values, max_length=50, num_beams=4)
        caption = tokenizer.decode(output_ids[0], skip_special_tokens=True)
        captions.append(caption)

    full_caption_en = " ".join(captions)

    # Перевод на русский
    try:
        full_caption_ru = GoogleTranslator(source='en', target='ru').translate(full_caption_en)
    except Exception as e:
        print("Ошибка перевода:", e)
        full_caption_ru = full_caption_en

    return jsonify({
        "caption_en": full_caption_en,
        "caption_ru": full_caption_ru
    })

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
