import json
import random

fixed_name = "web"
result = {
  "name": f"{fixed_name}-{random.randint(1, 99999999)}",
}

print(json.dumps(result))