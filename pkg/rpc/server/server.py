import grpc
import time
from concurrent import futures

from llama_index.core import KeywordTableIndex, Settings, SimpleDirectoryReader
from llama_index.core.base.embeddings.base import BaseEmbedding
from llama_index.core.node_parser import SemanticSplitterNodeParser, SentenceSplitter
from llama_index.core.extractors import KeywordExtractor
from llama_index.core.ingestion import IngestionPipeline
from llama_index.core.storage.docstore.simple_docstore import SimpleDocumentStore
from pkg.rpc.server import config
from pkg.rpc.server.llm import CustomLlamaCPP, EmbeddingModel
from pkg.rpc.server.services.mindpalace import MindPalaceService
from pkg.rpc.server.vdb import Milvus

import pkg.rpc.server.gen.Palace_pb2_grpc as grpcPalace

verbose = True
llm = CustomLlamaCPP(
    verbose=verbose,
    generate_kwargs={
        "top_k": 1,  # TODO: config
        "stop": ["<|endoftext|>", "</s>"],  # TODO: wtf
        # "seed": 4294967295,
        # "seed": -1,
    },
    # kwargs to pass to __init__()
    model_kwargs={
        "n_gpu_layers": -1,  # TODO: config
    },
)

client = Milvus(
    host=config.VDB_HOST,
    port=config.VDB_PORT,
)


def tmp():
    # Settings.llm = llm
    embed_model = EmbeddingModel()

    documents = SimpleDirectoryReader(
        input_files=["/home/jini/examples/example2.txt"]
    ).load_data()

    # after some trials and errors
    # buffer_size=2 and breakpoint_percentile_threshold=65 seems ok
    #
    # internally calls _get_text_embeddings
    splitter = SemanticSplitterNodeParser(
        buffer_size=2,
        breakpoint_percentile_threshold=65,
        embed_model=embed_model,
    )
    nodes = splitter.get_nodes_from_documents(documents)
    for node in nodes:
        embed_model.embeddings(node.get_content())


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    grpcPalace.add_PalaceServicer_to_server(
        MindPalaceService(llm=llm, client=client, verbose=verbose), server
    )
    server.add_insecure_port(f"[::]:{config.PYTHON_GRPC_SERVER_PORT}")
    server.start()
    print("Server started. Listen on port:", config.PYTHON_GRPC_SERVER_PORT)
    # server.wait_for_termination()
    try:
        while True:
            time.sleep(config.ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


def main():
    # server()
    tmp()


if __name__ == "__main__":
    main()
