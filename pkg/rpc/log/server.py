import grpc
import signal
from concurrent import futures

from pkg.rpc.services.log.service import LogService
from pkg.rpc.services.log.log import Log
import pkg.rpc.gen.Log_pb2_grpc as logService
from pkg.rpc import config


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    logService.add_LogServicer_to_server(LogService(log=Log), server)  # type:ignore
    server.add_insecure_port(f"[::]:{config.LOG_GRPC_SERVER_PORT}")
    server.start()
    Log.info(f"Log gRPC Server started on port: {config.LOG_GRPC_SERVER_PORT}")

    def handle_sigterm(*_):
        Log.warning("Log gRPC server graceful shutdown initiated...")
        Log.warning("Log gRPC Server stopped. Exiting.")
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
