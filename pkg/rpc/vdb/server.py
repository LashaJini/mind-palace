import grpc
import signal
from concurrent import futures

from pkg.rpc import config
from pkg.rpc.services.vdb.service import VDBService
from pkg.rpc.services.vdb.vdb import VDB
from pkg.rpc.loggers.vdb import log

import pkg.rpc.gen.VDB_pb2_grpc as vdbService

client = VDB(
    host=config.VDB_HOST,
    port=config.VDB_PORT,
    log=log,
)


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    vdbService.add_VDBServicer_to_server(VDBService(client=client), server)
    server.add_insecure_port(f"[::]:{config.VDB_GRPC_SERVER_PORT}")
    server.start()
    log.info(f"VDB gRPC Server started on port: {config.VDB_GRPC_SERVER_PORT}")

    def handle_sigterm(*_):
        log.warning("VDB gRPC server graceful shutdown initiated...")
        log.warning("VDB gRPC Server stopped. Exiting.")
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
