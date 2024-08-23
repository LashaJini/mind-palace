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
    log.info(f"Server started on port: {config.PALACE_GRPC_SERVER_PORT}")
    # server.wait_for_termination()
    try:
        while True:
            time.sleep(config.ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


def main():
    server()


if __name__ == "__main__":
    main()
