import os
from dotenv import load_dotenv

mp_env = os.environ.get("MP_ENV", "dev")
load_dotenv(f".env.{mp_env}")
MP_ENV = mp_env

VDB_HOST = os.getenv("VDB_HOST", "localhost")
VDB_PORT = int(os.getenv("VDB_PORT", 19530))
VDB_NAME = os.getenv("VDB_NAME", "_mind_palace_dev")
PALACE_GRPC_SERVER_PORT = int(os.getenv("PALACE_GRPC_SERVER_PORT", 50051))
VDB_GRPC_SERVER_PORT = int(os.getenv("VDB_GRPC_SERVER_PORT", 50052))
LOG_GRPC_SERVER_PORT = int(os.getenv("LOG_GRPC_SERVER_PORT", 50053))
ONE_DAY_IN_SECONDS = 60 * 60 * 24

# from logging package
#
# CRITICAL = 50
# FATAL = CRITICAL
# ERROR = 40
# WARNING = 30
# WARN = WARNING
# INFO = 20
# DEBUG = 10
# NOTSET = 0
LOG_LEVEL = (
    int(os.getenv("LOG_LEVEL", 0)) * 10
)  # logging package log levels are multiple of 10

LOG_FILEPATH = "logs/mindpalace-rpc-server.log"
