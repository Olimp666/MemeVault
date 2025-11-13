from flask import Flask, request, Response
import easyocr
import json

app = Flask(__name__)

# Инициализация OCR с поддержкой русского и английского языков
reader = easyocr.Reader(['ru', 'en'])

@app.route("/image_text_extraction", methods=["POST"])
def caption():
    if 'image' not in request.files:
        return Response(
            json.dumps({"error": "No image provided"}, ensure_ascii=False),
            mimetype="application/json",
            status=400
        )

    image_file = request.files['image']
    image_bytes = image_file.read()

    # Распознаём текст
    result = reader.readtext(image_bytes, detail=0)

    text = " ".join(result) if result else ""

    response_data = {"text": text}

    return Response(
        json.dumps(response_data, ensure_ascii=False),
        mimetype="application/json"
    )

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=5000)
