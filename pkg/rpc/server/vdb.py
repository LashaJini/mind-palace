from typing import List, TypedDict
from pymilvus import MilvusClient, connections, DataType, db

from pkg.rpc.server.llm import EmbeddingModel


class MilvusInsertData(TypedDict):
    id: str
    input: str


class Milvus:
    db_name = "mindpalace_vdb"

    def __init__(self, host: str, port: int):
        self.host = host
        self.port = port
        self.db_name = Milvus.db_name
        self.embedding_model = EmbeddingModel()

        connections.connect(host=self.host, port=self.port)
        if Milvus.db_name not in db.list_database():
            print(f"Database {Milvus.db_name} not found. Creating...")
            db.create_database(Milvus.db_name)

        self.client = MilvusClient(
            uri=f"http://{self.host}:{self.port}",
            db_name=self.db_name,
        )

    def _create_collection(self, user: str):
        collection_name = f"{user}_collection"
        if not self.client.has_collection(collection_name):
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

            print(f"Collection {collection_name} not found. Creating...")
            self.client.create_collection(
                collection_name=collection_name,
                schema=schema,
                index_params=index_params,
            )
        return collection_name

    def insert(self, user: str, data: MilvusInsertData):
        embedded_data = self._get_text_embedding([data])
        collection_name = self._create_collection(user)
        self.client.insert(collection_name=collection_name, data=embedded_data)

    def _get_text_embedding(self, data: List[MilvusInsertData]):
        embedded_data = []
        for item in data:
            new_item = {
                "id": item["id"],
                "vector": self.embedding_model.embeddings(item["input"]),
            }
            embedded_data.append(new_item)
        return embedded_data
