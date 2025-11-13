# Тестовый код для ИИ-модели по расшифровке видео
Код запускает локальный Flask веб-сервер с ИИ-моделью, которая описывает переданное видео.

В директории memes находятся видео для тестов (можно взять и свои).

## Запуск локального Flask веб-сервера с моделью

Из AI-video-description нужно выполнить следующие команды:

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

## Генерация описания к видео

Пример curl запроса:

```bash
curl -X POST -F "video=@<путь_к_видео>" http://127.0.0.1:5000/video_description
```
