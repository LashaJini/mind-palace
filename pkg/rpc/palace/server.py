import signal
import grpc
import time
from concurrent import futures

from pkg.rpc.services.llm.llm import CustomLlamaCPP, LLMConfig, EmbeddingModel
from pkg.rpc import config
from pkg.rpc.loggers.palace import log
from pkg.rpc.services.addon.service import AddonService
from pkg.rpc.services.llm.service import LLMService, EmbeddingModelService

import pkg.rpc.gen.Palace_pb2_grpc as grpcPalace


llm_config = LLMConfig(verbose=True)

llm = CustomLlamaCPP(
    config=llm_config,
    log=log,
    generate_kwargs={
        "top_k": 1,  # TODO: config
        "stop": ["<|endoftext|>", "</s>"],  # TODO: wtf
        # "seed": 4294967295,
        # "seed": -1,
    },
    # kwargs to pass to __init__()
    model_kwargs={
        "n_gpu_layers": -1,  # TODO: config
        "flash_attn": True,
    },
)

embedding_model = EmbeddingModel()


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    grpcPalace.add_AddonServicer_to_server(
        AddonService(llm=llm, log=log, verbose=True), server
    )
    grpcPalace.add_LLMServicer_to_server(LLMService(llm=llm), server)
    grpcPalace.add_EmbeddingModelServicer_to_server(
        EmbeddingModelService(embedding_model=embedding_model), server
    )

    server.add_insecure_port(f"[::]:{config.PALACE_GRPC_SERVER_PORT}")
    server.start()
    log.info(f"Palace gRPC Server started on port: {config.PALACE_GRPC_SERVER_PORT}")

    def handle_sigterm(*_):
        log.warning("Palace gRPC server graceful shutdown initiated...")
        log.warning("Palace gRPC Server stopped. Exiting.")
        exit(0)

    signal.signal(signal.SIGINT, handle_sigterm)  # Handle Ctrl+C
    signal.signal(signal.SIGTERM, handle_sigterm)  # Handle termination

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        handle_sigterm()


def main():
    server()


if __name__ == "__main__":
    main()
