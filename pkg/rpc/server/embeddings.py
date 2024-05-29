from llama_index.embeddings.huggingface import HuggingFaceEmbedding

EMBED_MODEL_NAME = "sentence-transformers/all-MiniLM-L6-v2"

embed_model = HuggingFaceEmbedding(model_name=EMBED_MODEL_NAME)
