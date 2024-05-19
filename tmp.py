from pymilvus import MilvusClient, connections, db, DataType
from sentence_transformers import SentenceTransformer

USER = "user1"
DB_NAME = USER + "_mind_palace"
PORT = 19530
HOST = "localhost"
BATCH_SIZE = 64
LIMIT = 1

conn = connections.connect(host=HOST, port=PORT)

if DB_NAME not in db.list_database():
    print(f"Database {DB_NAME} not found. Creating...")
    db.create_database(db_name=DB_NAME)

conn = connections.connect(host=HOST, port=PORT, db_name=DB_NAME)

client = MilvusClient(
    uri=f"http://localhost:{PORT}",
    db_name=DB_NAME,
)

## Addon 1

COLLECTION_NAME = "addon_1"
LLM_DIM = 384

### Schema
schema = MilvusClient.create_schema(
    auto_id=False,
    enable_dynamic_field=True,
)
schema.add_field(
    field_name="id", datatype=DataType.VARCHAR, is_primary=True, max_length=64
)
schema.add_field(field_name="vector", datatype=DataType.FLOAT_VECTOR, dim=LLM_DIM)

### index params
index_params = client.prepare_index_params()
index_params.add_index(field_name="id", index_type="INVERTED")
index_params.add_index(
    field_name="vector",
    index_type="IVF_FLAT",
    metric_type="L2",
    params={"nlist": 1536},
)

### create/load collection
if COLLECTION_NAME not in client.list_collections():
    print(f"Collection {COLLECTION_NAME} not found. Creating...")
    client.create_collection(
        collection_name=COLLECTION_NAME, schema=schema, index_params=index_params
    )

collection = client.load_collection(
    collection_name=COLLECTION_NAME,
    replica_number=1,  # standalone can be max 1 replica
)

### insert data
transformer = SentenceTransformer("all-MiniLM-L6-v2")
data = [
    {"id": "abcdefg", "input": "this is a user input text"},
    {"id": "abcdefgh", "input": "video games and pc related stuff"},
]


def convert_data(data):
    converted_data = []
    for item in data:
        new_item = {
            "id": item["id"],
            "vector": transformer.encode(item["input"]),
        }
        converted_data.append(new_item)
    return converted_data


ins = convert_data(data)

# res = client.insert(
#     collection_name=COLLECTION_NAME,
#     data=ins,
# )

### search
embeds = transformer.encode("computer")
search_data = [x for x in embeds]
res = client.search(
    collection_name=COLLECTION_NAME,
    data=[search_data],
    anns_field="vector",
    limit=LIMIT,
    output_fields=["id"],
)
print(res)
