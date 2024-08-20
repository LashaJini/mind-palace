import grpc
from typing import List
from pymilvus import MilvusClient, connections, DataType, db

from pkg.rpc.server import config, logger
from pkg.rpc.server.llm import EmbeddingModel
import pkg.rpc.server.gen.Palace_pb2_grpc as palaceService
import pkg.rpc.server.gen.Palace_pb2 as pbPalace


class MilvusInsertData:
    ids: List[str]
    inputs: List[str]

    def __init__(self, ids: List[str], inputs: List[str]):
        self.ids = ids
        self.inputs = inputs

    def __repr__(self):
        return f"MilvusInsertData({self.ids}, {self.inputs})"


class Milvus:
    db_name: str = ""

    def __init__(self, host: str, port: int):
        self.host = host
        self.port = port
        self.db_name = config.VDB_NAME

        channel = grpc.insecure_channel(f"localhost:{config.PALACE_GRPC_SERVER_PORT}")
        self.embedding_model = palaceService.EmbeddingModelStub(channel)

        connections.connect(host=self.host, port=self.port)
        if not self.db_exists():
            logger.log.info(f"Database '{config.VDB_NAME}' not found. Creating...")
            db.create_database(config.VDB_NAME)

        self.client = MilvusClient(
            uri=f"http://{self.host}:{self.port}",
            db_name=self.db_name,
        )

    def create_collection(self, collection_name: str) -> None:
        schema = self.client.create_schema(
            auto_id=False,
            enable_dynamic_field=True,
        )

        schema.add_field(
            field_name="id",
            datatype=DataType.VARCHAR,
            is_primary=True,
            max_length=64,
        )
        schema.add_field(
            field_name="vector",
            datatype=DataType.FLOAT_VECTOR,
            dim=EmbeddingModel.dimension,
        )

        index_params = self.client.prepare_index_params()
        index_params.add_index(field_name="id", index_type="INVERTED")
        index_params.add_index(
            field_name="vector",
            index_type="IVF_FLAT",
            metric_type=EmbeddingModel.metric_type,
            params={"nlist": 1536},
        )

        logger.log.info(f"Collection {collection_name} not found. Creating...")
        self.client.create_collection(
            collection_name=collection_name,
            schema=schema,
            index_params=index_params,
        )

    def collection(self, user: str) -> str:
        collection_name = f"{user}_collection"
        if not self.client.has_collection(collection_name):
            self.create_collection(collection_name)

        return collection_name

    def insert(self, user: str, data: MilvusInsertData):
        embedded_data = self._get_text_embedding(data)
        collection_name = self.collection(user)
        logger.log.info(f"Insert into {collection_name}: {data}")
        self.client.insert(collection_name=collection_name, data=embedded_data)

    def _get_text_embedding(self, data: MilvusInsertData):
        embedded_data = []

        for id, input in zip(data.ids, data.inputs):
            embeddings: pbPalace.Embeddings = self.embedding_model.CalculateEmbeddings(
                pbPalace.Text(text=input)
            )

            new_item = {
                "id": id,
                "vector": list(embeddings.embedding),
            }
            embedded_data.append(new_item)
        return embedded_data

    def drop(self):
        for collection in self.client.list_collections():
            self.client.drop_collection(collection)
            logger.log.info(f"Drop {collection} collection")

        db.drop_database(self.db_name)
        logger.log.info(f"Drop {self.db_name} database")

    def db_exists(self):
        return self.db_name in db.list_database()

    def ping(self):
        return self.db_exists()
