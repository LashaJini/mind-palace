import grpc
import time
from concurrent import futures

from pkg.rpc.server import config
from pkg.rpc.server.logger import log
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.services.mindpalace import MindPalaceService
from pkg.rpc.server.vdb import Milvus

import pkg.rpc.server.gen.Palace_pb2_grpc as grpcPalace

server_config = config.ServerConfig(verbose=True)

llm = CustomLlamaCPP(
    server_config=server_config,
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

client = Milvus(
    host=config.VDB_HOST,
    port=config.VDB_PORT,
)


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    grpcPalace.add_PalaceServicer_to_server(
        MindPalaceService(llm=llm, client=client, server_config=server_config), server
    )
    server.add_insecure_port(f"[::]:{config.PYTHON_GRPC_SERVER_PORT}")
    server.start()
    log.info(f"Server started on port: {config.PYTHON_GRPC_SERVER_PORT}")
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
