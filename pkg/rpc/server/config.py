import os
from dotenv import load_dotenv

mp_env = os.environ.get("MP_ENV", "dev")
load_dotenv(f".env.{mp_env}")

VDB_HOST = os.getenv("VDB_HOST", "localhost")
VDB_PORT = int(os.getenv("VDB_PORT", 19530))
PYTHON_GRPC_SERVER_PORT = os.getenv("PYTHON_GRPC_SERVER_PORT", 50051)
ONE_DAY_IN_SECONDS = 60 * 60 * 24

LOG_FILEPATH = "logs/mindpalace-rpc-server.log"
