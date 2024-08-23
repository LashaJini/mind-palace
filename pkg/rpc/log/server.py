import grpc
import time
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
    Log.info(f"Log Server started on port: {config.LOG_GRPC_SERVER_PORT}")
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
