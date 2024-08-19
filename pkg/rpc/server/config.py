import os
from dotenv import load_dotenv

mp_env = os.environ.get("MP_ENV", "dev")
load_dotenv(f".env.{mp_env}")
MP_ENV = mp_env

VDB_HOST = os.getenv("VDB_HOST", "localhost")
VDB_PORT = int(os.getenv("VDB_PORT", 19530))
VDB_NAME = os.getenv("VDB_NAME", "_mind_palace_dev")
PYTHON_GRPC_SERVER_PORT = os.getenv("PYTHON_GRPC_SERVER_PORT", 50051)
ONE_DAY_IN_SECONDS = 60 * 60 * 24

LOG_FILEPATH = "logs/mindpalace-rpc-server.log"


class ServerConfig:
    def __init__(self, verbose=False, **kwargs):
        self.verbose = verbose
        self.kwargs = kwargs

    def update(self, **kwargs):
        if len(kwargs.keys()):
            self.kwargs.update(kwargs)
        else:
            self.kwargs = {}

    def __getattr__(self, name):
        if name in self.kwargs:
            return self.kwargs[name]
        return None

    def __setattr__(self, name, value):
        if name in ["verbose", "kwargs"]:
            super().__setattr__(name, value)
        else:
            self.kwargs[name] = value

    def __repr__(self):
        return f"ServerConfig(verbose={self.verbose}, kwargs={self.kwargs})"
