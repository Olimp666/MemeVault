# Тестовый код для ИИ-модели, извлекающей текст из видео
Код запускает локальный Flask веб-сервер с ИИ-моделью, которая может извлечь текст из переданного видео.

В директории memes находятся видео для тестов (можно взять и свои).

## Запуск локального Flask веб-сервера с моделью

Из AI-video-text-extraction нужно выполнить следующие команды:

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

## Извлечение текста из видео

Пример curl запроса:

```bash
curl -X POST -F "video=@<путь_к_видео>" http://127.0.0.1:5000/video_text_extraction
```