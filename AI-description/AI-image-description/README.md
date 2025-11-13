# Тестовый код для ИИ-модели по расшифровке картинок
Код запускает локальный Flask веб-сервер с ИИ-моделью, которая описывает переданную картинку.

В директории memes находятся картинки для тестов (можно взять и свои).

## Запуск локального Flask веб-сервера с моделью

Из AI-image-description нужно выполнить следующие команды:

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

## Генерация описания к видео

Пример curl запроса:

```bash
curl -X POST -F "image=@<путь_к_картинке>" http://127.0.0.1:5000/image_description
```