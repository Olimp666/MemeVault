from flask import Flask, request, Response
import easyocr
import cv2
import tempfile
import os
import difflib
import json

app = Flask(__name__)

reader = easyocr.Reader(['ru', 'en'])

def is_similar(a, b, threshold=0.9):
    return difflib.SequenceMatcher(None, a, b).ratio() > threshold

@app.route("/video_text_extraction", methods=["POST"])
def video_ocr():
    if 'video' not in request.files:
        return Response(
            json.dumps({"error": "No video provided"}, ensure_ascii=False),
            mimetype="application/json",
            status=400
        )

    video_file = request.files['video']
    with tempfile.NamedTemporaryFile(delete=False, suffix=".mp4") as tmp:
        video_path = tmp.name
        video_file.save(video_path)

    cap = cv2.VideoCapture(video_path)
    if not cap.isOpened():
        os.remove(video_path)
        return Response(
            json.dumps({"error": "Failed to open video"}, ensure_ascii=False),
            mimetype="application/json",
            status=400
        )

    frame_rate = cap.get(cv2.CAP_PROP_FPS)
    frame_interval = int(frame_rate) if frame_rate > 0 else 1

    extracted_texts = []
    last_text = ""

    while True:
        ret, frame = cap.read()
        if not ret:
            break

        frame_no = int(cap.get(cv2.CAP_PROP_POS_FRAMES))
        if frame_no % frame_interval == 0:
            results = reader.readtext(frame, detail=0)
            if results:
                joined = " ".join(results).strip()

                if joined and not is_similar(joined, last_text):
                    extracted_texts.append(joined)
                    last_text = joined

    cap.release()
    os.remove(video_path)

    combined_text = "\n".join(extracted_texts)

    return Response(
        json.dumps({"text": combined_text}, ensure_ascii=False),
        mimetype="application/json"
    )

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
