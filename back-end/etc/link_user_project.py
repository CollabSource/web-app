import pymongo as pm
import json
import os
import random


config_path = os.path.join(os.path.dirname(__file__), 'db_config.json')

with open(config_path, 'r') as f:
    config = json.load(f)
    CONN_STR = config['CONN_STR']


if __name__ == '__main__':

    client = pm.MongoClient(CONN_STR)
    print('connected to mongodb')
    db = client['local']
    link_table = db['user_project_relationship']

    visited_proj_ids = set()
    n = 0
    for u in db['users'].find():
        uid = u['_id']
        n += 1
        n_proj = random.randint(0,6)

        random_ids = list(db['projects'].aggregate([
            {"$match": {"_id": {"$nin": list(visited_proj_ids)}}},  # Exclude visited IDs
            {"$sample": {"size": n_proj}},
            {"$project": {"_id": 1}}  # Only include the _id field
        ]))
        

        for d in random_ids:
            new_map = {
                'user_id': uid,
                'project_id': d['_id']
            }
            visited_proj_ids.add(d['_id'])
            link_table.insert_one(new_map)
        print(f'gave {uid} {len(random_ids)}')


