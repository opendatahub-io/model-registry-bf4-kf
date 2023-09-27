import requests
import json
from ml_metadata.proto import metadata_store_pb2
from ml_metadata.proto import metadata_store_service_pb2


grpc_gateway = "http://localhost:8081"
    
def main():
  artifact_type_request = {
    "artifact_type": {
      "name": "DataSetRest",
      "properties": {
        "day": metadata_store_pb2.INT,
        "split": metadata_store_pb2.STRING,
      }
    }
  }
  
  resp = putArtifactType(artifact_type_request)
  data_type_id = resp["type_id"]
  
  model_type_request = {
    "artifact_type": {
      "name": "SavedModelRest",
      "properties": {
        "version": metadata_store_pb2.INT,
        "name": metadata_store_pb2.STRING,
      }
    }
  }
  
  resp = putArtifactType(model_type_request)
  model_type_id = resp["type_id"]
  
  trainer_type_request = {
    "artifact_type": {
      "name": "TrainerRest",
      "properties": {
        "state": metadata_store_pb2.STRING,
      }
    }
  }
  
  resp = putArtifactType(trainer_type_request)
  trainer_type_id = resp["type_id"]
  
  data_artifact_request = {
    "artifacts": [
      {
        "name": "Train DataSet",
        "uri": "path/to/data",
        "type_id": data_type_id,
        "properties": {
          "day": {
            "int_value": 1
          },
          "split":  {
            "string_value": "train"
          }
        }
      }
    ]
  }
  
  putArtifacts(data_artifact_request)
    
def putArtifactType(req):
  resp = requests.post(f"{grpc_gateway}/ml_metadata.MetadataStoreService/PutArtifactType", json.dumps(req))
  print(resp.json())
  return resp.json()

def putArtifacts(req):
  resp = requests.post(f"{grpc_gateway}/ml_metadata.MetadataStoreService/PutArtifacts", json.dumps(req))
  print(resp.json())
  return resp.json()



if __name__ == '__main__':
  main()
