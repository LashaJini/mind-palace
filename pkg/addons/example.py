from transformers import pipeline

summarizer = pipeline("summarization", model="facebook/bart-large-cnn")


def Apply(data):
    summary = summarizer(data, max_length=len(data) / 3)
    return summary
