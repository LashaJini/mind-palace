import grpc
import time
from concurrent import futures

from pkg.rpc import config
from pkg.rpc.services.vdb.service import VDBService
from pkg.rpc.services.vdb.vdb import Milvus
from pkg.rpc.loggers.vdb import log

import pkg.rpc.gen.VDB_pb2_grpc as vdbService

client = Milvus(
    host=config.VDB_HOST,
    port=config.VDB_PORT,
    log=log,
)


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    vdbService.add_VDBServicer_to_server(VDBService(client=client), server)
    server.add_insecure_port(f"[::]:{config.VDB_GRPC_SERVER_PORT}")
    server.start()
    log.info(f"VDB Server started on port: {config.VDB_GRPC_SERVER_PORT}")
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
