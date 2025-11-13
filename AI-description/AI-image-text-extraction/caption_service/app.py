from flask import Flask, request, jsonify
import easyocr
from PIL import Image
import io

app = Flask(__name__)

reader = easyocr.Reader(['ru', 'en'])

@app.route("/caption", methods=["POST"])
def caption():
    if 'image' not in request.files:
        return jsonify({"error": "No image provided"}), 400

    image_file = request.files['image']
    image_bytes = image_file.read()

    result = reader.readtext(image_bytes, detail=0)

    text = " ".join(result) if result else ""

    return jsonify({"text": text})

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=5000)
