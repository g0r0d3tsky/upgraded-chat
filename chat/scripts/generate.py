import hashlib
import random
import uuid

import pandas as pd
from faker import Faker

FAKE = Faker(locale='ru_RU')

COUNT = 2_000_000

DATA = {
    'id': [],
    'content': [],
    'user_nickname': [],
    'time': [],
}


def generate():
    id = uuid.uuid4()
    content = FAKE.sentence()
    user_nickname = FAKE.user_name()
    time = FAKE.date_time_this_decade()

    DATA['id'].append(id)
    DATA['content'].append(content)
    DATA['user_nickname'].append(user_nickname)
    DATA['time'].append(time)


if __name__ == '__main__':
    for _ in range(COUNT):
        generate()

    df = pd.DataFrame(DATA)
    df.to_csv('user_data.csv', index=False)